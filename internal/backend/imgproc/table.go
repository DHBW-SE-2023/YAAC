package yaac_backend_imgproc

import (
	"image"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

type Table struct {
	Course string
	Image  gocv.Mat
	Rows   []TableRow
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

// Parses `img` to return a new table
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
