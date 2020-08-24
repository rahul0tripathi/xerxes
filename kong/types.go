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
	next string `json:"next"`
	data []struct{
		Id  string `json:"id"`
		Name string `json:"name"`
		HashOn string `json:"hash_on"`
	}  `mapstring:"data"`
}