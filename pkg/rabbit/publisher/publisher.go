package publisher

import (
	"github.com/wagslane/go-rabbitmq"
	"mailer/config"
)

//go:generate ifacemaker -f publisher.go -o interface.go -i Publisher -s RabbitMQPublisher -p publisher -y "Controller describes methods, implemented by the publisher package."
//go:generate mockgen -package mock -source interface.go -destination mock/mock_publisher.go
type RabbitMQPublisher struct {
	publisher *rabbitmq.Publisher
	queueName string
}

func New(params config.QueueConnection) (*RabbitMQPublisher, error) {
	conn, err := rabbitmq.NewConn(params.URL)
	if err != nil {
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQPublisher{
		publisher: publisher,
		queueName: params.QueueName,
	}, nil
}

func (p *RabbitMQPublisher) Publish(msg []byte) error {
	return p.publisher.Publish(
		msg,
		[]string{p.queueName},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
}
