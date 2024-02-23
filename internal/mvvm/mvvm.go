package yaac_mvvm

// Implements the shared.MVVM interface
type MVVM struct{}

func New() *MVVM {
	return &MVVM{}
}
