package sender

import (
	"bytes"
	"fmt"
	"github.com/toorop/go-dkim"
	ht "html/template"
	"io"
	"log"
	"mailer/config"
	"mailer/internal/consumer"
	"mailer/pkg/mail"
	"os"
	tt "text/template"
	"time"
)

// SingleSender represents the client to send emails.
// The client can send emails using the specified data.
//
// UseCase interface and mock.MockUC implementation are generated on the basis of SingleSender realization.
//
//go:generate ifacemaker -f *.go -o ../usecase.go -i UseCase -s SingleSender -p single_sender -y "Controller describes methods, implemented by the usecase package."
//go:generate mockgen -package mock -source ../usecase.go -destination mock/usecase_mock.go
type SingleSender struct {
	srv        *mail.SMTPServer
	clientPool chan *mail.SMTPClient // using not sync.Pool cuz chan has type definition
	isDkimSet  bool
	dkim       dkim.SigOptions
	createMsg  mail.CreateEmailMessage
}

func NewSingleSender(cfg config.EmailConnection) *SingleSender {
	res := SingleSender{
		srv:        mail.NewSMTPClient(cfg),
		clientPool: make(chan *mail.SMTPClient, 100),
		createMsg: mail.NewMSGCreator(
			fmt.Sprintf(`"%s" <%s>`, cfg.Name, cfg.Username),
			cfg.ErrorsTo,
			cfg.ReturnPath,
		),
	}

	// test client
	if client, err := getClient(res.srv); err != nil {
		panic(err)
	} else {
		res.clientPool <- client
	}

	// set dkim, if specified
	if cfg.PrivateKeyPath != "" {
		privateKey, err := os.ReadFile(cfg.PrivateKeyPath)
		if err != nil {
			panic(err)
		}

		res.dkim, res.isDkimSet = dkim.SigOptions{
			Version:          1,
			PrivateKey:       privateKey,
			Domain:           "_domainkey.crypto",
			Selector:         "mailru",
			Canonicalization: "relaxed/relaxed",
			Algo:             "rsa-sha256",
			Headers: []string{"date", "from", "to", "message-id", "subject",
				"mime-version", "content-type", "content-transfer-encoding"},
			QueryMethods:          []string{"dns/txt"},
			AddSignatureTimestamp: true,
			SignatureExpireIn:     7776000, // in seconds = 90 days
		}, true
	} else {
		log.Println("dkim is disabled")
	}

	return &res
}

// SendEmail to the specified receivers with given body data.
//
// Can also get templates from mongoDB, if found.
func (s *SingleSender) SendEmail(emailMsg *consumer.Email) error {
	email := s.createMsg().SetSubject(emailMsg.Subject)

	// SetDSN([]mail.DSN{mail.SUCCESS, mail.FAILURE}, false)

	if emailMsg.Sender != "" {
		email.SetSender(emailMsg.Sender)
	}

	if emailMsg.ReplyTo != "" {
		email.SetReplyTo(emailMsg.ReplyTo)
	}

	if len(emailMsg.To) != 0 {
		email.AddTo(emailMsg.To...)
	}

	if len(emailMsg.CopyTo) != 0 {
		email.AddCc(emailMsg.CopyTo...)
	}

	if len(emailMsg.BlindCopyTo) != 0 {
		email.AddBcc(emailMsg.BlindCopyTo...)
	}

	for _, file := range emailMsg.Files {
		f := mail.File(*file)
		email.Attach(&f)
	}

	email.Parts = make([]mail.Part, len(emailMsg.Parts))
	for i, part := range emailMsg.Parts {
		var (
			buf = new(bytes.Buffer)
			t   interface {
				Execute(wr io.Writer, data any) error
			}
			err error
		)

		switch part.ContentType {
		case consumer.TextHTML, consumer.TextAMP:
			t, err = ht.New("").Parse(part.Body)
		default:
			t, err = tt.New("").Parse(part.Body)
		}
		if err != nil {
			return err
		}

		if err = t.Execute(buf, emailMsg.PartValues); err != nil {
			return err
		}

		email.Parts[i] = mail.Part{
			ContentType: string(part.ContentType),
			Body:        buf,
		}
	}

	if s.isDkimSet {
		email.SetDkim(s.dkim)
	}

	if email.Error != nil {
		return email.Error
	}
	return s.send(email)
}

// send email message without error
func (s *SingleSender) send(email *mail.Email) error {
	var (
		client *mail.SMTPClient
		err    error
	)
	for i := 0; i < 10; i++ {
		select {
		case client = <-s.clientPool:
		case <-time.Tick(time.Millisecond * 250): // to avoid creating unnecessary clients
			client, err = getClient(s.srv)
			if err != nil {
				continue
			}
		}
		if err = email.Send(client); err == nil {
			s.clientPool <- client // client is healthy - insert in pool
			return nil
		}
	}
	return err
}

// getClient connects new smtp Client
func getClient(srv *mail.SMTPServer) (*mail.SMTPClient, error) {
	client, err := srv.Connect()
	if err != nil {
		return nil, err
	}
	go func() {
		for range time.Tick(time.Second * 30) {
			if err = client.Noop(); err != nil {
				return
			}
		}
	}()
	return client, nil
}
