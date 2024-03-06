package yaac_mvvm

import (
	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
)

// FIXME: The existance of this function should not be necessary
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

	// settings, err := m.Settings()
	// if err != nil {
	// 	log.Fatalln("Could not retrieve the application settings")
	// 	log.Fatalln("Resetting settings ...")
	// 	settings, _ = m.SettingsReset()
	// }

	// ms := settingsToMap(settings)

	// ms["MailServer"] = "email.server:80"
	// ms["UserEmail"] = "myemail@email.server"
	// ms["UserEmailPassword"] = "123"

	// err = m.NewMailBacked(yaac_shared.EmailData{MailServer: ms["MailServer"], Email: ms["UserEmail"], Password: ms["UserEmailPassword"]})
	// if err != nil {
	// 	log.Fatalf("Could not connect to email server")
	// }

	// m.StartDemon(5) // Refresh every 5 seconds

	// Needs to be the last step
	m.NewFrontendMain()
	m.OpenMainWindow()
}
