package yaac_frontend_pages

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gorm.io/gorm"
)

func StudentScreen(w fyne.Window) fyne.CanvasObject {
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
	studentInsertButton := ReturnStudentInsertButton(w)
	header := container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 200), title))
	dropdownArea := container.NewGridWrap(fyne.NewSize(200, 40), courseDropdown, studentDropdown, studentInsertButton)
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
ReturnStudentInsertButton returns the configured studentInsertButton
*/
func ReturnStudentInsertButton(w fyne.Window) *widget.Button {
	insertButton := widget.NewButton("Student hinzufügen", func() {
		DisplayInsertStudentDialog(w)
	})
	return insertButton
}

/*
ReturnStudentInsertButton returns the configured courseInsertButton
*/
func DisplayInsertStudentDialog(w fyne.Window) {
	content := container.NewVBox()
	form := dialog.NewCustomWithoutButtons("Student hinzufügen", content, w)
	studentEntry := widget.NewEntry()
	studentEntry.PlaceHolder = "Studentenname"
	courseEntry := widget.NewEntry()
	courseEntry.PlaceHolder = "Kursname"
	isImmatriculated := widget.NewCheck("Ist der Student immatrikuliert?", nil)
	exitButton := widget.NewButton("Zurück", func() {
		form.Hide()
	})
	confirmButton := widget.NewButton("Bestätigen", InsertStudent(w, form, studentEntry, courseEntry, isImmatriculated))
	confirmButton.Disable()
	ValidateInput(studentEntry, courseEntry, confirmButton)
	content.Add(studentEntry)
	content.Add(courseEntry)
	content.Add(isImmatriculated)
	content.Add(container.NewGridWithColumns(2, exitButton, confirmButton))
	form.Show()
}

/*
ValidateInput takes the two entries studentEntry, courseEntry as input and decides regarding their validity to enable the ConfirmButton
*/
func ValidateInput(studentEntry *widget.Entry, courseEntry *widget.Entry, confirmButton *widget.Button) {
	var validStudent bool
	var validCourse bool
	studentEntry.Validator = func(s string) error {
		re, _ := regexp.Compile(`^[a-zA-Z]+(\s[a-zA-Z]+)+$`)
		if !re.MatchString(s) {
			validStudent = false
			return errors.New("Die Eingabe entspricht nicht den Bedingungen(min. 1x Leerzeichen, nur Buchstaben!")
		} else {
			validStudent = true
		}
		return nil
	}
	studentEntry.OnChanged = func(s string) {
		if len(s) > 30 {
			s = s[0:10]
			studentEntry.SetText(s)
		}

		if validCourse && validStudent {
			confirmButton.Enable()
		} else {
			confirmButton.Disable()
		}
	}

	courseEntry.Validator = func(s string) error {
		re, _ := regexp.Compile(`\bT[A-Z]{2}\d{2}\b`)
		if !re.MatchString(s) {
			validCourse = false
			return errors.New("die Eingabe entspricht keinem validen Kurs")
		} else {
			validCourse = true
		}
		return nil
	}
	courseEntry.OnChanged = func(s string) {
		if len(s) > 10 {
			s = s[0:10]
			courseEntry.SetText(s)
		}
		if validCourse && validStudent {
			confirmButton.Enable()
		} else {
			confirmButton.Disable()
		}

	}
}

/*
InsertStudent collects all of the forms entries as input and builds student struct out of them to push on the db.
*/
func InsertStudent(w fyne.Window, form *dialog.CustomDialog, studentEntry *widget.Entry, courseEntry *widget.Entry, isImmatriculated *widget.Check) func() {
	insertStudent := func() {
		words := strings.Fields(studentEntry.Text)
		lastName := words[len(words)-1]
		firstName := strings.Join(words[:len(words)-1], " ")
		course, _ := myMVVM.CourseByName(courseEntry.Text)
		student := yaac_shared.Student{
			FirstName:        firstName,
			LastName:         lastName,
			CourseID:         course.ID,
			IsImmatriculated: isImmatriculated.Checked,
		}
		_, err := myMVVM.InsertStudent(student)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			form.Hide()
			dialog.ShowInformation("Student hinzufügen", fmt.Sprintf("%s %s %s", "Es wurde erfolgreich der Student", studentEntry.Text, "angelegt!"), w)
		}
	}
	return insertStudent
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
