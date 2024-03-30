package yaac_mvvm

import (
	yaac_backend_mail "github.com/DHBW-SE-2023/YAAC/internal/backend/mail"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

// Tries to connect to the mail server. Even if the mail server is unreachable (error), the mail backend is still valid
func (m *MVVM) NewMailBacked(loginCredentials yaac_shared.MailLoginData) error {
	backend, err := yaac_backend_mail.New(m, loginCredentials.MailServer, loginCredentials.Email, loginCredentials.Password)
	m.BackendMail = backend
	if err != nil {
		return err
	}

	return nil
}

func (m *MVVM) UpdateMailCredentials(credentials yaac_shared.MailLoginData) error {
	backend, err := yaac_backend_mail.New(m, credentials.MailServer, credentials.Email, credentials.Password)
	if err != nil {
		return err
	}

	m.BackendMail = backend

	return nil
}
