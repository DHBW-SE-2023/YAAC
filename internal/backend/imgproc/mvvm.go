package yaac_backend_imgproc

type mvvm interface {
}

type BackendOpenCV struct {
	MVVM mvvm
}

func New(mvvm mvvm) *BackendOpenCV {
	return &BackendOpenCV{
		MVVM: mvvm,
	}
}
