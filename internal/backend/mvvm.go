package yaac_backend

type mvvm interface {
}

type Backend struct {
	MVVM mvvm
}

func New(mvvm mvvm) *Backend {
	return &Backend{
		MVVM: mvvm,
	}
}
