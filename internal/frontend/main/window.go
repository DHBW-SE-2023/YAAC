package yaac_frontend_main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	resource "github.com/DHBW-SE-2023/YAAC/pkg/resource_manager"
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
	// Important setting to enable custom backgrounds without borders
	gv.Window.SetPadded(false)

	// handle main window
	gv.Window.SetContent(makeWindow(f))
	gv.Window.Resize(fyne.NewSize(1280, 720))
	gv.Window.Show()

	gv.App.Run()
}

func makeWindow(f *FrontendMain) fyne.CanvasObject {
	content := container.NewStack()
	title := widget.NewLabel("Component name")
	setPage := func(p pages.Page) {
		title.SetText(p.Title)

		content.Objects = []fyne.CanvasObject{p.View(gv.Window)}
		content.Refresh()
	}

	page := container.NewBorder(
		nil, nil, nil, nil, content)

	return container.NewBorder(nil, nil, makeNav(setPage, true), nil, page)
}

func makeNav(setPage func(page pages.Page), loadPrevious bool) fyne.CanvasObject {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return pages.PagesIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := pages.PagesIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			p, ok := pages.Pages[uid]
			if !ok {
				fyne.LogError("Missing Pages panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(p.Title)
		},
		OnSelected: func(uid string) {
			if p, ok := pages.Pages[uid]; ok {
				gv.App.Preferences().SetString(preferedStartPage, uid)
				setPage(p)
			}
		},
	}

	if loadPrevious {
		currentPref := gv.App.Preferences().StringWithFallback(preferedStartPage, "home")
		tree.Select(currentPref)
	}

	logo := canvas.NewImageFromFile("assets/DHBW.png")
	navFrame := canvas.NewRectangle(color.White)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(200, 200))

	return container.NewMax(navFrame, container.NewBorder(logo, nil, nil, nil, tree))
}
