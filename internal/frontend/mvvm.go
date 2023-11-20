package yaac_frontend

import (
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type mvvm interface {
	MailFormUpdated(data yaac_shared.EmailData)
	StartGoCV(img_path string)
}

type Frontend struct {
	MVVM mvvm
}

func New(mvvm mvvm) *Frontend {
	return &Frontend{
		MVVM: mvvm,
	}
}
