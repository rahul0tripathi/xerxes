package health
type Target struct {
	Id     string `json:"id"`
	Target string `json:"target"`
}
type TargetWithRoundRobin struct {
	RoundRobin int
	TargetHosts []Target
}