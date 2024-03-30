package yaac_backend_imgproc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"

	"golang.org/x/exp/slices"
)

// Find the attendance list in an image
// It expects a white paper on white background
// with the attendance list being black
// It performs a reverse perspective transform to
// show the table from an orthogonal view.
// This modified image is returned from the function.
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

// Review table shown in the image returned by `FindTable`.
// The image has to be prepared by `PrepareImage`.
// If it fails, e.g. if some information is missing, it returns an error
// It also takes in a Gosseract client to allow for reusing the Gosseract client.
func (table *Table) Review(tesseractClient *gosseract.Client) error {
	img := table.ImageWithoutTable.Clone()

	kernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))
	gocv.MorphologyEx(img, &img, gocv.MorphClose, kernel)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdBinary)
	gocv.MedianBlur(img, &img, 3)
	gocv.BitwiseNot(img, &img)

	gocv.CvtColor(img, &img, gocv.ColorGrayToBGRA)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 1.0, 0.0, gocv.BorderDefault)
	gocv.CvtColor(img, &img, gocv.ColorBGRAToGray)

	tesseractClient.SetLanguage("deu")

	err := table.studentNames(&img, tesseractClient)
	if err != nil {
		return err
	}

	newRows := make([]TableRow, 0, len(table.Rows))
	for _, r := range table.Rows {
		if r.FullName == "" || r.FirstName == "" || r.LastName == "" {
			continue
		}

		newRows = append(newRows, r)
	}

	// img = table.Image
	// gocv.MorphologyEx(img, &img, gocv.MorphClose, kernel)
	// gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdBinary)
	// gocv.MedianBlur(img, &img, 3)
	// gocv.BitwiseNot(img, &img)

	// gocv.CvtColor(img, &img, gocv.ColorGrayToBGRA)
	// gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 1.0, 0.0, gocv.BorderDefault)
	// gocv.CvtColor(img, &img, gocv.ColorBGRAToGray)

	table.Rows = newRows

	dyBot := 2
	dyTop := 2
	dxLeft := 2
	dxRight := 2
	for i, row := range table.Rows {
		r := image.Rect(row.SignatureROI.Min.X+dxLeft, row.SignatureROI.Min.Y+dyTop, row.SignatureROI.Max.X-dxRight, row.SignatureROI.Max.Y-dyBot)
		roi := img.Region(r)
		valid := ValidSignature(roi)
		table.Rows[i].Valid = valid
	}

	return nil
}

// Find the student names in `img` according to the rows in `table`.
// This function is called by `ReviewTable` and passes the Gosseract client
// it receives into this function.
// It returns a modified version of `table` with the names for the students filled in.
func (table *Table) studentNames(img *gocv.Mat, client *gosseract.Client) error {
	dyBot := -3
	dyTop := 0
	dxLeft := 0
	dxRight := 30

	for i, row := range table.Rows {
		nameROI := image.Rect(row.NameROI.Min.X+dxLeft, row.NameROI.Min.Y+dyTop, row.NameROI.Max.X-dxRight, row.NameROI.Max.Y-dyBot)
		nameROIImg := img.Region(nameROI)

		gocv.Line(img, row.NameROI.Min, row.NameROI.Min.Add(image.Pt(row.NameROI.Dx(), 0)), color.RGBA{255, 0, 0, 255}, 5)
		// gocv.Rectangle(img, nameROI, color.RGBA{255, 0, 0, 255}, 2)

		// Tesseract accepts (among others) the following formats: PNG, JPEG, ...
		// We choose PNG, because it is lossless and it doesn't have block artifacts
		roiPng, err := gocv.IMEncode(gocv.PNGFileExt, nameROIImg)
		if err != nil {
			return err
		}

		client.SetImageFromBytes(roiPng.GetBytes())
		name, err := client.Text()
		if err != nil {
			return err
		}

		table.Rows[i].RawName = name

		re := regexp.MustCompile(`([A-Z][a-z]+)[,.] (([A-Z][a-z\-]+ ?)+)`)
		nameParts := re.FindStringSubmatch(name)
		if nameParts == nil || len(nameParts) != 4 {
			continue
		}

		table.Rows[i].FullName = strings.TrimSpace(nameParts[0])
		table.Rows[i].LastName = strings.TrimSpace(nameParts[1])
		table.Rows[i].FirstName = strings.TrimSpace(nameParts[2])
	}

	if len(table.Rows) > 0 {
		c, err := extractCourseFromTitle(table.Rows[0].RawName)
		if err != nil {
			return err
		}

		table.Course = c

		// Skip the first row, as we have processed it here already
		table.Rows = table.Rows[1:]
	}

	return nil
}

// Check whether the signature in `img` is valid.
// This is done by checking that there is only one signature.
//
// It returns true if the signature is valid, false otherwise.
func ValidSignature(img PreparedImage) bool {
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(10, 3))

	// Inplace mutation not allowed for gocv.Canny
	canny := gocv.NewMat()
	gocv.Canny(img, &canny, 128.0, 255.0)
	gocv.MorphologyExWithParams(canny, &canny, gocv.MorphClose, kernel, 3, gocv.BorderDefault)
	gocv.Dilate(canny, &canny, kernel)

	cnts := gocv.FindContours(canny, gocv.RetrievalExternal, gocv.ChainApproxSimple).ToPoints()

	merged := make([]image.Rectangle, 0)
	for _, c := range cnts {
		r := gocv.BoundingRect(gocv.NewPointVectorFromPoints(c))
		merged = append(merged, r)
	}

	merged = merge(merged, 0.2*float64(img.Cols()), 0.2*float64(img.Rows()))
	// merged = merge(merged, 0.05*float64(img.Cols()), 0.1*float64(img.Rows()))

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

	for _, r := range filteredParts {
		gocv.Rectangle(&img, r, color.RGBA{255, 0, 0, 255}, 2)
	}

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

	slices.Compact(merged)

	return merged
}

// Extract the course from the title
// The title generally looks like this "<Department> <Course>"
// The course always only consists of upper case letters and numbers
// while the department name is written normaly.
func extractCourseFromTitle(title string) (string, error) {
	re := regexp.MustCompile("^[a-zA-Z ]* ([A-Z]+[0-9]+)")
	results := re.FindStringSubmatch(title)

	// First result should be the entire string,
	// second result is the capture gropu
	if len(results) != 2 {
		return "", errors.New("kurs konnte nicht erkannt werden")
	}

	return results[1], nil
}

func ParseTable(img PreparedImage) *Table {
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 2.0, 0.0, gocv.BorderDefault)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdOtsu)
	gocv.FastNlMeansDenoisingWithParams(img, &img, 11.0, 31, 9)

	gocv.BitwiseNot(img, &img)

	invImg := img.Clone()

	shape := img.Size()

	horizontalOpeningKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(shape[1]/10, 1))
	horizontalOpeningKernel2 := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(shape[1]/5, 1))
	horizontalClosingKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(shape[1]/5, 5))

	verticalOpeningKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(1, shape[0]/10))
	verticalOpeningKernel2 := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(1, shape[0]/5))
	verticalClosingKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, shape[0]/5))
	vhKernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))

	binaryHorizontal := gocv.NewMat()
	gocv.MorphologyExWithParams(img, &binaryHorizontal, gocv.MorphOpen, horizontalOpeningKernel, 1, gocv.BorderConstant)
	gocv.MorphologyExWithParams(binaryHorizontal, &binaryHorizontal, gocv.MorphClose, horizontalClosingKernel, 4, gocv.BorderConstant)
	gocv.MorphologyExWithParams(binaryHorizontal, &binaryHorizontal, gocv.MorphOpen, horizontalOpeningKernel2, 8, gocv.BorderConstant)

	binaryVertical := gocv.NewMat()
	gocv.MorphologyExWithParams(img, &binaryVertical, gocv.MorphOpen, verticalOpeningKernel, 1, gocv.BorderConstant)
	gocv.MorphologyExWithParams(binaryVertical, &binaryVertical, gocv.MorphClose, verticalClosingKernel, 4, gocv.BorderConstant)
	gocv.MorphologyExWithParams(binaryVertical, &binaryVertical, gocv.MorphOpen, verticalOpeningKernel2, 8, gocv.BorderConstant)

	binaryVerticalWide := binaryVertical.Clone()
	verticalClosingKernel2 := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(15, shape[0]/10))
	gocv.Dilate(binaryVertical, &binaryVerticalWide, verticalClosingKernel2)

	// Be conservative with the horizontal dilation
	binaryHorizontalWide := binaryHorizontal.Clone()
	horizontalClosingKernel2 := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(shape[1]/10, 5))
	gocv.Dilate(binaryHorizontal, &binaryHorizontalWide, horizontalClosingKernel2)
	// gocv.Dilate(binaryHorizontalWide, &binaryHorizontalWide, horizontalClosingKernel2)
	// gocv.Dilate(binaryHorizontalWide, &binaryHorizontalWide, horizontalClosingKernel2)

	vhLines := gocv.NewMat()
	gocv.AddWeighted(binaryVertical, 0.5, binaryHorizontal, 0.5, 0.0, &vhLines)

	gocv.BitwiseNot(vhLines, &vhLines)

	vhIters := 2
	gocv.ErodeWithParams(vhLines, &vhLines, vhKernel, image.Pt(-1, -1), vhIters, int(gocv.BorderConstant))

	gocv.Threshold(vhLines, &vhLines, 128.0, 255.0, gocv.ThresholdOtsu)

	contours := gocv.FindContours(vhLines, gocv.RetrievalTree, gocv.ChainApproxSimple).ToPoints()
	boundingRects := make([]image.Rectangle, len(contours))

	for _, contour := range contours {
		rect := gocv.BoundingRect(gocv.NewPointVectorFromPoints(contour))
		boundingRects = append(boundingRects, rect)
	}

	rows := gatherTableRows(boundingRects, shape)

	imageWithoutTable := invImg.Clone()
	gocv.Subtract(imageWithoutTable, binaryVerticalWide, &imageWithoutTable)
	gocv.Subtract(imageWithoutTable, binaryHorizontalWide, &imageWithoutTable)

	return &Table{
		Image:             invImg,
		ImageWithoutTable: imageWithoutTable,
		Rows:              rows,
	}
}

// imgShape: img.Rows(), img.Cols()
func gatherTableRows(boundingRects []image.Rectangle, imgShape []int) []TableRow {
	meanHeight := 0
	for _, rect := range boundingRects {
		meanHeight = meanHeight + rect.Dy()
	}
	meanHeight = meanHeight / len(boundingRects)

	sort.Slice(boundingRects, func(i, j int) bool {
		return boundingRects[i].Max.Y < boundingRects[j].Max.Y
	})

	rows := []TableRow{}

	newIdx := 0
	minDx := int(float32(imgShape[1]) * 0.04)
	minDy := int(float32(imgShape[0]) * 0.02)

	sigDx := int(float32(imgShape[1]) * 0.30)
	sigDy := int(float32(imgShape[0]) * 0.01)

	for i, rect := range boundingRects {
		if i < newIdx {
			continue
		}

		// Cell is to small, we skip it. Values were found empirically
		if rect.Dx() < minDx || rect.Dy() < minDy {
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
		if tableRow.NameROI.Dx() <= sigDx || tableRow.NameROI.Dy() <= sigDy {
			continue
		}

		if tableRow.SignatureROI.Dx() <= sigDx || tableRow.SignatureROI.Dy() <= sigDy {
			continue
		}

		rows = append(rows, tableRow)
	}

	return rows
}

func PrepareImage(img gocv.Mat) (PreparedImage, error) {
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
		return gocv.NewMat(), err
	}

	gocv.Filter2D(img, &img, -1, sharpeningKernel, image.Pt(-1, -1), 0, gocv.BorderDefault)

	return img, nil
}
