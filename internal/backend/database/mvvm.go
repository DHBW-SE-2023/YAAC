package yaac_backend_database

import (
	"gorm.io/gorm"
)

type mvvm interface {
}

type BackendDatabase struct {
	MVVM mvvm
	Path string
	DB   *gorm.DB
}

func NewBackend(mvvm mvvm, path string) *BackendDatabase {
	item := BackendDatabase{
		MVVM: mvvm,
		Path: path,
	}

	return &item
}
