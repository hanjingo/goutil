package network

import (
	"mime/multipart"
	"net/http"
)

//跨域发送数据
func WriteHttp(w http.ResponseWriter, data []byte) (int, error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return w.Write(data)
}

//返回
func SendBack(w http.ResponseWriter, content string, code int) {
	http.Error(w, content, code)
}

//解析表单数值
func GetMultPartFormValue(form *multipart.Form, key string) string {
	if form == nil || form.Value == nil {
		return ""
	}
	for k, v := range form.Value {
		if k == key && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

//解析表单文件
func GetMultPartFormData(form *multipart.Form, key string) []byte {
	if form == nil || form.File == nil {
		return nil
	}
	for k, v := range form.File {
		if k == key && len(v) > 0 {
			head := v[0]
			file, err := head.Open()
			if err != nil {
				return nil
			}
			back := make([]byte, head.Size)
			n, err := file.Read(back)
			if err != nil || n == 0 {
				return nil
			}
			return back
		}
	}
	return nil
}
