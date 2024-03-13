package settings

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var page int = 1
var pages, _ = os.ReadDir("assets/doku")

func wikiScreen() fyne.CanvasObject {
	title := canvas.NewText(" Nutzerdokumentation", color.Black)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	titleFrame := container.NewCenter(container.NewGridWrap(fyne.NewSize(800, 200), title))
	doku := canvas.NewImageFromFile(fmt.Sprintf("assets/doku/%d.png", page))
	doku.FillMode = canvas.ImageFillContain
	imageFrame := container.NewGridWrap(fyne.NewSize(1200, 1000), doku)
	nextButton := widget.NewButton("Weiter", func() {
		imageFrame.RemoveAll()
		page += 1
		print(len(pages))
		if page == len(pages) {
			page = 1
		}
		loadImage(fmt.Sprintf("assets/doku/%d.png", page), imageFrame)
	})
	backButton := widget.NewButton("Zurück", func() {
		imageFrame.RemoveAll()
		page -= 1
		if page == 0 {
			page = len(pages)
		}
		loadImage(fmt.Sprintf("assets/doku/%d.png", page), imageFrame)
	})
	buttonArea := container.NewCenter(container.NewHBox(container.NewAdaptiveGrid(3, backButton, layout.NewSpacer(), nextButton)))
	content := container.NewMax(container.NewVBox(titleFrame, container.NewCenter(imageFrame), buttonArea))
	return container.NewVScroll(content)
}

func loadImage(imagePath string, imageFrame *fyne.Container) {
	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain
	imageFrame.Add(image)
}
