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
)

/*
Delaration of custom Overwidget Struct
*/
type OverviewWidget struct {
	widget.BaseWidget
	frame   *canvas.Rectangle
	title   *fyne.Container
	content *fyne.Container
	button  *widget.Button
}

func NewOverviewWidget(w fyne.Window, title string, courseId int, nonAttending []string, totalStudents int) *OverviewWidget {
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
	item := &OverviewWidget{
		frame: &canvas.Rectangle{
			FillColor:    color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeColor:  color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeWidth:  4.0,
			CornerRadius: 20,
		},
		title:   container.NewVBox(titleLabel),
		content: contentFrame,
		button: widget.NewButtonWithIcon("", yaac_shared.ResourceImageIconPng, func() {
			v := VerificationScreen(w, GetImageByDate(title, time.Now()), courseId, container.NewWithoutLayout())
			lastView = w.Content()
			w.SetContent(v)
		}),
	}
	item.ExtendBaseWidget(item)
	return item
}
func (item *OverviewWidget) CreateRenderer() fyne.WidgetRenderer {
	item.frame.Resize(fyne.NewSize(250, 250))
	header := container.NewPadded(container.NewPadded(container.NewGridWithColumns(2, item.title, item.button)))
	c := container.NewPadded(
		item.frame,
		container.NewBorder(header, nil, nil, nil, container.NewVScroll(item.content)),
	)
	return widget.NewSimpleRenderer(c)
}

/*
Delaration of custom AttendanceRow Struct for course and student view
*/
type AttendanceRow struct {
	widget.BaseWidget
	frame    *canvas.Rectangle
	date     *widget.Label
	state    *widget.Label
	content  *fyne.Container
	OnTapped func()
}

func NewAttendanceRow(dateText string, stateText string) *AttendanceRow {
	item := &AttendanceRow{
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
	DetectAttendancyState(item)
	return item
}
func (item *AttendanceRow) Tapped(ev *fyne.PointEvent) { item.OnTapped() }

func (item *AttendanceRow) CreateRenderer() fyne.WidgetRenderer {
	item.frame.Resize(fyne.NewSize(250, 250))
	item.content.Add(item.date)
	item.content.Add(item.state)
	c := container.NewPadded(
		item.frame,
		item.content,
	)
	return widget.NewSimpleRenderer(c)
}

/*
Delaration of custom VerificationWidget Struct for the verification view
*/
type VerificationWidget struct {
	widget.BaseWidget
	frame     *canvas.Rectangle
	attending *widget.Check
	content   *fyne.Container
}

func NewVerificationWidget(student string, attendance bool, students []string, attendances []bool) *VerificationWidget {
	item := &VerificationWidget{
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

func (item *VerificationWidget) CreateRenderer() fyne.WidgetRenderer {
	item.content.Add(item.attending)
	c := container.NewPadded(
		item.frame,
		item.content,
	)
	c.Resize(fyne.NewSize(200, 200))
	return widget.NewSimpleRenderer(c)
}

/*
DetectAttendancyState will detect the current attendancies state and change frame color accordingly for each AttendanceRow
*/
func DetectAttendancyState(item *AttendanceRow) {
	if item.state.Text != "Anwesend" {
		item.frame.FillColor = color.RGBA{227, 0, 27, 255}
		item.frame.StrokeColor = color.RGBA{227, 0, 27, 255}
		item.frame.StrokeWidth = 2.0
	}
}

/*
Declaration of SelectionTracker struct this will track the current selections for courses and students
*/
type SelectionTracker struct {
	courseName *widget.Label
	secondary  *widget.Label
}
