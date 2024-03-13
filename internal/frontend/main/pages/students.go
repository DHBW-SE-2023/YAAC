package pages

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gorm.io/gorm"
)

type students struct {
	name   *widget.Label
	course *widget.Label
}

func studentScreen(_ fyne.Window) fyne.CanvasObject {
	title := canvas.NewText(" Studenten", color.Black)
	title.Alignment = fyne.TextAlignLeading
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	student := &students{
		name:   widget.NewLabel(""),
		course: widget.NewLabel(""),
	}
	var studentNames []string
	tableHeader := container.NewGridWithColumns(2)
	tableHeader.Add(widget.NewLabel("Datum"))
	tableHeader.Add(widget.NewLabel("Anwesenheit"))
	table := container.NewVBox()

	selection := widget.NewLabel("")

	studentDropdown := widget.NewSelect(studentNames, func(s string) {
		student.name.SetText(s)
		selection.SetText(updateSelection(student))
		table.RemoveAll()
		refreshStudentAttendancyList(table, student.course.Text, s)
	})
	studentDropdown.Disable()
	courseDropdown := widget.NewSelect([]string{
		"TIK22",
		"TIT22",
		"TIS22",
		"TIM22",
	}, func(s string) {
		student.course.SetText(s)
		selection.SetText(updateSelection(student))
		refreshStudentDropdown(studentDropdown, s)
		studentDropdown.Enable()
	})
	courseDropdown.Selected = "Kursauswahl"

	dropdownArea := container.NewHBox(courseDropdown, studentDropdown)
	selectionArea := container.NewVBox(selection, widget.NewSeparator(), tableHeader)
	header := container.NewGridWrap(fyne.NewSize(400, 200), title)
	studentView := container.NewBorder(container.NewVBox(header, dropdownArea), nil, nil, nil, container.NewBorder(selectionArea, nil, nil, nil, table))
	return studentView
}

func updateSelection(student *students) string {
	return fmt.Sprintf("%s - %s", student.course.Text, student.name.Text)
}

func refreshStudentDropdown(studentDropdown *widget.Select, course string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	students, _ := myMVVM.CourseStudents(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}})
	var studentNames []string
	for _, studElement := range students {
		studentNames = append(studentNames, fmt.Sprintf("%s %s", studElement.FirstName, studElement.LastName))
	}
	studentDropdown.Options = studentNames
}

func refreshStudentAttendancyList(table *fyne.Container, course string, student string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	selectedStudent, _ := myMVVM.Students(yaac_shared.Student{LastName: strings.Split(student, " ")[1]})
	attendances, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
	for _, attendance := range attendances {
		for _, attendancy := range attendance.Attendancies {
			if attendancy.StudentID == selectedStudent[0].ID {
				table.Add(NewAttendanceRow(attendance.ReceivedAt.Format("2006-01-02"), MapBooleans(attendancy.IsAttending)))
			}
		}
	}
}

type attendanceRow struct {
	widget.BaseWidget
	frame    *canvas.Rectangle
	date     *widget.Label
	state    *widget.Label
	content  *fyne.Container
	OnTapped func()
}

func NewAttendanceRow(dateText string, stateText string) *attendanceRow {
	item := &attendanceRow{
		frame: &canvas.Rectangle{
			FillColor:    color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeColor:  color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeWidth:  4.0,
			CornerRadius: 10,
		},
		date:    widget.NewLabel(dateText),
		state:   widget.NewLabel(stateText),
		content: container.NewGridWithColumns(2),
		OnTapped: func() {
			yaac_shared.App.SendNotification(fyne.NewNotification("Weiterleitung", fmt.Sprintf("%s %s", dateText, stateText)))
		},
	}
	item.ExtendBaseWidget(item)
	detectAttendancyState(item)
	return item
}
func (item *attendanceRow) Tapped(ev *fyne.PointEvent) { item.OnTapped() }

func (item *attendanceRow) CreateRenderer() fyne.WidgetRenderer {
	item.frame.Resize(fyne.NewSize(250, 250))
	item.content.Add(item.date)
	item.content.Add(item.state)
	c := container.NewPadded(
		item.frame,
		item.content,
	)
	return widget.NewSimpleRenderer(c)
}

func detectAttendancyState(item *attendanceRow) {
	if item.state.Text != "Anwesend" {
		item.frame.FillColor = color.RGBA{227, 0, 27, 255}
		item.frame.StrokeColor = color.RGBA{227, 0, 27, 255}
		item.frame.StrokeWidth = 2.0
	}
}
