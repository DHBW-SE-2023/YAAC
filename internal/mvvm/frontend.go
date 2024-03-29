package yaac_mvvm

import (
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

var frontendMain *yaac_frontend_main.FrontendMain = nil

func (m *MVVM) NewFrontendMain() {
	m.FrontendMain = yaac_frontend_main.New(m)
}

func (m *MVVM) NotifyError(source string, err error) {
	frontendMain.ReceiveError(source, err)
}

func (m *MVVM) NotifyNewList(list yaac_shared.AttendanceList) {
	frontendMain.ReceiveNewTable(list)
}
