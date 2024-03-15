package pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func LandingScreen(_ fyne.Window) fyne.CanvasObject {
	yaacLogo := canvas.NewImageFromFile("assets/Icon.png")
	yaacLogo.FillMode = canvas.ImageFillContain
	yaacLogo.SetMinSize(fyne.NewSize(350, 350))
	title := canvas.NewText("YAAC", color.Black)
	title.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	title.TextSize = 150
	description := canvas.NewText("Ihre Plattform zur Studenten Verwaltung", color.Black)
	description.TextStyle = fyne.TextStyle{Bold: true}
	description.TextSize = 60

	return container.NewCenter(container.NewVBox(container.NewCenter(container.NewHBox(
		yaacLogo,
		title,
	)), description))
}
