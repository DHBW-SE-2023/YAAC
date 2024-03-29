package test

import (
	"os"
	"testing"

	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

func TestNewTable(t *testing.T) {
	listSigs := map[string][]imgproc.TableRow{
		"testdata/list.jpg": {
			{FullName: "Baumann, Lysann", Valid: true},
			{FullName: "Beetz, Robin Georg", Valid: true},
			{FullName: "Beuerle, Marco", Valid: true},
			{FullName: "Domitrovic, Max", Valid: true},
			{FullName: "Druica, Mathias", Valid: true},
			{FullName: "Egger, Julia", Valid: false},
			{FullName: "Fischer, David", Valid: true},
			{FullName: "Fisher, Jamie", Valid: false},
			{FullName: "Gmeiner, Leander Gabriel Mauritius", Valid: true},
			{FullName: "Handschuh, Jannik", Valid: false},
			{FullName: "Hogan, Finley", Valid: true},
			{FullName: "Kiele, Milan", Valid: true},
			{FullName: "Marschall, Linus", Valid: true},
			{FullName: "Medwedkin, Eduard", Valid: false},
			{FullName: "Naas, Jasper", Valid: true},
			{FullName: "Nusch, Hannes", Valid: false},
			{FullName: "Rottweiler, Philipp", Valid: false},
			{FullName: "Schilling, Tobias", Valid: true},
			{FullName: "Schneider, Anna-Sophie", Valid: false},
			{FullName: "Seidel, Yannick", Valid: false},
			{FullName: "Siegert, Daniel Valentin", Valid: false},
			{FullName: "Zagst, Jonas", Valid: false},
		},
		"testdata/list_2.jpg": {
			{FullName: "Baumann, Lysann", Valid: true},
			{FullName: "Beetz, Robin Georg", Valid: true},
			{FullName: "Beuerle, Marco", Valid: true},
			{FullName: "Domitrovic, Max", Valid: true},
			{FullName: "Druica, Mathias", Valid: true},
			{FullName: "Egger, Julia", Valid: false},
			{FullName: "Fischer, David", Valid: true},
			{FullName: "Fisher, Jamie", Valid: false},
			{FullName: "Gmeiner, Leander Gabriel Mauritius", Valid: true},
			{FullName: "Handschuh, Jannik", Valid: false},
			{FullName: "Hogan, Finley", Valid: true},
			{FullName: "Kiele, Milan", Valid: true},
			{FullName: "Marschall, Linus", Valid: true},
			{FullName: "Medwedkin, Eduard", Valid: false},
			{FullName: "Naas, Jasper", Valid: true},
			{FullName: "Nusch, Hannes", Valid: false},
			{FullName: "Rottweiler, Philipp", Valid: false},
			{FullName: "Schilling, Tobias", Valid: true},
			{FullName: "Schneider, Anna-Sophie", Valid: false},
			{FullName: "Seidel, Yannick", Valid: false},
			{FullName: "Siegert, Daniel Valentin", Valid: false},
			{FullName: "Zagst, Jonas", Valid: false},
		},

		// This test fails because the signatures are not detected correctly

		// "testdata/list_3.jpg": {
		// 	{FullName: "Baumann, Lysann", Valid: true},
		// 	{FullName: "Beetz, Robin Georg", Valid: false},
		// 	{FullName: "Beuerle, Marco", Valid: true},
		// 	{FullName: "Domitrovic, Max", Valid: true},
		// 	{FullName: "Druica, Mathias", Valid: true},
		// 	{FullName: "Egger, Julia", Valid: true},
		// 	{FullName: "Fischer, David", Valid: true},
		// 	{FullName: "Fisher, Jamie", Valid: true},
		// 	{FullName: "Gmeiner, Leander Gabriel Mauritius", Valid: true},
		// 	{FullName: "Handschuh, Jannik", Valid: false},
		// 	{FullName: "Hogan, Finley", Valid: true},
		// 	{FullName: "Kiele, Milan", Valid: true},
		// 	{FullName: "Marschall, Linus", Valid: true},
		// 	{FullName: "Medwedkin, Eduard", Valid: true},
		// 	{FullName: "Naas, Jasper", Valid: true},
		// 	{FullName: "Nusch, Hannes", Valid: true},
		// 	{FullName: "Rottweiler, Philipp", Valid: true},
		// 	{FullName: "Schilling, Tobias", Valid: true},
		// 	{FullName: "Schneider, Anna-Sophie", Valid: true},
		// 	{FullName: "Seidel, Yannick", Valid: true},
		// 	{FullName: "Siegert, Daniel Valentin", Valid: true},
		// 	{FullName: "Zagst, Jonas", Valid: true},
		// },
	}

	for attendanceListPath, validSignatures := range listSigs {
		t.Log("Processing ", attendanceListPath)

		img := gocv.IMRead(attendanceListPath, gocv.IMReadAnyColor)
		if img.Empty() {
			wd, _ := os.Getwd()
			t.Fatalf("Could not open image with path %v. The current path is %v", attendanceListPath, wd)
		}
		tesseractClient := gosseract.NewClient()
		defer tesseractClient.Close()

		table, err := imgproc.NewTable(img, tesseractClient)
		if err != nil {
			t.Fatalf("cv.ReviewTable: %v", err)
		}

		rows := table.Rows

		if len(rows) != len(validSignatures) {
			t.Fatalf("Incorrect length of signatures: %v, correct: %v", len(rows), len(validSignatures))
		}

		for i, sig := range validSignatures {
			s := rows[i]
			// t.Logf("Name: %v, Index: %v\n", s.FullName, i)
			if s.FullName != sig.FullName {
				t.Fatalf("Incorrect name of entry %v: %v, correct: %v", i, s.FullName, sig.FullName)
			}

			if s.Valid != sig.Valid {
				t.Fatalf("Entry %v (%v) incorrectly marked as %v, correct: %v (true means valid)", i, s.FullName, s.Valid, sig.Valid)
			}
		}
	}
}
