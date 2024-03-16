package yaac_frontend_pages

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

func VerificationScreen(w fyne.Window, img []byte, course int, courseTable *fyne.Container, optional ...time.Time) fyne.CanvasObject {
	header := ReturnVerificationHeader()
	description := canvas.NewText("Überprüfen sie die dargestellt Liste und wählen gegebenfalls Anwesende Studenten aus:", color.Black)
	description.TextSize = 16
	description.TextStyle = fyne.TextStyle{Bold: true}
	image := RotateImage(img)

	widgetList := container.NewVBox()
	if len(optional) > 0 {
		widgetList = LoadVerificationWidgets(course, optional[0])
	} else {
		widgetList = LoadVerificationWidgets(course, time.Now())
	}

	confirmButton := ReturnConfirmButton(w, widgetList, optional, course, courseTable)
	exitButton := ReturnExitButton(w)

	verificationList := container.NewVBox(description, widgetList, container.NewCenter(container.NewGridWithRows(1, confirmButton, exitButton)))
	contentBox := container.NewAdaptiveGrid(2, image, container.NewPadded(container.NewPadded(container.NewVScroll(verificationList))))
	return container.NewBorder(header, nil, nil, nil, contentBox)
}

/*
ReturnVerificationHeader returns the configured VerificationHeader including navBar extern layout
*/
func ReturnVerificationHeader() *fyne.Container {
	headerFrame := canvas.NewRectangle(color.White)
	logo := canvas.NewImageFromFile("assets/DHBW.png")
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 200))
	title := ReturnHeader("Anwesenheitsprüfung")
	header := container.NewGridWrap(fyne.NewSize(200, 200), logo, title)
	return container.NewMax(headerFrame, header)
}

/*
ReturnConfirmButton returns the configured VerificationHeader passing the window,widgetList, optional params []time.Time,courseID and courseTable.
These values will be used to update the attendancees when the button gets clicked
*/
func ReturnConfirmButton(w fyne.Window, widgetList *fyne.Container, optional []time.Time, course int, courseTable *fyne.Container) *widget.Button {
	confirmButton := widget.NewButton("Bestätigen", func() {
		var attendances []bool
		for _, obj := range widgetList.Objects {
			if widget, ok := obj.(*VerificationWidget); ok {
				attendances = append(attendances, widget.attending.Checked)
			}
		}
		if len(optional) > 0 {
			UpdateList(attendances, course, optional[0])
		} else {
			UpdateList(attendances, course, time.Now())
		}
		ReturnToPreviousPage(w, course, optional, courseTable)
	})
	return confirmButton
}

/*
UpdateList updates all new attendancies by initializing a AttedanceList Object as soon as the confirmButton gets clicked.
*/
func UpdateList(attendances []bool, course int, date time.Time) {
	list, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: uint(course)}}, date.AddDate(0, 0, -31), date.Add(24*time.Hour))
	attendanceList := ReturnUpdatedAttendancies(attendances, list)
	_, err := myMVVM.UpdateList(yaac_shared.AttendanceList{
		ID:           list[0].ID,
		CreatedAt:    list[0].CreatedAt,
		CourseID:     uint(course),
		Attendancies: attendanceList,
		Image:        list[0].Image,
		ReceivedAt:   list[0].ReceivedAt,
	})

	if err != nil {
		yaac_shared.App.SendNotification(fyne.NewNotification("Fehler bei Listenaktualisierung", err.Error()))
		return
	} else {
		yaac_shared.App.SendNotification(fyne.NewNotification("Ihre Liste wurde erfolgreich aktualisiert", ""))
	}
}

/*
ReturnUpdatedAttendancies provides UpdateList all updatedAttendancies by using the VerificationWidgets Checkbox values and construting
an Attendance object for each.
*/
func ReturnUpdatedAttendancies(attendances []bool, list []yaac_shared.AttendanceList) []yaac_shared.Attendance {
	var attendanceList []yaac_shared.Attendance
	for i := 0; i < len(list[0].Attendancies); i++ {
		att := yaac_shared.Attendance{
			StudentID:        list[0].Attendancies[i].StudentID,
			IsAttending:      attendances[i],
			AttendanceListID: list[0].ID,
		}
		attendanceList = append(attendanceList, att)
	}
	return attendanceList
}

/*
ReturnToPreviousPage executes a command to return to the previousPage passing the window, courseId,optional params []time.Time
and courseTable. Depending on the fact if optional params are passed or not the command decides where to navigate to.
*/
func ReturnToPreviousPage(w fyne.Window, course int, optional []time.Time, courseTable *fyne.Container) {
	if lastView != nil {
		w.SetContent(lastView)
		if len(optional) > 0 {
			courses, _ := myMVVM.Courses()
			var selectedCourse string
			for _, element := range courses {
				if element.ID == uint(course) {
					selectedCourse = element.Name
				}
			}
			courseTable.RemoveAll()
			RefreshCourseAttendancy(courseTable, selectedCourse, optional[0].Format("2006-01-02"))
		} else {
			LoadOverviewWidgets(w, overviewGrid)
		}
	}
}

/*
ReturnExit returns the configured exitButton to exit the current Page and return to lastView.
*/
func ReturnExitButton(w fyne.Window) *widget.Button {
	exitButton := widget.NewButton("Zurück zur Startseite", func() {
		if lastView != nil {
			w.SetContent(lastView)
		}
	})
	return exitButton
}

/*
GetAttendancies returns all attendancies using the courseID and date for the selected List. Returning the studentNames and isAttending
values in two seperate lists for LoadVerificationWidgets to use.
*/
func GetAttendancies(course int, date time.Time) ([]string, []bool) {
	attendancies, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: uint(course)}}, date.AddDate(0, 0, -31), date.Add(24*time.Hour))
	var students []string
	var attendance []bool
	for _, element := range attendancies[0].Attendancies {
		student, _ := myMVVM.Students(yaac_shared.Student{Model: gorm.Model{ID: element.StudentID}})
		students = append(students, fmt.Sprintf("%s %s", student[0].FirstName, student[0].LastName))
		attendance = append(attendance, element.IsAttending)
	}
	return students, attendance
}

/*
LoadVerificationWidgets loads all attendancies as VerificationWidgets into a container.
*/
func LoadVerificationWidgets(course int, date time.Time) *fyne.Container {
	students, attendance := GetAttendancies(course, date)
	studentList := container.NewVBox()
	for i := 0; i < len(students); i++ {
		studentList.Add(NewVerificationWidget(students[i], attendance[i], students, attendance))
	}
	return studentList
}
