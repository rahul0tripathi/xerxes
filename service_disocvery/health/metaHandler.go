package health

import (
	"fmt"
	"github.com/docker/docker/client"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore"
	"github.com/rahultripathidev/docker-utility/service_disocvery/discovery"
	definitions "github.com/rahultripathidev/docker-utility/types"
)

var (
	Services       []definitions.ServiceDef
	DockerConnPool map[string]*client.Client
)

func LoadMetaData() {
	err := config.LoadConfig.ServiceDef()
	if err != nil {
		panic(err)
	}
	err = config.LoadConfig.Cache()
	if err != nil {
		panic(err)
	}
	err = config.LoadConfig.Nodes()
	if err != nil {
		panic(err)
	}
	Services = []definitions.ServiceDef{}
	for id := range config.ServicesDec.Def {
		func(id string) {
			services := datastore.GetAllServicesPods(id)
			Services = append(Services, services...)
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
	InitSub()
}
func LoadTargets() {
	discovery.Targets = make(map[string]struct {
		TargetList []string
		RoundRobin int
	})
	for _, service := range Services {
		func(service definitions.ServiceDef) {
			target := discovery.Targets["/"+service.Service]
			target.TargetList = append(target.TargetList, fmt.Sprintf("http://%s:%s", service.Ip, service.Port))
			discovery.Targets["/"+service.Service] = target
		}(service)
	}
	for key, value := range config.ServiceDiscovery {
		discovery.Targets["/"+key] = struct {
			TargetList []string
			RoundRobin int
		}{TargetList: []string{value}, RoundRobin: 0}
	}
}
