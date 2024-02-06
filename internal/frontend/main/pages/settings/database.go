package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func databaseScreen() fyne.CanvasObject {
	title := canvas.NewText(" Datenbank", color.Black)
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	header := container.NewGridWrap(fyne.NewSize(800, 200), title)

	labelDatabaseConnection := widget.NewLabel("Datenbank Status")
	labelDatabaseConnection.TextStyle = fyne.TextStyle{Bold: true}
	labelDatabaseConnectionStatus := widget.NewLabel("Aktiv")
	databaseConnectionStatus := container.NewVBox(labelDatabaseConnection, container.NewGridWrap(fyne.NewSize(400, 40), labelDatabaseConnectionStatus))
	labelDatabaseConnectionString := widget.NewLabel("Datenbank Verbindung")
	labelDatabaseConnectionString.TextStyle = fyne.TextStyle{Bold: true}
	entryDatabaseConnectionString := widget.NewEntry()
	databaseConnection := container.NewVBox(labelDatabaseConnectionString, container.NewGridWrap(fyne.NewSize(400, 40), entryDatabaseConnectionString))

	buttonFrame := canvas.NewRectangle(color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	pullDatabaseButton := widget.NewButton("Aktualisiere Daten", func() { println("Fresh Pull") })
	resetDatabaseButton := widget.NewButton("Setze Datenbank zurück", func() { println("Restart") })
	buttons := container.NewGridWrap(fyne.NewSize(250, 50), container.NewMax(buttonFrame, pullDatabaseButton), container.NewMax(buttonFrame, resetDatabaseButton))
	return container.NewVBox(header, databaseConnectionStatus, databaseConnection, buttons)
}
