package yaac_mvvm

import (
	yaac_backend_mail "github.com/DHBW-SE-2023/YAAC/internal/backend/mail"
	yaac_frontend_mail "github.com/DHBW-SE-2023/YAAC/internal/frontend/mail"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func (m *MVVM) MailFormUpdated(data yaac_shared.EmailData) {
	var mail_window = yaac_frontend_mail.New(m)
	var backend, _ = yaac_backend_mail.New(m, "webserver:port", "mail-adresse", "password")

	resp := backend.GetResponse()
	mail_window.UpdateResultLabel(resp)
}
