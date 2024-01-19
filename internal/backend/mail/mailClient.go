package yaac_backend_mail

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
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
	_, err = c.Select("INBOX", false)
	if err != nil {
		return "", err
	}

	// Get first unseen message
	//firstUnseen := mbox.UnseenSeqNum
	//seqset := new(imap.SeqSet)
	//seqset.AddRange(firstUnseen, mbox.Messages)
	//seqset.AddNum(firstUnseen)
	//seqset.AddRange(1, mbox.Messages)

	// Get the whole message body
	//items := []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}

	// Get only unseen Messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		return "", err
	}

	// Create a sequenz for the found mails
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	// Channel for mail messages
	messages := make(chan *imap.Message)

	// Fetching mails
	go func() {
		if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages); err != nil {
			log.Fatal(err)
		}
	}()

	// Check all unread messages
	for msg := range messages {
		if msg == nil {
			log.Fatal("Error: Could not read Mail")
			continue
		}

		mailstring, err := getMailasString(msg)
		if err == nil {
			log.Fatal("Error decoding mail")
		} else {
			subject, err := getSubject(mailstring)
			if err == nil {
				log.Fatal("Error decoding mail")
			} else {
				log.Printf("Mail detected with Subject: %v\n", subject)

				// Check if mail Subject is correct

				image, err := getBase64AttachmentFromMail(mailstring)
				if err == nil {
					log.Println("Image Attachment identified")
					err := writeToDatabase(image)
					if err != nil {
						log.Println(("Error writing to Database"))
					} else {
						// Mark Mail as read
					}
				}
			}
		}
	}
	return "Sucessfull!", err

}

func getBoundary(contentType string) (string, error) {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}
	return params["boundary"], err
}

func getMailasString(msg *imap.Message) (string, error) {
	// Get Mail Literal
	var section imap.BodySectionName
	mailLiteral := msg.GetBody(&section)
	if mailLiteral == nil {
		err := errors.New("No litteral in mail body found")
		return "", err
	}

	// Get Mail as String
	mailString, err := imap.ParseString(mailLiteral)
	if err != nil {
		return "", err
	}

	return mailString, err
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
			//fmt.Println("----------")
			//fmt.Println(string(body))
			//fmt.Println("----------")
			if strings.HasPrefix(part.Header.Get("Content-Type"), "image/jpeg") {
				return (string(body)), err
			}
		}
	}
	// Return Error
	err = errors.New("Found no attached image in mail")
	return "", err
}

func getSubject(mailString string) (string, error) {
	// Read Mail
	message, err := mail.ReadMessage(strings.NewReader(mailString))
	if err != nil {
		return "", err
	}

	// Read content type from Mail Header
	contentType := message.Header.Get("Subject")

	return contentType, err
}

// demo func, need to ne replaced
func writeToDatabase(base64data string) error {
	f, err := os.Create("/tmp/image.jpg")
	if err != nil {
		log.Println("Error creating file")
		return err
	}

	defer f.Close()

	data, err := decodeBase64(base64data)
	if err != nil {
		return err
	}

	bytes, err := f.Write(data)
	if err != nil {
		return err
	}
	log.Printf(("Wrote data to file: %v Bytes"), bytes)
	return nil
}

func decodeBase64(base64data string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
