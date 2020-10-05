package logger

import (
	"github.com/alex60217101990/vse-instrumenty-bst/external/configs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Strings []string

func (ss Strings) ToInterfaceSlice() []interface{} {
	iface := make([]interface{}, len(ss))
	for i := range ss {
		iface[i] = ss[i]
	}
	return iface
}

type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewZapLogger() Logger {
	var (
		cfg zap.Config
		l   ZapLogger
		err error
	)
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	cfg.Encoding = "json"
	if len(configs.Conf.Logger.LogsPath) > 0 {
		cfg.OutputPaths = []string{configs.Conf.Logger.LogsPath} // "/tmp/logs"
	} else {
		cfg.OutputPaths = []string{"stdout"}
	}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.InitialFields = map[string]interface{}{
		"service": configs.Conf.ServiceName,
	}
	if configs.Conf.IsDebug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	l.logger, err = cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		DefaultLogger.Fatal(err)
	}

	l.sugar = l.logger.Sugar()
	defer func() {
		l.sugar.Sync()
		l.logger.Sync()
	}()

	return &l
}

func (l *ZapLogger) GetNativeLogger() interface{} {
	return l.logger
}

func (l *ZapLogger) Close() {
	l.sugar.Sync()
	l.logger.Sync()
}

func (l *ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Infof(tpl string, args ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Infof(tpl, args...)
}

func (l *ZapLogger) Errorw(err string, keysAndValues ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Errorw(err, keysAndValues...)
}

func (l *ZapLogger) Errorf(tpl string, args ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Infof(tpl, args...)
}

func (l *ZapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Warnf(tpl string, args ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Warnf(tpl, args...)
}

func (l *ZapLogger) Fatalw(err string, keysAndValues ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Warnw(err, keysAndValues...)
}

func (l *ZapLogger) Fatalf(tpl string, args ...interface{}) {
	defer func() {
		_ = l.sugar.Sync()
	}()

	l.sugar.Warnf(tpl, args...)
}
