package usecase

import (
	"bytes"
	"fmt"
	"github.com/toorop/go-dkim"
	ht "html/template"
	"io"
	"mailer/config"
	"mailer/internal/consumer"
	"mailer/internal/single_sender"
	"mailer/pkg/logger"
	"mailer/pkg/mail"
	"mailer/pkg/ptr"
	"sync"
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
	srv                        *mail.SMTPServer
	client                     *mail.SMTPClient
	repo                       single_sender.Repository
	returnPath, from, errorsTo string
	isDkimSet                  bool
	dkim                       dkim.SigOptions
	log                        logger.Logger
	m                          sync.Mutex
}

func NewSingleSender(cfg config.EmailConnection, repo single_sender.Repository, log logger.Logger) (*SingleSender, error) {
	srv := mail.NewSMTPClient()
	srv.Host = cfg.Host
	srv.Port = int(cfg.Port)
	srv.Username = cfg.Username
	srv.Password = cfg.Password

	client, err := srv.Connect()
	if err != nil {
		return nil, err
	}

	res := SingleSender{
		srv:        srv,
		client:     client,
		repo:       repo,
		from:       fmt.Sprintf(`"%s" <%s>`, cfg.Name, cfg.Username),
		returnPath: cfg.ReturnPath,
		errorsTo:   cfg.ErrorsTo,
		log:        log,
	}

	if len(cfg.PrivateKey) != 0 {
		res.dkim, res.isDkimSet = dkim.SigOptions{
			Version:          1,
			PrivateKey:       cfg.PrivateKey,
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
	}

	go func() {
		for range time.Tick(time.Second * 30) {
			res.clientKeepAlive()
		}
	}()

	return &res, nil
}

// SendEmail to the specified receivers with given body data.
//
// Can also get templates from mongoDB, if found.
func (s *SingleSender) SendEmail(emailMsg consumer.Email) (recipients []string, err error) {
	email := mail.NewHighPriorityMSG(s.from, s.errorsTo, s.returnPath)
	defer func() {
		recipients = email.GetRecipients()
	}()

	email.SetSubject(emailMsg.Subject)
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
		email.Attach(ptr.Get(mail.File(*file)))
	}

	email.Parts = make([]mail.Part, len(emailMsg.Parts))
	for i, part := range emailMsg.Parts {
		var (
			buf = new(bytes.Buffer)
			t   interface {
				Execute(wr io.Writer, data any) error
			}
		)

		switch part.ContentType {
		case consumer.TextHTML, consumer.TextAMP:
			t, err = ht.New("").Parse(part.Body)
		default:
			t, err = tt.New("").Parse(part.Body)
		}
		if err != nil {
			return
		}

		if err = t.Execute(buf, emailMsg.PartValues); err != nil {
			return
		}

		email.Parts[i] = mail.Part{
			ContentType: string(part.ContentType),
			Body:        buf,
		}
	}

	if s.isDkimSet {
		email.SetDkim(s.dkim)
	}

	if err = email.Error; err != nil {
		return
	}

	s.m.Lock()
	defer s.m.Unlock()
	if err = email.SendEnvelopeFrom(s.from, s.client); err == io.EOF {
		s.clientKeepAlive()
		err = email.SendEnvelopeFrom(s.from, s.client)
	}
	return
}

func (s *SingleSender) clientKeepAlive() {
	s.m.Lock()
	defer s.m.Unlock()
	if err := s.client.Noop(); err != nil {
		s.log.Warnf("email connection was closed due %v", err)
		var newClient *mail.SMTPClient
		newClient, err = s.srv.Connect()
		if err != nil {
			s.log.Warnf("failed to create new connection due %v", err)
			return
		}

		// TODO do with channels or sync
		*s.client = *newClient
		s.log.Info("successfully created new connection")
	}
}
