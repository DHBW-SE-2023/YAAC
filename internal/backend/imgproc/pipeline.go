package yaac_backend_imgproc

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"

	"golang.org/x/exp/slices"
)

func FindTable(img gocv.Mat) gocv.Mat {
	origImg := img.Clone()
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 2.0, 0.0, gocv.BorderDefault)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdOtsu)
	gocv.FastNlMeansDenoisingWithParams(img, &img, 11.0, 31, 9)
	gocv.Canny(img, &img, 50.0, 150.0)

	hierachy := gocv.NewMat()
	contours := gocv.FindContoursWithParams(img, &hierachy, gocv.RetrievalExternal, gocv.ChainApproxSimple).ToPoints()

	maxArea := 0.0
	maxRect := gocv.NewPointVector()
	for _, contour := range contours {
		hull := gocv.NewMat()

		// This sorts the hull counter-clockwise
		gocv.ConvexHull(gocv.NewPointVectorFromPoints(contour), &hull, false, true)

		hullPoints := gocv.NewPointVectorFromMat(hull)

		arcLen := gocv.ArcLength(hullPoints, true)
		points := gocv.ApproxPolyDP(hullPoints, 0.001*arcLen, true)

		if points.Size() != 4 {
			continue
		}

		polyArea := gocv.ContourArea(points)
		if polyArea < maxArea {
			continue
		}

		maxArea = polyArea
		maxRect = points
	}

	lt, lb, rb, rt := maxRect.At(0), maxRect.At(1), maxRect.At(2), maxRect.At(3)

	lDiff := lt.Sub(lb)
	rDiff := rt.Sub(rb)
	bDiff := rb.Sub(lb)
	tDiff := rt.Sub(lt)

	leftHeight := math.Sqrt(float64(lDiff.Y*lDiff.Y + lDiff.X*lDiff.X))
	rightHeight := math.Sqrt(float64(rDiff.Y*rDiff.Y + rDiff.X*rDiff.X))
	bottomWidth := math.Sqrt(float64(bDiff.Y*bDiff.Y + bDiff.X*bDiff.X))
	topWidth := math.Sqrt(float64(tDiff.Y*tDiff.Y + tDiff.X*tDiff.X))

	maxHeight := int(max(leftHeight, rightHeight))
	maxWidth := int(max(bottomWidth, topWidth))

	origRect := maxRect
	destRect := gocv.NewPointVectorFromPoints([]image.Point{
		image.Pt(0, 0),
		image.Pt(maxWidth-1, 0),
		image.Pt(maxWidth-1, maxHeight-1),
		image.Pt(0, maxHeight-1),
	})

	transform := gocv.GetPerspectiveTransform(origRect, destRect)
	gocv.WarpPerspective(origImg, &img, transform, image.Pt(maxWidth, maxHeight))

	return img
}

// Expects an image which is made up of the table in question.
func ReviewTable(img gocv.Mat, tesseractClient *gosseract.Client) (Table, error) {
	// We now have the warped image, where the table is front and center
	// Now lets convert it to binary
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 2.0, 0.0, gocv.BorderDefault)
	gocv.FastNlMeansDenoisingWithParams(img, &img, 10.0, 7, 21)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdOtsu)

	k := [3][3]int8{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	// memcpy(binaryK, sizeof(k), k)
	binaryK := bytes.NewBuffer([]byte{})
	binary.Write(binaryK, binary.NativeEndian, k)

	sharpeningKernel, err := gocv.NewMatFromBytes(3, 3, gocv.MatTypeCV8S, binaryK.Bytes())
	if err != nil {
		return Table{}, err
	}

	gocv.Filter2D(img, &img, -1, sharpeningKernel, image.Pt(-1, -1), 0, gocv.BorderDefault)

	table := NewTable(img)
	img = table.Image.Clone()

	kernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))
	gocv.MorphologyEx(img, &img, gocv.MorphClose, kernel)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdBinary)
	gocv.MedianBlur(img, &img, 3)
	gocv.BitwiseNot(img, &img)

	gocv.CvtColor(img, &img, gocv.ColorGrayToBGRA)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 1.0, 0.0, gocv.BorderDefault)

	tesseractClient.SetLanguage("deu")

	table, err = StudentNames(img, table, tesseractClient)
	if err != nil {
		return Table{}, err
	}

	dyBot := 2
	dyTop := 2
	dxLeft := 2
	dxRight := 2
	for i, row := range table.Rows {
		r := image.Rect(row.SignatureROI.Min.X+dxLeft, row.SignatureROI.Min.Y+dyTop, row.SignatureROI.Max.X-dxRight, row.SignatureROI.Max.Y-dyBot)
		roi := img.Region(r)
		valid := ValidSignature(roi.Clone())
		table.Rows[i].Valid = valid
	}

	return table, nil
}

func StudentNames(img gocv.Mat, table Table, client *gosseract.Client) (Table, error) {
	dyBot := 2
	dyTop := 4
	dxLeft := 2
	dxRight := 30

	for i, row := range table.Rows {
		nameROI := image.Rect(row.NameROI.Min.X+dxLeft, row.NameROI.Min.Y+dyTop, row.NameROI.Max.X-dxRight, row.NameROI.Max.Y-dyBot)
		nameROIImg := img.Region(nameROI)

		// Tesseract accepts (among others) the following formats: PNG, JPEG, ...
		// We choose PNG, because it is lossless and it doesn't have block artifacts
		roiPng, err := gocv.IMEncode(gocv.PNGFileExt, nameROIImg)
		if err != nil {
			return Table{}, err
		}

		client.SetImageFromBytes(roiPng.GetBytes())
		name, err := client.Text()
		if err != nil {
			return Table{}, err
		}

		table.Rows[i].FullName = name

		nameParts := strings.Split(name, ",")
		if len(nameParts) != 2 {
			continue
		}

		table.Rows[i].LastName = nameParts[0]
		table.Rows[i].FirstName = nameParts[1]
	}

	return table, nil
}

func ValidSignature(img gocv.Mat) bool {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(10, 5))

	// Inplace mutation not allowed for gocv.Canny
	canny := gocv.NewMat()
	gocv.Canny(img, &canny, 128.0, 255.0)
	gocv.MorphologyExWithParams(canny, &canny, gocv.MorphClose, kernel, 3, gocv.BorderDefault)

	cnts := gocv.FindContours(canny, gocv.RetrievalExternal, gocv.ChainApproxSimple).ToPoints()

	merged := make([]image.Rectangle, 0)
	for _, c := range cnts {
		r := gocv.BoundingRect(gocv.NewPointVectorFromPoints(c))
		merged = append(merged, r)
	}

	merged = merge(merged, 0.01*float64(img.Cols()), 0.01*float64(img.Rows()))

	filteredParts := make([]image.Rectangle, 0)
	minWidth := img.Cols() / 10
	minHeight := img.Rows() / 2
	for _, r := range merged {
		if r.Dx() < minWidth || r.Dy() < minHeight {
			continue
		}

		filteredParts = append(filteredParts, r)
	}

	valid := len(filteredParts) == 1
	return valid
}

func merge(rects []image.Rectangle, deltaX float64, deltaY float64) []image.Rectangle {
	merged := make([]image.Rectangle, 0, len(rects))

	for i, r1 := range rects {
		mr := r1
		for j, r2 := range rects {
			if i == j {
				continue
			}

			dx12 := math.Abs(float64(r1.Max.X - r2.Min.X))
			dx21 := math.Abs(float64(r2.Max.X - r1.Min.X))
			dx := min(dx12, dx21)

			dy12 := math.Abs(float64(r1.Max.Y - r2.Min.Y))
			dy21 := math.Abs(float64(r2.Max.Y - r1.Min.Y))
			dy := min(dy12, dy21)

			closeEnough := dx < deltaX && dy < deltaY
			if r1.Overlaps(r2) || closeEnough {
				mr = mr.Union(r2)
			}
		}
		merged = append(merged, mr)
	}

	slices.Compact[[]image.Rectangle, image.Rectangle](merged)

	return merged
}
