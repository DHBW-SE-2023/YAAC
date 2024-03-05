package yaac_mvvm

import (
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
)

var frontendMain *yaac_frontend_main.FrontendMain = nil

func (m *MVVM) NewFrontendMain() {
	frontendMain = yaac_frontend_main.New(m)
}

func (m *MVVM) OpenMainWindow() {
	frontendMain.OpenMainWindow()
}

func (m *MVVM) NotifyError(err error) {
	panic("Notification service not implemented")
}
