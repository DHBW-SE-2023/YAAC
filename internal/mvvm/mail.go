package yaac_mvvm

import (
	yaac_backend_mail "github.com/DHBW-SE-2023/YAAC/internal/backend/mail"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

var mailBackend *yaac_backend_mail.BackendMail = nil

func (m *MVVM) NewMailBacked(loginCredentials yaac_shared.EmailData) error {
	b, err := yaac_backend_mail.New(m, loginCredentials.MailServer, loginCredentials.Email, loginCredentials.Password)
	if err != nil {
		return err
	}

	mailBackend = b

	return nil
}

func (m *MVVM) GetMailsToday() ([]yaac_backend_mail.MailData, error) {
	return mailBackend.GetMailsToday()
}

func (m *MVVM) UpdateMailCredentials(credentials yaac_shared.EmailData) error {
	return m.NewMailBacked(credentials)
}

func (m *MVVM) CheckMailConnection() bool {
	return mailBackend.CheckMailConnection()
}
