package yaac_frontend_settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

var page int = 1
var pages int = 8

func wikiScreen() fyne.CanvasObject {
	title := ReturnHeader("Nutzerdokumentation")

	// This is neccecary because the frame wont accept an uninitialized image
	image := canvas.NewImageFromResource(yaac_shared.ResourceDoku1Png)
	image.FillMode = canvas.ImageFillContain
	imageFrame := container.NewGridWrap(fyne.NewSize(1200, 1000), image)

	nextButton := createReturnNextButton(imageFrame)
	backButton := createReturnBackButton(imageFrame)

	buttonArea := container.NewCenter(container.NewHBox(container.NewAdaptiveGrid(3, backButton, layout.NewSpacer(), nextButton)))
	content := container.NewStack(container.NewVBox(title, container.NewCenter(imageFrame), buttonArea))

	return container.NewVScroll(content)
}

/*
ReturnNextButton returns the fully configured nextButton which is responsible for switch to the next page
*/
func createReturnNextButton(imageFrame *fyne.Container) *widget.Button {
	nextButton := widget.NewButton("Weiter", func() {
		imageFrame.RemoveAll()
		page += 1
		if page >= pages+1 {
			page = 1
		}
		loadCurrentPageImage(imageFrame)
	})
	return nextButton
}

/*
ReturnBackButton returns the fully configured backButton which is responsible for switch to the last page
*/

func createReturnBackButton(imageFrame *fyne.Container) *widget.Button {
	backButton := widget.NewButton("Zurück", func() {
		imageFrame.RemoveAll()
		page -= 1
		if page <= 0 {
			page = pages
		}
		loadCurrentPageImage(imageFrame)
	})
	return backButton
}

/*
Load the image matching the current page into the frame
*/
func loadCurrentPageImage(imageFrame *fyne.Container) {
	var image *canvas.Image

	// Ingenious solution
	switch page {
	case 1:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku1Png)
	case 2:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku2Png)
	case 3:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku3Png)
	case 4:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku4Png)
	case 5:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku5Png)
	case 6:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku6Png)
	case 7:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku7Png)
	case 8:
		image = canvas.NewImageFromResource(yaac_shared.ResourceDoku8Png)
	}

	image.FillMode = canvas.ImageFillContain
	imageFrame.Add(image)
}
