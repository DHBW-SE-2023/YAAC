package yaac_mvvm

import (
	"log"

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
	m.ImgprocBackendStart()
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
	_, _, _ = ReturnMailSettings(settings)
	ms["MailServer"] = "secureimap.t-online.de:993"
	ms["UserEmail"] = "dhbw.rust.sweng@t-online.de"
	ms["UserEmailPassword"] = "Hallo123"
	// ms["MailServer"] = mailConnection
	// ms["UserEmail"] = mailUser
	// ms["UserEmailPassword"] = mailPassword

	err = m.NewMailBacked(yaac_shared.EmailData{MailServer: ms["MailServer"], Email: ms["UserEmail"], Password: ms["UserEmailPassword"]})
	if err != nil {
		// log.Fatalf("Could not connect to email server")
		log.Default()
	}

	m.StartDemon(10) // Refresh every 5 seconds

	// Needs to be the last step
	m.NewFrontendMain()
	m.OpenMainWindow()
}

func ReturnMailSettings(setting []yaac_shared.Setting) (string, string, string) {
	var mailConnection string
	var mailUser string
	var mailPassword string
	for _, element := range setting {
		if element.Setting == "mailConnection" {
			mailConnection = element.Value
		} else if element.Setting == "mailUser" {
			mailUser = element.Value
		} else {
			mailPassword = element.Value
		}
	}
	return mailConnection, mailUser, mailPassword
}
