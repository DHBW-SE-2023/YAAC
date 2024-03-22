package yaac_frontend_settings

import (
	"errors"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func emailScreen(w fyne.Window) fyne.CanvasObject {
	title := ReturnHeader("Email")
	form := ReturnForm(w)
	content := container.NewCenter(container.NewVBox(container.NewGridWrap(fyne.NewSize(600, 40), form)))
	return container.NewVBox(container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 200), title)), widget.NewSeparator(), content)
}

/*
ReturnForm returns the fully configured Email Form responsible for managing and updating mail settings
*/
func ReturnForm(w fyne.Window) *widget.Form {
	mailConnection, mailUser, mailPassword := ReturnMailSettings()
	// var serverStatus *widget.Label
	serverStatus := container.NewHBox()
	if mailConnection != "" && mailUser != "" && mailPassword != "" {
		alive := myMVVM.CheckMailConnection()
		serverStatusText := widget.NewLabel(MapMailBooleans(alive))
		serverStatusImage := loadImage("assets/alive.png")
		serverStatus.Add(serverStatusImage)
		serverStatus.Add(serverStatusText)
	} else {
		serverStatusText := widget.NewLabel("Keine Daten hinterlegt")
		serverStatusImage := loadImage("assets/down.png")
		serverStatus.Add(serverStatusImage)
		serverStatus.Add(serverStatusText)
	}
	form := ConfigureForm(w, mailConnection, mailUser, mailPassword, serverStatus)
	return form
}

/*
UpdateSetting collects the currently selected setting values, intialize a Setting Object
for each and updates the changes on the database
*/
func UpdateSetting(w fyne.Window, key string, value string) {
	var settings []yaac_shared.Setting
	setting := yaac_shared.Setting{
		Setting: key,
		Value:   value,
	}
	settings = append(settings, setting)
	_, err := myMVVM.SettingsUpdate(settings)
	if err != nil {
		dialog.ShowError(err, w)
	} else {
		dialog.ShowInformation("Erfolgreiche Aktualisierung", "Ihre Mail Daten wurden erfolgreich aktualisiert", w)
	}
}

/*
ReturnMailSettings makes a lookup in the setting DB to assign each key to it's respective value returning all assigned values
*/
func ReturnMailSettings() (string, string, string) {
	setting, _ := myMVVM.Settings()
	var mailConnection string
	var mailUser string
	var mailPassword string
	for _, element := range setting {
		if element.Setting == "mailConnection" {
			mailConnection = element.Value
		} else if element.Setting == "mailUser" {
			mailUser = element.Value
		} else {
			mailPassword = element.Value
		}
	}
	return mailConnection, mailUser, mailPassword
}

/*
ConfigureForm configures the mailForm regarding input as well as functionalities passing mailConnection, mailUser, mailPassword, serverStatus
returning the fully configure form.
*/
func ConfigureForm(w fyne.Window, mailConnection string, mailUser string, mailPassword string, serverStatus *fyne.Container) *widget.Form {
	server := widget.NewEntry()
	server.SetText(mailConnection)
	username := widget.NewEntry()
	username.SetText(mailUser)
	password := widget.NewPasswordEntry()
	password.SetText(mailPassword)
	ConfigureInputValidation(server, username, password)
	restartButton := widget.NewButton("Zurücksetzen", func() {
		_, err := myMVVM.SettingsReset()
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Email Einstellungen wurden erfolgreich zurückgesetzt", "", w)
		}
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
	ConfigureSubmitButton(w, server, username, password, submitButton)
	return form
}

/*
ConfigureSubmitButton configures the submitButton responsible for udpating the settings DB passing sever(Entry),username(Entry),password(Entry)as well as the sumbitButton itself.
returning the fully configure form.
*/
func ConfigureSubmitButton(w fyne.Window, server *widget.Entry, username *widget.Entry, password *widget.Entry, submitButton *widget.Button) {
	submitButton.OnTapped = func() {
		if CheckInputValidity([]*widget.Entry{server, username, password}) != nil {
			dialog.ShowInformation("Überprüfen sie ihre Eingaben", "Fehler bei Udpate", w)
		} else {
			serverConnection := server.Text
			serverUser := username.Text
			password := password.Text
			UpdateSetting(w, "mailConnection", serverConnection)
			UpdateSetting(w, "mailUser", serverUser)
			UpdateSetting(w, "mailPassword", password)
			myMVVM.UpdateMailCredentials(yaac_shared.MailLoginData{MailServer: serverConnection, Email: serverUser, Password: password})
		}
	}
}

/*
ConfigureInputValidation takes every widget.Entry and assigns a specific Inputvalidater reducing the max input as well
as the input regarding syntax passing all entries
*/
func ConfigureInputValidation(server *widget.Entry, username *widget.Entry, password *widget.Entry) {
	InputValidater(server, `^[a-zA-Z0-9-._]+:[0-9]+$`, 30, "Geben sie einen validen Servernamen ein.")
	InputValidater(username, `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`, 30, "Geben sie einen validen Usernamen ein.")
	InputValidater(password, `.*`, 30, "")
}

/*
InputValidater defines the actual Vaildatorfunctions for ConfigureInputValidation passing the entry, the regex matching schema, the maxLen
as well as a message for failure.
*/
func InputValidater(entry *widget.Entry, regex string, maxLen int, failure string) {
	entry.Validator = func(s string) error {
		re, _ := regexp.Compile(regex)
		if !re.MatchString(s) {
			return errors.New(failure)
		}
		return nil
	}
	entry.OnChanged = func(s string) {
		if len(s) > maxLen {
			s = s[0:maxLen]
			entry.SetText(s)
		}

	}
}

/*
CheckInputValidity checks the validity of each entry before Updating the setting, passing all the entries as []*widget.Entry
*/
func CheckInputValidity(entries []*widget.Entry) error {
	for _, element := range entries {
		if element.Validate() != nil {
			return element.Validate()
		}
	}
	return nil
}
