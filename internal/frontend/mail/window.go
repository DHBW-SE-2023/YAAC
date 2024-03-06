package yaac_frontend_mail

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	resource "github.com/DHBW-SE-2023/YAAC/pkg/resource_manager"
)

var gv GlobalVars

type GlobalVars struct {
	App         fyne.App
	Window      fyne.Window
	ResultLabel *widget.Label
}

func (f *WindowMail) Open() {
	gv = GlobalVars{}
	gv.App = *yaac_shared.GetApp()
	gv.Window = gv.App.NewWindow("Mail Demo")

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	gv.Window.SetIcon(r)

	gv.Window.SetContent(makeFormTab(gv.Window, f))
	gv.Window.Show()
}

func (f *WindowMail) UpdateResultLabel(content string) {
	gv.ResultLabel.SetText(content)
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title: content,
	})
}

func makeFormTab(_ fyne.Window, f *WindowMail) fyne.CanvasObject {
	mailServer := widget.NewEntry()
	mailServer.SetPlaceHolder("John Smith")

	email := widget.NewEntry()
	email.SetPlaceHolder("test@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	gv.ResultLabel = widget.NewLabel("")

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

			// FIXME: Error handling here
			err := f.MVVM.UpdateMailCredentials(formStruct)
			if err != nil {
				log.Fatalf("Could not update E-Mail credentials")
			}
		},
	}
	form.Append("Password", password)
	form.Append("Your first unread message:", gv.ResultLabel)

	return form
}
