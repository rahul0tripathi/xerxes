package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
)

func pullNewImagesInAllNodes(nodes *map[string]*client.Client, targetImage string, tag string) {
	for id, nodeConn := range *(nodes) {
		func(id string) {
			ctx := context.Background()
			err := Pullimage(&ctx, nodeConn, targetImage)
			if err != nil {
				fmt.Print("Failed to pull image ", err, " nodeId: ", id)
			}
		}(id)
	}
}
func Update(serviceId string, nodeId string) {
	runningFlakes, err := bitcask.GetAllServiceFlakes(serviceId)
	if err != nil {
		return
	}
	instancesToupdate := make(map[string]bool)
	connPool := make(map[string]*client.Client)
	if len(runningFlakes) > 0 {
		targetImage := config.ServicesDec.Def[serviceId].ImageUri
		tag := config.ServicesDec.Def[serviceId].Image
		if nodeId == "" {
			for _, flake := range runningFlakes {
				instancesToupdate[flake.HostId] = true
			}
		} else {
			instancesToupdate[nodeId] = true
		}

		for id, _ := range config.Nodes.NodeList {
			if instancesToupdate[id] == true {
				func(id string) {
					newDockerClient, err := dockerClient.NewDockerClient(id)
					if err != nil {
					}
					connPool[id] = newDockerClient
				}(id)
			}
		}
		pullNewImagesInAllNodes(&connPool, targetImage, tag)
		for _, flake := range runningFlakes {
			fmt.Println("Shutting flake ", flake.Id)
			err := shutdownPod(flake)
			if err != nil {
				fmt.Print("Failed to shutdown flake ", flake.Id)
				continue
			}
			fmt.Println("scaling service back ", serviceId)
			err = ScaleUp(serviceId, nodeId)
			if err != nil {
				fmt.Print("Failed to scale up service ", serviceId)
				continue
			}
		}
	}
}
