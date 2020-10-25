package mysql

type MsDbConfig struct {
	UserName string
	PassWord string
	Addr     string
	DataBase string
	Charset  string
}

func (conf *MsDbConfig) Name() string {
	return ""
}
