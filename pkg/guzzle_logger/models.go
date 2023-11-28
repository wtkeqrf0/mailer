package guzzle_logger

const (
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"
)

type Log struct {
	LogMessage       interface{} `json:"log_message"`
	LogLevel         *string     `json:"log_level"`
	RemoteService    *string     `json:"remote_service"`
	RemoteSubService *string     `json:"remote_sub_service"`
	LogDescription   *string     `json:"log_description"`
	MessageType      *string     `json:"message_type"`
}

type LogMessage struct {
	Body *string `json:"body"`
}
