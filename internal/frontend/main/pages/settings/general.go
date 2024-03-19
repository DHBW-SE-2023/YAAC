package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func generalScreen() fyne.CanvasObject {
	title := canvas.NewText(" Allgemein", color.Black)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	description := canvas.NewText(" Das Team YAAC präsentiert stolz Yet Another Attendance Checker", color.Black)
	link := canvas.NewText(" Git: https://github.com/DHBW-SE-2023", color.Black)
	description.TextSize = 16
	link.TextSize = 16
	content := container.NewVBox(description, link)
	logo := canvas.NewImageFromResource(yaac_shared.ResourceIconPng)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(600, 600))
	contentFrame := container.NewVBox(content, logo)
	return container.NewGridWrap(fyne.NewSize(800, 200), title, contentFrame)
}
