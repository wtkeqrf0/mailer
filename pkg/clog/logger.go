package clog

import (
	"github.com/goccy/go-json"
	"log"
	"mailer/pkg/rabbit"
	"time"
)

type Logger struct {
	produce     rabbit.Produce
	serviceName string
	location    *time.Location
}

type LogMessage struct {
	ServiceName string    `json:"service_name"`
	Message     string    `json:"message"`
	MessageType string    `json:"message_type"`
	TimeDate    time.Time `json:"time_date"`
}

func New(produce rabbit.Produce, serviceName string) *Logger {
	return &Logger{
		produce:     produce,
		serviceName: serviceName,
		location:    time.FixedZone("UTC+3", 3*60*60),
	}
}

func (l *Logger) SendLog(msg string, msgLevel Level) {
	log.Println(msg)
	jsonMsg, _ := json.Marshal(LogMessage{
		ServiceName: l.serviceName,
		Message:     msg,
		MessageType: msgLevel.String(),
		TimeDate:    time.Now().In(l.location),
	})
	if err := l.produce(jsonMsg); err != nil {
		log.Printf("failed to publish, %v", err)
	}
}
