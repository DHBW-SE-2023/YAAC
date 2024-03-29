package yaac_mvvm

import (
	"log"
	"time"

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

	err = m.NewMailBacked(yaac_shared.MailLoginData{MailServer: ms["MailServer"], Email: ms["UserEmail"], Password: ms["UserEmailPassword"]})
	if err != nil {
		log.Println("ERROR: Could not connect to email server")
		log.Println("ERROR: Please set your email credentails in the settings")
	}

	m.StartDemon(30 * time.Minute) // Refresh every 30 minutes

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
