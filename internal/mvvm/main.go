package yaac_mvvm

import (
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
)

func (m *MVVM) OpenMainWindow() {
	var frontend = yaac_frontend_main.New(m)
	frontend.OpenMainWindow()
}
