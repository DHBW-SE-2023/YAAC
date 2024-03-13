package pages

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
	"gorm.io/gorm"
)

var LastView fyne.CanvasObject
var overviewGrid *fyne.Container

func overviewScreen(w fyne.Window) fyne.CanvasObject {
	insertListIcon, _ := fyne.LoadResourceFromPath("assets/imageUpload2.png")
	image := canvas.NewImageFromResource(insertListIcon)
	image.FillMode = canvas.ImageFillOriginal // Ensure the original size of the image
	insertList := widget.NewButton("", func() {
		openImagUpload(w)
	})

	buttonIconContainer := container.NewCenter(image)
	title := canvas.NewText(" Anwesenheitsliste der Kurse - Heute", color.Black)
	title.Alignment = fyne.TextAlignLeading
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}

	overviewGrid = container.NewGridWrap(fyne.NewSize(250, 250))
	loadOverviewWidgets(w, overviewGrid)
	header := container.NewVBox(container.NewBorder(nil, nil, nil, container.NewPadded(container.NewPadded(container.NewPadded(container.NewPadded(buttonIconContainer, insertList)))), container.NewGridWrap(fyne.NewSize(400, 200), title)), widget.NewSeparator())
	return container.NewBorder(header, nil, nil, nil, container.NewVScroll(overviewGrid))
}

type overviewWidget struct {
	widget.BaseWidget
	frame   *canvas.Rectangle
	title   *fyne.Container
	content *fyne.Container
	button  *widget.Button
}

func NewOverviewWidget(w fyne.Window, title string, courseId int, nonAttending []string, totalStudents int) *overviewWidget {
	imageResource, _ := fyne.LoadResourceFromPath("assets/imageIcon.png")
	titleLabel := widget.NewLabel(title)
	contentFrame := container.NewVBox()
	if len(nonAttending) == 0 {
		contentFrame.Add(widget.NewLabelWithStyle(fmt.Sprintf("%d/%d %s", totalStudents, totalStudents, "Anwesend"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: false, Monospace: false}))
	} else {
		for _, student := range nonAttending {
			contentFrame.Add(widget.NewLabelWithStyle(student, fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: false, Monospace: false}))
		}
	}
	titleLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: false, Monospace: false}
	titleLabel.Alignment = fyne.TextAlignCenter
	item := &overviewWidget{
		frame: &canvas.Rectangle{
			FillColor:    color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeColor:  color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeWidth:  4.0,
			CornerRadius: 20,
		},
		title:   container.NewVBox(titleLabel),
		content: contentFrame,
		button: widget.NewButtonWithIcon("", imageResource, func() {
			v := verificationScreen(w, getImage(title, time.Now()), courseId)
			LastView = w.Content()
			w.SetContent(v)
		}),
	}
	item.ExtendBaseWidget(item)
	return item
}

func (item *overviewWidget) CreateRenderer() fyne.WidgetRenderer {
	item.frame.Resize(fyne.NewSize(250, 250))
	header := container.NewPadded(container.NewPadded(container.NewGridWithColumns(2, item.title, item.button)))
	c := container.NewPadded(
		item.frame,
		container.NewBorder(header, nil, nil, nil, container.NewVScroll(item.content)),
	)
	return widget.NewSimpleRenderer(c)
}

func getImage(title string, date time.Time) []byte {
	selectedCourse, _ := myMVVM.CourseByName(title)
	list, _ := myMVVM.LatestList(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, date)
	return list.Image
}

func returnNonAttending(attendance []yaac_shared.Attendance) []string {
	var returnNonAttending []string
	for _, element := range attendance {
		if !element.IsAttending {
			students, _ := myMVVM.Students(yaac_shared.Student{Model: gorm.Model{ID: element.StudentID}})
			returnNonAttending = append(returnNonAttending, fmt.Sprintf("%s %s", students[0].FirstName, students[0].LastName))
		}
	}
	return returnNonAttending
}

func openImagUpload(w fyne.Window, optional ...string) {
	courseEntry := widget.NewEntry()
	if len(optional) != 0 {
		courseEntry.Text = optional[0]
	} else {
		courseEntry.Text = "TIK22,TIT22...."
	}

	fileUpload := widget.NewButton("Load Image", func() {
		if len(optional) > 1 {
			showFileDialog(w, courseEntry.Text, optional[1])
		} else {
			showFileDialog(w, courseEntry.Text)
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

func showFileDialog(w fyne.Window, course string, optional ...string) {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			log.Println("Cancelled")
			return
		}
		img := loadImage(reader)
		if len(optional) != 0 {
			insertList(img, course, optional[0])
		} else {
			insertList(img, course)
		}
	}, w)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	fd.Show()
}

func loadImage(f fyne.URIReadCloser) []byte {
	data, err := io.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	return data
}

func insertList(img []byte, course string, optional ...string) {
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

func loadOverviewWidgets(w fyne.Window, grid *fyne.Container) {
	grid.RemoveAll()
	courses, _ := myMVVM.Courses()
	for _, element := range courses {
		var students []string
		lists, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: element.ID}}, time.Now().AddDate(0, 0, -30), time.Now())
		if len(lists[0].Attendancies) > 0 {
			students = returnNonAttending(lists[0].Attendancies)
		} else {
			students = append(students, "Kein Listeingang")
		}
		grid.Add(NewOverviewWidget(w, element.Name, int(element.ID), students, len(lists[0].Attendancies)))
	}
}
