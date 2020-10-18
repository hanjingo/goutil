package goutil

import (
	"crypto/md5"
	"io"
	"os"
	"errors"
)

//FSIZE 文件大小
type FSIZE int64

//定义文件单位
const (
	BYTE FSIZE = 1
	KB   FSIZE = 2 ^ 10
	MB   FSIZE = 2 ^ 20
	GB   FSIZE = 2 ^ 30
	TB   FSIZE = 2 ^ 40
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
func GetSize(filePathName string) (FSIZE, error) {
	info, err := os.Stat(filePathName)
	if err != nil {
		return FSIZE(0), errors.New("get file stat fail")
	}
	return FSIZE(info.Size()), nil
}
