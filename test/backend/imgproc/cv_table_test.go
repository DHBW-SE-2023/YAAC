package test

import (
	"bytes"
	"encoding/binary"
	"image"
	"os"
	"testing"

	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
	"github.com/otiai10/gosseract"
	"gocv.io/x/gocv"
)

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

	img = imgproc.FindTable(img)

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

	table := imgproc.NewTable(img)
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

	table, err = imgproc.StudentNames(img, table, tesseractClient)
	if err != nil {
		t.Fatalf("cv.StudentNames: %v", err)
	}

	// FIXME: ReviewTables takes the whole name column
	// StudentNames takes the whole name column
	// The first two and last two rows are uninteresting to us
	// len(correctNames) = 26
	rows := table.Rows
	rows = rows[2:24]

	names := make([]string, 0, len(rows))
	for _, name := range rows {
		names = append(names, name.Name)
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

func TestReviewTable(t *testing.T) {
	validSignatures := []imgproc.TableRow{
		{Name: "Baumann, Lysann", Valid: true},
		{Name: "Beetz, Robin Georg", Valid: true},
		{Name: "Beuerle, Marco", Valid: true},
		{Name: "Domitrovic, Max", Valid: true},
		{Name: "Druica, Mathias", Valid: true},
		{Name: "Egger, Julia", Valid: false},
		{Name: "Fischer, David", Valid: true},
		{Name: "Fisher, Jamie", Valid: false},
		{Name: "Gmeiner, Leander Gabriel Mauritius", Valid: true},
		{Name: "Handschuh, Jannik", Valid: false},
		{Name: "Hogan, Finley", Valid: true},
		{Name: "Kiele, Milan", Valid: true},
		{Name: "Marschall, Linus", Valid: true},
		{Name: "Medwedkin, Eduard", Valid: false},
		{Name: "Naas, Jasper", Valid: true},
		{Name: "Nusch, Hannes", Valid: false},
		{Name: "Rottweiler, Philipp", Valid: false},
		{Name: "Schilling, Tobias", Valid: true},
		{Name: "Schneider, Anna-Sophie", Valid: false},
		{Name: "Seidel, Yannick", Valid: false},
		{Name: "Siegert, Daniel Valentin", Valid: false},
		{Name: "Zagst, Jonas", Valid: false},
	}

	attendanceListPath := "testdata/list.jpg"
	img := gocv.IMRead(attendanceListPath, gocv.IMReadAnyColor)
	if img.Empty() {
		wd, _ := os.Getwd()
		t.Fatalf("Could not open image with path %v. The current path is %v", attendanceListPath, wd)
	}

	img = imgproc.FindTable(img)
	table, err := imgproc.ReviewTable(img)
	if err != nil {
		t.Fatalf("cv.ReviewTable: %v", err)
	}

	rows := table.Rows

	// FIXME: ReviewTables takes the whole name column
	// The first two and last two rows are uninteresting to us
	// len(correctNames) = 26
	rows = rows[2:24]

	if len(rows) != len(validSignatures) {
		t.Fatalf("Incorrect length of signatures: %v, correct: %v", len(rows), len(validSignatures))
	}

	for i, sig := range validSignatures {
		s := rows[i]
		t.Logf("Name: %v, Index: %v\n", s.Name, i)
		if s.Name != sig.Name {
			t.Fatalf("Incorrect name of entry %v: %v, correct: %v", i, s.Name, sig.Name)
		}

		if s.Valid != sig.Valid {
			t.Fatalf("Entry %v (%v) incorrectly marked as %v, correct: %v (true means valid)", i, s.Name, s.Valid, sig.Valid)
		}
	}
}
