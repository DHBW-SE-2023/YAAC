package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func settingsScreen(_ fyne.Window) fyne.CanvasObject {
	spacy := layout.NewSpacer()
	spacy.Resize(fyne.NewSquareSize(2))
	return container.NewGridWrap(fyne.NewSize(90, 90),
		canvas.NewImageFromResource(data.FyneLogo),
		spacy,
		canvas.NewImageFromResource(data.FyneLogo),
	)
}
