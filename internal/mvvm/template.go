package yaac_mvvm

import (
	yaac_backend_template "github.com/DHBW-SE-2023/YAAC/internal/backend/template"
	yaac_frontend_template "github.com/DHBW-SE-2023/YAAC/internal/frontend/template"
)

// Template function in the MVVM.
// Packages can include this function in their MVVM interfaces
func (m *MVVM) TemplateFuncSend() {
	// Do something
	backend := yaac_backend_template.New(m)
	frontend := yaac_frontend_template.New(m)
	backend.Foo()
	frontend.Open()
}

func (m *MVVM) TemplateFuncRecive(data any) {
	// Do something with data
}
