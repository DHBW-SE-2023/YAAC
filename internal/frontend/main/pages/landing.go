package yaac_frontend_pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func LandingScreen(_ fyne.Window) fyne.CanvasObject {
	yaacLogo := canvas.NewImageFromResource(yaac_shared.ResourceIconPng)
	yaacLogo.FillMode = canvas.ImageFillContain
	yaacLogo.SetMinSize(fyne.NewSize(350, 350))
	title := canvas.NewText("YAAC", color.Black)
	title.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	title.TextSize = 150
	description := canvas.NewText("Ihre Plattform zur Studentenverwaltung", color.Black)
	description.TextStyle = fyne.TextStyle{Bold: true}
	description.TextSize = 60
	spacer := container.NewGridWrap(fyne.NewSize(50, 50), layout.NewSpacer())
	return container.NewBorder(spacer, spacer, spacer, spacer, container.NewPadded(container.NewPadded(container.NewCenter(container.NewVBox(container.NewCenter(container.NewHBox(
		yaacLogo,
		title,
	)), description)))))
}
