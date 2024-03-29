package yaac_mvvm

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

	// err = m.NewMailBacked(yaac_shared.MailLoginData{MailServer: ms["MailServer"], Email: ms["UserEmail"], Password: ms["UserEmailPassword"]})
	// if err != nil {
	// 	log.Fatalf("Could not connect to email server")
	// }

	// m.StartDemon(5) // Refresh every 5 seconds

	// Needs to be the last step
	m.NewFrontendMain()
	m.OpenMainWindow()
}
