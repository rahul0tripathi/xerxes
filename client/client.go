package client

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/rahultripathidev/docker-utility/config"
)

var DockerClient *client.Client

func init() {
	var err error
	err = config.LoadHosts(config.ConfigDir)
	DockerClient, err = client.NewClient(config.Nodelist.Master.Host, config.Nodelist.Master.Version, nil, nil)
	if err != nil {
		fmt.Println("Unable to Create Client to docker api ", err)
	}
	swarmList, _ := DockerClient.SwarmInspect(context.Background())
	if swarmList.ClusterInfo.ID == "" {
		fmt.Println("NO SWARM FOUND >> Creating a new one")
		_, err := DockerClient.SwarmInit(context.Background(), swarm.InitRequest{
			ListenAddr: "0.0.0.0:2337",
		})
		if err != nil {
			fmt.Println("Unable to Create Swarm  ", err)
		}
	}
}
