package discovery

type (
	TargetResp struct {
		Host string `json:"host"`
	}
)
var (
	Targets  map[string]struct {
		TargetList []string
		RoundRobin int
	}
)

