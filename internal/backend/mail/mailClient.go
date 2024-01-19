package yaac_backend_mail

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func (b *BackendMail) GetResponse(input yaac_shared.EmailData) string {
	msg, err := getLatestMessage(input.MailServer, input.Email, input.Password)
	if err != nil {
		return ("Something went wrong please try again:")
	}
	return msg
}

func getLatestMessage(serverAddr string, username string, password string) (string, error) {
	// Connect to server
	c, err := client.DialTLS(serverAddr, nil)
	if err != nil {
		return "", err
	}
	log.Println("Connected to server at: ", serverAddr)

	// close connection to server when data was recieved
	defer c.Logout()

	// login to the email server
	if err := c.Login(username, password); err != nil {
		return "", err
	}
	log.Println("Logged in to server at: ", serverAddr)

	// Select to default INBOX, would not work if renamed
	// TODO: in actual implementation return an error for this and show a filed for inbox selection
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return "", err
	}

	// Get first unseen message
	firstUnseen := mbox.UnseenSeqNum
	seqset := new(imap.SeqSet)
	seqset.AddNum(firstUnseen)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}

	// Channes for messages
	messages := make(chan *imap.Message, 1)

	go func() {
		if err := c.Fetch(seqset, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	msg := <-messages
	if msg == nil {
		return "", err
	}

	// Get Mail Literal
	mailLiteral := msg.GetBody(&section)
	if mailLiteral == nil {
		return "", err
	}

	// Get Mail Body as String
	mailString, err := imap.ParseString(mailLiteral)
	if err != nil {
		return "", err
	}

	// Get Base64Image from Body
	base64Image, err := getBase64AttachmentFromMail(mailString)
	if err != nil {
		return "", err
	}

	return base64Image, err

}

func getBoundary(contentType string) (string, error) {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}
	return params["boundary"], err
}

func getBase64AttachmentFromMail(mailString string) (string, error) {

	// Read Mail
	message, err := mail.ReadMessage(strings.NewReader(mailString))
	if err != nil {
		return "", err
	}

	// Read content type from Mail Header
	contentType := message.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {

		// Read Body Boundary from Mail Header
		boundary, err := getBoundary(contentType)
		if err != nil {
			return "", err
		}

		// Divide Mail into body parts
		mr := multipart.NewReader(message.Body, boundary)

		for {

			// End of File
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}

			// Other Error
			if err != nil {
				return "", err
			}

			// Read one part
			body, err := io.ReadAll(part)
			if err != nil {
				return "", err
			}

			fmt.Printf("Teil mit Content-Type %s:\n", part.Header.Get("Content-Type"))
			fmt.Println("----------")
			fmt.Println(string(body))
			fmt.Println("----------")
			if strings.HasPrefix(part.Header.Get("Content-Type"), "image/jpeg") {
				return (string(body)), err
			}
		}
	}
	// Return Error
	err = errors.New("Found no attached image in mail")
	return "", err
}
