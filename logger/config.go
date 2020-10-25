package logger

//新建配置
func NewConfig() *LogConfig {
	return &LogConfig{
		Console: &ConsoleWriterConfig{},
		File:    &FileWriterConfig{},
		Net:     &NetWriterConfig{},
	}
}

/*日志配置*/
type LogConfig struct {
	DoPrintConsole bool                 `json:"DoPrintConsole"` //是否打印终端
	DoPrintFile    bool                 `json:"DoPrintFile"`    //是否打印文件
	DoPrintNet     bool                 `json:"DoPrintNet"`     //是否打印网路
	Console        *ConsoleWriterConfig `json:"Console"`        //终端打印机配置
	File           *FileWriterConfig    `json:"File"`           //文件打印机配置
	Net            *NetWriterConfig     `json:"Net"`            //网络打印机配置
}

/*这些参数从配置文件中读取
符号解释：
↑：值越大，性能越好
↓：值越大，性能越差
=：对性能没影响
？：对性能有影响，但不是线性关系，即无法确定影响的好坏*/
type ConsoleWriterConfig struct {
	WriterName string `json:"WriterName"` //名字
	CacheCapa  int    `json:"CacheCapa"`  //日志缓冲队列容量 ↑
	OpenLevels []int  `json:"OpenLevels"` //要打开的等级
}

func (c *ConsoleWriterConfig) Check() bool {
	if c.CacheCapa <= 0 {
		return false
	}
	return true
}

func (c *ConsoleWriterConfig) Type() uint32 {
	return CONSOLE
}

func (c *ConsoleWriterConfig) Name() string {
	return c.WriterName
}

/*这些参数从配置文件中读取
符号解释：
↑：值越大，性能越好
↓：值越大，性能越差
=：对性能没影响
？：对性能有影响，但不是线性关系，即无法确定影响的好坏*/
type FileWriterConfig struct {
	WriterName string `json:"WriterName"` //打印机名字
	Capacity   int64  `json:"Capacity"`   //文件容量(单位: MB) ↑
	ExpireDur  int    `json:"ExpireDur"`  //文件生命周期(单位: Hour) ↑
	DelDur     int    `json:"DelDur"`     //文件删除周期(单位: Day)
	CheckDur   int    `json:"CheckDur"`   //写入多少次检查一次打印机
	BufCapa    int64  `json:"BufCapa"`    //文件缓冲区容量 ↑
	CacheCapa  int    `json:"CacheCapa"`  //日志缓冲队列容量 ↑
	OpenLevels []int  `json:"OpenLevels"` //要打开的等级
	MaxFileNum int    `json:"MaxFileNum"` //同时打开最大文件数量
}

func (c *FileWriterConfig) Check() bool {
	if c.Capacity <= 0 {
		return false
	}
	if c.ExpireDur <= 0 {
		return false
	}
	if c.CheckDur <= 0 {
		return false
	}
	if c.DelDur < 0 {
		return false
	}
	if c.BufCapa <= 0 {
		return false
	}
	if c.CacheCapa <= 0 {
		return false
	}
	if c.MaxFileNum <= 0 {
		return false
	}
	return true
}

func (c *FileWriterConfig) Type() uint32 {
	return FILE
}

func (c *FileWriterConfig) Name() string {
	return c.WriterName
}

/*这些参数从配置文件中读取
符号解释：
↑：值越大，性能越好
↓：值越大，性能越差
=：对性能没影响
？：对性能有影响，但不是线性关系，即无法确定影响的好坏*/
type NetWriterConfig struct {
	WriterName string `json:"WriterName"` //名字
	CacheCapa  int    `json:"CacheCapa"`  //日志缓冲队列容量 ↑
	CheckDur   int    `json:"CheckDur"`   //检查周期(写入多少次检查一次打印机)
	OpenLevels []int  `json:"OpenLevels"` //要打开的等级
}

func (c *NetWriterConfig) Check() bool {
	if c.CacheCapa <= 0 {
		return false
	}
	if c.CheckDur <= 0 {
		return false
	}
	return true
}

func (c *NetWriterConfig) Type() uint32 {
	return NET
}
