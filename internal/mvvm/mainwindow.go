package yaac_mvvm

import yaac_frontend "github.com/DHBW-SE-2023/YAAC/internal/frontend"

func (m *MVVM) OpenMainWindow() {
	var frontend = yaac_frontend.New(m)
	frontend.OpenMainWindow()
}
