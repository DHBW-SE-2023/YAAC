package yaac_backend_imgproc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"math"
	"regexp"
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

// Parse table shown in the image returned by `FindTable`.
// If it fails, e.g. if some information is missing, it returns an error
// It also takes in a Gosseract client to allow for reusing the Gosseract client.
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

	return table, nil
}

// Find the student names in `img` according to the rows in `table`.
// This function is called by `ReviewTable` and passes the Gosseract client
// it receives into this function.
// It returns a modified version of `table` with the names for the students filled in.
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

		name = strings.TrimSpace(name)
		table.Rows[i].FullName = name

		nameParts := strings.Split(name, ", ")
		if len(nameParts) != 2 {
			continue
		}

		table.Rows[i].LastName = nameParts[0]
		table.Rows[i].FirstName = nameParts[1]
	}

	if len(table.Rows) > 0 {
		c, err := extractCourseFromTitle(table.Rows[0].FullName)
		if err != nil {
			return Table{}, err
		}

		table.Course = c

		// Skip the first row, as we have processed it here already
		table.Rows = table.Rows[1:]
	}

	return table, nil
}

// Check whether the signature in `img` is valid.
// This is done by checking that there is only one signature.
//
// It returns true if the signature is valid, false otherwise.
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

	slices.Compact(merged)

	return merged
}

// Extract the course from the title
// The title generally looks like this "<Department> <Course>"
// The course always only consists of upper case letters and numbers
// while the department name is written normaly.
func extractCourseFromTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	re := regexp.MustCompile("^[a-zA-Z ]* ([A-Z]+[0-9]+)$")
	results := re.FindStringSubmatch(title)

	// First result should be the entire string,
	// second result is the capture gropu
	if len(results) != 2 {
		return "", errors.New("could not identify course label")
	}

	return results[1], nil
}
