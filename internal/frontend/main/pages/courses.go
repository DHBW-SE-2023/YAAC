package pages

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gorm.io/gorm"
)

var courseTable *fyne.Container

type courses struct {
	name *widget.Label
	date *widget.Label
}

func coursesScreen(w fyne.Window) fyne.CanvasObject {
	title := canvas.NewText(" Kursansicht", color.Black)
	title.Alignment = fyne.TextAlignLeading
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	course := &courses{
		name: widget.NewLabel(""),
		date: widget.NewLabel(""),
	}

	var dates []string

	tableHeader := container.NewGridWithColumns(2)
	tableHeader.Add(widget.NewLabel("Name"))
	tableHeader.Add(widget.NewLabel("Status"))
	courseTable = container.NewVBox()

	selection := widget.NewLabel("")

	dateDropdown := widget.NewSelectEntry(dates)
	dateDropdown.PlaceHolder = "Type or select date"
	dateDropdown.Disable()
	courseDropdown := widget.NewSelect([]string{
		"TIK22",
		"TIT22",
		"TIS22",
		"TIM22",
	}, func(s string) {
		course.name.SetText(s)
		selection.SetText(refreshSelection(course))
		refreshDateDropdown(dateDropdown, s)
		dateDropdown.SetText("")
		dateDropdown.Enable()

	})
	courseDropdown.Selected = "Kursauswahl"

	editDropdown := widget.NewSelect([]string{
		"Liste bearbeiten",
		"Liste anzeigen",
		"Liste hochladen",
	}, func(s string) {
		if s == "Liste bearbeiten" {
			LastView = w.Content()
			verifyList(w, courseDropdown.Selected, dateDropdown.Text)
		} else if s == "Liste anzeigen" {
			showImage(w, courseDropdown.Selected, dateDropdown.Text)
		} else {
			openImagUpload(w, courseDropdown.Selected, dateDropdown.Text)
		}
	})
	editDropdown.Selected = "Listenkonfiguration"
	editDropdown.Disable()
	dateDropdown.OnChanged = func(s string) {
		course.date.SetText(s)
		selection.SetText(refreshSelection(course))
		courseTable.RemoveAll()
		refreshCourseAttendancy(courseTable, course.name.Text, s)
		editDropdown.Enable()
	}

	dropdownArea := container.NewGridWrap(fyne.NewSize(200, 40), courseDropdown, dateDropdown, layout.NewSpacer(), layout.NewSpacer(), editDropdown)
	selectionArea := container.NewVBox(selection, widget.NewSeparator(), tableHeader)
	header := container.NewGridWrap(fyne.NewSize(400, 200), title)
	studentView := container.NewBorder(container.NewVBox(header, dropdownArea), nil, nil, nil, container.NewBorder(selectionArea, nil, nil, nil, container.NewVScroll(courseTable)))
	return studentView
}

func refreshSelection(course *courses) string {
	return fmt.Sprintf("%s - %s", course.name.Text, course.date.Text)
}

func refreshDateDropdown(dateDropdown *widget.SelectEntry, course string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	lists, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
	var dates []string
	for _, element := range lists {
		dates = append(dates, element.ReceivedAt.Format("2006-01-02"))
	}
	dateDropdown.SetOptions(dates)
}

func refreshCourseAttendancy(table *fyne.Container, course string, date string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	lists, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, time.Now().AddDate(0, 0, -180), time.Now())
	for _, list := range lists {
		if list.ReceivedAt.Format(("2006-01-02")) == date {
			for _, attendancies := range list.Attendancies {
				student, _ := myMVVM.Students(yaac_shared.Student{Model: gorm.Model{ID: attendancies.StudentID}})
				table.Add(NewAttendanceRow(fmt.Sprintf("%s %s", student[0].FirstName, student[0].LastName), MapBooleans(attendancies.IsAttending)))
			}
		}
	}
}

func MapBooleans(b bool) string {
	if bool(b) {
		return "Anwesend"
	}
	return "Abwesend"
}

func verifyList(w fyne.Window, course string, date string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	parsedTime, _ := time.Parse("2006-01-02", date)
	w.SetContent(verificationScreen(w, getImage(course, parsedTime.Add(24*time.Hour)), int(selectedCourse.ID), parsedTime))
}

func showImage(w fyne.Window, course string, date string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	parsedTime, _ := time.Parse("2006-01-02", date)
	list, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, parsedTime, parsedTime.Add(24*time.Hour))
	img := RotateImage(list[0].Image)
	img.FillMode = canvas.ImageFillOriginal
	customDialog := dialog.NewCustom(fmt.Sprintf("%s %s", "Listen vom", date), "Beenden", container.NewVScroll(container.NewGridWrap(fyne.NewSize(800, 1000), img)), w)
	customDialog.Resize(fyne.NewSize(800, 800))
	customDialog.Show()
}
