package yaac_mvvm

import (
	"time"

	demon "github.com/DHBW-SE-2023/YAAC/internal/backend/demon"
)

// Duation in seconds to wait between refreshes
// The demon is run in a goroutine
func (m *MVVM) StartDemon(duration time.Duration) {
	go demon.StartDemon(m, duration)
}
