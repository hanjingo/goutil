package logger

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileManager1 struct {
	maxFileNum int //最大同时打开文件数量
	files      map[string]*FileInfo
}

func NewFileManager1(maxFileNum int) *FileManager1 {
	return &FileManager1{
		maxFileNum: maxFileNum,
		files:      make(map[string]*FileInfo),
	}
}

//打开文件 args:文件容量 dur:文件寿命 s
func (fm *FileManager1) OpenFile(filePathName string, args ...interface{}) (*FileInfo, error) {
	if len(fm.files) >= fm.maxFileNum {
		return nil, errors.New("达到最大文件数量")
	}
	if fi, ok := fm.files[filePathName]; ok {
		return fi, nil
	}
	if filePathName == "" {
		return nil, errors.New("文件路径或文件名不能为空")
	}
	info := &FileInfo{mu: new(sync.Mutex)}
	info.Path, info.Name = filepath.Split(filePathName)
	info.Type = filepath.Ext(filePathName)
	info.CreateTime = time.Now()
	if err := info.CreateFile(); err != nil {
		return nil, err
	}
	wfd, err := os.OpenFile(filePathName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	info.WFd = wfd
	rfd, err := os.OpenFile(filePathName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	info.RFd = rfd
	fm.files[filePathName] = info
	return info, nil
}

func (fm *FileManager1) Write(filePathName string, data []byte) (int, error) {
	if err := fm.checkFile(filePathName); err != nil {
		return 0, err
	}
	w := bufio.NewWriter(fm.files[filePathName].WFd)
	n, err := w.Write(data)
	if err != nil {
		return n, err
	}
	w.Flush()
	return n, nil
}
func (fm *FileManager1) checkFile(filePathName string) error {
	if _, ok := fm.files[filePathName]; !ok {
		if _, err := fm.OpenFile(filePathName); err != nil {
			return err
		}
	}
	return nil
}

func (fm *FileManager1) Read(filePathName string, length int) ([]byte, error) {
	if filePathName == "" {
		return nil, errors.New("文件路径为空")
	}
	if _, ok := fm.files[filePathName]; !ok {
		if _, err := fm.OpenFile(filePathName); err != nil {
			return nil, err
		}
	}
	info := fm.files[filePathName]
	buf := make([]byte, length)
	n, err := info.RFd.Read(buf)
	info.RFd.Read(buf)
	if err != nil {
		if err == io.EOF {
			return buf[:n], io.EOF
		}
		return nil, err
	}
	return buf[:n], nil
}

func (fm *FileManager1) CloseFile(filePathName string) {
	fi, ok := fm.files[filePathName]
	if !ok {
		return
	}
	fi.WFd.Close()
	fi.RFd.Close()
	delete(fm.files, filePathName)
}

func (fm *FileManager1) GetFileInfo(filePathName string) *FileInfo {
	if fi, ok := fm.files[filePathName]; ok {
		fi.Info()
		return fi
	}
	return nil
}

func (fm *FileManager1) CleanFile(filePathName string) error {
	_, ok := fm.files[filePathName]
	if !ok {
		return errors.New("文件不存在或没有注册")
	}
	_, err := os.OpenFile(filePathName, os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	return nil
}
