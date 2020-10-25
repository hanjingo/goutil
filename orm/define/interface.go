package define

//数据模型接口
type ModelI interface {
	Do(...interface{}) error  //处理命令
	Create(interface{}) error //创建
	Update(interface{}) error //更新
	Del() error               //删除
	Get(interface{}) error    //查询
	Close() error 			  //关闭模型
}

//数据库接口
type DBI interface {
	Name() string					//数据库名字
	Where(...interface{}) ModelI 	//查询操作
	Close() error 					//关闭操作
}