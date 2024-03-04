package yaac_mvvm

import (
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
)

// var frontend yaac_frontend_main.FrontendMain

func (m *MVVM) StartApplication() {
	m.ConnectDatabase("data/data.db")
	m.StartDemon(5) // Refresh every 5 seconds

	// Needs to be the last step
	m.OpenMainWindow()
}

func (m *MVVM) OpenMainWindow() {
	var frontend = yaac_frontend_main.New(m)
	frontend.OpenMainWindow()
}
