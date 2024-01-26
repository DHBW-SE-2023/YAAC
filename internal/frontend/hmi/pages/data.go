package pages

import (
	"fyne.io/fyne/v2"
)

// Tutorial defines the data structure for a tutorial
type Page struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Tutorials defines the metadata for each tutorial
	Pages = map[string]Page{
		"home":     {"Home", "", landingScreen, true},
		"overview": {"Übersicht", "Anwesenheitsliste der Kurse - Heute", overviewScreen, true},
		"courses":  {"Kurse", "Kursansicht", coursesScreen, true},
		"students": {"Studenten", "Studenten", studentScreen, true},
		"settings": {"Einstellungen", "Einstellungen", settingsScreen, true},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	PagesIndex = map[string][]string{
		"": {"home", "overview", "courses", "students", "settings"},
	}
)
