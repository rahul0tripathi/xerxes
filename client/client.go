package dockerClient

import (
	"errors"
	"github.com/docker/docker/client"
	"github.com/rahultripathidev/docker-utility/config"
)

func init() {
	err := config.LoadConfig.Nodes()
	if err != nil {
		panic(err)
	}
}

func NewDockerClient(nodeId string) (*client.Client, error) {
	if _, ok := config.Nodes.NodeList[nodeId]; ok {
		dockerClient, err := client.NewClient(config.Nodes.NodeList[nodeId].Host, config.Nodes.NodeList[nodeId].Version, nil, nil)
		if err != nil {
			return nil, err
		}
		return dockerClient, nil
	} else {
		return nil, errors.New("unable to find node")
	}
}
