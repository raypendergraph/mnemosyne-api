package system

import (
	"github.com/rs/zerolog"
	"os"
)

type LogLevel string

var (
	DebugLevel LogLevel = "DEBUG"
	ErrorLevel LogLevel = "ERROR"
	WarnLevel  LogLevel = "WARN"
)

type LogController interface {
	SetLogLevel(level LogLevel)
}

type LoggerImpl struct {
	logger zerolog.Logger
}

func (r LoggerImpl) WithAdditionalMetaData(values MetaData) Logger {
	return LoggerImpl{
		logger: r.logger.With().Fields(values).Logger(),
	}
}

func (r LoggerImpl) LogAlways(msg string) {
	r.logger.Info().Msg(msg)
}

func (r LoggerImpl) LogDebug(msg string) {
	r.logger.Debug().Msg(msg)
}

func (r LoggerImpl) LogError(e Error) {
	r.logger.Error().Err(e).Msg(e.GetCause().GetDefinition().Message)
}

func (r LoggerImpl) LogWarning(msg string, errors ...error) {
	if len(errors) == 0 {
		r.logger.Warn().Msg(msg)
	}
	event := r.logger.Log()
	for _, e := range errors {
		event = event.Err(e)
	}
	event.Msg(msg)
}

func NewLogger() LoggerImpl {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
	logger := zerolog.New(multi).With().Timestamp().Logger()
	return LoggerImpl{logger: logger}
}
