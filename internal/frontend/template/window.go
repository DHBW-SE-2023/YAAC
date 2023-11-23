package yaac_frontend_template

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	resource "github.com/DHBW-SE-2023/YAAC/pkg/resource_manager"
)

var gv GlobalVars

// Put all variables you want to expose to the package here
type GlobalVars struct {
	App    fyne.App
	Window fyne.Window
}

// This function MUST be the first access to this package
func (f *WindowTemplate) Open() {
	gv = GlobalVars{}
	gv.App = *yaac_shared.GetApp()

	// setuping window
	gv.Window = gv.App.NewWindow("TEMPLATE")

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	gv.Window.SetIcon(r)

	// handle main window
	gv.Window.SetContent(makeWindow(f))
	//gv.Window.Resize(fyne.NewSize(800, 600))
	gv.Window.Show()

	gv.App.Run()
}

// Outsources UI building
func makeWindow(f *WindowTemplate) *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("TEMPLATE"),
	)
}

// Accessing content in the MVVM
func (f *WindowTemplate) AccessMVVM() {
	f.MVVM.TemplateFuncRecive(false)
}

/*
# Functions

func (f *WindowTemplate) ExposedFunction() {
	// This function is Exposed to the MVVM using WindowTemplate 'f'
	// Functions exposed to this package, see the mvvm interface in the mvvm.go file,
	// can be accessed via f.MVVM.funcname()
}

func PublicFunction() {
	// ONLY use when neccecary
	// This function can be accessed by all packages importing this one
}

func privateFunction() {
	// DEFAULT
	// This function is invisible to packages importing this one
	// Decalring an unused private function will result in a warning
}
*/
