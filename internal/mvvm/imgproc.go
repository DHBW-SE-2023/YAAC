package yaac_mvvm

import (
	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
)

func (m *MVVM) ImgprocBackendStart() {
	m.BackendImgproc = imgproc.NewBackend(m)
	m.NewTable([]byte{})
}
