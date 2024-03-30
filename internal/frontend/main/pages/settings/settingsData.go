package yaac_frontend_settings

import (
	"fyne.io/fyne/v2"
)

/*
Defines the Setting struct which will hold each settingPages properties
*/
type Setting struct {
	Title      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
}

/*
Defines the actual Mapping from each Setting struct to it's respective index
*/
var (
	SettingPages = map[string]Setting{
		"email":     {"Email", emailScreen, true},
		"wiki":      {"WIKI", wikiScreen, true},
		"impressum": {"Impressum", impressumScreen, true},
	}

	SettingPagesIndex = map[string][]string{
		"": {"email", "wiki", "impressum"},
	}
)
