package yaac_frontend_settings

import (
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type FrontendMain struct {
	MVVM shared.MVVM
}

var myMVVM shared.MVVM = nil

func New(mvvm shared.MVVM) *FrontendMain {
	myMVVM = mvvm
	return &FrontendMain{
		MVVM: mvvm,
	}
}
