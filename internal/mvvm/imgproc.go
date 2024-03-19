package yaac_mvvm

import (
	"log"

	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
)

var imgprocBackend *imgproc.BackendImgproc = nil

func (m *MVVM) ImgprocBackendStart() {
	imgprocBackend = imgproc.NewBackend(m)
}

func (m *MVVM) NewTable(img []byte) (*imgproc.Table, error) {

	table, err := imgprocBackend.NewTable(img)
	if err != nil {
		log.Fatalf("backend.ValidateTable: %v", err)
		return nil, err
	}

	return table, nil
}
