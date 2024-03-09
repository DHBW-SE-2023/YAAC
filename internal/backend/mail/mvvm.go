package yaac_backend_mail

import "time"

type mvvm interface {
}

type BackendMail struct {
	MVVM       mvvm
	serverAddr string
	username   string
	password   string
}

type MailData struct {
	Image      []byte
	Course     string
	ReceivedAt time.Time
}

// Create a new backend_mail struct
// Paramter: mvvm, serverAddress, username, password
// returns an error if it is not possible to connect and login to the server and NO Mailservice
func New(mvvm mvvm, serverAddr string, username string, password string) (*BackendMail, error) {
	mailservice := BackendMail{
		MVVM:       mvvm,
		serverAddr: serverAddr,
		username:   username,
		password:   password,
	}

	c, _, err := mailservice.setupMail()
	if err != nil {
		return nil, err
	}

	c.Logout()
	return &mailservice, nil
}
