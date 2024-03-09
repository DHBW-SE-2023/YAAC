package yaac_frontend_mail

import (
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type WindowMail struct {
	MVVM yaac_shared.MVVM
}

func New(mvvm yaac_shared.MVVM) *WindowMail {
	return &WindowMail{
		MVVM: mvvm,
	}
}
