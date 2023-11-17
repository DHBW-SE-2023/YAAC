package yaac_frontend

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/yaac-go-prototype/internal/shared"
	resource "github.com/DHBW-SE-2023/yaac-go-prototype/pkg/resource_manager"
)

var mailWindow fyne.Window
var result_label *widget.Label

func (f *Frontend) OpenMailWindow() {
	mailWindow = App.NewWindow("Mail Demo")

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	mailWindow.SetIcon(r)

	mailWindow.SetContent(makeFormTab(mailWindow, f))
	mailWindow.Show()
}

func (f *Frontend) UpdateResultLabel(content string) {
	result_label.SetText(content)
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title: content,
	})
}

func makeFormTab(_ fyne.Window, f *Frontend) fyne.CanvasObject {
	mailServer := widget.NewEntry()
	mailServer.SetPlaceHolder("John Smith")

	email := widget.NewEntry()
	email.SetPlaceHolder("test@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	result_label = widget.NewLabel("")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "MailServer", Widget: mailServer, HintText: "Specify the address of your mail server with ':port'"},
			{Text: "Email", Widget: email, HintText: "Your email address"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {

			formStruct := yaac_shared.EmailData{
				MailServer: mailServer.Text,
				Email:      email.Text,
				Password:   password.Text,
			}
			f.MVVM.MailFormUpdated(formStruct)
		},
	}
	form.Append("Password", password)
	form.Append("Your first unread message:", result_label)

	return form
}
