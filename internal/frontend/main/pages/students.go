package yaac_frontend_pages

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gorm.io/gorm"
)

func StudentScreen(_ fyne.Window) fyne.CanvasObject {
	student := &SelectionTracker{
		courseName: widget.NewLabel(""),
		secondary:  widget.NewLabel(""),
	}
	var studentNames []string
	selection := widget.NewLabel("")

	title := ReturnHeader("Studenten")
	tableHeader, attendanceTable := ReturnAttendanceTable("Datum", "Anwesenheit")
	studentDropdown := ReturnStudentDropdown(studentNames, student, selection, attendanceTable)
	courseDropdown := ReturnCourseDropdown(student, selection, studentDropdown, "student")

	header := container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 200), title))
	dropdownArea := container.NewGridWrap(fyne.NewSize(200, 40), courseDropdown, studentDropdown)
	selectionArea := container.NewVBox(selection, tableHeader)
	studentView := container.NewBorder(container.NewVBox(header, widget.NewSeparator(), dropdownArea), nil, nil, nil, container.NewBorder(selectionArea, nil, nil, nil, attendanceTable))
	return studentView
}

/*
ReturnStudentDropdown returns the configured studentDropdown passing the studentNames list,selectionTracker,selection label and the attendance Table
*/
func ReturnStudentDropdown(studentNames []string, student *SelectionTracker, selection *widget.Label, attendanceTable *fyne.Container) *widget.SelectEntry {
	studentDropdown := widget.NewSelectEntry(studentNames)
	studentDropdown.OnChanged = func(s string) {
		if len(s) > 30 {
			s = s[0:30]
		}
		re := regexp.MustCompile(`[^a-zA-Z-ä-ö-ü\s]`)
		s = re.ReplaceAllString(s, "")
		studentDropdown.SetText(s)
		attendanceTable.RemoveAll()
		RefreshStudentAttendancyList(attendanceTable, student.courseName.Text, s, student)
		selection.SetText(RefreshSelection(student))
	}
	studentDropdown.PlaceHolder = "Type or select student"
	studentDropdown.Disable()
	return studentDropdown
}

/*
RefreshStudentDropdown is responsible for refreshing the studentDropdown list as soon as a course has been selected
*/
func RefreshStudentDropdown(studentDropdown *widget.SelectEntry, course string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	students, _ := myMVVM.CourseStudents(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}})
	if len(students) != 0 {
		var studentNames []string
		for _, studElement := range students {
			studentNames = append(studentNames, fmt.Sprintf("%s %s", studElement.FirstName, studElement.LastName))
		}
		studentDropdown.SetOptions(studentNames)
	}
}

/*
RefreshStudentAttendancyList is responsible for refreshing the attedanceTable list as soon as a course and a student has been selected
*/
func RefreshStudentAttendancyList(table *fyne.Container, course string, student string, selection *SelectionTracker) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	if strings.Contains(student, " ") {
		selectedStudent, _ := myMVVM.Students(yaac_shared.Student{LastName: strings.Split(student, " ")[1]})
		if len(selectedStudent) != 0 {
			selection.secondary.SetText(student)
			attendances, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
			for _, attendance := range attendances {
				for _, attendancy := range attendance.Attendancies {
					if attendancy.StudentID == selectedStudent[0].ID {
						table.Add(NewAttendanceRow(attendance.ReceivedAt.Format("2006-01-02"), MapBooleans(attendancy.IsAttending)))
					}
				}
			}
		}
	}
}
