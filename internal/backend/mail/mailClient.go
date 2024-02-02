package yaac_backend_mail

import (
	"encoding/base64"
	"errors"
	"io"
	"log"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func (b *BackendMail) GetResponse(input yaac_shared.EmailData) string {
	go MailService(input.MailServer, input.Email, input.Password)
	return "Please read log!"
}

func (b *BackendMail) GetMailsToday() ([]MailData, error) {

	//channel for mails
	messages := make(chan *imap.Message)

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
			log.Println("all mails processed")
			break
		}
		b.processMail(msg)

	}
	//close connection and wait one hour to fetch mails again
	b.closeConnection(c)

	return nil, nil
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

// getBase64AttachmentFromMailString needs the mail as a string
// returns the base64 encode image file that is attached in the mail
// returns an error if there is a problem extracting the image or if there is no image
func (b *BackendMail) getBase64AttachmentFromMailString(mailString string) (string, error) {

	// Read Mail
	message, err := mail.ReadMessage(strings.NewReader(mailString))
	if err != nil {
		return "", err
	}

	// Read content type from Mail Header
	contentType := message.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {

		// Read Body Boundary from Mail Header
		boundary, err := b.getBoundary(contentType)
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

			if strings.HasPrefix(part.Header.Get("Content-Type"), "image/jpeg") {
				return (string(body)), err
			}
		}
	}
	// Return Error
	err = errors.New("Found no attached image in mail")
	return "", err
}

// getSubject needs the mail as string and return the subject of the mail
func (b *BackendMail) getSubject(mailString string) (string, error) {
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
func (b *BackendMail) writeToDatabase(base64data string) error {
	path := "/Users/vinzent/image_" + strconv.Itoa(rand.Intn(1000)) + ".jpg"
	f, err := os.Create(path)
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

// decodeBase64 needs the base64 encoded string and return the binary data of it
// returns an error if it is not possible to decode
func (b *BackendMail) decodeBase64(base64data string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

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

func (b *BackendMail) processImage() bool {
	//get the base64 encoded image from the mail, continue if there is no
	image, err := b.getBase64AttachmentFromMailString(mailstring)
	if err != nil {
		return false
	}
	log.Println("Image Attachment identified")

	//write that to database
	err = b.writeToDatabase(image)
	if err != nil {
		log.Println(("Error writing to Database"))
		return false
	}

	//mark mail as read in order to show that it has been proceed
	err = b.markMailAsRead(c, msg)
	if err != nil {
		log.Println(("Error marking mail as read!"))
		log.Print(err)
		return false
	}
	log.Println(("Data successfully stored in data base"))

}

func (b *BackendMail) processMail(msg *imap.Message) bool {
	//get mail message as string in order to prcess the mail further
	mailstring, err := b.getMailasString(msg)
	if err != nil {
		log.Println("Error decoding mail")
		return false
	}

	//check the subject of the mail and process it only if the subject contains "Anwesendtheitsliste"
	if b.checkMailSubject(mailstring) {

		//get the base64 encoded image from the mail, continue if there is no
		image, err := b.getBase64AttachmentFromMailString(mailstring)
		if err != nil {
			return
		}
		log.Println("Image Attachment identified")

		//write that to database
		err = b.writeToDatabase(image)
		if err != nil {
			log.Println(("Error writing to Database"))
			return
		}

		//mark mail as read in order to show that it has been proceed
		err = b.markMailAsRead(c, msg)
		if err != nil {
			log.Println(("Error marking mail as read!"))
			log.Print(err)
			return
		}
		log.Println(("Data successfully put forward"))
	}
	return false
}

// MailService needs the user crudentials and checks the mails every hour
// contains an endless loop, if there is an error with the login data, it returns an error
func MailService(serverAddr string, username string, password string) error {

	for {

		//channel for mails
		messages := make(chan *imap.Message)

		//setup mail client
		c, seqset, err := setupMail(serverAddr, username, password)
		if err != nil {
			return nil
		}

		//start fetchings mails
		go b.fetchMails(c, seqset, messages)

		//process fetched mails
		for msg := range messages {

			//no new unread mail
			if msg == nil {
				log.Println("all mails processed")
				break
			}
			pro

		}
		//close connection and wait one hour to fetch mails again
		CloseConnection(c)
		time.Sleep(1 * time.Hour)
	}
}

// can be deleted
func fetchMails_old(serverAddr string, username string, password string, messages chan *imap.Message) error {
	// init all
	// Connect to server
	c, err := client.DialTLS(serverAddr, nil)
	if err != nil {
		return err
	}
	log.Println("Connected to server at: ", serverAddr)

	// close connection to server when data was recieved
	defer c.Logout()

	// login to the email server
	if err := c.Login(username, password); err != nil {
		return err
	}
	log.Println("Logged in to server at: ", serverAddr)

	// Select to default INBOX, would not work if renamed
	// TODO: in actual implementation return an error for this and show a filed for inbox selection
	_, err = c.Select("INBOX", false)
	if err != nil {
		return err
	}

	// Get only unseen Messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		return err
	}

	// Create a sequenz for the found mails
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	// Fetching mails
	log.Println(("start fetching"))
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages); err != nil {
		log.Println("no unread mails:")
	}
	log.Println(("end fetching"))
	return nil
}

// can be deleted
func checkMails(messages chan *imap.Message) {
	// Check all unread messages
	for msg := range messages {
		if msg == nil {
			log.Println("Error reading mail")
			return
		}

		mailstring, err := getMailasString(msg)
		if err != nil {
			log.Println("Error decoding mail")
		} else {
			subject, err := getSubject(mailstring)
			if err != nil {
				log.Println("Error decoding mail subject")
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
	log.Println(("End checking mails"))
}

// checkMailSubject needs the mail string
// returns true if the subject of the mail contains ""
// returns true if the subject of the mail contains not "" or there is an error reading the mail subject
func (b *BackendMail) checkMailSubject(mailstring string) bool {
	subject, err := getSubject(mailstring)
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

// CloseConnection needs the client and closes the connection
func (b *BackendMail) closeConnection(c *client.Client) {
	// close connection to server
	c.Logout()
	log.Println("logt out")
}
