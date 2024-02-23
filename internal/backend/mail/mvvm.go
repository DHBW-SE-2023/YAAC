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
	image_data []byte
	course     string
	datetime   time.Time
}

// try to create a new backend_mail struct
// Paramter: mvvm, serverAddress, username, password
// returns an error if it is not possible to connect to the server
func New(mvvm mvvm, serverAddr string, username string, password string) (*BackendMail, error) {
	mailservice := BackendMail{
		MVVM:       mvvm,
		serverAddr: serverAddr,
		username:   username,
		password:   password,
	}
	_, _, err := mailservice.setupMail()
	if err != nil {
		return nil, err
	}
	return &mailservice, nil
}
