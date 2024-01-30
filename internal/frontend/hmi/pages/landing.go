package pages

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func landingScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("data/Icon.png")
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Willkommen zurück XD", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewHBox(
			widget.NewHyperlink("User Dokumentation", parseURL("https://developer.fyne.io/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("Sponsoren", parseURL("https://www.sap.com/germany/index.html?url_id=auto_hp_redirect_germany")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))
}
