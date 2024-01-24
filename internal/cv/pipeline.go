package cv

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"math"

	"github.com/otiai10/gosseract"
	"gocv.io/x/gocv"
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

type ValidNames struct {
	string
	bool
}

// Expects an image which is made up of the table in question.
func ReviewTable(img gocv.Mat) ([]ValidNames, error) {
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
		return nil, err
	}

	gocv.Filter2D(img, &img, -1, sharpeningKernel, image.Pt(-1, -1), 0, gocv.BorderDefault)

	table := NewTable(img)
	shape := table.Image.Size()

	scaling := 2

	// Shove table.Image into img for ease of use
	gocv.Resize(table.Image, &img, image.Pt(scaling*shape[0], scaling*shape[1]), 0.0, 0.0, gocv.InterpolationLinear)

	kernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))
	gocv.MorphologyEx(img, &img, gocv.MorphClose, kernel)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdBinary)
	gocv.MedianBlur(img, &img, 3)

	tesseractClient := gosseract.NewClient()
	defer tesseractClient.Close()

	tesseractClient.SetLanguage("deu")

	dyBot := 2
	dyTop := 4
	dxLeft := 1
	dxRight := 10

	gocv.CvtColor(table.Image, &img, gocv.ColorGrayToBGRA)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 1.0, 0.0, gocv.BorderDefault)

	for _, row := range table.Rows {
		for _, col := range row {
			gocv.Rectangle(&img, col, color.RGBA{255, 0, 0, 255}, 1)
		}
	}

	nameCells := make([]image.Rectangle, 0, len(table.Rows))
	namesigs := make([]struct {
		name   string
		sigROI image.Rectangle
	}, 0, len(nameCells))
	for _, row := range table.Rows {
		// If len(row) > 10, then we assume the row is deformed
		if len(row) > 10 {
			continue
		}

		// We need the columns num, name, signature
		if len(row) < 2 {
			continue
		}

		var nameCell image.Rectangle
		nextIdx := 0
		for _, r := range row {
			// The name and signature column have at least a width of 30%
			nextIdx = nextIdx + 1
			if r.Dx() <= int(0.30*float32(shape[0])) || r.Dy() <= int(0.01*float32(shape[1])) {
				continue
			}

			nameCell = r
			break
		}

		// We need at least two columns left
		if nextIdx > len(row)-1 {
			nextIdx = 0
			continue
		}

		sigCell := row[nextIdx]
		roi := image.Rect(nameCell.Min.X+dxLeft, nameCell.Min.Y+dyTop, nameCell.Max.X-dxRight, nameCell.Max.Y-dyBot)
		sigROI := image.Rect(sigCell.Min.X+dxLeft, sigCell.Min.Y+dyTop, sigCell.Max.X-dxRight, sigCell.Max.Y-dyBot)

		// Show the rectangles of the students names on the image
		// gocv.Rectangle(&img, roi, color.RGBA{255, 0, 0, 255}, 1)

		// Tesseract accepts (among others) the following formats: PNG, JPEG, ...
		// We choose PNG, because it is lossless and it doesn't have block artifacts
		roiPng, err := gocv.IMEncode(gocv.PNGFileExt, img.Region(roi))
		if err != nil {
			return nil, err
		}

		tesseractClient.SetImageFromBytes(roiPng.GetBytes())
		text, err := tesseractClient.Text()
		if err != nil {
			return nil, err
		}

		namesigs = append(namesigs, struct {
			name   string
			sigROI image.Rectangle
		}{text, sigROI})
	}

	results := make([]ValidNames, 0)

	return results, nil
}
