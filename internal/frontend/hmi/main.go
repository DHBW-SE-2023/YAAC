// Package main provides various examples of Fyne API capabilities.
package main

import (
	"TestFyne/pages"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

func main() {
	a := app.NewWithID("yaac.proto")
	logLifecycle(a)
	w := a.NewWindow("YAAC")
	topWindow = w
	w.SetMaster()

	dialogTitle := canvas.NewText("Willkommen", color.Black)
	dialogTitle.TextSize = 20
	dialogTitle.Alignment = fyne.TextAlignCenter
	dialogTitle.TextStyle = fyne.TextStyle{Bold: true, Italic: false, Monospace: false}
	dialogContent := canvas.NewText("Es sind x neue Listen gekommen", color.Black)
	dialogContent.TextSize = 14
	dialogContent.Alignment = fyne.TextAlignCenter

	data := container.NewBorder(dialogTitle, nil, nil, nil, dialogContent)
	dialogWindow := dialog.NewCustom("", "Okay", data, w)
	dialogWindow.Resize(fyne.NewSize(400, 200))
	dialogWindow.Show()

	content := container.NewStack()
	title := widget.NewLabel("Component name")
	setPage := func(p pages.Page) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(p.Title)
			topWindow = child
			child.SetContent(p.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(p.Title)

		content.Objects = []fyne.CanvasObject{p.View(w)}
		content.Refresh()
	}
	page := container.NewBorder(
		nil, nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setPage, false))
	} else {
		split := container.NewHSplit(makeNav(setPage, true), page)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(1280, 920))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeNav(setPage func(page pages.Page), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

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
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setPage(p)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	logo := canvas.NewImageFromFile("data/DHBW.png")
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(200, 200))
	}

	return container.NewBorder(logo, nil, nil, nil, tree)
}
