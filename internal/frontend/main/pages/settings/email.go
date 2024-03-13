package settings

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func emailScreen() fyne.CanvasObject {
	title := canvas.NewText(" Email", color.Black)
	title.TextSize = 28
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	header := container.NewCenter(container.NewGridWrap(fyne.NewSize(800, 200), title))

	form := generateForm()
	content := container.NewCenter(container.NewVBox(container.NewGridWrap(fyne.NewSize(600, 40), form)))
	return container.NewVBox(header, content)
}

func generateForm() *widget.Form {
	serverStatus := widget.NewLabel("Läuft Sascha")
	server := widget.NewEntry()
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()
	restartButton := widget.NewButton("Zurücksetzen", func() {
		println("")
	})
	submitButton := widget.NewButton("Bestätigen", nil)
	buttonArea := container.NewAdaptiveGrid(2, submitButton, restartButton)
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Status", Widget: serverStatus},
			{Text: "E-Mail Server", Widget: server},
			{Text: "E-Mail User", Widget: username},
			{Text: "E-Mail Password", Widget: password},
			{Widget: buttonArea}},
	}
	submitButton.OnTapped = func() {
		serverConnection := server.Text
		serverUser := username.Text
		password := password.Text
		setSetting("mailConnection", serverConnection)
		setSetting("mailUser", serverUser)
		setSetting("mailPassword", password)
	}
	return form
}

func setSetting(key string, value string) {
	var settings []yaac_shared.Setting
	setting := yaac_shared.Setting{
		Setting: key,
		Value:   value,
	}
	settings = append(settings, setting)
	_, err := myMVVM.SettingsUpdate(settings)
	if err != nil {
		yaac_shared.App.SendNotification(fyne.NewNotification("Fehler beim Aktualisieren", err.Error()))
	} else {
		yaac_shared.App.SendNotification(fyne.NewNotification("Erfolgreiche Aktualisierung", "Ihre Mail Daten wurden erfolgreich aktualisiert"))
	}
}
