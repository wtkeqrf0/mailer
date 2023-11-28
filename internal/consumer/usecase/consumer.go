package usecase

import (
	"encoding/json"
	"fmt"
	"mailer/internal/consumer"
	"mailer/internal/single_sender"
	"mailer/pkg/guzzle_logger"
	"mailer/pkg/logger"
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
	tgLogger       guzzle_logger.API
	logger         logger.Logger
	repo           consumer.Repository
	singleSenderUC single_sender.UseCase
}

func NewConsumer(logger logger.Logger, tgLogger guzzle_logger.API, repo consumer.Repository, singleSenderUC single_sender.UseCase) *Consumer {
	return &Consumer{
		tgLogger:       tgLogger,
		logger:         logger,
		repo:           repo,
		singleSenderUC: singleSenderUC,
	}
}

// ProcessEmail marshal email to the consumer.Email struct and
// send it by other email packages.
//
// All method should be done by one goroutine.
func (s *Consumer) ProcessEmail(b []byte) {
	var emailMsg consumer.Email
	if err := json.Unmarshal(b, &emailMsg); err != nil {
		s.logger.Warnf("failed to unmarshal message %s due %v", string(b), err)
		return
	}

	switch emailMsg.Settings.Locale {
	case consumer.LocaleRu:
	default:
		emailMsg.Settings.Locale = consumer.LocaleEn
	}

	if err := s.repo.GetTemplateByName(&emailMsg); err != nil {
		s.tgLogger.SendLog(guzzle_logger.LevelWarning, fmt.Sprintf("failed to get template due %v", err), nil, emailMsg)
		return
	}

	if len(emailMsg.Parts) != 0 && len(emailMsg.Files) != 0 && emailMsg.Subject != "" {
		s.logger.Warn("email body doesn't have any part OR file OR subject")
		return
	}

	switch emailMsg.Settings.From {
	case consumer.SingleSenderAddress, "":
		if recipients, err := s.singleSenderUC.SendEmail(emailMsg); err != nil {
			s.tgLogger.SendLog(guzzle_logger.LevelError, fmt.Sprintf("failed to send email to %s due %v", strings.Join(recipients, ", "), err), nil, emailMsg)
		} else {
			s.logger.Infof("email was sent to %s", strings.Join(recipients, ", "))
		}
	default:
		s.logger.Warnf("from (%s) is not valid", emailMsg.Settings.From)
	}
}
