package yaac_frontend_template

// Expose functions from the MVVM here
type mvvm interface {
	// Template function within the MVVM exposed to this package
	TemplateFuncRecive(data any)
}

// Rename Me!
type WindowTemplate struct {
	MVVM mvvm
}

func New(mvvm mvvm) *WindowTemplate {
	return &WindowTemplate{
		MVVM: mvvm,
	}
}
