package yaac_frontend_settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

/*
ReturnHeader will return the canvas.Text objet of each page
*/
func ReturnHeader(pageTitle string) *fyne.Container {
	title := canvas.NewText(pageTitle, color.Black)
	title.Alignment = fyne.TextAlignLeading
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	return container.NewCenter(container.NewGridWrap(fyne.NewSize(800, 200), title))
}

func MapMailBooleans(b bool) string {
	if bool(b) {
		return "Aktiv"
	}
	return "Nicht erreichbar"
}
