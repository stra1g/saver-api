package logging

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(err error, msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type logger struct {
	zlog zerolog.Logger
}

func (l *logger) Debug(msg string, fields map[string]interface{}) {
	l.zlog.Debug().Fields(fields).Msg(msg)
}

func (l *logger) Info(msg string, fields map[string]interface{}) {
	l.zlog.Info().Fields(fields).Msg(msg)
}
func (l *logger) Warn(msg string, fields map[string]interface{}) {
	l.zlog.Warn().Fields(fields).Msg(msg)
}

func (l *logger) Error(err error, msg string, fields map[string]interface{}) {
	l.zlog.Error().Err(err).Fields(fields).Msg(msg)
}

func (l *logger) Fatal(msg string, fields map[string]interface{}) {
	l.zlog.Fatal().Fields(fields).Msg(msg)
}

func Initialize(output io.Writer, isDebug bool) Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	
	level := zerolog.InfoLevel
	if isDebug {
		level = zerolog.DebugLevel
	}
	
	zl := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
		
	return &logger{
		zlog: zl,
	}
}
