package definitions

type BitConf struct {
	Dbpath string `mapstructure:"dbpath",json:"dbpath"`
	MaxWriteSize int `mapstructure:"max_write_size"`
}