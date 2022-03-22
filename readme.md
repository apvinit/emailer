# Emailer

A pretty wrapper around go `net/smtp`
send multiple email using single client

```go
func main() {
  var client *emailer.Client
  client, err := emailer.NewClient(Options{
		Host: "smtp-test.example.com",
		Port: "1025",
		User: "username",
		Pass: "password",
		Name: "Example User",
		From: "user@example.com",
	})
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