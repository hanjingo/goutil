package goutil

import (
	"strings"
	"strconv"
)

//CmpVersion 对比版本; flag:版本分割符 返回: 0:相等 1:v1高于v2 -1:v1低于v2 -2:无法比较
func CmpVersion(v1, v2 string, flag string) int {
	tmp1 := strings.Split(v1, flag)
	tmp2 := strings.Split(v2, flag)
	if len(tmp1) > len(tmp2) {
		for i, v := range tmp2 {
			n2, err := strconv.Atoi(v)
			if err != nil {
				return -2
			}
			n1, err := strconv.Atoi(tmp1[i])
			if err != nil {
				return -2
			}
			if n1 > n2 {
				return 1
			}
			if n2 > n1 {
				return -1
			}
		}
		return 0
	}
	for i, v := range tmp1 {
		n1, err := strconv.Atoi(v)
		if err != nil {
			return -2
		}
		n2, err := strconv.Atoi(tmp2[i])
		if err != nil {
			return -2
		}
		if n1 > n2 {
			return 1
		}
		if n2 > n1 {
			return -1
		}
	}
	return 0
} 