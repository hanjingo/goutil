package logger

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	file "github.com/hanjingo/golib/file"
)

//文件信息
type FileInfo struct {
	mu         *sync.Mutex
	Path       string    //路径名称
	Name       string    //文件名称(文件名+类型："007.txt")
	Type       string    //类型(".txt")
	Size       file.SIZE //实际大小(单位：MB)(不准确)
	CreateTime time.Time //文件创建日期(不准确)
	WFd        *os.File  //写文件描述符
	RFd        *os.File  //读文件描述符
}

//获得文件的绝对路径
func (fi *FileInfo) GetAbsPath() string {
	if fi.Path != "" && fi.Name != "" {
		return filepath.Join(fi.Path, fi.Name)
	}
	return ""
}

func (fi *FileInfo) Info() error {
	size, err := getSize(fi.GetAbsPath())
	if err != nil {
		return err
	}
	fi.Size = size
	return nil
}

func (fi *FileInfo) IsFull(sz file.SIZE) bool {
	fi.Info()
	return fi.Size >= sz
}

func (fi *FileInfo) IsExpired(dur time.Duration) bool {
	fi.Info()
	expireTime := fi.CreateTime.Add(dur)
	return time.Now().After(expireTime)
}

func (fi *FileInfo) CreateFile() error {
	if !isExist(fi.Path) { //先创建路径
		err := os.MkdirAll(fi.Path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if !isExist(fi.GetAbsPath()) { //再创建文件
		fd, err := createFile(fi.GetAbsPath())
		if err != nil {
			return err
		}
		if err := fd.Close(); err != nil {
			return err
		}
	}
	return nil
}

//判断文件/路径是否存在
func isExist(arg string) bool {
	_, err := os.Stat(arg) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//获取文件的实际大小; filePathName:"c:\\007.log"; 返回值单位：MB
func getSize(filePathName string) (file.SIZE, error) {
	if !isExist(filePathName) {
		return 0, errors.New("文件不存在")
	}
	info, err := os.Stat(filePathName)
	if err != nil {
		return 0, errors.New("获得文件尺寸失败")
	}
	return file.SIZE(info.Size()).MB(), nil
}

/*创建文件*/
func createFile(arg string) (*os.File, error) {
	if !isExist(arg) { //如果文件不存在
		fd, err := os.Create(arg)
		if err != nil {
			return nil, err
		}
		return fd, nil
	}
	return nil, errors.New("文件已经存在")
}
