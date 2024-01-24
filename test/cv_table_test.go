package test

import (
	"bytes"
	"encoding/binary"
	"image"
	"os"
	"testing"

	"github.com/DHBW-SE-2023/YAAC/internal/cv"
	"github.com/otiai10/gosseract"
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

func TestStudentNameRecognition(t *testing.T) {
	correctNames := []string{
		"Baumann, Lysann",
		"Beetz, Robin Georg",
		"Beuerle, Marco",
		"Domitrovic, Max",
		"Druica, Mathias",
		"Egger, Julia",
		"Fischer, David",
		"Fisher, Jamie",
		"Gmeiner, Leander Gabriel Mauritius",
		"Handschuh, Jannik",
		"Hogan, Finley",
		"Kiele, Milan",
		"Marschall, Linus",
		"Medwedkin, Eduard",
		"Naas, Jasper",
		"Nusch, Hannes",
		"Rottweiler, Philipp",
		"Schilling, Tobias",
		"Schneider, Anna-Sophie",
		"Seidel, Yannick",
		"Siegert, Daniel Valentin",
		"Zagst, Jonas",
	}

	attendanceListPath := "testdata/list.jpg"
	img := gocv.IMRead(attendanceListPath, gocv.IMReadAnyColor)
	if img.Empty() {
		wd, _ := os.Getwd()
		t.Fatalf("Could not open image with path %v. The current path is %v", attendanceListPath, wd)
	}

	img = cv.FindTable(img)

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
		t.Fatalf("gocv.NewMatFromBytes: %v", err)
	}

	gocv.Filter2D(img, &img, -1, sharpeningKernel, image.Pt(-1, -1), 0, gocv.BorderDefault)

	table := cv.NewTable(img)
	img = table.Image.Clone()

	kernel := gocv.GetStructuringElement(gocv.MorphCross, image.Pt(3, 3))
	gocv.MorphologyEx(img, &img, gocv.MorphClose, kernel)
	gocv.Threshold(img, &img, 128.0, 255.0, gocv.ThresholdBinary)
	gocv.MedianBlur(img, &img, 3)
	gocv.BitwiseNot(img, &img)

	gocv.CvtColor(img, &img, gocv.ColorGrayToBGRA)
	gocv.GaussianBlur(img, &img, image.Point{X: 3, Y: 3}, 1.0, 0.0, gocv.BorderDefault)

	tesseractClient := gosseract.NewClient()
	defer tesseractClient.Close()
	tesseractClient.SetLanguage("deu")

	namesROI, err := cv.StudentNames(img, table, tesseractClient)
	if err != nil {
		t.Fatalf("cv.StudentNames: %v", err)
	}

	// StudentNames takes the whole name column
	// The first two and last two rows are uninteresting to us
	// len(correctNames) = 26
	namesROI = namesROI[2:24]

	names := make([]string, 0, len(namesROI))
	for _, name := range namesROI {
		names = append(names, name.Name())
	}

	t.Logf("Recognised names: %v", names)

	if len(correctNames) != len(names) {
		t.Fatalf("Name list has incorrect length. Length %v, Correct length: %v", len(correctNames), len(names))
	}

	for i, name := range names {
		if name != correctNames[i] {
			t.Fatalf("Name not equal: \"%v\" != \"%v\"", name, correctNames[i])
		}
	}
}
