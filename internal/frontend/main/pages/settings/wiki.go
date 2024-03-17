package yaac_frontend_settings

import (
	"fmt"
	"io/fs"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func wikiScreen() fyne.CanvasObject {
	title := ReturnHeader("Nutzerdokumentation")
	var page int = 1
	var pages, _ = os.ReadDir("assets/doku")

	doku := canvas.NewImageFromFile(fmt.Sprintf("assets/doku/%d.png", page))
	doku.FillMode = canvas.ImageFillContain
	imageFrame := container.NewGridWrap(fyne.NewSize(1200, 1000), doku)

	nextButton := ReturnNextButton(imageFrame, page, pages)
	backButton := ReturnBackButton(imageFrame, page, pages)

	buttonArea := container.NewCenter(container.NewHBox(container.NewAdaptiveGrid(3, backButton, layout.NewSpacer(), nextButton)))
	content := container.NewMax(container.NewVBox(title, container.NewCenter(imageFrame), buttonArea))
	return container.NewVScroll(content)
}

/*
ReturnNextButton returns the fully configured nextButton which is responsible for switch to the next page
*/
func ReturnNextButton(imageFrame *fyne.Container, page int, pages []fs.DirEntry) *widget.Button {
	nextButton := widget.NewButton("Weiter", func() {
		imageFrame.RemoveAll()
		page += 1
		if page == len(pages) {
			page = 1
		}
		LoadImage(fmt.Sprintf("assets/doku/%d.png", page), imageFrame)
	})
	return nextButton
}

/*
ReturnBackButton returns the fully configured backButton which is responsible for switch to the last page
*/
func ReturnBackButton(imageFrame *fyne.Container, page int, pages []fs.DirEntry) *widget.Button {
	backButton := widget.NewButton("Zurück", func() {
		imageFrame.RemoveAll()
		page -= 1
		if page == 0 {
			page = len(pages)
		}
		LoadImage(fmt.Sprintf("assets/doku/%d.png", page), imageFrame)
	})
	return backButton
}

/*
LoadImage refreshes the currently displayed image in the respective imageFrame on backButton|nextButton Clicked Events
*/
func LoadImage(imagePath string, imageFrame *fyne.Container) {
	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillContain
	imageFrame.Add(image)
}
