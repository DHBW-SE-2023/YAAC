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
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// GetMailsToday fetches all unread mails from today
// and checks the mails with the subject containing "Anwesenheitsliste".
// It extracts the attached image as binary data.
// Returns an array with the maildata from the mails
func (b *BackendMail) GetMailsToday() ([]MailData, error) {

	// array for the maildata
	var maildata []MailData

	//setup mail client
	c, ids, err := b.setupMail(true)
	if err != nil {
		return nil, err
	}

	// logout before function returns
	defer c.Logout()

	for _, id := range ids {

		messages := make(chan *imap.Message, 1)
		header := make(chan *imap.Message, 1)

		log.Printf("ID: %v\n", id)

		seqset := new(imap.SeqSet)
		seqset.AddNum(id)

		// Nachrichten-Header abrufen
		if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, header); err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		head := <-header

		log.Printf("Subject: %v\n", head.Envelope.Subject)
		log.Printf("DAte: %v\n", head.Envelope.Date)

		if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages); err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		msg := <-messages

		binary_image, err := b.getBinaryImageFromMail(msg)
		if err == nil {
			maildata = append(maildata, MailData{Image: binary_image, ReceivedAt: head.Envelope.Date, ID: id})
			log.Println("Successfully added")
		}
	}
	return maildata, nil
}

// Checks the mail connection to the server and the login credentials. Returns true if the connection and authentication is fine otherwise false
func (b *BackendMail) CheckMailConnection() bool {
	c, _, err := b.setupMail(true)
	if err != nil {
		return false
	}
	c.Logout()
	return true
}

// getBoundary needs the contentType part of the mail and returns the value of boundary
func (b *BackendMail) getBoundary(contentType string) (string, error) {
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}
	if params["boundary"] == "" {
		err := errors.New("no boundary found")
		return "", err
	}
	return params["boundary"], nil
}

// getMailAsString gets a pointer to a imap.message and returns the mail as string
// returns an error if it is not possible to convert the mail to a string
func (b *BackendMail) getMailAsString(msg *imap.Message) (string, error) {
	// Get Mail Literal
	var section imap.BodySectionName
	mailLiteral := msg.GetBody(&section)
	if mailLiteral == nil {
		err := errors.New("no litteral in mail body found")
		log.Printf("Error: %v", err)
		return "", err
	}

	// Get Mail as String
	mailString, err := imap.ParseString(mailLiteral)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	return mailString, err
}

// getBinaryImageFromMailString needs the mail as a string
// returns the binary image file that is attached in the mail
// returns an error if there is a problem extracting the image or if there is no image
func (b *BackendMail) getBinaryImageFromMail(msg *imap.Message) ([]byte, error) {

	mailString, err := b.getMailAsString(msg)
	if err != nil {
		return nil, err
	}

	// Read Mail
	message, err := mail.ReadMessage(strings.NewReader(mailString))
	if err != nil {
		log.Printf("Error: %v", err)
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
				log.Printf("Error: %v", err)
				return nil, err
			}

			// Read one part
			body, err := io.ReadAll(part)
			if err != nil {
				log.Printf("Error: %v", err)
				return nil, err
			}

			//Check if the part contains a jpeg image
			if strings.HasPrefix(part.Header.Get("Content-Type"), "image/jpeg") {
				binaryData, err := base64.StdEncoding.DecodeString(string(body))
				log.Printf("Error: %v", err)
				return binaryData, err
			}
		}
	}
	// Return Error if no image found
	err = errors.New("found no attached image in mail")
	log.Printf("Error: %v", err)
	return nil, err
}

// setupMail sets up the completly mail setup and returns the client and the ids of the unread mails
// returns an error if something went wrong
func (b *BackendMail) setupMail(onlyReadable bool) (*client.Client, []uint32, error) {
	//connect to the mail server
	c, err := b.connectToServer(b.serverAddr)
	if err != nil {
		return nil, nil, err
	}

	//login with the user credentials
	c, err = b.logInToInbox(c, b.username, b.password, onlyReadable)
	if err != nil {
		return nil, nil, err
	}

	//focus only on unread mails
	c, ids, err := b.getIDsOfMails(c)
	if err != nil {
		return nil, nil, err
	}
	return c, ids, nil
}

func (b *BackendMail) getIDsOfMails(c *client.Client) (*client.Client, []uint32, error) {
	// Search for unread messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, nil, err
	}

	// sort them in reverse order to get the latest message first
	var sortedUIDs []uint32
	for _, uid := range ids {
		sortedUIDs = append([]uint32{uint32(uid)}, sortedUIDs...)
	}

	return c, sortedUIDs, nil
}

/*
func (b *BackendMail) processMail(mail *imap.Message) (MailData, error) {

	var mailData = MailData{}

	var err error

	mailstring, err := b.getMailAsString(mail)
	if err != nil {
		return mailData, err
	}

	//get the base64 encoded image from the mail
	mailData.Image, err = b.getBinaryImageFromMailString(mailstring)
	if err != nil {
		return mailData, err
	}

	mailData.ReceivedAt, err = b.getDatetime(mailstring)
	if err != nil {
		return mailData, err
	}

	return mailData, err
}
*/

// checkDatetime checks if the mail is from today. So it checks if the date from the mail is from today.
// returns true if the mail is from today otherwise false
func (b *BackendMail) checkDatetime(mailTime time.Time) bool {
	// Get Year, month and day from mail
	mail_year, mail_month, mail_day := mailTime.Local().Date()
	// Get Year, month and day from now
	current_year, current_month, current_day := time.Now().Date()

	return (mail_year == current_year && mail_month == current_month && mail_day == current_day)
}

// checkMailSubject needs the mail subject
// returns if the subject of the mail contains "Anwesenheitsliste" or not
func (b *BackendMail) checkMailSubject(subject string) bool {
	return strings.Contains(subject, "Anwesenheitsliste")
}

// connectToServer needs the server address and conects to the server
// returns the client and an error
func (b *BackendMail) connectToServer(serverAddr string) (*client.Client, error) {
	// Connect to server
	c, err := client.DialTLS(serverAddr, nil)
	if err != nil {
		log.Printf("Info: %v", err)
		log.Print("Trying connect to mail server without TLS ...")
		c, err = client.Dial(serverAddr)
		if err != nil {
			log.Printf("Error: %v", err)
			return nil, err
		}
	}
	return c, err
}

// LogInToInbox needs the client, username and passwort and logs in into INBOX
// returns the new client and an error
func (b *BackendMail) logInToInbox(c *client.Client, username string, password string, onlyReadable bool) (*client.Client, error) {
	// login to the email server
	if err := c.Login(username, password); err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	// Select to default INBOX
	_, err := c.Select("INBOX", onlyReadable)
	if err != nil {
		log.Printf("Error: %v", err)
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
		log.Printf("Error: %v", err)
		return nil, nil, err
	}
	// Create a sequenz for the found mails
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)
	return c, seqset, nil
}

/*
// fetchMails needs the client, seqset and a imap message channel
// It fetches the mails and send them to the channel
func (b *BackendMail) fetchMails(c *client.Client, seqset *imap.SeqSet, messages chan *imap.Message) {
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages); err != nil {
		log.Printf("Info: No unread mails: %v", err)
	}
}
*/

// Marks the mail with the given ID as read
// returns an error if there is an error with that
func (b *BackendMail) MarkMailsAsRead(ids []uint32) error {
	//connect to server
	c, err := b.connectToServer(b.serverAddr)
	if err != nil {
		return err
	}

	//login with the user credentials
	c, err = b.logInToInbox(c, b.username, b.password, false)
	if err != nil {
		return err
	}

	//mark all mails with the given ids as read
	for _, id := range ids {

		seqset := new(imap.SeqSet)
		seqset.AddNum(id)

		err = c.Store(seqset, "+FLAGS.SILENT", []interface{}{imap.SeenFlag}, nil)
		if err != nil {
			log.Printf("Error: %v", err)
			return err
		}
	}
	return nil
}
