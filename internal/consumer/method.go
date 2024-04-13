package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"mailer/pkg/clog"
	"strings"
)

// Consumer represents the client to send emails.
// The client can send emails using the specified credentials.
//
// UseCase interface and mock.MockUC implementation are generated on the basis of Consumer realization.
//
//go:generate ifacemaker -f *.go -o ../usecase.go -i UseCase -s Consumer -p consumer -y "Controller describes methods, implemented by the usecase package."
//go:generate mockgen -package mock -source ../usecase.go -destination mock/usecase_mock.go
type Consumer struct {
	logger         *clog.Logger
	repo           repo
	singleSenderUC sender.UseCase
	ch             <-chan amqp.Delivery
}

func NewConsumer(logger *clog.Logger, repo repo, singleSenderUC sender.UseCase, ch <-chan amqp.Delivery) *Consumer {
	return &Consumer{
		logger:         logger,
		repo:           repo,
		singleSenderUC: singleSenderUC,
		ch:             ch,
	}
}

func (s *Consumer) ProcessEmails() {
	for msg := range s.ch {
		go s.check(msg)
	}
}

func (s *Consumer) check(msg amqp.Delivery) {
	var (
		resend bool
		err    error
	)
	defer func() {

		if r := recover(); r != nil {
			resend = true
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}

		if resend {
			if err != nil {
				s.logger.SendLog(err.Error(), clog.LevelError)
			}
			err = msg.Nack(false, false)
		} else {
			if err != nil {
				log.Println(err.Error())
			}
			err = msg.Ack(false)
		}

		if err != nil {
			s.logger.SendLog(fmt.Sprintf("failed to proceed queue delivery, %v", err), clog.LevelFatal)
		}
	}()
	resend, err = s.processEmail(msg)
}

// processEmail marshal email to the consumer.Email struct and
// send it by other email packages.
//
// All method should be done by one goroutine.
func (s *Consumer) processEmail(body []byte) (bool, error) {
	emailMsg := new(Email)
	if err := json.Unmarshal(body, emailMsg); err != nil {
		return false, fmt.Errorf("failed to unmarshal message %s due %v", string(body), err)
	}

	switch emailMsg.Settings.Locale {
	case LocaleRu:
	default:
		emailMsg.Settings.Locale = LocaleEn
	}

	if err := s.repo.GetTemplateByName(emailMsg); err != nil {
		return true, err
	}

	if (len(emailMsg.Parts) == 0 && len(emailMsg.Files) == 0) || emailMsg.Subject == "" {
		return false, errors.New("email body doesn't have any part, file or subject")
	}

	switch emailMsg.Settings.From {
	case SingleSenderAddress, "":
		recipients, err := s.singleSenderUC.SendEmail(emailMsg)
		if err != nil {
			err = fmt.Errorf("failed to send email to %s due %v", strings.Join(recipients, ", "), err)
			return true, err
		}
		log.Printf("email was sent to %s", strings.Join(recipients, ", "))
	default:
		return false, fmt.Errorf("from (%s) is not valid", emailMsg.Settings.From)
	}
	return false, nil
}
