package pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingsScreen(_ fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Einstellungen", color.Black)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	settingNav := canvas.NewRectangle(color.NRGBA{R: 230, G: 233, B: 235, A: 255})
	settingNav.Resize(fyne.NewSize(400, 400))
	var settingOptions = []string{"Allgemein", "Datenbank", "Email", "Wiki", "Impressum"}
	settingList := widget.NewList(
		func() int {
			return len(settingOptions)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(li widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(settingOptions[li])
		},
	)
	navBar := container.NewGridWrap((fyne.NewSize(300, 200)), title, settingList)
	navFrame := container.NewHBox(container.NewStack(settingNav, navBar))
	settingsContent := canvas.NewRectangle(color.NRGBA{R: 125, G: 136, B: 142, A: 255})
	logo := canvas.NewImageFromFile("assets/Icon.png")
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 200))
	contentFrame := container.NewStack(settingsContent, logo)
	content := container.NewBorder(nil, nil, navFrame, nil, contentFrame)

	/*
		header := widget.NewLabel("Select an action:")
		mail_button := widget.NewButton(
			"Open Mail Window",
			yaac_frontend_mail.New(f.MVVM).Open,
		)
		opencv_button := widget.NewButton(
			"Open OpenCV Demo Window",
			yaac_frontend_opencv.New(f.MVVM).Open,
		)

		return container.NewVBox(
			header,
			mail_button,
			opencv_button,
		)
	*/

	return content
}
