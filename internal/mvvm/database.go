package yaac_mvvm

import (
	database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
)

func (m *MVVM) ConnectDatabase(dbPath string) error {
	m.BackendDatabase = database.NewBackend(m, dbPath)
	return m.BackendDatabase.ConnectDatabase()
}
