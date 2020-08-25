package kong

type Upstream struct {
	Name string `json:"name"`
	HashOn string `json:"hash_on"` //hash_on should be ip only
}

type UpstreamTarget struct {
	Target string `json:"target"`
	Weight int `json:"weight"`
}

type  Service struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Path string `json:"path"`
	Port int `json:"port"`
}

type Route struct {
	Paths []string `json:"paths"`
	Name string `json:"name"`
}

type UpstreamResp struct {
	Next string `json:"next"`
	Data []struct{
		Id  string `json:"id"`
		Name string `json:"name"`
		HashOn string `json:"hash_on"`
	}  `mapstring:"data"`
}
type ServiceResp struct {
	Next string `json:"next"`
	Data []struct{
		Id  string `json:"id"`
		Host string `json:"host"`
		Protocol string `json:"protocol"`
		Port int `json:"port"`
		Name string `json:"name"`
		Path string `json:"path"`
	}  `mapstring:"data"`
}
type RouteResp struct{
	Next string `json:"next"`
	Data []struct{
		Id  string `json:"id"`
		Paths []string `json:"paths"`
		Name string `json:"name"`
	}  `mapstring:"data"`
}
type TargetResp struct {
	Next string `json:"next"`
	Data []struct {
		Id     string `json:"id"`
		Weight int    `json:"weight"`
		Target string `json:"target"`
	} `mapstring:"data"`
}