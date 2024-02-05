package pages

import (
	"fyne.io/fyne/v2"
)

// Tutorial defines the data structure for a tutorial
type Page struct {
	Title      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
}

var (
	// Tutorials defines the metadata for each tutorial
	Pages = map[string]Page{
		"home":     {"Home", landingScreen, true},
		"overview": {"Übersicht", overviewScreen, true},
		"courses":  {"Kurse", coursesScreen, true},
		"students": {"Studenten", studentScreen, true},
		"settings": {"Einstellungen", settingsScreen, true},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	PagesIndex = map[string][]string{
		"": {"home", "overview", "courses", "students", "settings"},
	}
)
