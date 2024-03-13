package pages

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
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
	logo := canvas.NewImageFromFile("assets/Icon.png")
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(256, 256))
	button := widget.NewButton("Insert List", func() {
		testInsertAttendanceList()
	})

	button2 := widget.NewButton("Bulk Insert Students TIK", func() {
		testInsertStudent()
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
		button2,
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
		CourseID:   1,
		Image:      imageBytes,
	}
	myMVVM.InsertList(attendanceList)
}

func testInsertStudent() {
	students := []string{"Robin Beats", "Phillip Rotweiler", "Finley Hogan", "Julia Egger", "Yannick Seidel", "David Fischer", "Linus Marshall", "Hannes Nusch", "Edward Medwedkin", "Marco Beuerle", "Lysann Baumann", "Milan Kiele", "Tobias y"}
	for _, element := range students {
		myMVVM.InsertStudent(yaac_shared.Student{FirstName: strings.Split(element, " ")[0], LastName: strings.Split(element, " ")[1], IsImmatriculated: true, CourseID: 1})
	}
}
