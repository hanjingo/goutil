package file

type SIZE int64

//定义文件单位
const (
	BYTE SIZE = 1
	KB   SIZE = 2 ^ 10
	MB   SIZE = 2 ^ 20
	GB   SIZE = 2 ^ 30
	TB   SIZE = 2 ^ 40
)

func (s SIZE) Add(arg SIZE) SIZE {
	return s + arg
}

func (s SIZE) Del(arg SIZE) SIZE {
	return s - arg
}

func (s SIZE) TB() SIZE {
	return s / TB
}

func (s SIZE) GB() SIZE {
	return s / GB
}

func (s SIZE) MB() SIZE {
	return s / MB
}

func (s SIZE) KB() SIZE {
	return s / KB
}
