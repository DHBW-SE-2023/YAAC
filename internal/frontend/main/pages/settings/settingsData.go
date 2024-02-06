package settings

import (
	"fyne.io/fyne/v2"
)

// Tutorial defines the data structure for a tutorial
type Setting struct {
	Title      string
	View       func() fyne.CanvasObject
	SupportWeb bool
}

var (
	// Tutorials defines the metadata for each tutorial
	SettingPages = map[string]Setting{
		"general":   {"Allgemein", generalScreen, true},
		"database":  {"Datenbank", databaseScreen, true},
		"email":     {"Email", emailScreen, true},
		"wiki":      {"WIKI", wikiScreen, true},
		"impressum": {"Impressum", impressumScreen, true},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	SettingPagesIndex = map[string][]string{
		"": {"general", "database", "email", "wiki", "impressum"},
	}
)
