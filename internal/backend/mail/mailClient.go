package yaac_backend_mail

import (
	"encoding/base64"
	"errors"
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
	//go MailService(input.MailServer, input.Email, input.Password)
	return "Please read log!"
}

// GetMailsToday fetches all unread mails from today
// and checks the mails with the subject containing "Anwesenheitsliste".
// It extracts the attached image as binary data.
// Returns an array with the maildata from the mails
func (b *BackendMail) GetMailsToday() ([]MailData, error) {

	//channel for mails
	messages := make(chan *imap.Message)

	// array for the maildata
	var maildata []MailData

	//setup mail client
	c, seqset, err := b.setupMail()
	if err != nil {
		return nil, err
	}

	//start fetchings mails
	go b.fetchMails(c, seqset, messages)

	//process fetched mails
	for msg := range messages {

		//no new unread mail
		if msg == nil {
			log.Println("all mails fetched")
			break
		}
		if b.checkSubjekt() && b.mailToday() { //need to be added
			maildata_temp, err := b.processMail(msg)
			if err == nil {
				maildata = append(maildata, maildata_temp)
			}
		}
	}
	//close connection
	c.Logout()

	return maildata, nil
}

// kann gelöscht werden
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
			log.Println(err)
		}
	}()

	// Check all unread messages
	for msg := range messages {
		if msg == nil {
			log.Println("Error: Could not read Mail")
			continue
		}

		mailstring, err := getMailasString(msg)
		if err == nil {
			log.Println("Error decoding mail")
		} else {
			subject, err := getSubject(mailstring)
			if err == nil {
				log.Println("Error decoding mail")
			} else {
				log.Printf("Mail detected with Subject: %v\n", subject)

				// Check if mail Subject is correct

				image, err := getBase64AttachmentFromMailString(mailstring)
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

// getBoundary needs the contentType part of the mail and returns the value of boundary
func (b *BackendMail) getBoundary(contentType string) (string, error) {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}
	return params["boundary"], nil
}

// getMailasString gets a pointer to a imap.message and returns the mail as string
// returns an error if it is not possible to convert the mail to a string
func (b *BackendMail) getMailasString(msg *imap.Message) (string, error) {
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

// getBinaryImageFromMailString needs the mail as a string
// returns the binary image file that is attached in the mail
// returns an error if there is a problem extracting the image or if there is no image
func (b *BackendMail) getBinaryImageFromMailString(mailString string) ([]byte, error) {

	// Read Mail
	message, err := mail.ReadMessage(strings.NewReader(mailString))
	if err != nil {
		return nil, err
	}

	// Read content type from Mail Header
	contentType := message.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {

		// Read Body Boundary from Mail Header
		boundary, err := b.getBoundary(contentType)
		if err != nil {
			return nil, err
		}

		// Divide Mail into body parts
		mr := multipart.NewReader(message.Body, boundary)

		for {

			// Read next part
			part, err := mr.NextPart()

			// End of file
			if err == io.EOF {
				break
			}

			// Other Error
			if err != nil {
				return nil, err
			}

			// Read one part
			body, err := io.ReadAll(part)
			if err != nil {
				return nil, err
			}

			//Check if the part contains a jpeg image
			if strings.HasPrefix(part.Header.Get("Content-Type"), "image/jpeg") {
				binaryData, err := base64.StdEncoding.DecodeString(string(body))
				return binaryData, err
			}
		}
	}
	// Return Error
	err = errors.New("Found no attached image in mail")
	return nil, err
}

// getSubject needs the mail as string and return the subject of the mail
func (b *BackendMail) getSubject(mailString string) (string, error) {
	// Read Mail
	message, err := mail.ReadMessage(strings.NewReader(mailString))
	if err != nil {
		return "", err
	}

	// Read content type from Mail Header
	subject := message.Header.Get("Subject")

	return subject, err
}

// setupMail sets up the completly mail setup and returns the client and the seqset
// returns an error if
func (b *BackendMail) setupMail() (*client.Client, *imap.SeqSet, error) {
	//connect to the mail server
	c, err := b.connectToServer(b.serverAddr)
	if err != nil {
		return nil, nil, err
	}

	//login with the user crudentials
	c, err = b.logInToInbox(c, b.username, b.password)
	if err != nil {
		return nil, nil, err
	}

	//focus only on unread mails
	c, seqset, err := b.getOnlyUnreadMails(c)
	if err != nil {
		return nil, nil, err
	}
	return c, seqset, nil
}

// processMail processes the whole imap message object
// returns the MailData Struct with the data
// returns an error if there went something wrong with decoding and parsing the mail
func (b *BackendMail) processMail(msg *imap.Message) (MailData, error) {

	var mailData = MailData{}

	//get mail message as string in order to prcess the mail further
	mailstring, err := b.getMailasString(msg)
	if err != nil {
		log.Println("Error decoding mail")
		return mailData, err
	}

	//get the base64 encoded image from the mail
	mailData.image_data, err = b.getBinaryImageFromMailString(mailstring)
	if err != nil {
		log.Println("Error Image Attachment Extraction")
		return mailData, err
	}

	mailData.course, err = b.getCourse(mailstring)
	if err != nil {
		log.Println("Error detecting course")
		return mailData, err
	}

	mailData.datetime, err = b.getDatetime(mailstring)
	if err != nil {
		log.Println("Error detecting course")
		return mailData, err
	}

	return mailData, err
}

// checkMailSubject needs the mail string
// returns true if the subject of the mail contains ""
// returns true if the subject of the mail contains not "" or there is an error reading the mail subject
func (b *BackendMail) checkMailSubject(mailstring string) bool {
	subject, err := b.getSubject(mailstring)
	if err != nil {
		log.Println("Error decoding mail subject")
		return false
	}
	log.Printf("Mail detected with Subject: %v\n", subject)
	return strings.Contains(subject, "Anwesenheitsliste")
}

// markMailAsRead needs the mail client and the imap message and marks the message as read
// returns an error if it didn't work otherwise nil
func (b *BackendMail) markMailAsRead(c *client.Client, msg *imap.Message) error {
	seqnums := new(imap.SeqSet)
	seqnums.AddNum(msg.SeqNum)
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	if err := c.Store(seqnums, item, flags, nil); err != nil {
		return err
	}
	return nil
}

// connectToServer needs the server address and conects to the server
// returns the client and an error
func (b *BackendMail) connectToServer(serverAddr string) (*client.Client, error) {
	// Connect to server
	c, err := client.DialTLS(serverAddr, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to server at: ", serverAddr)
	return c, err
}

// LogInToInbox needs the client, username and passwort and logs in into INBOX
// returns the new client and an error
func (b *BackendMail) logInToInbox(c *client.Client, username string, password string) (*client.Client, error) {
	// login to the email server
	if err := c.Login(username, password); err != nil {
		return nil, err
	}
	log.Println("Logged in to server")

	// Select to default INBOX
	_, err := c.Select("INBOX", false)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// getOnlyUnreadMails needs the client and only selects unread mails
// returns the new client, the imap Seqset, that is need for fetching these unread mails and an error
func (b *BackendMail) getOnlyUnreadMails(c *client.Client) (*client.Client, *imap.SeqSet, error) {
	// Get only unseen Messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		return nil, nil, err
	}
	// Create a sequenz for the found mails
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)
	return c, seqset, nil
}

// fetchMails needs the client, seqset and a imap message channel
// It fetches the mails and send them to the channel
func (b *BackendMail) fetchMails(c *client.Client, seqset *imap.SeqSet, messages chan *imap.Message) {
	// Fetching mails
	log.Println(("start fetching"))
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages); err != nil {
		log.Println("no unread mails:")
	}
	log.Println(("end fetching"))
}
