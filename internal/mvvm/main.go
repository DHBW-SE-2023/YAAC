package yaac_mvvm

import (
	yaac_backend_database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
	"log"
	"os"
)

var databaseinst *yaac_backend_database.BackendDatabase

func (m *MVVM) StartApplication() {
	m.CreateOrConnectDatabase()
	m.OpenMainWindow()
}

func (m *MVVM) OpenMainWindow() {
	var frontend = yaac_frontend_main.New(m)
	frontend.OpenMainWindow()
}

func (m *MVVM) CreateOrConnectDatabase() {
	databaseinst = yaac_backend_database.New(m, "./data/", "data.db")

	if _, err := os.Stat("./data/data.db"); err != nil {
		databaseinst.CreateDatabase()
	} else {
		log.Println("Found existing database")
	}

	databaseinst.ConnectDatabase()
}

func (m *MVVM) CloseDatabase() {
	databaseinst.DisconnectDatabase()
}
