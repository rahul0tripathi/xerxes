package definitions

type Nodes struct {
	NodeList map[string]Daemon `mapstructure:"nodes"`
}
type Daemon struct {
	Host    string `mapstructure:"host"`
	Version string `mapstructure:"version"`
	Ip      string `mapstructure:"ip"`
}
