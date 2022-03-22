package emailer

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEmailer(t *testing.T) {
	var client *Client
	client, err := NewClient(Options{
		Host: "smtp-test.kause.in",
		Port: "1025",
		User: "user",
		Pass: "pass",
		Name: "Mailhog 1",
		From: "mailhog1@mailhog.com",
	})
	if err != nil {
		t.Error(err)
	}

	// send multiple emails from same client
	for i := 0; i < 100; i++ {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprint("Test Email Content ", i+1))

		mail := Mail{
			To:      []string{"vinit@crowdpouch.com", "vinit@kreateworld.in"},
			Cc:      []string{"test@email.com", "test-01@email.com"},
			Bcc:     []string{"some@email.com"},
			Subject: fmt.Sprint("Test Email Subject ", i+1),
			Body:    buf,
		}
		mail.SendWith(client)
		// time.Sleep(time.Second * 1)
	}
}

/*
func TestEmailerMT(t *testing.T) {
	var client *Client

	client, err := NewClient(
		Options{
			Host: "smtp.mailtrap.io",
			Port: "587",
			User: "5b3e8dc2b0b914",
			Pass: "0a67a6dd68737a",
			Name: "Mailtrap 1",
			From: "mailtrap1@mailtrap.io",
		})
	if err != nil {
		t.Error(err)
	}
	var buf bytes.Buffer
	buf.WriteString("Mailtrap content 1")

	mail := Mail{
		To:      []string{"vinit@crowdpouch.com"},
		Subject: "Mailtrap subject 1",
		Body:    buf,
	}

	mail.SendWith(client)
}
*/
