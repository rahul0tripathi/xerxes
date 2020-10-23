package health

import (
	"fmt"
	"github.com/docker/docker/client"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
	"github.com/rahultripathidev/docker-utility/service_disocvery/discovery"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"time"
)

var (
	Flakes         []definitions.FlakeDef
	DockerConnPool map[string]*client.Client
)

func LoadMetaData() {
	err := config.LoadConfig.Bit()
	if err != nil {
		panic(err)
	}
	bitcask.InitClient()
	if bitcask.BitClient.Locked() {
		<-time.After(5 * time.Second)
	}
	defer func() {
		bitcask.BitClient.Flock.Unlock()
		bitcask.GracefulClose()
	}()
	err = config.LoadConfig.ServiceDef()
	if err != nil {
		panic(err)
	}
	err = config.LoadConfig.Nodes()
	if err != nil {
		panic(err)
	}
	Flakes = []definitions.FlakeDef{}
	for id := range config.ServicesDec.Def {
		func(id string) {
			serviceFlakes, err := bitcask.GetAllServiceFlakes(id)
			if err == nil {
				Flakes = append(Flakes, serviceFlakes...)
			}
		}(id)
	}
	DockerConnPool = make(map[string]*client.Client)
	for id := range config.Nodes.NodeList {
		func(id string) {
			conn, err := dockerClient.NewDockerClient(id)
			if err == nil {
				DockerConnPool[id] = conn
			}
		}(id)
	}
	LoadTargets()
}

func init() {
	config.LoadConfig.ServiceDiscovery()
	LoadMetaData()
}
func LoadTargets() {
	discovery.Targets = make(map[string]struct {
		TargetList []string
		RoundRobin int
	})
	for _, flake := range Flakes {
		func(flake definitions.FlakeDef) {
			target := discovery.Targets["/"+flake.Service]
			target.TargetList = append(target.TargetList, fmt.Sprintf("http://%s:%s", flake.Ip, flake.Port))
			discovery.Targets["/"+flake.Service] = target
		}(flake)
	}
	for key, value := range config.ServiceDiscovery {
		discovery.Targets["/"+key] = struct {
			TargetList []string
			RoundRobin int
		}{TargetList: []string{value}, RoundRobin: 0}
	}
}
