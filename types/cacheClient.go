package definitions

type CacheClient struct {
	Addr string `json:"addr"`
	Password string `json:"password"`
	DB int `json:"db"`
}
