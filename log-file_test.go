package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestFileLogger(t *testing.T)  {
	logger := NewFileLogger(LogLevelDebug,"./","test")
	logger.Debug("test Debug\n")
	logger.Trace("test Trace\n")
	logger.Info("test Info\n")
	logger.Warn("test Warn\n")
	logger.Error("test Error\n")
	logger.Fatal("test Fatal\n")
	logger.close()
	fmt.Println(time.Now())
}

func TestConsoleLogger(t *testing.T)  {
	logger := NewConsoleLogger(LogLevelDebug)
	logger.Debug("test Debug\n")
	logger.Trace("test Trace\n")
	logger.Info("test Info\n")
	logger.Warn("test Warn\n")
	logger.Error("test Error\n")
	logger.Fatal("test Fatal\n")
	logger.close()
	fmt.Println(time.Now())
}
