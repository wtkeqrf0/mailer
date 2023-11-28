package logger

import (
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"mailer/config"
	"os"
	"runtime"
)

//go:generate ifacemaker -f zerolog.go -o interface.go -i Logger -s ApiLogger -p logger -y "Controller describes methods, implemented by the logger package."
//go:generate mockgen -package mock -source interface.go -destination mock/mock_logger.go
type ApiLogger struct {
	loki   promtail.Client
	logger zerolog.Logger
	level  string
}

func NewApiLogger(logCfg config.Logger, serviceName string) (*ApiLogger, error) {
	a := ApiLogger{level: logCfg.Level}

	var w zerolog.LevelWriter
	a.logger = log.With().Caller().Logger()
	if logCfg.InFile {
		logFile, err := os.OpenFile(logCfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			return nil, err
		}
		w = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)
	} else {
		w = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	a.logger = zerolog.New(w).Level(loggerLevelMap[logCfg.Level]).With().Timestamp().Logger().Hook(copyLogger{&a})
	return &a, nil
}

var loggerLevelMap = map[string]zerolog.Level{
	"debug":    zerolog.DebugLevel,
	"info":     zerolog.InfoLevel,
	"warn":     zerolog.WarnLevel,
	"error":    zerolog.ErrorLevel,
	"panic":    zerolog.PanicLevel,
	"fatal":    zerolog.FatalLevel,
	"noLevel":  zerolog.NoLevel,
	"disabled": zerolog.Disabled,
}

type copyLogger struct {
	*ApiLogger
}

func (c copyLogger) Run(_ *zerolog.Event, _ zerolog.Level, _ string) {}

func (a *ApiLogger) Debug(msg string) {
	a.logger.Debug().Msg(msg)
}

func (a *ApiLogger) Debugf(template string, args ...interface{}) {
	a.logger.Debug().Msgf(template, args...)
}

func (a *ApiLogger) Info(msg string) {
	a.logger.Info().Msg(msg)
}

func (a *ApiLogger) Infof(template string, args ...interface{}) {
	a.logger.Info().Msgf(template, args...)
}

func (a *ApiLogger) Warn(msg string) {
	a.logger.Warn().Msg(msg)
}

func (a *ApiLogger) Warnf(template string, args ...interface{}) {
	a.logger.Warn().Msgf(template, args...)
}

func (a *ApiLogger) Error(err error) {
	a.logger.Error().Msg(err.Error())
}

func (a *ApiLogger) Errorf(template string, args ...interface{}) {
	a.logger.Error().Msgf(template, args...)
}

func (a *ApiLogger) Panic(msg string) {
	a.logger.Panic().Msg(msg)
}

func (a *ApiLogger) Panicf(template string, args ...interface{}) {
	a.logger.Panic().Msgf(template, args...)
}

func (a *ApiLogger) Fatal(msg string) {
	a.logger.Fatal().Msg(msg)
}

func (a *ApiLogger) Fatalf(template string, args ...interface{}) {
	a.logger.Fatal().Msgf(template, args...)
}

func (a *ApiLogger) Tracef(s string, i ...interface{}) {
	go a.logger.Trace().Msgf(s, i...)
}

func (a *ApiLogger) ErrorFull(error error) {
	pc, _, line, _ := runtime.Caller(1)
	det := runtime.FuncForPC(pc)
	msg := fmt.Sprintf("ERROR:\n%s :: %d :: %s", det.Name(), line, error.Error())
	a.logger.Error().Stack().Err(error).Msg(msg)
}
