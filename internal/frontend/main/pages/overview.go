package yaac_frontend_pages

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gorm.io/gorm"
)

var overviewGrid *fyne.Container

func OverviewScreen(w fyne.Window) fyne.CanvasObject {
	title := ReturnHeader("Anwesenheitsliste der Kurse - Heute")
	buttonImageContainer := ReturnVerifyImageContainer(w)
	buttonRefreshContainer := ReturnMailRefreshContainer(w)
	overviewGrid = container.NewGridWrap(fyne.NewSize(250, 250))
	LoadOverviewWidgets(w, overviewGrid)

	// header := container.NewVBox(container.NewBorder(nil, nil, nil, container.NewPadded(container.NewPadded(container.NewPadded(container.NewPadded(buttonImageContainer)))), container.NewGridWrap(fyne.NewSize(400, 200), title)), widget.NewSeparator())
	header := container.NewVBox(container.NewBorder(nil, nil, nil, container.NewPadded(container.NewPadded(container.NewPadded(container.NewGridWithColumns(2, container.NewPadded(buttonImageContainer), buttonRefreshContainer)))), container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 200), title))), widget.NewSeparator())
	return container.NewBorder(header, nil, nil, nil, container.NewVScroll(overviewGrid))
}

/*
ReturnVerifyImageContainer returns the buttonImageContaier containing the image for insertList Button.
*/
func ReturnVerifyImageContainer(w fyne.Window) *tappableImage {
	image := canvas.NewImageFromResource(yaac_shared.ResourceImageUploadPng)
	buttonImageContainer := newTappableImage(image, func() {
		OpenImageUpload(w)
	})
	return buttonImageContainer
}

/*
ReturnVerifyImageContainer returns the buttonImageContaier containing the image for insertList Button.
*/
func ReturnMailRefreshContainer(w fyne.Window) *tappableImage {
	image := canvas.NewImageFromResource(yaac_shared.ResourceRefreshPng)
	buttonRefreshContainer := newTappableImage(image, func() {
		myMVVM.SingleDemonRunthrough()
	})
	return buttonRefreshContainer
}

/*
LoadOverviewWidgets loads all OverviewWidgets for each course and adds them to the overviewGrid
*/
func LoadOverviewWidgets(w fyne.Window, grid *fyne.Container) {
	grid.RemoveAll()
	var frameColor color.NRGBA
	var hidden bool
	courses, _ := myMVVM.Courses()
	for _, element := range courses {
		var students []string
		var totalStudents int
		lists, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: element.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
		if len(lists) > 0 {
			if len(lists[0].Attendancies) > 0 {
				students = ReturnNonAttending(lists[0].Attendancies)
				totalStudents = len(lists[0].Attendancies)
				if len(students) > 0 {
					frameColor = color.NRGBA{227, 0, 27, 255}
					hidden = false
				} else {
					hidden = false
					frameColor = color.NRGBA{R: 209, G: 209, B: 209, A: 255}
				}
			} else {
				students = append(students, "Keine Anwesenheiten")
				totalStudents = 0
				frameColor = color.NRGBA{227, 0, 27, 255}
				hidden = false
			}
		} else {
			students = append(students, "Kein Listeingang")
			totalStudents = 0
			frameColor = color.NRGBA{241, 230, 60, 200}
			hidden = true
		}
		widget := NewOverviewWidget(w, element.Name, int(element.ID), students, totalStudents)
		widget.frame.FillColor = frameColor
		widget.button.Hidden = hidden
		grid.Add(widget)
	}
}
