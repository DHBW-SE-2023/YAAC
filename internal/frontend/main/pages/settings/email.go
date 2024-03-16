package yaac_frontend_settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func emailScreen() fyne.CanvasObject {
	title := ReturnHeader("Email")
	form := ReturnForm()
	content := container.NewCenter(container.NewVBox(container.NewGridWrap(fyne.NewSize(600, 40), form)))
	return container.NewVBox(title, content)
}

/*
ReturnForm returns the fully configured Email Form responsible for managing and updating mail settings
*/
func ReturnForm() *widget.Form {
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
		UpdateSetting("mailConnection", serverConnection)
		UpdateSetting("mailUser", serverUser)
		UpdateSetting("mailPassword", password)
	}
	return form
}

/*
UpdateSetting collects the currently selected setting values, intialize a Setting Object
for each and updates the changes on the database
*/
func UpdateSetting(key string, value string) {
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
