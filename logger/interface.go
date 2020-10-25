package logger

type LogWriterI interface {
	Type() uint32              //返回打印机类型
	Name() string              //返回打印机名字
	SetTarget(target string)   //设置打印地址
	SetLevel(lvl LevelInfoI)   //设置打印机可打印的日志等级
	CloseLevel(lvl int)        //关闭某个等级的日志
	Write(lvl int, str string) //写
	Flush()                    //刷入
	Close()                    //关闭
	SetValid(isvalid bool)     //设置是否可用
}

type LevelInfoI interface {
	IsValid() bool         //当前等级是否可用
	SetValid(isValid bool) //设置是否可用
	Level() int            //返回等级num
	Head() string          //返回日志头字符串
}

type ConfigI interface {
	Check() bool
	Type() uint32
}
