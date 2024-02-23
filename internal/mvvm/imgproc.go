package yaac_mvvm

import (
	"log"

	imgproc "github.com/DHBW-SE-2023/YAAC/internal/backend/imgproc"
	frontend "github.com/DHBW-SE-2023/YAAC/internal/frontend/main"
)

func (m *MVVM) ValidateTable(img []byte) {
	b := imgproc.NewBackend(m)
	f := frontend.New(m)

	go func() {
		table, err := b.ValidateTable(img)
		if err != nil {
			log.Fatalf("backend.ValidateTable: %v", err)
		}

		f.ReceiveNewTable(table)
	}()
}
