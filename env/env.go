package env

import (
	"os"
	"path/filepath"
	"runtime"
)

/*获得当前系统类型*/
func GetOsType() string {
	return runtime.GOOS
}

/*获得当前程序的运行路径*/
func GetCurrPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

//获得当前runtime stack
func GetRuntimeStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
