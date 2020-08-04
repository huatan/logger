package logger

type LogInterface interface {
	SetLevel(level int)
	Init()
	Debug(formart string, args ...interface{})
	Trace(formart string, args ...interface{})
	Info(formart string, args ...interface{})
	Warn(formart string, args ...interface{})
	Error(formart string, args ...interface{})
	Fatal(formart string, args ...interface{})
	close()
}
