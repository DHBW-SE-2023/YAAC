package yaac_frontend_mail

import (
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type mvvm interface {
	MailFormUpdated(data yaac_shared.EmailData)
	ValidateTable(img []byte)
}

type WindowMail struct {
	MVVM mvvm
}

func New(mvvm mvvm) *WindowMail {
	return &WindowMail{
		MVVM: mvvm,
	}
}
