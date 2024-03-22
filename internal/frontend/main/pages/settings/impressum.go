package yaac_frontend_settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func impressumScreen() fyne.CanvasObject {
	title := ReturnHeader("Impressum")
	doku := canvas.NewImageFromFile("assets/YAACImpressum.png")
	doku.FillMode = canvas.ImageFillContain
	imageFrame := container.NewGridWrap(fyne.NewSize(1200, 1000), doku)
	content := container.NewMax(container.NewVBox(container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 200), title)), widget.NewSeparator(), container.NewCenter(imageFrame)))
	return container.NewVScroll(content)
}
