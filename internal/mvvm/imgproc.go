package yaac_mvvm

import (
	"log"

	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
)

var imgprocBackend *imgproc.BackendImgproc = nil

func (m *MVVM) ImgprocBackendStart() {
	imgprocBackend = imgproc.NewBackend(m)
}

func (m *MVVM) ValidateTable(img []byte) (imgproc.Table, error) {

	table, err := imgprocBackend.ValidateTable(img)
	if err != nil {
		log.Fatalf("backend.ValidateTable: %v", err)
		return imgproc.Table{}, err
	}

	return table, nil
}
