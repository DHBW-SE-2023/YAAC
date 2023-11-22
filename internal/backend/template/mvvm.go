package yaac_backend_template

// Expose functions from the MVVM here
type mvvm interface {
	// Template function within the MVVM exposed to this package
	TemplateFuncRecive(data any)
}

// Rename Me!
type BackendTemplate struct {
	MVVM mvvm
}

func New(mvvm mvvm) *BackendTemplate {
	return &BackendTemplate{
		MVVM: mvvm,
	}
}
