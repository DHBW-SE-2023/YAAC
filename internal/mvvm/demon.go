package yaac_mvvm

import (
	"time"

	demon "github.com/DHBW-SE-2023/YAAC/internal/backend/demon"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

// Duation in seconds to wait between refreshes
// The demon is run in a goroutine
func (m *MVVM) StartDemon(duration time.Duration) {
	go demon.StartDemon(m, duration)
}

func (m *MVVM) UploadImage(img []byte, course *yaac_shared.Course) (*yaac_shared.AttendanceList, error) {
	return demon.UploadImage(m, img, course)
}
