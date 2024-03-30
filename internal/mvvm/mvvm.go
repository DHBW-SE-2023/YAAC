package yaac_mvvm

import (
	yaac_backend_database "github.com/DHBW-SE-2023/YAAC/internal/backend/database"
	yaac_backend_imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
	yaac_backend_mail "github.com/DHBW-SE-2023/YAAC/internal/backend/mail"
	yaac_frontend_main "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
)

// Implements the shared.MVVM interface
type MVVM struct {
	*yaac_backend_imgproc.BackendImgproc
	*yaac_backend_mail.BackendMail
	*yaac_frontend_main.FrontendMain
	*yaac_backend_database.BackendDatabase
}

func New() *MVVM {
	return &MVVM{}
}
