package logger

import (
	"fmt"
	"os"
	"time"
)

type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(config Config) (log LogInterface, err error) {

	logLevel := config.LogLevel

	log = &ConsoleLogger{
		level: logLevel,
	}
	return
}

func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	c.level = level
}
func (c *ConsoleLogger) Init() {}
func (c *ConsoleLogger) Debug(formart string, args ...interface{}) {
	c.writeLog(LogLevelDebug, formart, args...)
}
func (c *ConsoleLogger) Trace(formart string, args ...interface{}) {
	c.writeLog(LogLevelTrace, formart, args...)
}
func (c *ConsoleLogger) Info(formart string, args ...interface{}) {
	c.writeLog(LogLevelInfo, formart, args...)
}
func (c *ConsoleLogger) Warn(formart string, args ...interface{}) {
	c.writeLog(LogLevelWarn, formart, args...)
}
func (c *ConsoleLogger) Error(formart string, args ...interface{}) {
	c.writeLog(LogLevelError, formart, args...)
}
func (c *ConsoleLogger) Fatal(formart string, args ...interface{}) {
	c.writeLog(LogLevelFatal, formart, args...)
}
func (c *ConsoleLogger) close() {}
func (c *ConsoleLogger) writeLog(level int, format string, args ...interface{}) {
	if c.level > level {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05.000000000")
	levelString := getLevelString(level)
	filename, funcName, lineNo := GetLineInfo()

	msg := fmt.Sprintf(format, args...)

	fmt.Fprintf(os.Stdout, "%s %s (%s %s %d) %s\n", now, levelString, filename, funcName, lineNo, msg)

}
