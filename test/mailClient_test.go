package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/smtp"
	"os"
	"testing"

	yaac_backend_mail "github.com/DHBW-SE-2023/YAAC/internal/backend/mail"
	yaac_mvvm "github.com/DHBW-SE-2023/YAAC/internal/mvvm"
)

type Credentials struct {
	ServerIMAP string `json:"serverIMAP"`
	ServerSMTP string `json:"serverSMTP"`
	Mail       string `json:"mail"`
	Password   string `json:"password"`
}

type testData struct {
	subject    string
	attachment []byte
}

// returns (IMAP Server:Port) (SMTP Server:Port) (Mailadresse) (Password) (error)
func getTestLoginData() (string, string, string, string, error) {
	mailLoginDataPath := "testdata/mailLoginData.json"
	file, err := os.Open(mailLoginDataPath)
	if err != nil {
		return "", "", "", "", err
	}
	defer file.Close()

	var creds Credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	if err != nil {
		return "", "", "", "", err
	}
	return creds.ServerIMAP, creds.ServerSMTP, creds.Mail, creds.Password, nil
}

func setupMail() (*yaac_backend_mail.BackendMail, error) {
	mvvm := yaac_mvvm.New()
	serverIMAP, _, mail, password, err := getTestLoginData()
	if err != nil {
		return nil, err
	}
	return yaac_backend_mail.New(mvvm, serverIMAP, mail, password)
}

// Send test mails
func sendTestMails() error {
	_, serverSMTP, mail, password, err := getTestLoginData()
	if err != nil {
		return err
	}

}

func TestCheckMailConnection(t *testing.T) {
	mail_backend, err := setupMail()
	if err != nil {
		t.Fatalf("Error setup Mail client: Is login data correct? Error: %v", err)
	}

	if !mail_backend.CheckMailConnection() {
		t.Fatal("Mail Connection ist not active")
	}
}

func buildEmail_Attachment(from string, to string, subject string, body string, imageData []byte) []byte {
	// mail header
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/mixed; boundary=mime-boundary"

	// mail boy
	var emailBody bytes.Buffer

	for key, value := range headers {
		emailBody.WriteString(key + ": " + value + "\r\n")
	}

	emailBody.WriteString("\r\n")
	emailBody.WriteString(body + "\r\n")

	// image attachment
	emailBody.WriteString("\r\n--mime-boundary\r\n")
	emailBody.WriteString("Content-Type: image/jpeg\r\n")
	emailBody.WriteString("Content-Disposition: attachment; filename=\"testfile.jpg\"\r\n")
	emailBody.WriteString("\r\n")
	emailBody.WriteString(base64.StdEncoding.EncodeToString(imageData))
	emailBody.WriteString("\r\n--mime-boundary--")

	return emailBody.Bytes()
}

func sendEmail(from, password, smtpHost, smtpPort, to string, msg []byte) error {
	// authenticate
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}

func TestGetMailsToday(t *testing.T) {

	mail_backend, err := setupMail()
	if err != nil {
		t.Fatalf("Error setup Mail client: Is the login data correct? Error: %v", err)
	}

	/*
		err = sendTestMails()
		if err != nil {
			t.Fatalf("Error sending test Mails: %v", err)
		}
	*/

	_, _, mail, password, err := getTestLoginData()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	//----------

	// E-Mail Konfiguration
	from := mail
	to := mail
	smtpHost := "smtp.mail.de"
	smtpPort := "587"

	// E-Mail Inhalt
	subject := "Anwesenheitsliste TIT22 Test-E-Mail mit Bildanhang"
	body := "Hallo, hier ist eine E-Mail mit einem Bildanhang."

	var testdata []testData

	// E-Mail zusammenbauen
	for i := range testdata {
		msg := buildEmail_Attachment(from, to, subject, body, imageData)
		err = sendEmail(from, password, smtpHost, smtpPort, to, msg)
		if err != nil {
			return err
		}
	}

	// E-Mail senden

	//--------

	println("E-Mail wurde erfolgreich gesendet!")

	mailData, err := mail_backend.GetMailsToday()
	if err != nil {
		t.Fatalf("Error getting test Mails: %v", err)
	}
	log.Printf("maildata %v", string(mailData[0].Image))
}
