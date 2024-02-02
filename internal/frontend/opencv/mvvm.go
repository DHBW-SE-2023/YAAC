package yaac_frontend_opencv

import (
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type mvvm interface {
	MailFormUpdated(data yaac_shared.EmailData)
	StartGoCV(img_path string)
}

type WindowOpenCV struct {
	MVVM mvvm
}

func New(mvvm mvvm) *WindowOpenCV {
	return &WindowOpenCV{
		MVVM: mvvm,
	}
}
