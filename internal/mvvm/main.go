package yaac_mvvm

import (
	yaac_backend_database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
	"log"
	"os"
)

var databaseinst *yaac_backend_database.BackendDatabase

func (m *MVVM) StartApplication() {
	m.CreateOrConnectDatabase("./data/", "data.db")
	m.OpenMainWindow()
}

func (m *MVVM) OpenMainWindow() {
	var frontend = yaac_frontend_main.New(m)
	frontend.OpenMainWindow()
}

func (m *MVVM) CreateOrConnectDatabase(path string, dbName string) {
	databaseinst = yaac_backend_database.New(m, path, dbName)

	if _, err := os.Stat(path + dbName); err != nil {
		databaseinst.CreateDatabase()
	} else {
		log.Println("Found existing database")
	}

	databaseinst.ConnectDatabase()
}

func (m *MVVM) CloseDatabase() {
	databaseinst.DisconnectDatabase()
}
