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
	"gocv.io/x/gocv"
	"gorm.io/gorm"
)

type verificationWidget struct {
	widget.BaseWidget
	frame     *canvas.Rectangle
	attending *widget.Check
	content   *fyne.Container
}

func NewVerificationWidget(student string, attendance bool, students []string, attendances []bool) *verificationWidget {
	item := &verificationWidget{
		frame: &canvas.Rectangle{
			FillColor:    color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeColor:  color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeWidth:  4.0,
			CornerRadius: 10,
		},
		attending: widget.NewCheck(student, func(value bool) {
			for index, name := range students {
				if name == student {
					attendances[index] = value
				}
			}
		}),
		content: container.NewGridWithColumns(2),
	}
	item.ExtendBaseWidget(item)
	item.attending.Checked = attendance
	item.attending.OnChanged = func(b bool) { item.attending.Checked = b }

	return item
}

func (item *verificationWidget) CreateRenderer() fyne.WidgetRenderer {
	item.content.Add(item.attending)
	c := container.NewPadded(
		item.frame,
		item.content,
	)
	c.Resize(fyne.NewSize(200, 200))
	return widget.NewSimpleRenderer(c)
}

func verificationScreen(w fyne.Window, img []byte, course int, optional ...time.Time) fyne.CanvasObject {
	//Define Header
	headerFrame := canvas.NewRectangle(color.White)
	logo := canvas.NewImageFromFile("assets/DHBW.png")
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 200))
	title := canvas.NewText("Anwesenheitsprüfung", color.Black)
	title.Alignment = fyne.TextAlignLeading
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	header := container.NewGridWrap(fyne.NewSize(200, 200), logo, title)

	description := canvas.NewText("Überprüfen sie die dargestellt Liste und wählen gegebenfalls Anwesende Studenten aus:", color.Black)
	description.TextSize = 16
	description.TextStyle = fyne.TextStyle{Bold: true}

	widgetList := container.NewVBox()
	if len(optional) > 0 {
		widgetList = loadData(course, optional[0])
	} else {
		widgetList = loadData(course, time.Now())
	}
	image := RotateImage(img)

	confirmButton := widget.NewButton("Bestätigen", func() {
		var attendances []bool
		for _, obj := range widgetList.Objects {
			if widget, ok := obj.(*verificationWidget); ok {
				attendances = append(attendances, widget.attending.Checked)
			}
		}
		if len(optional) > 0 {
			updateList(attendances, course, optional[0])
		} else {
			updateList(attendances, course, time.Now())
		}
		if LastView != nil {
			w.SetContent(LastView)
			if len(optional) > 0 {
				courses, _ := myMVVM.Courses()
				var selectedCourse string
				for _, element := range courses {
					if element.ID == uint(course) {
						selectedCourse = element.Name
					}
				}
				courseTable.RemoveAll()
				refreshCourseAttendancy(courseTable, selectedCourse, optional[0].Format("2006-01-02"))
			} else {
				loadOverviewWidgets(w, overviewGrid)
			}
		}
	})

	exitButton := widget.NewButton("Zurück zur Startseite", func() {
		if LastView != nil {
			w.SetContent(LastView)
		}
	})
	verificationList := container.NewVBox(description, widgetList, container.NewCenter(container.NewGridWithRows(1, confirmButton, exitButton)))
	contentBox := container.NewAdaptiveGrid(2,
		image,
		container.NewPadded(container.NewPadded(container.NewVScroll(verificationList))),
	)
	return container.NewBorder(container.NewMax(headerFrame, header), nil, nil, nil, contentBox)
}

func RotateImage(img []byte) *canvas.Image {
	rotated := gocv.NewMat()
	imgMat, _ := gocv.IMDecode(img, gocv.IMReadAnyColor)
	gocv.Rotate(imgMat, &rotated, gocv.Rotate90Clockwise)
	imgNew, _ := imgMat.ToImage()
	image := canvas.NewImageFromImage(imgNew)
	return image
}

func getAttendancies(course int, date time.Time) ([]string, []bool) {
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

func loadData(course int, date time.Time) *fyne.Container {
	students, attendance := getAttendancies(course, date)
	studentList := container.NewVBox()
	for i := 0; i < len(students); i++ {
		studentList.Add(NewVerificationWidget(students[i], attendance[i], students, attendance))
	}
	return studentList
}

func updateList(attendances []bool, course int, date time.Time) {
	list, _ := myMVVM.AllAttendanceListInRangeByCourse(yaac_shared.Course{Model: gorm.Model{ID: uint(course)}}, date.AddDate(0, 0, -31), date.Add(24*time.Hour))
	var attendanceList []yaac_shared.Attendance
	for i := 0; i < len(list[0].Attendancies); i++ {
		att := yaac_shared.Attendance{
			StudentID:        list[0].Attendancies[i].StudentID,
			IsAttending:      attendances[i],
			AttendanceListID: list[0].ID,
		}
		attendanceList = append(attendanceList, att)
	}
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
