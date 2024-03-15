package test

import (
	"encoding/json"
	"os"
	"testing"

	yaac_backend_mail "github.com/DHBW-SE-2023/YAAC/internal/backend/mail"
	yaac_mvvm "github.com/DHBW-SE-2023/YAAC/internal/mvvm"
)

type Credentials struct {
	Server   string `json:"server"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

func getTestLoginData() (string, string, string, error) {
	file, err := os.Open("testdata/mailLoginData.json")
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	var creds Credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	if err != nil {
		return "", "", "", err
	}
	return creds.Server, creds.Mail, creds.Password, nil
}

func setupMail() (*yaac_backend_mail.BackendMail, error) {
	mvvm := yaac_mvvm.New()
	server, mail, password, err := getTestLoginData()
	if err != nil {
		return nil, err
	}
	return yaac_backend_mail.New(mvvm, server, mail, password)
}

func TestCheckMailConnection(t *testing.T) {
	mail_backend, err := setupMail()
	if err != nil {
		t.Fatalf("Error setup Mail client: Is the server online and the login data correct? Error: %v", err)
	}

	if !mail_backend.CheckMailConnection() {
		t.Fatal("Mail Connection ist not active")
	}
}
