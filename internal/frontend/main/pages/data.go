package yaac_frontend_pages

import (
	"fyne.io/fyne/v2"
)

/*
Defines the Page struct which will hold each Pages properties
*/
type Page struct {
	Title      string
	View       func(w fyne.Window) fyne.CanvasObject
	SupportWeb bool
}

/*
Defines the actual Mapping from each Page struct to it's respective index
*/
var (
	Pages = map[string]Page{
		"home":     {"Home", LandingScreen, true},
		"overview": {"Übersicht", OverviewScreen, true},
		"courses":  {"Kurse", CoursesScreen, true},
		"students": {"Studenten", StudentScreen, true},
		"settings": {"Einstellungen", SettingsScreen, true},
	}

	PagesIndex = map[string][]string{
		"": {"home", "overview", "courses", "students", "settings"},
	}
)
