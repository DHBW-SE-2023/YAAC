package yaac_frontend_main

import (
	pages "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages"
	shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

type FrontendMain struct {
	MVVM shared.MVVM
}

func New(mvvm shared.MVVM) *FrontendMain {
	pages.New(mvvm)
	return &FrontendMain{
		MVVM: mvvm,
	}
}

func (*FrontendMain) ReceiveNewTable(table shared.AttendanceList) {
	panic("Not yet implemented")
}
