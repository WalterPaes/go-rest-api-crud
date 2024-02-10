package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(message string, tags ...zap.Field)
	Error(message string, err error, tags ...zap.Field)
}

type logger struct {
	log *zap.Logger
}

func NewLogger(level, output string) *logger {
	logConfig := zap.Config{
		OutputPaths: []string{
			getLogOutput(output),
		},
		Level:    zap.NewAtomicLevelAt(getLogLevel(level)),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "message",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	log, _ := logConfig.Build()

	return &logger{
		log: log,
	}
}

func (l *logger) Info(message string, tags ...zap.Field) {
	l.log.Info(message, tags...)
	l.log.Sync()
}

func (l *logger) Error(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Info(message, tags...)
	l.log.Sync()
}

func getLogOutput(output string) string {
	if strings.ToLower(strings.TrimSpace(output)) == "" {
		return "stdout"
	}
	return output
}

func getLogLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "error":
		return zapcore.ErrorLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}
