package yaac_mvvm

import (
	yaac_backend_database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
)

var databaseinst *yaac_backend_database.BackendDatabase

func (m *MVVM) StartApplication() {
	m.ConnectDatabase("data/data.db")
	m.OpenMainWindow()

	print(m.Courses())
}

func (m *MVVM) OpenMainWindow() {
	var frontend = yaac_frontend_main.New(m)
	frontend.OpenMainWindow()
}
