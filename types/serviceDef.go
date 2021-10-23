package definitions

type FlakeDef struct {
	Id string `json:"id,omitempty"`
	ContainerId string `json:"container_id"`
	HostId string `json:"host_id"`
	Service string `json:"service"`
	Ip string `json:"ip"`
	Port string `json:"port"`
}
type FlakeStats struct {
	MemUsage string
	CpuPer string
    Network string
}
type HealthDef struct {
	Endpoint string `mapstructure:"endpoint"`
}
type ServiceDeclaration struct {
	Image string `mapstructure:"image"`
	ImageUri string `mapstructure:"image_uri"`
	ContainerPort string `mapstructure:"container_port"`
	BasePort int `mapstructure:"base_port"`
	MaxPort int `mapstructure:"max_port"`
	KongConf Kongdef `mapstructure:"kong_conf"`
	Health HealthDef `mapstructure:"health"`
}
