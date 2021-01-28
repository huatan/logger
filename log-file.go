package logger

import (
	"fmt"
	"github.com/huatan/logger/utils"
	"os"
	"time"
)

type FileLogger struct {
	level         int
	logPath       string
	logName       string
	file          *os.File
	errorFile     *os.File
	logDataChan   chan *utils.LogData
	logSplitType  int
	logSplitSize  int64
	lastSplitTime time.Time
}

func NewFileLogger(config Config) (log LogInterface, err error) {
	logPath := config.LogPath
	if logPath == "" {
		err = fmt.Errorf("not found logPath ")
		return
	}
	logName := config.LogName
	if logName == "" {
		err = fmt.Errorf("not found logName ")
		return
	}
	logChanSize := config.LogChanSize
	if logChanSize == 0 {
		logChanSize = 60000
	}

	//日志文件切分
	if config.LogSplitType == LogSplitTypeSize {
		if config.LogSplitSize == 0 {
			config.LogSplitSize = 100 << 20
		}
	} else {
		config.LogSplitType = LogSplitTypeHour
		if config.LogSplitSize == 0 {
			config.LogSplitSize = 1
		}
	}

	log = &FileLogger{
		level:        config.LogLevel,
		logPath:      logPath,
		logName:      logName,
		logDataChan:  make(chan *utils.LogData, logChanSize),
		logSplitType: config.LogSplitType,
		logSplitSize: config.LogSplitSize,
	}
	log.Init()
	return
}
func (f *FileLogger) Init() {
	//如果路径不存在，创造路径
	exists, err := utils.PathExists(f.logPath)
	if !exists {
		os.MkdirAll(f.logPath,0755)
	}
	//打开普通日志文件
	filename := fmt.Sprintf("%s%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
	}
	f.file = file

	//打开Error以上日志文件
	filename = fmt.Sprintf("%s%serror.log", f.logPath, f.logName)
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
	}
	f.errorFile = file
	go f.writeLogAsync()
}

func (f *FileLogger) splitFileHour(errorAndFatal bool) {
	now := time.Now()
	hour := now.Truncate(time.Hour).Sub(f.lastSplitTime.Truncate(time.Hour)).Hours()
	if now.Day() == f.lastSplitTime.Day() && hour < float64(f.logSplitSize) {
		return
	}
	var backupFilename string
	if errorAndFatal {
		backupFilename = fmt.Sprintf("%s%serror.log_bak%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(),
			now.Hour())
	} else {
		backupFilename = fmt.Sprintf("%s%s.log_bak%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(),
			now.Hour())
	}
	f.lastSplitTime = now
	file := f.file
	if errorAndFatal {
		file = f.errorFile
	}
	filename := file.Name()
	file.Close()
	os.Rename(filename, backupFilename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	for err != nil {
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		fmt.Println(err)
		time.Sleep(time.Second)
	}

	if errorAndFatal {
		f.errorFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) splitFileSize(errorAndFatal bool) {
	file := f.file
	if errorAndFatal {
		file = f.errorFile
	}
	statInfo, err := file.Stat()
	if err != nil {
		return
	}
	if statInfo.Size() < f.logSplitSize {
		return
	}

	now := time.Now()
	var backupFilename string
	if errorAndFatal {
		backupFilename = fmt.Sprintf("%s%serror.log_bak%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second())
	} else {
		backupFilename = fmt.Sprintf("%s%s.log_bak%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second())
	}
	f.lastSplitTime = now
	filename := file.Name()

	file.Close()
	os.Rename(filename, backupFilename)

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	for err != nil {
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		fmt.Println(err)
		time.Sleep(time.Second)
	}

	if errorAndFatal {
		f.errorFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) checkSplitFile(errorAndFatal bool) {
	if f.logSplitType == LogSplitTypeHour {
		f.splitFileHour(errorAndFatal)
	} else {
		f.splitFileSize(errorAndFatal)
	}
}

func (f *FileLogger) writeLogAsync() {
	for logData := range f.logDataChan {
		var file *os.File = f.file
		if logData.ErrorAndFatal {
			file = f.errorFile
		}
		f.checkSplitFile(logData.ErrorAndFatal)
		fmt.Fprintf(file, "%s %s (%s %s %d) %s\n",
			logData.TimeStr, logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)

	}
}

func (f *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}
func (f *FileLogger) Debug(formart string, args ...interface{}) {
	f.writeLog(LogLevelDebug, formart, args...)
}
func (f *FileLogger) Trace(formart string, args ...interface{}) {
	f.writeLog(LogLevelTrace, formart, args...)
}
func (f *FileLogger) Info(formart string, args ...interface{}) {
	f.writeLog(LogLevelInfo, formart, args...)
}
func (f *FileLogger) Warn(formart string, args ...interface{}) {
	f.writeLog(LogLevelWarn, formart, args...)
}
func (f *FileLogger) Error(formart string, args ...interface{}) {
	f.writeLog(LogLevelError, formart, args...)
}
func (f *FileLogger) Fatal(formart string, args ...interface{}) {
	f.writeLog(LogLevelFatal, formart, args...)
}
func (f *FileLogger) Close() {
	f.file.Close()
	f.errorFile.Close()
}

func (f *FileLogger) writeLog(level int, format string, args ...interface{}) {
	if f.level > level {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05.000000000")
	levelString := getLevelString(level)
	filename, funcName, lineNo := utils.GetLineInfo()

	msg := fmt.Sprintf(format, args...)

	logData := &utils.LogData{
		Message:       msg,
		TimeStr:       now,
		LevelStr:      levelString,
		Filename:      filename,
		FuncName:      funcName,
		LineNo:        lineNo,
		ErrorAndFatal: false,
	}

	if level == LogLevelError || level == LogLevelFatal {
		logData.ErrorAndFatal = true
	}
	select {
	case f.logDataChan <- logData:
	default:
	}

	//if level <= LogLevelWarn {
	//	fmt.Fprintf(f.file, "%s %s (%s %s %d) %s\n", now, levelString, filename, funcName, lineNo, msg)
	//} else {
	//	fmt.Fprintf(f.errorFile, "%s %s (%s %s %d) %s\n", now, levelString, filename, funcName, lineNo, msg)
	//}

}
