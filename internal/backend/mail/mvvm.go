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
	ReceivedAt time.Time
	ID         uint32
}

// Create a new backend_mail struct
// Paramter: mvvm, serverAddress, username, password
// returns an error if it is not possible to connect and login to the server
func New(mvvm mvvm, serverAddr string, username string, password string) (*BackendMail, error) {
	mailservice := BackendMail{
		MVVM:       mvvm,
		serverAddr: serverAddr,
		username:   username,
		password:   password,
	}

	c, err := mailservice.setupMail(true)
	if err != nil {
		return nil, err
	}

	c.Logout()
	return &mailservice, nil
}
