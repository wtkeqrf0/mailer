package guzzle_logger

import (
	"encoding/json"
	"mailer/config"
	"mailer/pkg/logger"
	"mailer/pkg/rabbit/publisher"
)

//go:generate ifacemaker -f api.go -o interface.go -i API -s GuzzleAPI -p guzzle_logger -y "Controller describes methods, implemented by the guzzle_logger package."
//go:generate mockgen -package mock -source interface.go -destination mock/mock_guzzle_logger.go
type GuzzleAPI struct {
	remoteService    *string
	remoteSubService *string
	logPublisher     publisher.Publisher
	logger           logger.Logger
}

func New(service, subService string, logger logger.Logger, params config.QueueConnection) (*GuzzleAPI, error) {
	queuePublisher, err := publisher.New(params)
	if err != nil {
		return nil, err
	}

	return &GuzzleAPI{
		remoteService:    &service,
		remoteSubService: &subService,
		logPublisher:     queuePublisher,
		logger:           logger,
	}, nil
}

func (g *GuzzleAPI) SendLog(level string, description string, messageType *string, msg interface{}) {
	switch level {
	case LevelInfo:
		g.logger.Info(description)
	case LevelWarning:
		g.logger.Warn(description)
	case LevelError:
		g.logger.Errorf(description)
	}

	log := Log{
		LogMessage:       msg,
		LogLevel:         &level,
		RemoteService:    g.remoteService,
		RemoteSubService: g.remoteSubService,
		LogDescription:   &description,
	}

	if messageType != nil {
		log.MessageType = messageType
	}

	data, err := json.Marshal(log)
	if err != nil {
		g.logger.ErrorFull(err)
		return
	}

	if err = g.logPublisher.Publish(data); err != nil {
		g.logger.ErrorFull(err)
	}
}
