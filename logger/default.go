package logger

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	env "github.com/hanjingo/golib/env"
)

//默认终端打印机
var DefaultConsoleWriterConfig *ConsoleWriterConfig
var DefaultConsoleWriter *ConsoleWriter

//默认文件打印机
var DefaultFileWriterConfig *FileWriterConfig
var DefaultFileWriter *FileWriter

//默认网络打印机
var DefaultNetConfig *NetWriterConfig
var DefaultNetWriter *NetWriter

func initLogger() {
	//终端打印机
	DefaultConsoleWriterConfig = &ConsoleWriterConfig{
		WriterName: DefaultConsoleWriterName,
		CacheCapa:  100,
		OpenLevels: []int{1, 2, 3, 4, 5, 6},
	}
	//文件打印机
	DefaultFileWriterConfig = &FileWriterConfig{
		WriterName: DefaultFileWriterName,
		Capacity:   2,
		ExpireDur:  24,
		CheckDur:   100,
		BufCapa:    40960,
		CacheCapa:  100,
		OpenLevels: []int{1, 2, 3, 4, 5, 6},
		MaxFileNum: 1,
	}
	//网络打印机
	DefaultNetConfig = &NetWriterConfig{
		WriterName: DefaultNetWriterName,
		CacheCapa:  100,
		CheckDur:   100,
		OpenLevels: []int{1, 2, 3, 4, 5, 6},
	}
}

var defaultLog *Logger
var logOnce = new(sync.Once)

func GetDefaultLogger() *Logger {
	logOnce.Do(func() {
		defaultLog = NewLogger()
		initLogger()
		DefaultConsoleWriter = NewConsoleWriter(DefaultConsoleWriterConfig)
		defaultLog.SetWriter(DefaultConsoleWriter)

		DefaultFileWriter = NewFileWriter(DefaultFileWriterConfig)
		DefaultFileWriter.SetNewFileFunc(defaultNewFileFunc)
		defaultLog.SetWriter(DefaultFileWriter)

		DefaultNetWriter = NewNetWriter(DefaultNetConfig)
		defaultLog.SetWriter(DefaultNetWriter)
		DefaultNetWriter.SetValid(false)
	})
	return defaultLog
}

//默认日志头函数
var defaultHeadFunc = func(lvl int) string {
	tag := ""
	switch lvl {
	case FATAL:
		tag = " [FATAL]"
	case ERROR:
		tag = " [ERROR]"
	case WARNING:
		tag = " [WARNING]"
	case NOTICE:
		tag = " [NOTICE]"
	case DEBUG:
		tag = " [DEBUG]"
	case INFO:
		tag = " [INFO]"
	}
	now := time.Now()
	ms := strconv.FormatInt(now.UnixNano()/1e6, 10) //获得当前的毫秒数
	ms_str := ms[len(ms)-3:]                        //取最后三位
	back := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d:%s%s> ", now.Year(), now.Month(),
		now.Day(), now.Hour(), now.Minute(), now.Second(), ms_str, tag)
	return back
}

//默认新建文件函数
var defaultNewFileFunc = func() string {
	now := time.Now()
	fName := fmt.Sprintf("%d-%02d-%02d %02d.%02d.%02d.log", now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	fdatePath := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	fpath := filepath.Join(env.GetCurrPath(), "log", fdatePath)
	return filepath.Join(fpath, fName)
}
