package yaac_frontend_main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	yaac_frontend_pages "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

const preferedStartPage = "home"

var gv GlobalVars

type GlobalVars struct {
	App    fyne.App
	Window fyne.Window
}

func (f *FrontendMain) OpenMainWindow() {
	gv = GlobalVars{}
	gv.App = *yaac_shared.GetApp()
	gv.App.Settings().SetTheme(ytheme)

	// setuping window
	gv.Window = gv.App.NewWindow(yaac_shared.APP_NAME)

	// set icon
	gv.Window.SetIcon(yaac_shared.ResourceIconPng)

	// setup systray
	if desk, ok := gv.App.(desktop.App); ok {
		m := fyne.NewMenu(yaac_shared.APP_NAME,
			fyne.NewMenuItem("Show", func() {
				gv.Window.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(yaac_shared.ResourceIconPng)
	}
	gv.Window.SetCloseIntercept(func() {
		gv.Window.Hide()
	})
	// Important setting to enable custom backgrounds without borders
	gv.Window.SetPadded(false)
	gv.Window.Show()
	gv.Window.SetContent(makeWindow())
	gv.Window.Resize(fyne.NewSize(1280, 720))
	gv.Window.Show()

	gv.App.Run()
}

func makeWindow() fyne.CanvasObject {
	content := container.NewStack()
	title := widget.NewLabel("Component name")
	setPage := func(p yaac_frontend_pages.Page) {
		title.SetText(p.Title)

		content.Objects = []fyne.CanvasObject{p.View(gv.Window)}
		content.Refresh()
	}

	page := makeNavPage(setPage, true, content)
	return page
}

func makeNavPage(setPage func(page yaac_frontend_pages.Page), loadPrevious bool, content *fyne.Container) fyne.CanvasObject {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return yaac_frontend_pages.PagesIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := yaac_frontend_pages.PagesIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			p, ok := yaac_frontend_pages.Pages[uid]
			if !ok {
				fyne.LogError("Missing Pages panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(p.Title)
		},
		OnSelected: func(uid string) {
			if p, ok := yaac_frontend_pages.Pages[uid]; ok {
				gv.App.Preferences().SetString(preferedStartPage, uid)
				setPage(p)
			}
		},
	}
	if loadPrevious {
		currentPref := gv.App.Preferences().StringWithFallback(preferedStartPage, "home")
		tree.Select(currentPref)
	}

	logo := canvas.NewImageFromResource(yaac_shared.ResourceDHBWPng)
	navFrame := canvas.NewRectangle(color.White)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 200))
	return container.NewBorder(nil, nil, container.NewStack(navFrame, container.NewBorder(logo, nil, nil, nil, tree)), nil, content)
}
