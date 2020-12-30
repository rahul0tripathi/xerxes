package definitions

type Kongdef struct {
	Service  KongService  `mapstructure:"service"`
	Upstream KongUpstream `mapstructure:"upstream"`
}

type KongService struct {
	Name       string `mapstructure:"name"`
	Route      string `mapstructure:"route"`
	TaregtPath string `mapstructure:"target_path"`
}

type KongUpstream struct {
	Name   string `mapstructure:"name"`
	Hashon string `mapstructure:"hashon"`
}

type KongConn struct {
	Host string `mapstructure:"host"`
	Admin string `mapstructure:"admin"`
}
type InternalGateway struct {
	Host string `mapstructure:"host"`
	Admin string `mapstructure:"admin_port"`
}