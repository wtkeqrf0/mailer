package consumer

import (
	"github.com/wagslane/go-rabbitmq"
	"mailer/config"
)

func StartConsuming(
	params config.QueueConnection,
	consumeFunc func(d []byte),
) (*rabbitmq.Consumer, error) {
	conn, err := rabbitmq.NewConn(params.URL)
	if err != nil {
		return nil, err
	}

	consumer, err := rabbitmq.NewConsumer(
		conn, func(d rabbitmq.Delivery) (action rabbitmq.Action) {
			go consumeFunc(d.Body)
			return rabbitmq.Ack
		}, params.QueueName,
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsConcurrency(10),
		rabbitmq.WithConsumerOptionsQueueDurable,
	)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
