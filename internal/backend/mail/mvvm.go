package yaac_backend_mail

type mvvm interface {
}

type BackendMail struct {
	MVVM mvvm
}

func New(mvvm mvvm) *BackendMail {
	return &BackendMail{
		MVVM: mvvm,
	}
}
