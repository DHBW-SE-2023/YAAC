package test

import (
	"os"
	"testing"

	"github.com/DHBW-SE-2023/YAAC/internal/cv"
	"gocv.io/x/gocv"
)

func TestTableColumnCount(t *testing.T) {
	attendanceListPath := "testdata/list.jpg"
	img := gocv.IMRead(attendanceListPath, gocv.IMReadAnyColor)
	if img.Empty() {
		wd, _ := os.Getwd()
		t.Fatalf("Could not open image with path %v. The current path is %v", attendanceListPath, wd)
	}

	img = cv.FindTable(img)
	table := cv.NewTable(img)

	for idx, row := range table.Rows {
		// A standard attendance list always has 3 columns
		l := len(row)
		if l != 3 {
			t.Fatalf("The row %v is expected to have 3 columns, but it has %v", idx, l)
		}
	}

}
