package yaac_mvvm

import (
	"log"

	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

func settingsToMap(settings []yaac_shared.Setting) map[string]string {
	ms := make(map[string]string)
	for _, s := range settings {
		ms[s.Setting] = s.Value
	}
	return ms
}

func (m *MVVM) StartApplication() {
	err := m.ConnectDatabase("data/data.db")
	if err != nil {
		panic("Could not connect to database")
	}

	settings, err := m.Settings()
	if err != nil {
		log.Fatalln("Could not retrieve the application settings")
		log.Fatalln("Resetting settings ...")
		settings, _ = m.SettingsReset()
	}

	ms := settingsToMap(settings)

	err = m.NewMailBacked(yaac_shared.MailLoginData{MailServer: ms["MailServer"], Email: ms["UserEmail"], Password: ms["UserEmailPassword"]})
	if err != nil {
		log.Println("ERROR: Could not connect to email server")
		log.Println("ERROR: Please set your email credentails in the settings")
	}

	m.StartDemon(5 * 1000) // Refresh every 5 seconds

	// Needs to be the last step
	m.NewFrontendMain()
	m.OpenMainWindow()
}
