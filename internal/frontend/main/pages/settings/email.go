package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func emailScreen() fyne.CanvasObject {
	title := canvas.NewText(" Email", color.Black)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	header := container.NewGridWrap(fyne.NewSize(800, 200), title)

	labelMailConnection := widget.NewLabel("Mail Server Status")
	labelMailConnection.TextStyle = fyne.TextStyle{Bold: true}
	labelMailConnectionStatus := widget.NewLabel("Aktiv")
	mailConnectionStatus := container.NewVBox(labelMailConnection, container.NewGridWrap(fyne.NewSize(400, 40), labelMailConnectionStatus))
	labelMailConnectionString := widget.NewLabel("Mail Server Verbindung")
	labelMailConnectionString.TextStyle = fyne.TextStyle{Bold: true}
	entryMailConnectionString := widget.NewEntry()
	mailConnection := container.NewVBox(labelMailConnectionString, container.NewGridWrap(fyne.NewSize(400, 40), entryMailConnectionString))

	buttonFrame := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	pullMailButton := widget.NewButton("Aktualisiere Daten", func() { println("Fresh Pull") })
	resetMailButton := widget.NewButton("Mail Sever Neustart", func() { println("Restart") })
	buttons := container.NewGridWrap(fyne.NewSize(250, 50), container.NewMax(buttonFrame, pullMailButton), container.NewMax(buttonFrame, resetMailButton))
	return container.NewVBox(header, mailConnectionStatus, mailConnection, buttons)
}
