package yaac_frontend_main

import (
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type mvvm interface {
	MailFormUpdated(data yaac_shared.EmailData)
	StartGoCV(img_path string)
}

type FrontendMain struct {
	MVVM mvvm
}

func New(mvvm mvvm) *FrontendMain {
	return &FrontendMain{
		MVVM: mvvm,
	}
}
