package emailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"sync"
)

// A Client represent a connection to smtp server with identity
type Client struct {
	// Name of the client
	name     string
	username string
	password string
	// smtp sever without port
	host string
	// smtp sever address with port
	addr   string
	from   string
	sender string
	smtp   *smtp.Client
	mu     *sync.Mutex
}

type Options struct {
	Host string // smtp host
	Port string // smtp port
	User string // smtp username
	Pass string // smtp password
	Name string // sender name for the client
	From string // sender email for the client
}

// NewClient create a emailer client for sending emails using provided
// smtp sevice provider
func NewClient(o Options) (*Client, error) {
	var mu sync.Mutex
	c := &Client{
		name:     o.Name,
		username: o.User,
		password: o.Pass,
		host:     o.Host,
		addr:     fmt.Sprintf("%s:%s", o.Host, o.Port),
		from:     o.From,
		sender:   fmt.Sprintf("%s <%s>", o.Name, o.From),
		mu:       &mu,
	}
	err := c.connnect()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (mail *Mail) SendWith(c *Client) {
	c.sendMail(mail)
}

// SendMail sends the mail
func (c *Client) SendMail(mail *Mail) error {
	return c.sendMail(mail)
}

type Mail struct {
	sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    bytes.Buffer
}

func (mail *Mail) BuildMail() string {
	msg := strings.Builder{}
	msg.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n")
	msg.WriteString(fmt.Sprintf("From: %s\r\n", mail.sender))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";")))
	msg.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(mail.Bcc, ";")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))
	msg.WriteString(fmt.Sprintf("\r\n%s\r\n", mail.Body.String()))
	return msg.String()
}

func (c *Client) sendMail(mail *Mail) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	mail.sender = c.sender
	err := c.smtp.Noop()
	if err != nil {
		err = c.connnect()
		if err != nil {
			log.Fatal("Noop Reconnectiong Failed", err)
			return err
		}
	}
	err = c.smtp.Mail(c.from)
	if err != nil {
		log.Fatal("Error creating mail from", err)
	}
	for _, v := range mail.To {
		c.smtp.Rcpt(v)
	}
	for _, v := range mail.Cc {
		c.smtp.Rcpt(v)
	}
	msg := mail.BuildMail()
	w, err := c.smtp.Data()
	if err != nil {
		log.Printf("[Failed] \"%s\" to %s\n. %s\n", mail.Subject, mail.To, err.Error())
		c.smtp.Reset()
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Println("Error writing mail bytes", err)
		return err
	}
	err = w.Close()
	if err != nil {
		log.Println("Error closing mail writer", err)
		return err
	}
	fmt.Printf("[SENT] \"%s\" to %s\n", mail.Subject, mail.To)
	return nil
}

func (c *Client) connnect() (err error) {
	if c.smtp != nil {
		c.smtp.Close()
	}
	c.smtp, err = smtp.Dial(c.addr)
	if err != nil {
		return err
	}
	c.smtp.StartTLS(&tls.Config{InsecureSkipVerify: true})
	auth := smtp.CRAMMD5Auth(c.username, c.password)
	err = c.smtp.Auth(auth)
	if err != nil {
		log.Fatal("error authenticating ", err)
		return err
	}
	return
}
