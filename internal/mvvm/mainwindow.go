package yaac_mvvm

import yaac_frontend "github.com/DHBW-SE-2023/yaac-go-prototype/internal/frontend"

func (m *MVVM) OpenMainWindow() {
	var frontend = yaac_frontend.New(m)
	frontend.OpenMainWindow()
}
