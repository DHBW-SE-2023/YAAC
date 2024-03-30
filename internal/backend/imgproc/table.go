package yaac_backend_imgproc

import (
	"image"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

type Table struct {
	Course            string
	Image             gocv.Mat
	ImageWithoutTable gocv.Mat
	Rows              []TableRow
}

type TableRow struct {
	FirstName    string
	LastName     string
	RawName      string
	FullName     string
	Valid        bool
	NameROI      image.Rectangle
	SignatureROI image.Rectangle
	TotalROI     image.Rectangle
}

// An image that has been prepared with to PrepareImage
type PreparedImage = gocv.Mat

// Parses `img` to return a new table
//
// This extracts the students names from the image and validates the signatures.
func NewTable(img gocv.Mat, client *gosseract.Client) (*Table, error) {
	img = FindTable(img)
	warpedImg := img.Clone()

	// We now have the warped image, where the table is front and center
	// Now lets convert it to binary
	img, err := PrepareImage(img)
	if err != nil {
		return nil, err
	}

	table := ParseTable(img)
	if err := table.Review(client); err != nil {
		return nil, err
	}

	table.Image = warpedImg

	return table, nil
}
