package cv

import (
	"image"
	"sort"

	"gocv.io/x/gocv"
)

type Table struct {
	Image gocv.Mat
	Rows  [][]image.Rectangle
}

func NewTable(img gocv.Mat) Table {
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 2.0, 0.0, gocv.BorderDefault)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdOtsu)
	gocv.FastNlMeansDenoisingWithParams(img, &img, 11.0, 31, 9)

	gocv.BitwiseNot(img, &img)

	invImg := img.Clone()

	size := img.Size()

	// TODO: This needs some more work, why 50 and 35?
	horizontalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(size[0]/50, 1))
	verticalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(1, size[1]/35))
	vhKernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))

	// TODO: Tweek this number or it takes longer than necessary
	iters := 8
	binaryHorizontal := gocv.NewMat()
	gocv.MorphologyExWithParams(img, &binaryHorizontal, gocv.MorphOpen, horizontalKernel, iters, gocv.BorderConstant)

	binaryVertical := gocv.NewMat()
	gocv.MorphologyExWithParams(img, &binaryVertical, gocv.MorphOpen, verticalKernel, iters, gocv.BorderConstant)

	vhLines := gocv.NewMat()
	gocv.AddWeighted(binaryVertical, 0.5, binaryHorizontal, 0.5, 0.0, &vhLines)

	gocv.BitwiseNot(vhLines, &vhLines)

	vhIters := 2
	gocv.ErodeWithParams(vhLines, &vhLines, vhKernel, image.Pt(-1, -1), vhIters, int(gocv.BorderConstant))

	gocv.Threshold(vhLines, &vhLines, 128.0, 255.0, gocv.ThresholdOtsu)

	contours := gocv.FindContours(vhLines, gocv.RetrievalTree, gocv.ChainApproxSimple).ToPoints()

	boundingRects := make([]image.Rectangle, len(contours))
	meanHeight := 0

	for _, contour := range contours {
		rect := gocv.BoundingRect(gocv.NewPointVectorFromPoints(contour))
		boundingRects = append(boundingRects, rect)
		meanHeight = meanHeight + rect.Dy()
	}

	meanHeight = meanHeight / len(contours)

	sort.Slice(boundingRects, func(i, j int) bool {
		return boundingRects[i].Max.Y < boundingRects[j].Max.Y
	})

	rows := [][]image.Rectangle{}

	newIdx := 0
	for i, rect := range boundingRects {
		if i < newIdx {
			continue
		}

		if rect.Dx() == 0 || rect.Dy() == 0 {
			continue
		}

		// We'll probably never have more than 10 rects detected in one line
		row := make([]image.Rectangle, 0, 10)
		currentMaxHeight := rect.Min.Y + meanHeight/2

		for _, next := range boundingRects[i:] {
			if next.Min.Y > currentMaxHeight {
				break
			}

			row = append(row, next)
		}

		newIdx = i + len(row)

		sort.Slice(row, func(i, j int) bool {
			return row[i].Min.X < row[j].Min.X
		})

		rows = append(rows, row)
	}

	return Table{
		Image: invImg,
		Rows:  rows,
	}
}
