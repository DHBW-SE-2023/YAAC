package yaac_backend_imgproc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
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
	img := table.Image.Clone()

	kernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))
	gocv.MorphologyEx(img, &img, gocv.MorphClose, kernel)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdBinary)
	gocv.MedianBlur(img, &img, 3)
	gocv.BitwiseNot(img, &img)

	gocv.CvtColor(img, &img, gocv.ColorGrayToBGRA)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 1.0, 0.0, gocv.BorderDefault)

	tesseractClient.SetLanguage("deu")

	err := table.studentNames(img, tesseractClient)
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
	table.Rows = newRows

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

	return nil
}

// Find the student names in `img` according to the rows in `table`.
// This function is called by `ReviewTable` and passes the Gosseract client
// it receives into this function.
// It returns a modified version of `table` with the names for the students filled in.
func (table *Table) studentNames(img gocv.Mat, client *gosseract.Client) error {
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
			return err
		}

		client.SetImageFromBytes(roiPng.GetBytes())
		name, err := client.Text()
		if err != nil {
			return err
		}

		table.Rows[i].FullName = name

		nameParts := strings.Split(name, ",")
		if len(nameParts) != 2 {
			continue
		}

		table.Rows[i].LastName = nameParts[0]
		table.Rows[i].FirstName = nameParts[1]
	}

	if len(table.Rows) > 0 {
		c, err := extractCourseFromTitle(table.Rows[0].FullName)
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

	slices.Compact(merged)

	return merged
}

// Extract the course from the title
// The title generally looks like this "<Department> <Course>"
// The course always only consists of upper case letters and numbers
// while the department name is written normaly.
func extractCourseFromTitle(title string) (string, error) {
	re := regexp.MustCompile("^[a-zA-Z ]* ([A-Z]+[0-9]+)$")
	results := re.FindStringSubmatch(title)

	// First result should be the entire string,
	// second result is the capture gropu
	if len(results) != 2 {
		return "", errors.New("could not identify course label")
	}

	return results[0], nil
}

func ParseTable(img PreparedImage) *Table {
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 2.0, 0.0, gocv.BorderDefault)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdOtsu)
	gocv.FastNlMeansDenoisingWithParams(img, &img, 11.0, 31, 9)

	gocv.BitwiseNot(img, &img)

	invImg := img.Clone()

	shape := img.Size()

	horizontalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(shape[0]/50, 1))
	verticalKernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(1, shape[1]/35))
	vhKernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))

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

	return &Table{
		Image: invImg,
		Rows:  rows,
	}
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
