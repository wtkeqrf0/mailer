package sender

import (
	"context"
	"fmt"
	"github.com/toorop/go-dkim"
	"log"
	"mailer/config"
	"mailer/pkg/mail"
	"os"
	"time"
)

//go:generate ifacemaker -f *.go -o sender_if.go -i Sender -s sender -p sender -y "Sender represents the email client."
type sender struct {
	srv        *mail.SMTPServer
	clientPool chan *mail.SMTPClient // using not sync.Pool cuz chan has type definition
	isDkimSet  bool
	dkim       dkim.SigOptions
	createMsg  mail.CreateEmailMessage
}

func New(ctx context.Context, cfg config.Email) Sender {
	s := sender{
		srv:        mail.NewSMTPClient(cfg),
		clientPool: make(chan *mail.SMTPClient, 100),
		createMsg: mail.NewMSGCreator(
			fmt.Sprintf(`"%s" <%s>`, cfg.Name, cfg.Username),
			cfg.ErrorsTo,
			cfg.ReturnPath,
		),
	}

	// test client
	if client, err := getClient(s.srv); err != nil {
		panic(err)
	} else {
		s.clientPool <- client
	}

	context.AfterFunc(ctx, s.clean())

	// set dkim, if specified
	if cfg.PrivateKeyPath != "" {
		privateKey, err := os.ReadFile(cfg.PrivateKeyPath)
		if err != nil {
			panic(err)
		}

		s.dkim, s.isDkimSet = dkim.SigOptions{
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

	return &s
}

// Send to the specified receivers with given body data.
//
// Can also get templates from mongoDB, if found.
func (s *sender) Send(receivedEmail *mail.Parsable) error {
	email := receivedEmail.ToEmail(s.createMsg())

	if s.isDkimSet {
		email.SetDkim(s.dkim)
	}

	if email.Error != nil {
		return email.Error
	}
	return s.send(email)
}

// send email message without error
func (s *sender) send(email *mail.Email) error {
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
			s.clientPool <- client // client is healthy - insert into pool
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

func (s *sender) clean() func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		for {
			select {
			case client := <-s.clientPool:
				_ = client.Quit()
			case <-ctx.Done():
				return
			}
		}
	}
}
