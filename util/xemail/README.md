email
=====

The ```email``` package currently supports the following:
*  From, To, Bcc, and Cc fields
*  Email addresses in both "test@example.com" and "First Last &lt;test@example.com&gt;" format
*  Text and HTML Message Body
*  Attachments
*  Read Receipts
*  Custom headers
*  More to come!

### Examples
#### Sending email using Gmail
```go
e := email.NewEmail()
e.From = "Jordan Wright <test@gmail.com>"
e.To = []string{"test@example.com"}
e.Bcc = []string{"test_bcc@example.com"}
e.Cc = []string{"test_cc@example.com"}
e.Subject = "Awesome Subject"
e.Text = []byte("Text Body is, of course, supported!")
e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"))
```

#### Another Method for Creating an Email
You can also create an email directly by creating a struct as follows:
```go
e := &email.Email {
	To: []string{"test@example.com"},
	From: "Jordan Wright <test@gmail.com>",
	Subject: "Awesome Subject",
	Text: []byte("Text Body is, of course, supported!"),
	HTML: []byte("<h1>Fancy HTML is supported, too!</h1>"),
	Headers: textproto.MIMEHeader{},
}
```

#### Creating an Email From an io.Reader
You can also create an email from any type that implements the ```io.Reader``` interface by using ```email.NewEmailFromReader```.

#### Attaching a File
```go
e := NewEmail()
e.AttachFile("test.txt")
```

#### A Pool of Reusable Connections
```go
(var ch <-chan *email.Email)
p := email.NewPool(
	"smtp.gmail.com:587",
	4,
	smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"),
)
for i := 0; i < 4; i++ {
	go func() {
		for e := range ch {
			p.Send(e, 10 * time.Second)
		}
	}()
}
```
