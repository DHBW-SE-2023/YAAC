package pages

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func landingScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(yaac_shared.ResourceIconPng)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(256, 256))
	button := widget.NewButton("Insert List", func() {
		testInsertAttendanceList()
	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Willkommen zurück XD", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewHBox(
			widget.NewHyperlink("User Dokumentation", parseURL("https://developer.fyne.io/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("Sponsoren", parseURL("https://www.sap.com/germany/index.html?url_id=auto_hp_redirect_germany")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
		button,
	))
}

func testInsertAttendanceList() {
	imageFilePath := "assets/list.jpg"
	var testTime = time.Now()
	// Read the image file
	imageBytes, err := ioutil.ReadFile(imageFilePath)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return
	}
	attendanceList := yaac_shared.AttendanceList{
		ReceivedAt: testTime,
		CourseID:   2,
		Image:      imageBytes,
	}
	myMVVM.InsertList(attendanceList)
}
