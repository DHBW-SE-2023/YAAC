package yaac_frontend_settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func impressumScreen() fyne.CanvasObject {
	title := ReturnHeader("Impressum")
	doku := canvas.NewImageFromResource(yaac_shared.ResourceYAACImpressumPng)
	doku.FillMode = canvas.ImageFillContain
	imageFrame := container.NewGridWrap(fyne.NewSize(1200, 1000), doku)
	content := container.NewMax(container.NewVBox(title, container.NewCenter(imageFrame)))
	return container.NewVScroll(content)
}
