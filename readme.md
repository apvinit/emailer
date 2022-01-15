# Emailer

A pretty wrapper around go `net/smtp`
send multiple email using single client

```go
func main() {
  var sarah *emailers.Client
  var sarahX sync.Mutex
  sarah, err = emailers.NewClient("Sender Name", "user@email.com", "user@email.com", "password", "smtp.gmail.com", "smtp.gmail.com:587", &sarahX)
	if err != nil {
		log.Fatal("Could not initialize Sarah Emailers")
	}

  var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("email-template.tmpl.html"))
	tmpl.Execute(&buf, nil)

	mail := emailers.Mail{
		To:      []string{"test@user.com"},
		Subject: "Test subject",
		Body:    buf,
	}
	mail.SendWith(sarah)

}
```