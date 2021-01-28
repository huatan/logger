package logger

import "fmt"

type Config struct {
	//日志方法
	Method string
	//日志存放路径
	LogPath string
	//日志名
	LogName string
	//日志过滤级别
	LogLevel int
	//文件日志缓存大小
	LogChanSize int
	//文件日志切割类型
	LogSplitType int
	//文件日志切割单位数量
	LogSplitSize int64

}

var log LogInterface

//file,初始化一个文件日志实例
//console,初始化console日志实例
func InitLogger(config Config)(err error){
	switch config.Method {
	case "file":
		log,err = NewFileLogger(config)
	case "console":
		log,err = NewConsoleLogger(config)
	default:
		err = fmt.Errorf("unsupport logger method:%s",config.Method)
	}
	if err != nil {
		return
	}
	Info("init logger success")
	return
}

func Debug(format string,args ...interface{})  {
	log.Debug(format,args...)
}
func Trace(format string,args ...interface{})  {
	log.Trace(format,args...)
}
func Info(format string,args ...interface{})  {
	log.Info(format,args...)
}
func Warn(format string,args ...interface{})  {
	log.Warn(format,args...)
}
func Error(format string,args ...interface{})  {
	log.Error(format,args...)
}
func Fatal(format string,args ...interface{})  {
	log.Fatal(format,args...)
}
func Close()  {
	log.Close()
}

