package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func impressumScreen() fyne.CanvasObject {
	title := canvas.NewText("Impressum", color.Black)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	titleFrame := container.NewCenter(container.NewGridWrap(fyne.NewSize(800, 200), title))
	doku := canvas.NewImageFromFile("assets/YAACImpressum.png")
	doku.FillMode = canvas.ImageFillContain
	imageFrame := container.NewGridWrap(fyne.NewSize(1200, 1000), doku)
	content := container.NewMax(container.NewVBox(titleFrame, container.NewCenter(imageFrame)))
	return container.NewVScroll(content)
}
