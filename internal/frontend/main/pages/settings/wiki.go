package yaac_frontend_settings

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

var page int = 1
var pages int = 24

func wikiScreen(_ fyne.Window) fyne.CanvasObject {
	title := ReturnHeader("Nutzerdokumentation")

	// This is neccecary because the frame wont accept an uninitialized image
	image := canvas.NewImageFromResource(yaac_shared.ResourceYaacManual01Png)
	image.FillMode = canvas.ImageFillContain
	imageFrame := container.NewGridWrap(fyne.NewSize(1200, 1000), image)

	currentPage := widget.NewLabel(fmt.Sprintf("%s %d", "Seite", page))
	nextButton := createReturnNextButton(imageFrame, currentPage)
	backButton := createReturnBackButton(imageFrame, currentPage)

	currentPage.Importance = widget.MediumImportance
	currentPage.TextStyle = fyne.TextStyle{Bold: true}

	buttonArea := container.NewCenter(container.NewHBox(container.NewAdaptiveGrid(5, backButton, layout.NewSpacer(), currentPage, layout.NewSpacer(), nextButton)))
	content := container.NewStack(container.NewVBox(container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 200), title)), widget.NewSeparator(), container.NewCenter(imageFrame), buttonArea))
	return container.NewVScroll(content)
}

/*
ReturnNextButton returns the fully configured nextButton which is responsible for switch to the next page
*/
func createReturnNextButton(imageFrame *fyne.Container, currentPage *widget.Label) *widget.Button {
	nextButton := widget.NewButton("Weiter", func() {
		imageFrame.RemoveAll()
		page += 1
		if page >= pages+1 {
			page = 1
		}
		currentPage.SetText(fmt.Sprintf("%s %d", "Seite", page))
		loadCurrentPageImage(imageFrame)
	})
	return nextButton
}

/*
ReturnBackButton returns the fully configured backButton which is responsible for switch to the last page
*/

func createReturnBackButton(imageFrame *fyne.Container, currentPage *widget.Label) *widget.Button {
	backButton := widget.NewButton("Zurück", func() {
		imageFrame.RemoveAll()
		page -= 1
		if page <= 0 {
			page = pages
		}
		currentPage.SetText(fmt.Sprintf("%s %d", "Seite", page))
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
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual01Png)
	case 2:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual02Png)
	case 3:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual03Png)
	case 4:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual04Png)
	case 5:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual05Png)
	case 6:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual06Png)
	case 7:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual07Png)
	case 8:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual08Png)
	case 9:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual09Png)
	case 10:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual10Png)
	case 11:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual11Png)
	case 12:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual12Png)
	case 13:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual13Png)
	case 14:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual14Png)
	case 15:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual15Png)
	case 16:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual16Png)
	case 17:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual17Png)
	case 18:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual18Png)
	case 19:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual19Png)
	case 20:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual20Png)
	case 21:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual21Png)
	case 22:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual22Png)
	case 23:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual23Png)
	case 24:
		image = canvas.NewImageFromResource(yaac_shared.ResourceYaacManual24Png)
	}

	image.FillMode = canvas.ImageFillContain
	imageFrame.Add(image)
}
