package goutil

import (
	"crypto/md5"
	"io"
	"os"
	"errors"
)

//SIZE 文件大小
type SIZE int64

//定义文件单位
const (
	BYTE SIZE = 1
	KB   SIZE = 2 ^ 10
	MB   SIZE = 2 ^ 20
	GB   SIZE = 2 ^ 30
	TB   SIZE = 2 ^ 40
)

//ComputeMD5 计算md5
func ComputeMD5(filePathName string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePathName)
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

//GetSize 获取文件的实际大小;
func GetSize(filePathName string) (SIZE, error) {
	info, err := os.Stat(filePathName)
	if err != nil {
		return SIZE(0), errors.New("get file stat fail")
	}
	return SIZE(info.Size()), nil
}
