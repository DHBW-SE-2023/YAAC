package yaac_mvvm

import (
	yaac_backend_mail "github.com/DHBW-SE-2023/YAAC/internal/backend/mail"
	yaac_frontend_mail "github.com/DHBW-SE-2023/YAAC/internal/frontend/mail"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

// can be deleted
func (m *MVVM) MailFormUpdated(data yaac_shared.EmailData) {
	var mail_window = yaac_frontend_mail.New(m)

	//if New() returns an error do NEVER use the yaac_backend_mail object!!!
	backend, err := yaac_backend_mail.New(m, "webserver:port", "mailadress", "password")
	if err != nil {
		return
	}

	// can be deleted
	test_mail_data := yaac_shared.EmailData{
		MailServer: "webserver:port",
		Email:      "mailadress",
		Password:   "password",
	}
	resp := backend.GetResponse(test_mail_data)
	mail_window.UpdateResultLabel(resp)
}

// TODO: Refresh mails, extract mails that we are interessted in
// extract the image and pass it to the imgproc
func (m *MVVM) MailUpdateData() {
	panic("Not implemented")
}

func (m *MVVM) MailRestartServer() {
	panic("Not implemented")
}

func (m *MVVM) MailsRefresh() []yaac_shared.Email {
	m.MailUpdateData()
	panic("Not implemented")
}
