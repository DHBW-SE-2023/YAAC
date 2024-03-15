package pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages/settings"
)

func SettingsScreen(_ fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Einstellungen", color.Black)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	settingNav := canvas.NewRectangle(color.NRGBA{R: 230, G: 233, B: 235, A: 255})
	settingNav.Resize(fyne.NewSize(400, 400))

	content := container.NewStack()

	setContent := func(s settings.Setting) {
		content.Objects = []fyne.CanvasObject{s.View()}
		content.Refresh()
	}

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return settings.SettingPagesIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := settings.SettingPagesIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			p, ok := settings.SettingPages[uid]
			if !ok {
				fyne.LogError("Missing Pages panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(p.Title)
		},
		OnSelected: func(uid string) {
			if p, ok := settings.SettingPages[uid]; ok {
				setContent(p)
			}
		},
	}
	navBar := container.NewGridWrap((fyne.NewSize(300, 200)), title, tree)
	navFrame := container.NewHBox(container.NewStack(settingNav, navBar))
	contentFrame := canvas.NewRectangle(color.NRGBA{R: 125, G: 136, B: 142, A: 255})
	page := container.NewBorder(nil, nil, navFrame, nil, container.NewMax(contentFrame, content))
	return page
}
