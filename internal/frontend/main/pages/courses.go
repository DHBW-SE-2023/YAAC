package pages

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gorm.io/gorm"
)

type courses struct {
	name *widget.Label
	date *widget.Label
}

func coursesScreen(_ fyne.Window) fyne.CanvasObject {
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
	table := container.NewVBox()

	selection := widget.NewLabel("")

	dateDropdown := widget.NewSelect(dates, func(s string) {
		course.date.SetText(s)
		selection.SetText(refreshSelection(course))
		table.RemoveAll()
		refreshCourseAttendancy(table, course.name.Text, s)
	})

	courseDropdown := widget.NewSelect([]string{
		"TIK22",
		"TIT22",
		"TIS22",
		"TIM22",
	}, func(s string) {
		course.name.SetText(s)
		selection.SetText(refreshSelection(course))
		refreshDateDropdown(dateDropdown, s)
	})

	dropdownArea := container.NewHBox(courseDropdown, dateDropdown)
	selectionArea := container.NewVBox(selection, widget.NewSeparator(), tableHeader)
	header := container.NewGridWrap(fyne.NewSize(400, 200), title)
	studentView := container.NewBorder(container.NewVBox(header, dropdownArea), nil, nil, nil, container.NewBorder(selectionArea, nil, nil, nil, table))
	return studentView
}

func refreshSelection(course *courses) string {
	return fmt.Sprintf("%s - %s", course.name.Text, course.date.Text)
}

func refreshDateDropdown(dateDropdown *widget.Select, course string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	lists, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
	var dates []string
	for _, element := range lists {
		dates = append(dates, element.ReceivedAt.Format("2006-01-02"))
	}
	dateDropdown.Options = dates
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
