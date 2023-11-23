package yaac_frontend_main

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

var gv GlobalVars

type GlobalVars struct {
	App    fyne.App
	Window fyne.Window
}

func (f *FrontendMain) OpenMainWindow() {
	gv = GlobalVars{}
	gv.App = *yaac_shared.GetApp()

	// setuping window
	gv.Window = gv.App.NewWindow(yaac_shared.APP_NAME)

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	gv.Window.SetIcon(r)

	// setup systray
	if desk, ok := gv.App.(desktop.App); ok {
		m := fyne.NewMenu(yaac_shared.APP_NAME,
			fyne.NewMenuItem("Show", func() {
				gv.Window.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(r)
	}
	gv.Window.SetCloseIntercept(func() {
		gv.Window.Hide()
	})

	// handle main window
	gv.Window.SetContent(makeWindow(f))
	gv.Window.Show()

	gv.App.Run()
}

func makeWindow(f *FrontendMain) *fyne.Container {
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
