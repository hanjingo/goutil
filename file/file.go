package file

import (
	"crypto/md5"
	"errors"
	"io"
	"os"
	"path/filepath"

	core "github.com/hanjingo/gocore"
)

//计算md5
func MD5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}
	return hash.Sum(result), nil
}

//判断文件是否存在
func Exist(arg string) bool {
	_, err := os.Stat(arg) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//获取文件的实际大小;
func Size(filePathName string) (core.FSIZE, error) {
	if !Exist(filePathName) {
		return core.FSIZE(0), errors.New("文件不存在")
	}
	info, err := os.Stat(filePathName)
	if err != nil {
		return core.FSIZE(0), errors.New("获得文件尺寸失败")
	}
	return core.FSIZE(info.Size()), nil
}

//获得文件名
func FullName(filePathName string) string {
	_, name := filepath.Split(filePathName)
	return name
}

//获得文件名 和 类型
func NameAndExt(filePathName string) (string, string) {
	_, name := filepath.Split(filePathName)
	file_name := filepath.Base(name)
	file_type := filepath.Ext(filePathName)
	return file_name, file_type
}

//创建文件
func Create(filePathName string) (*os.File, error) {
	if !Exist(filePathName) {
		return os.Create(filePathName)
	}
	return nil, errors.New("文件已经存在")
}
