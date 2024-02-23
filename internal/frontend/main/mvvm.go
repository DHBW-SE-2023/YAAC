package yaac_frontend_main

import (
	"log"

	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type mvvm interface {
	MailFormUpdated(data shared.EmailData)
	ValidateTable(img []byte)
}

type FrontendMain struct {
	MVVM mvvm
}

func New(mvvm mvvm) *FrontendMain {
	return &FrontendMain{
		MVVM: mvvm,
	}
}

func (*FrontendMain) ReceiveNewTable(table shared.Table) {
	log.Fatal("Not yet implemented")
}
