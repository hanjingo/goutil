package logger

type LevelInfo struct {
	isValid bool                 //是否可用
	lvl     int                  //日志等级
	fun     func(lvl int) string //日志头设置函数
}

func NewLevel(lvl int, headFunc func(lvl int) string) *LevelInfo {
	back := &LevelInfo{
		isValid: true,
		lvl:     lvl,
		fun:     headFunc,
	}
	return back
}

func (li *LevelInfo) SetHeadFunc(f func(lvl int) string) {
	li.fun = f
}

func (li *LevelInfo) Level() int {
	return li.lvl
}
func (li *LevelInfo) Head() string {
	if li.fun == nil {
		return ""
	}
	return li.fun(li.lvl)
}
func (li *LevelInfo) IsValid() bool {
	return li.isValid
}
func (li *LevelInfo) SetValid(isValid bool) {
	li.isValid = isValid
}
