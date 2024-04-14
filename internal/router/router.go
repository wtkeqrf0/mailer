package router

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"mailer/internal/sender"
	"mailer/pkg/clog"
	"mailer/pkg/mail"
)

// router represents the client to send emailr.
// The client can send emails using the specified credentialr.
//
// UseCase interface and mock.MockUC implementation are generated on the basis of router realization.
//
//go:generate ifacemaker -f *.go -o router_if.go -i Router -s router -p router -y "Router represents the message router."
type router struct {
	logger      *clog.Logger
	repo        Repository
	emailSender sender.Sender
	ch          <-chan amqp.Delivery
}

func New(logger *clog.Logger, repo Repository, sender sender.Sender, ch <-chan amqp.Delivery) Router {
	return &router{
		logger:      logger,
		repo:        repo,
		emailSender: sender,
		ch:          ch,
	}
}

func (r *router) ProcessEmails() {
	r.logger.SendLog("server started", clog.LevelInfo)
	for msg := range r.ch {
		go r.check(msg)
	}
}

func (r *router) check(msg amqp.Delivery) {
	var (
		resend bool
		cause  string
	)
	defer func() {
		re := recover()
		var err error
		switch {
		case re != nil:
			if err, _ = re.(error); err == nil {
				err = fmt.Errorf("%v", re)
			}
			r.logger.SendLog(err.Error(), clog.LevelFatal)
			err = msg.Nack(false, false)
		case resend:
			r.logger.SendLog(cause, clog.LevelError)
			err = msg.Nack(false, false)
		default:
			log.Println(cause)
			err = msg.Ack(false)
		}
		if err != nil {
			r.logger.SendLog(fmt.Sprintf("failed to proceed queue delivery, %v", cause), clog.LevelFatal)
		}
	}()
	resend, cause = r.processEmail(msg.Body)
}

// processEmail marshal email to the router.Email struct and
// send it by other email packager.
//
// All method should be done by one goroutine.
func (r *router) processEmail(body []byte) (bool, string) {
	emailMsg := new(mail.Parsable)
	if err := json.Unmarshal(body, emailMsg); err != nil {
		return false, fmt.Sprintf("failed to unmarshal message %s due %v", string(body), err)
	}

	if err := r.repo.GetTemplateByName(emailMsg); err != nil {
		return true, err.Error()
	}

	if (len(emailMsg.Parts) == 0 && len(emailMsg.Files) == 0) || emailMsg.Subject == "" {
		return false, "email body doesn't have any part, file or subject"
	}

	if err := r.emailSender.Send(emailMsg); err != nil {
		return true, fmt.Sprintf("failed to send email to %s, %v", emailMsg.Recipients(", "), err)
	}
	return false, fmt.Sprintf("email was sent to %s", emailMsg.Recipients(", "))
}
