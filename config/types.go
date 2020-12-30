package config

type Services struct {
	Type string `mapstructure:"type"`
	Image string `mapstructure:"image"`
	Containerport uint32 `mapstructure:"containerPort"`
	Hostport uint32 `mapstructure:"hostPort"`
	Dev Dev `mapstructure:"dev"`
	HostPortRange []uint32 `mapstructure:"hostportRange"`
	ImageUri string `mapstring:"imageUri"`
	Init uint64 `mapstring:"init"`
	KongConf KongDesc `mapstructure:"kongConfig"`
}

type config struct {
	Services map[string]Services `mapstructure:"services"`
}

type Kong struct {
	Host string `mapstructure:"host"`
	Admin string `mapstructure:"admin"`
}

type KongDesc struct {
	Upstream bool `mapstructure:"upstream"`
	Route string `mapstructure:"route"`
	UpstreamService string `mapstructure:"upstreamService"`
	Service string `mapstructure:"service"`
	ServicePath string `mapstructure:"service_path"`
}
type NodeList struct{
	Master Daemon `mapstructure:"master"`
	Nodes  []Daemon `mapst ructure:"nodes"`
}

type Daemon struct {
	Host string `mapstructure:"host"`
	Version string `mapstructure:"version"`
	Ip string `mapstructure:"ip"`
	Id int `mapstrucutr:"Id"`
}
type Dev struct {
	Manager string `mapstructure:"manager"`
	Cmd DevCmd `mapstructure:"cmd"`
	Port string `mapstructure:"port"`
}

type DevCmd struct {
	Start string `mapstructure:"start"`
	Stop string `mapstructure:"stop"`
	Reload string `mapstructure:"reload"`
}
