package v1

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	file "github.com/hanjingo/golib/file"
)

const FILE_CAPA_UNLIMIT int64 = 0
const FILE_LIFE_UNLIMIT int64 = 0

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
	size, err := file.GetSize(fi.GetAbsPath())
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
	if !file.IsExist(fi.Path) { //先创建路径
		err := os.MkdirAll(fi.Path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if !file.IsExist(fi.GetAbsPath()) { //再创建文件
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

/*创建文件*/
func createFile(arg string) (*os.File, error) {
	return file.Create(arg)
}
