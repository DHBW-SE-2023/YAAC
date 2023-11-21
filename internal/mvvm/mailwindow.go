package yaac_mvvm

import (
	yaac_backend "github.com/DHBW-SE-2023/YAAC/internal/backend"
	yaac_frontend_mail "github.com/DHBW-SE-2023/YAAC/internal/frontend/mail"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func (m *MVVM) MailFormUpdated(data yaac_shared.EmailData) {
	var mail_window = yaac_frontend_mail.New(m)
	var backend = yaac_backend.New(m)

	resp := backend.GetResponse(data)
	mail_window.UpdateResultLabel(resp)
}
