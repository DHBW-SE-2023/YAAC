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

func getBase64AttachmentFromMailString(mailString string) (string, error) {

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
	path := "/Users/$USER/image_" + strconv.Itoa(rand.Intn(1000)) + ".jpg"
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

func decodeBase64(base64data string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func MailService(serverAddr string, username string, password string) error {

	for {

		//channel for mails
		messages := make(chan *imap.Message)

		//connect to the mail server
		c, err := connectToServer(serverAddr)
		if err != nil {
			return err
		}

		//close connection when the function returns (in case of an error)
		defer CloseConnection(c)

		//login with the user crudentials
		c, err = LogInToInbox(c, username, password)
		if err != nil {
			return err
		}

		//focus only on unread mails
		c, seqset, err := getOnlyUnreadMails(c)
		if err != nil {
			return err
		}
		//start fetchings mails
		go fetchMails(c, seqset, messages)

		//process fetched mails
		for msg := range messages {

			//no new unread mail
			if msg == nil {
				log.Println("all mails processed")
				break
			}

			//get mail message as string in order to prcess the mail further
			mailstring, err := getMailasString(msg)
			if err != nil {
				log.Println("Error decoding mail")
				continue
			}

			//check the subject of the mail and process it only if the subject contains "Anwesendtheitsliste"
			if checkMailSubject(mailstring) {

				//get the base64 encoded image from the mail, continue if there is no
				image, err := getBase64AttachmentFromMailString(mailstring)
				if err != nil {
					continue
				}
				log.Println("Image Attachment identified")

				//write that to database
				err = writeToDatabase(image)
				if err != nil {
					log.Println(("Error writing to Database"))
					continue
				}

				//mark mail as read in order to show that it has been proceed
				err = markMailAsRead(c, msg)
				if err != nil {
					log.Println(("Error marking mail as read!"))
					log.Print(err)
					continue
				}
				log.Println(("Data successfully stored in data base"))
			}
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

func checkMailSubject(mailstring string) bool {
	subject, err := getSubject(mailstring)
	if err != nil {
		log.Println("Error decoding mail subject")
		return false
	} else {
		log.Printf("Mail detected with Subject: %v\n", subject)
		// Check if mail Subject is correct
	}
	return true
}

func markMailAsRead(c *client.Client, msg *imap.Message) error {
	seqnums := new(imap.SeqSet)
	seqnums.AddNum(msg.SeqNum)
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	if err := c.Store(seqnums, item, flags, nil); err != nil {
		return err
	}
	return nil
}

func connectToServer(serverAddr string) (*client.Client, error) {
	// Connect to server
	c, err := client.DialTLS(serverAddr, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to server at: ", serverAddr)
	return c, err
}

func LogInToInbox(c *client.Client, username string, password string) (*client.Client, error) {
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

func getOnlyUnreadMails(c *client.Client) (*client.Client, *imap.SeqSet, error) {
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

func fetchMails(c *client.Client, seqset *imap.SeqSet, messages chan *imap.Message) {
	// Fetching mails
	log.Println(("start fetching"))
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchItem("BODY.PEEK[]")}, messages); err != nil {
		log.Println("no unread mails:")
	}
	log.Println(("end fetching"))
}

func CloseConnection(c *client.Client) {
	// close connection to server
	c.Logout()
	log.Println("logt out")
}
