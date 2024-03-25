package yaac_mvvm

import (
	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
)

<<<<<<< Updated upstream
func (m *MVVM) ImgprocBackendStart() {
	m.BackendImgproc = imgproc.NewBackend(m)
	m.NewTable([]byte{})
=======
func (m *MVVM) ValidateTable(img []byte) {
	b := imgproc.NewBackend(m)
	f := frontend.New(m)

	go func() {
		table, err := b.ValidateTable(img)
		if err != nil {
			log.Println("backend.ValidateTable: ", err)
		}

		f.ReceiveNewTable(table)
	}()
>>>>>>> Stashed changes
}
