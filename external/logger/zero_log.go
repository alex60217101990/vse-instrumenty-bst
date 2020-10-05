package logger

import (
	"fmt"
	"time"

	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"

	"github.com/ilya1st/rotatewriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type ZeroLogger struct {
	logger zerolog.Logger
	writer *rotatewriter.RotateWriter
}

func NewZeroLogger() Logger {
	var (
		l   ZeroLogger
		err error
	)

	if len(configs.Conf.Logger.LogsPath) > 0 {
		l.writer, err = rotatewriter.NewRotateWriter(configs.Conf.Logger.LogsPath, 8)
		if err != nil {
			DefaultLogger.Fatal(err)
		}

		l.logger = zerolog.New(l.writer).With().Caller().Str("service", configs.Conf.ServiceName).Timestamp().Stack().Logger()
	} else {
		l.logger = log.With().Caller().Str("service", configs.Conf.ServiceName).Timestamp().Stack().Logger()
	}

	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return &l
}

func (l *ZeroLogger) Close() {
	if l.writer != nil {
		l.writer.Rotate(nil)
		l.writer.CloseWriteFile()
	}
}

func (l *ZeroLogger) GetNativeLogger() interface{} {
	return l.logger
}

func (l *ZeroLogger) Infow(msg string, keysAndValues ...interface{}) {
	var (
		key string
	)

	e := l.logger.Info()
	for i, val := range keysAndValues {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", val)
			continue
		}
		e.Str(key, fmt.Sprintf("%v", val))
		key = ""
	}
	e.Msg(msg)
}

func (l *ZeroLogger) Infof(tpl string, args ...interface{}) {
	l.logger.Info().Msgf(tpl, args...)
}

func (l *ZeroLogger) Errorw(err string, keysAndValues ...interface{}) {
	var (
		key string
	)

	e := l.logger.Error().Stack()
	for i, val := range keysAndValues {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", val)
			continue
		}
		e.Str(key, fmt.Sprintf("%v", val))
		key = ""
	}
	e.Msg(err)
}

func (l *ZeroLogger) Errorf(tpl string, args ...interface{}) {
	l.logger.Error().Stack().Msgf(tpl, args...)
}

func (l *ZeroLogger) Warnw(msg string, keysAndValues ...interface{}) {
	var (
		key string
	)

	e := l.logger.Warn()
	for i, val := range keysAndValues {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", val)
			continue
		}
		e.Str(key, fmt.Sprintf("%v", val))
		key = ""
	}
	e.Msg(msg)
}

func (l *ZeroLogger) Warnf(tpl string, args ...interface{}) {
	l.logger.Warn().Msgf(tpl, args...)
}

func (l *ZeroLogger) Fatalw(err string, keysAndValues ...interface{}) {
	var (
		key string
	)

	e := l.logger.Fatal()
	for i, val := range keysAndValues {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", val)
			continue
		}
		e.Str(key, fmt.Sprintf("%v", val))
		key = ""
	}
	e.Msg(err)
}

func (l *ZeroLogger) Fatalf(tpl string, args ...interface{}) {
	l.logger.Fatal().Msgf(tpl, args...)
}
