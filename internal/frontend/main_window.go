package yaac_frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	yaac_frontend_mail "github.com/DHBW-SE-2023/YAAC/internal/frontend/mail"
	yaac_frontend_opencv "github.com/DHBW-SE-2023/YAAC/internal/frontend/opencv"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	resource "github.com/DHBW-SE-2023/YAAC/pkg/resource_manager"
)

var App fyne.App
var mainWindow fyne.Window

func (f *Frontend) OpenMainWindow() {
	App = *yaac_shared.GetApp()

	// setuping window
	mainWindow = App.NewWindow(yaac_shared.APP_NAME)

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	mainWindow.SetIcon(r)

	// setup systray
	if desk, ok := App.(desktop.App); ok {
		m := fyne.NewMenu(yaac_shared.APP_NAME,
			fyne.NewMenuItem("Show", func() {
				mainWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(r)
	}
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	// handle main window
	mainWindow.SetContent(makeMainWindow(f))
	mainWindow.Show()

	App.Run()
}

func makeMainWindow(f *Frontend) *fyne.Container {
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
}
