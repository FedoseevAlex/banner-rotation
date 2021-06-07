package logger

import (
	"log"
	"os"
	"strings"

	"github.com/FedoseevAlex/banner-rotation/internal/types"
	"github.com/rs/zerolog"
)

type (
	LogArgs types.LogFields
	Logger  struct {
		logger zerolog.Logger
	}
)

func New(level, file string) *Logger {
	logfile, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o640)
	if err != nil {
		log.Println("can't open file ", file, err)
		return nil
	}

	logger := zerolog.New(logfile).With().Timestamp().Logger()
	numLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		logger.Warn().
			Str("level", level).
			Msg("Unknown log level. Using DEBUG.")
	} else {
		logger = logger.Level(numLevel)
	}

	return &Logger{logger: logger}
}

func (l Logger) Trace(msg string, fields types.LogFields) {
	l.logWithLevel(zerolog.TraceLevel, msg, fields)
}

func (l Logger) Debug(msg string, fields types.LogFields) {
	l.logWithLevel(zerolog.DebugLevel, msg, fields)
}

func (l Logger) Info(msg string, fields types.LogFields) {
	l.logWithLevel(zerolog.InfoLevel, msg, fields)
}

func (l Logger) Warn(msg string, fields types.LogFields) {
	l.logWithLevel(zerolog.WarnLevel, msg, fields)
}

func (l Logger) Error(msg string, fields types.LogFields) {
	l.logWithLevel(zerolog.ErrorLevel, msg, fields)
}

func (l Logger) logWithLevel(level zerolog.Level, msg string, fields types.LogFields) {
	record := l.logger.WithLevel(level)
	record = record.Fields(fields)
	record.Msg(msg)
}

func (l Logger) ChildLogger(name string) types.Logger {
	log := l.logger.
		With().
		Str("name", name).
		Logger()

	return Logger{logger: log}
}
