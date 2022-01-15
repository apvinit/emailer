# Emailer

A pretty wrapper around go `net/smtp`
send multiple email using single client

```go
func main() {
  var client *emailer.Client
  var clientX sync.Mutex
  client, err = emailer.NewClient("Sender Name", "user@email.com", "user@email.com", "password", "smtp.gmail.com", "smtp.gmail.com:587", &clientX)
	if err != nil {
		log.Fatal("Could not initialize client Emailer")
	}

  var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("email-template.tmpl.html"))
	tmpl.Execute(&buf, nil)

	mail := emailer.Mail{
		To:      []string{"test@user.com"},
		Subject: "Test subject",
		Body:    buf,
	}
	mail.SendWith(client)

}
```