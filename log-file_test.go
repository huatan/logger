package logger

import (
	"fmt"
	logold "log"
	"testing"
	"time"
)

func init() {
	//初始化日志
	err := InitLogger(Config{
		Method:       "file",
		LogPath:      "log/", //需要预先建立log文件夹
		LogName:      "neighbor",
		LogLevel:     LogLevelInfo,
		LogSplitType: LogSplitTypeSize,
		LogSplitSize: 50 << 20,
	})
	if err != nil {
		logold.Fatal(err)
	}
}
func TestFileLogger(t *testing.T) {

	Debug("test Debug\n")
	Trace("test Trace\n")
	Info("test Info\n")
	Warn("test Warn\n")
	Error("test Error\n")
	Fatal("test Fatal\n")
	Close()
	fmt.Println(time.Now())
}

func TestConsoleLogger(t *testing.T) {
	Debug("test Debug\n")
	Trace("test Trace\n")
	Info("test Info\n")
	Warn("test Warn\n")
	Error("test Error\n")
	Fatal("test Fatal\n")
	Close()
	fmt.Println(time.Now())
}
