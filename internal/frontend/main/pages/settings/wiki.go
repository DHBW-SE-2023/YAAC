package settings

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
)

func wikiScreen() fyne.CanvasObject {
	gradient := canvas.NewHorizontalGradient(color.NRGBA{0x80, 0, 0, 0xff}, color.NRGBA{0, 0x80, 0, 0xff})
	go func() {
		for {
			time.Sleep(time.Second)

			gradient.Angle += 45
			if gradient.Angle >= 360 {
				gradient.Angle -= 360
			}
			canvas.Refresh(gradient)
		}
	}()

	return container.NewGridWrap(fyne.NewSize(90, 90),
		canvas.NewImageFromResource(data.FyneLogo),
		&canvas.Rectangle{FillColor: color.NRGBA{0x80, 0, 0, 0xff},
			StrokeColor: color.NRGBA{R: 255, G: 120, B: 0, A: 255},
			StrokeWidth: 1},
		&canvas.Rectangle{
			FillColor:    color.NRGBA{R: 255, G: 200, B: 0, A: 180},
			StrokeColor:  color.NRGBA{R: 255, G: 120, B: 0, A: 255},
			StrokeWidth:  4.0,
			CornerRadius: 20},
		&canvas.Line{StrokeColor: color.NRGBA{0, 0, 0x80, 0xff}, StrokeWidth: 5},
		&canvas.Circle{StrokeColor: color.NRGBA{0, 0, 0x80, 0xff},
			FillColor:   color.NRGBA{0x30, 0x30, 0x30, 0x60},
			StrokeWidth: 2},
		canvas.NewText("Text", color.NRGBA{0, 0x80, 0, 0xff}),
		gradient,
		canvas.NewRadialGradient(color.NRGBA{0x80, 0, 0, 0xff}, color.NRGBA{0, 0x80, 0x80, 0xff}),
	)
}
