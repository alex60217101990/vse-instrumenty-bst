package logger

type Logger interface {
	GetNativeLogger() interface{}
	Infow(msg string, keysAndValues ...interface{})
	Infof(tpl string, args ...interface{})
	Errorw(err string, keysAndValues ...interface{})
	Errorf(tpl string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Warnf(tpl string, args ...interface{})
	Fatalw(err string, keysAndValues ...interface{})
	Fatalf(tpl string, args ...interface{})
	Close()
}
