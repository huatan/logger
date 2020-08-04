package logger

import "runtime"

type LogData struct {
	Message  string
	TimeStr  string
	LevelStr string
	Filename string
	FuncName string
	LineNo   int
	ErrorAndFatal bool
}

func GetLineInfo() (filename string, funcName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		filename = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}
	return
}
