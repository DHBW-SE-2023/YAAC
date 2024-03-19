package yaac_frontend_pages

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"gocv.io/x/gocv"
	"gorm.io/gorm"
)

var lastView fyne.CanvasObject

/*
RefreshSelection is responsible for refreshing the selection text on changes by using the SelectionTracker struct for student and course view
*/
func RefreshSelection(selection *SelectionTracker) string {
	return fmt.Sprintf("%s - %s", selection.courseName.Text, selection.secondary.Text)
}

/*
ReturnHeader will return the canvas.Text objet of each page
*/
func ReturnHeader(pageTitle string) *canvas.Text {
	title := canvas.NewText(pageTitle, color.Black)
	title.Alignment = fyne.TextAlignLeading
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	return title
}

/*
ReturnAttendanceTable returns the configured attendanceTable as well the tableHeader for the courses and student view
*/
func ReturnAttendanceTable(header1 string, header2 string) (*fyne.Container, *fyne.Container) {
	tableHeader := container.NewGridWithColumns(2)
	tableHeader.Add(widget.NewLabel(header1))
	tableHeader.Add(widget.NewLabel(header2))
	attendanceTable := container.NewVBox()
	return tableHeader, attendanceTable
}

/*
ReturnCourseDropdown returns the configured courseDropdown passing the SelectionTracker Struct, selection Label, dependingDropdown
and source(course|student) since they will be necessary for change handling for use in course and student view.
*/
func ReturnCourseDropdown(selectionTracker *SelectionTracker, selection *widget.Label, dependingDropdown *widget.SelectEntry, source string) *widget.Select {
	courseDropdown := widget.NewSelect([]string{
		"TIK22",
		"TIT22",
		"TIS22",
		"TIM22",
	}, func(s string) {
		fmt.Println(s)
		selectionTracker.courseName.SetText(s)
		selection.SetText(RefreshSelection(selectionTracker))
		if source == "course" {
			RefreshDateDropdown(dependingDropdown, s)
		} else {
			RefreshStudentDropdown(dependingDropdown, s)
		}
		dependingDropdown.SetText("")
		dependingDropdown.Enable()
	})
	courseDropdown.Selected = "Kursauswahl"
	return courseDropdown
}

/*
MapBooleans will map the booleans to be displayed accordingly on the frontend
*/
func MapBooleans(b bool) string {
	if bool(b) {
		return "Anwesend"
	}
	return "Abwesend"
}

/*
GetImageByDate requests the latest list regarding a defined date from the dataset
*/
func GetImageByDate(title string, date time.Time) []byte {
	selectedCourse, _ := myMVVM.CourseByName(title)
	list, _ := myMVVM.LatestList(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, date)
	return list.Image
}

/*
ReturnNonAttending extracts the non attending students given a certain list of attendances(List).
*/
func ReturnNonAttending(attendance []yaac_shared.Attendance) []string {
	var returnNonAttending []string
	for _, element := range attendance {
		if !element.IsAttending {
			students, _ := myMVVM.Students(yaac_shared.Student{Model: gorm.Model{ID: element.StudentID}})
			returnNonAttending = append(returnNonAttending, fmt.Sprintf("%s %s", students[0].FirstName, students[0].LastName))
		}
	}
	return returnNonAttending
}

/*
OpenImageUpload opens a customized FileDialog Window with the possibility to pass optional parameters.
- Optional[0]=Course Name
- Optional[1]=List Date
These values will be passed to the function ShowFileDialog to handle the actual file upload.
*/
func OpenImageUpload(w fyne.Window, optional ...string) {
	courseEntry := widget.NewEntry()
	if len(optional) != 0 {
		courseEntry.Text = optional[0]
	} else {
		courseEntry.Text = "TIK22,TIT22...."
	}

	fileUpload := widget.NewButton("Load Image", func() {
		if len(optional) > 1 {
			ShowFileDialog(w, courseEntry.Text, optional[1])
		} else {
			ShowFileDialog(w, courseEntry.Text)
		}
	})
	fileUpload.Disable()
	courseEntry.OnSubmitted = func(text string) {
		fileUpload.Enable()
	}
	content := container.NewVBox(
		widget.NewLabel("Geben sie das Kürzel des betroffenen Kurses ein:"),
		courseEntry,
		fileUpload,
	)
	customDialog := dialog.NewCustom("Listen Upload", "Beenden", content, w)
	customDialog.Show()
}
func ShowFileDialog(w fyne.Window, course string, optional ...string) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			log.Println("Cancelled")
			return
		}
		img := LoadImage(reader)
		if len(optional) != 0 {
			InsertList(img, course, optional[0])
		} else {
			InsertList(img, course)
		}
	}, w)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	fd.Show()
}
func LoadImage(f fyne.URIReadCloser) []byte {
	data, err := io.ReadAll(f)
	if err != nil {
		yaac_shared.App.SendNotification(fyne.NewNotification("Failed to load image data", err.Error()))
		return nil
	}
	return data
}

/*
InsertList will finally construct a yaac_shared.AttendanceList object from the previously gathered functions and push it on the database.
*/
func InsertList(img []byte, course string, optional ...string) {
	var testTime time.Time
	if len(optional) != 0 {
		testTime, _ = time.Parse("2006-01-02", optional[0])
	} else {
		testTime = time.Now()
	}

	selectedCourse, err := myMVVM.CourseByName(course)
	if err != nil {
		yaac_shared.App.SendNotification(fyne.NewNotification("Fehler bei Listen Uplaod", err.Error()))
		return
	}

	attendanceList := yaac_shared.AttendanceList{
		ReceivedAt: testTime,
		CourseID:   selectedCourse.ID,
		Image:      img,
	}
	_, err = myMVVM.InsertList(attendanceList)
	if err != nil {
		yaac_shared.App.SendNotification(fyne.NewNotification("Fehler bei Listen Uplaod", err.Error()))
	} else {
		yaac_shared.App.SendNotification(fyne.NewNotification("Liste erfolgreich hochgeladen", fmt.Sprintf("%s %s %s", "Ihre Liste für den Kurs", course, "wurde erfolgreich hochgeladen!")))
	}
}

/*
RotateImage will rotate images if necessary so they can be displayed vertically
*/
func RotateImage(img []byte) *canvas.Image {
	rotated := gocv.NewMat()
	imgMat, _ := gocv.IMDecode(img, gocv.IMReadAnyColor)
	gocv.Rotate(imgMat, &rotated, gocv.Rotate90Clockwise)
	imgNew, _ := imgMat.ToImage()
	image := canvas.NewImageFromImage(imgNew)
	return image
}