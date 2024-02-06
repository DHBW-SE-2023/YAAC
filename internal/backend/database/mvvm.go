package yaac_backend_database

import "database/sql"

type mvvm interface {
}

type BackendDatabase struct {
	MVVM     mvvm
	path     string
	dbName   string
	database *sql.DB
}

func New(mvvm mvvm, path string, dbName string) *BackendDatabase {
	item := BackendDatabase{
		MVVM:   mvvm,
		path:   path,
		dbName: dbName,
	}

	return &item
}
