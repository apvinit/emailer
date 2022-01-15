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

type Client struct {
	name     string
	username string
	password string
	host     string
	addr     string
	from     string
	sender   string
	smtp     *smtp.Client
	m        *sync.Mutex
}

// NewClient create a smtp client using the credentials passed
func NewClient(name, from, username, password, host, addr string, m *sync.Mutex) (*Client, error) {
	c := &Client{name, username, password, host, addr, from, fmt.Sprintf("%s <%s>", name, from), nil, m}
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
	Subject string
	Body    bytes.Buffer
}

func (mail *Mail) BuildMail() string {
	msg := strings.Builder{}
	msg.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n")
	msg.WriteString(fmt.Sprintf("From: %s\r\n", mail.sender))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))
	msg.WriteString(fmt.Sprintf("\r\n%s\r\n", mail.Body.String()))
	return msg.String()
}

func (c *Client) sendMail(mail *Mail) error {
	c.m.Lock()
	defer c.m.Unlock()
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
	auth := smtp.PlainAuth("", c.username, c.password, c.host)
	err = c.smtp.Auth(auth)
	if err != nil {
		log.Fatal("error authenticating ", err)
		return err
	}
	return
}
