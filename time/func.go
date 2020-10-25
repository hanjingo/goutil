package time

import (
	"strconv"
	"time"
)

/*将string转为time (cst时间)*/
func StrToCstTime(timeStr string) time.Time {
	t, _ := time.ParseInLocation(FmtTimeStr, timeStr, time.Local)
	return t
}

/*time转string（标准格式）*/
func TimeToStdStr(t time.Time) string {
	return t.Format(FmtTimeStr)
}

//获得unix毫秒数
func UnixMs(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

//time dur转毫秒
func DurToMs(dur time.Duration) int64 {
	return dur.Nanoseconds() / 1e6
}

//time转string (时间戳)
func TimeToStamp(t time.Time) string {
	return strconv.FormatInt(UnixMs(t), 10)
}

//string(时间戳) 转 time
func StampToTime(TimeStamp string) time.Time {
	ms, _ := strconv.ParseInt(TimeStamp, 10, 64)
	second := ms / 1000
	ns := (ms % 1000) * 1000000
	return time.Unix(second, ns)
}
