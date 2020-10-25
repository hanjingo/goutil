package logger

const (
	UNLIMIT             int64 = 0
	UNLIMIT_EXPIRE_TIME int   = 0
)

/*打印机类型*/
const (
	NET     uint32 = 100 //网络打印机
	CONSOLE uint32 = 200 //终端打印机
	FILE    uint32 = 300 //文件打印机
)

//等级
const (
	UNKNOWN int = 0
	FATAL   int = 1
	ERROR   int = 2
	WARNING int = 3
	NOTICE  int = 4
	INFO    int = 5
	DEBUG   int = 6
)

//默认文件打印机名字
const DefaultFileWriterName string = "defaultFileWriter"

//默认网络打印机名字
const DefaultNetWriterName string = "defaultNetWriter"

//默认终端打印机名字
const DefaultConsoleWriterName string = "defaultConsoelWriter"
