package yaac_backend_mail

import (
	"log"

	yaac_shared "github.com/DHBW-SE-2023/YAAC/internal/shared"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func (b *BackendMail) GetResponse(input yaac_shared.EmailData) string {
	msg, err := getLatestMessage(input.MailServer, input.Email, input.Password)
	if err != nil {
		return "Something went wrong please try again"
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
		log.Fatal(err)
	}

	firstUnseen := mbox.UnseenSeqNum
	seqset := new(imap.SeqSet)
	seqset.AddNum(firstUnseen)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	if err := <-done; err != nil {
		return "", err
	}

	var firstUnseenMsg string
	for msg := range messages {
		firstUnseenMsg = msg.Envelope.Subject
	}

	return firstUnseenMsg, nil
}
