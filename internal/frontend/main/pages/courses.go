package yaac_frontend_pages

import (
	"errors"
	"fmt"
	"regexp"
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

func CoursesScreen(w fyne.Window) fyne.CanvasObject {
	course := &SelectionTracker{
		courseName: widget.NewLabel(""),
		secondary:  widget.NewLabel(""),
	}
	var dates []string
	selection := widget.NewLabel("")

	title := ReturnHeader("Kursansicht")
	tableHeader, courseTable := ReturnAttendanceTable("Name", "Status")
	dateDropdown := ReturnDateDropdown(dates)
	courseDropdown := ReturnCourseDropdown(course, selection, dateDropdown, "course")
	editDropdown := ReturnEditDropdown(w, courseDropdown, dateDropdown, courseTable)
	ConfigureDateDropdownVerification(dateDropdown, course, selection, courseTable, editDropdown)
	header := container.NewGridWrap(fyne.NewSize(400, 200), title)
	dropdownArea := container.NewGridWrap(fyne.NewSize(200, 40), courseDropdown, dateDropdown, layout.NewSpacer(), layout.NewSpacer(), editDropdown)
	selectionArea := container.NewVBox(selection, widget.NewSeparator(), tableHeader)
	studentView := container.NewBorder(container.NewVBox(header, dropdownArea), nil, nil, nil, container.NewBorder(selectionArea, nil, nil, nil, container.NewVScroll(courseTable)))
	return studentView
}

/*
ReturnDateDropdown returns the configured dateDropdown passing the dates list
*/
func ReturnDateDropdown(dates []string) *widget.SelectEntry {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	dateDropdown := widget.NewSelectEntry(dates)
	dateDropdown.PlaceHolder = "YYYY-MM-DD"
	dateDropdown.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
	dateDropdown.Scroll = container.ScrollNone
	dateDropdown.Validator = func(s string) error {
		if !re.MatchString(s) {
			return errors.New("failure")
		}
		return nil
	}
	dateDropdown.Disable()

	return dateDropdown
}

/*
ReturnEditDropdown returns the configured editDropdown passing the fyne.Window(for redirection),courseDropdown,dateDropdown, courseTable
since they will be necessary for change handling.
*/
func ReturnEditDropdown(w fyne.Window, courseDropdown *widget.Select, dateDropdown *widget.SelectEntry, courseTable *fyne.Container) *widget.Select {
	editDropdown := widget.NewSelect([]string{
		"Liste bearbeiten",
		"Liste anzeigen",
		"Liste hochladen",
	}, func(s string) {
		if s == "Liste bearbeiten" {
			lastView = w.Content()
			VerifyList(w, courseDropdown.Selected, dateDropdown.Text, courseTable)
		} else if s == "Liste anzeigen" {
			ShowImage(w, courseDropdown.Selected, dateDropdown.Text)
			courseTable.RemoveAll()
		} else {
			OpenImageUpload(w, courseDropdown.Selected, dateDropdown.Text)
		}
		courseTable.RemoveAll()
	})
	editDropdown.Selected = "Listenkonfiguration"
	editDropdown.Disable()
	return editDropdown
}

/*
RefreshDateDropdown is responsible for refreshing the date dropdown as soon as course has been selected
*/
func RefreshDateDropdown(dateDropdown *widget.SelectEntry, course string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	lists, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
	var dates []string
	for _, element := range lists {
		dates = append(dates, element.ReceivedAt.Format("2006-01-02"))
	}
	dateDropdown.SetOptions(dates)
}

/*
ConfigureDateDropdownVerification is responsible for configuring the verification process for text inputs passing the dateDropdown,
course SelectionTracker, selection Label, courseTable Container as well as the editDropdown for further processing.
*/
func ConfigureDateDropdownVerification(dateDropdown *widget.SelectEntry, course *SelectionTracker, selection *widget.Label, courseTable *fyne.Container, editDropdown *widget.Select) {
	dateDropdown.OnChanged = func(s string) {
		if len(s) > 10 {
			s = s[0:10]
		}
		re := regexp.MustCompile(`[^\d-]`)
		s = re.ReplaceAllString(s, "")
		dateDropdown.SetText(s)
		if dateDropdown.Validate() != nil {
			course.secondary.SetText("Falsches Datumsformat")
		} else {
			course.secondary.SetText(s)
		}
		selection.SetText(RefreshSelection(course))
		courseTable.RemoveAll()
		RefreshCourseAttendancy(courseTable, course.courseName.Text, s)
		editDropdown.Enable()
	}
}

/*
RefreshCourseAttendancy is responsible for refreshing the course attendancy list as soon as a date and a course has been selected
*/
func RefreshCourseAttendancy(table *fyne.Container, course string, date string) {
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

/*
VerifyList will redirect the user to the Verification Page passing the currently selected course,date, courseTable
*/
func VerifyList(w fyne.Window, course string, date string, courseTable *fyne.Container) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	parsedTime, _ := time.Parse("2006-01-02", date)
	w.SetContent(VerificationScreen(w, GetImageByDate(course, parsedTime.Add(24*time.Hour)), int(selectedCourse.ID), courseTable, parsedTime))
}

/*
ShowImage will display the currently selected list in a seperate window passing course and date.
*/
func ShowImage(w fyne.Window, course string, date string) {
	selectedCourse, _ := myMVVM.CourseByName(course)
	parsedTime, _ := time.Parse("2006-01-02", date)
	list, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, parsedTime, parsedTime.Add(24*time.Hour))
	img := RotateImage(list[0].Image)
	img.FillMode = canvas.ImageFillOriginal
	customDialog := dialog.NewCustom(fmt.Sprintf("%s %s", "Listen vom", date), "Beenden", container.NewVScroll(container.NewGridWrap(fyne.NewSize(800, 1000), img)), w)
	customDialog.Resize(fyne.NewSize(800, 800))
	customDialog.Show()
}
