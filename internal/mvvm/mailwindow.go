package yaac_mvvm

import (
	yaac_backend "github.com/DHBW-SE-2023/YAAC/internal/backend"
	yaac_frontend "github.com/DHBW-SE-2023/YAAC/internal/frontend"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func (m *MVVM) MailFormUpdated(data yaac_shared.EmailData) {
	var frontend = yaac_frontend.New(m)
	var backend = yaac_backend.New(m)

	resp := backend.GetResponse(data)
	frontend.UpdateResultLabel(resp)
}
