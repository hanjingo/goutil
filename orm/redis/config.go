package redis

type RdsConfig struct {
	PoolCapa int `json:"PoolCapa"`
	InitNum int `json:"InitNum"`
	Addr string `json:"Addr"`
	PassWord string `json:"PassWord"`
}

func (conf *RdsConfig) Name() string {
	return ""
}
