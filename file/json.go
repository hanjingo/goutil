package file

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

/*加载json格式的配置文件*/
func LoadJsonConfig(file_path string, conf interface{}) error {
	if file_path == "" {
		return errors.New("配置文件地址不能为空")
	}
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	return nil
}
