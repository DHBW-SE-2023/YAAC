package yaac_backend_imgproc

import (
	"image"
	"sort"

	"gocv.io/x/gocv"
)

type Table struct {
	Image gocv.Mat
	Rows  []TableRow
}

type TableRow struct {
	FirstName    string
	LastName     string
	FullName     string
	Valid        bool
	NameROI      image.Rectangle
	SignatureROI image.Rectangle
	TotalROI     image.Rectangle
}

// Expects a grayscale image
func NewTable(img gocv.Mat) Table {
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 2.0, 0.0, gocv.BorderDefault)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdOtsu)
	gocv.FastNlMeansDenoisingWithParams(img, &img, 11.0, 31, 9)

	gocv.BitwiseNot(img, &img)

	invImg := img.Clone()

	shape := img.Size()

	// TODO: This needs some more work, why 50 and 35?
	horizontalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(shape[0]/50, 1))
	verticalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(1, shape[1]/35))
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

	rows := []TableRow{}

	newIdx := 0
	for i, rect := range boundingRects {
		if i < newIdx {
			continue
		}

		// Cell is to small, we skip it. Values were found empirically
		if rect.Dx() < int(float32(img.Cols())*0.04) || rect.Dy() < int(float32(img.Rows())*0.02) {
			continue
		}

		// We'll probably never have more than 10 rects detected in one line
		row := make([]image.Rectangle, 0, 10)
		currentMaxHeight := rect.Min.Y + meanHeight/2

		minX, minY, maxX, maxY := 0, 0, 0, 0
		for _, next := range boundingRects[i:] {
			if next.Min.Y > currentMaxHeight {
				break
			}

			row = append(row, next)

			minX = min(minX, next.Min.X)
			minY = min(minY, next.Min.Y)
			maxX = max(maxX, next.Max.X)
			maxY = max(maxY, next.Max.Y)
		}

		newIdx = i + len(row)

		sort.Slice(row, func(i, j int) bool {
			return row[i].Min.X < row[j].Min.X
		})

		// FIXME: Should we fail if we have a misformed row?
		// We need the columns num, name, signature
		// If not, we assume the column is malformed and skip it
		if len(row) != 3 {
			continue
		}

		tableRow := TableRow{
			// IndexROI: row[0], but it is unused
			NameROI:      row[1],
			SignatureROI: row[2],
			TotalROI:     image.Rect(minX, minY, maxX, maxY),
		}

		// The name and signature column have at least a width of 30%
		if tableRow.NameROI.Dx() <= int(0.30*float32(shape[0])) || tableRow.NameROI.Dy() <= int(0.01*float32(shape[1])) {
			continue
		}

		if tableRow.SignatureROI.Dx() <= int(0.30*float32(shape[0])) || tableRow.SignatureROI.Dy() <= int(0.01*float32(shape[1])) {
			continue
		}

		rows = append(rows, tableRow)
	}

	return Table{
		Image: invImg,
		Rows:  rows,
	}
}
