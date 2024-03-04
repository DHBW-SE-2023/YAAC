package pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	for i := 0; i <= 10; i++ {
		grid.Add(NewOverviewWidget("TIK22", "Max Alberti"))
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

func NewOverviewWidget(title string, attendance string) *overviewWidget {
	imageResource, _ := fyne.LoadResourceFromPath("assets/imageIcon.png")
	titleLabel := widget.NewLabel(title)
	contentLabel := widget.NewLabel(attendance)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: false, Monospace: false}
	titleLabel.Alignment = fyne.TextAlignCenter
	contentLabel.TextStyle = fyne.TextStyle{Bold: true, Italic: false, Monospace: false}
	contentLabel.Alignment = fyne.TextAlignCenter
	item := &overviewWidget{
		frame: &canvas.Rectangle{
			FillColor:    color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeColor:  color.NRGBA{R: 209, G: 209, B: 209, A: 255},
			StrokeWidth:  4.0,
			CornerRadius: 20,
		},
		title:   container.NewVBox(titleLabel),
		content: container.NewVBox(contentLabel),
		button: widget.NewButtonWithIcon("", imageResource, func() {
			println("Show Image!")
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
