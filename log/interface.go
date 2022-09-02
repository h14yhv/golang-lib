package log

type Logger interface {
	Log(lvl LogLevel, message string)
	Fatal(message string)
	Fatalf(format string, a ...interface{})
	Panic(message string)
	Panicf(format string, a ...interface{})
	Critical(message string)
	Criticalf(format string, a ...interface{})
	Error(message string)
	Errorf(format string, a ...interface{})
	Warning(message string)
	Warningf(format string, a ...interface{})
	Notice(message string)
	Noticef(format string, a ...interface{})
	Info(message string)
	Infof(format string, a ...interface{})
	Debug(message string)
	Debugf(format string, a ...interface{})
}
