package main

import (
	"github.com/huatan/logger"
	"log"
)


func Run() {
	for {
		logger.Info("user server is running")
		//time.Sleep(time.Second)
	}
}
func main() {
	err := logger.InitLogger(logger.Config{
		Method:       "file",
		LogPath:      "log/",//需要预先建立log文件夹
		LogName:      "test",
		LogLevel:     logger.LogLevelTrace,
		LogSplitType: logger.LogSplitTypeSize,
		LogSplitSize: 50 << 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	Run()
	return
}
