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

func rgbGradient(x, y, w, h int) color.Color {
	g := int(float32(x) / float32(w) * float32(255))
	b := int(float32(y) / float32(h) * float32(255))

	return color.NRGBA{uint8(255 - b), uint8(g), uint8(b), 0xff}
}

func overviewScreen(_ fyne.Window) fyne.CanvasObject {
	title := canvas.NewText(" Anwesenheitsliste der Kurse - Heute", color.Black)
	title.Alignment = fyne.TextAlignLeading
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	grid := container.NewGridWrap(fyne.NewSize(250, 250))
	courses, _ := myMVVM.Courses()
	for _, element := range courses {
		// element is the element from someSlice for where we are
		students, _ := myMVVM.CourseStudents(yaac_shared.Course{Model: gorm.Model{ID: element.ID}})
		grid.Add(NewOverviewWidget(element.Name, students))
	}
	header := container.NewVBox(container.NewGridWrap(fyne.NewSize(400, 200), title), widget.NewSeparator())
	return container.NewBorder(header, nil, nil, nil, container.NewVScroll(grid))
}

type overviewWidget struct {
	widget.BaseWidget
	frame   *canvas.Rectangle
	title   *fyne.Container
	content *fyne.Container
	button  *widget.Button
}

func NewOverviewWidget(title string, attendance []yaac_shared.Student) *overviewWidget {
	titleLabel := widget.NewLabel(title)
	contentFrame := container.NewVBox()
	for _, student := range attendance {
		contentFrame.Add(widget.NewLabelWithStyle(fmt.Sprintf("%s %s", student.FirstName, student.LastName), fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: false, Monospace: false}))
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
		content: container.NewVBox(contentFrame),
		button: widget.NewButtonWithIcon("", yaac_shared.ResourceIconPng, func() {
			_ = getImage(title)
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
		container.NewBorder(header, nil, nil, nil, item.content),
	)
	return widget.NewSimpleRenderer(c)
}
func getImage(title string) []byte {
	selectedCourse, _ := myMVVM.CourseByName(title)
	list, _ := myMVVM.LatestList(yaac_shared.Course{Model: gorm.Model{ID: selectedCourse.ID}}, time.Now())
	return list.Image
}
