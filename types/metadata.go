package definitions

import (
	"math/rand"
	"time"
)

type NodeMeta struct {
	Id              string   `json:"id,omitempty"`
	TotalRegistered int      `json:"total_registered"`
	Healthy         int      `json:"healthy"`
	Unhealthy       int      `json:"unhealthy"`
	UnhealthyPods   []string `json:"unhealthy_pods"`
}
type ServiceMeta struct {
	BasePort  int        `json:"base_port"`
	MaxPort   int        `json:"max_port"`
	UsedPorts []UsedPort `json:"used_ports"`
}

type UsedPort struct {
	Port string `json:"port"`
	Node string `json:"node"`
}

func GenRandPort(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

