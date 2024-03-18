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

	//channel for mails
	messages := make(chan *imap.Message)

	// array for the maildata
	var maildata []MailData

	//setup mail client
	c, seqset, err := b.setupMail()
	defer c.Logout()

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

		//get mail message as string in order to prcess the mail further
		mailstring, err := b.getMailAsString(msg)
		if err != nil {
			log.Println("Error decoding mail")
			continue
		}

		if b.checkMailSubject(mailstring) && b.checkDatetime(mailstring) {
			maildata_temp, err := b.processMail(mailstring)
			if err == nil {
				log.Println("Successfully added")
				maildata = append(maildata, maildata_temp)
			}
		}
	}

	return maildata, nil
}

// Checks the mail connection to the server and the login credentials. Returns true if the connection and authentication is fine otherwise false
func (b *BackendMail) CheckMailConnection() bool {
	c, _, err := b.setupMail()
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
	// Return Error if no image found
	err = errors.New("found no attached image in mail")
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

	//login with the user credentials
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
func (b *BackendMail) processMail(mailstring string) (MailData, error) {

	var mailData = MailData{}

	var err error

	//get the base64 encoded image from the mail
	mailData.Image, err = b.getBinaryImageFromMailString(mailstring)
	if err != nil {
		log.Println("Error Image Attachment Extraction")
		return mailData, err
	}

	mailData.Course, err = b.getCourse(mailstring)
	if err != nil {
		log.Println("Error detecting course")
		return mailData, err
	}

	mailData.ReceivedAt, err = b.getDatetime(mailstring)
	if err != nil {
		log.Println("Error detecting course")
		return mailData, err
	}

	return mailData, err
}

// getCourse extractes Course from mail subject and returns the course as string
// returns an error if it is not possilble to read the mail or if there is no course found in the subject
func (b *BackendMail) getCourse(mailstring string) (string, error) {
	subject, err := b.getSubject(mailstring)
	if err != nil {
		return "", err
	}
	words := strings.Fields(subject)
	for _, word := range words {
		if strings.Contains(word, "TI") {
			return word, nil
		}
	}
	err = errors.New("no course found in subject")
	return "", err
}

// getDatetime extraces the date and time of the mail and returns it as time struct
// returns an error if it is not possilble to read the mail or if there occurs an error parsing the date or time
func (b *BackendMail) getDatetime(mailstring string) (time.Time, error) {
	// Read Mail
	message, err := mail.ReadMessage(strings.NewReader(mailstring))
	if err != nil {
		log.Println("Not able to read mail")
		return time.Now(), err
	}

	// Read datetime from Mail Header
	// These two formats are mostly used in mails
	datestring := message.Header.Get("Date")
	datetime, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", datestring)
	if err != nil {
		datetime, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700 (MST)", datestring)
		if err != nil {
			log.Printf("Error: %v", err)
			return time.Now(), err
		}
	}
	return datetime, nil
}

// checkDatetime checks if the mail is from today. So it checks if the date from the mail is from today.
// returns true if the mail is from today otherwise false
func (b *BackendMail) checkDatetime(mailsting string) bool {
	maildate, err := b.getDatetime(mailsting)
	if err != nil {
		log.Printf("Error parsing date from mail: %v", err)
		return false
	}
	mail_year, mail_month, mail_day := maildate.Local().Date()
	current_year, current_month, current_day := time.Now().Date()
	if mail_year == current_year && mail_month == current_month && mail_day == current_day {
		log.Printf("Fiting Date: %v", maildate.Local())
		return true
	}
	log.Printf("Not Fiting Date: %v", maildate.Local())
	return false
}

// checkMailSubject needs the mail string
// returns true if the subject of the mail contains "Anwesenheitsliste"
// returns true if the subject of the mail contains not "Anwesenheitsliste" or there is an error reading the mail subject
func (b *BackendMail) checkMailSubject(mailstring string) bool {
	subject, err := b.getSubject(mailstring)
	if err != nil {
		log.Println("Error decoding mail subject")
		return false
	}
	if strings.Contains(subject, "Anwesenheitsliste") {
		log.Printf(("Fiting Subject: %v"), subject)
		return true
	}
	return false
}

// connectToServer needs the server address and conects to the server
// returns the client and an error
func (b *BackendMail) connectToServer(serverAddr string) (*client.Client, error) {
	// Connect to server
	c, err := client.DialTLS(serverAddr, nil)
	if err != nil {
		log.Printf("Error: %v. Trying without TLS ...", err)
		c, err = client.Dial(serverAddr)
		if err != nil {
			log.Printf("Error: %v", err)
			return nil, err
		}
	}
	log.Println("Connected to server at: ", serverAddr)
	return c, err
}

// LogInToInbox needs the client, username and passwort and logs in into INBOX
// returns the new client and an error
func (b *BackendMail) logInToInbox(c *client.Client, username string, password string) (*client.Client, error) {
	// login to the email server
	if err := c.Login(username, password); err != nil {
		log.Printf("Error: %v", err)
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
