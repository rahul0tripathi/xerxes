package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore"
)
func pullNewImagesInAllNodes(nodes *map[string]*client.Client , targetImage string , tag string){
	for id , nodeConn := range *(nodes) {
		func(id string) {
			ctx := context.Background()
			err := Pullimage(&ctx, nodeConn , targetImage)
			if err != nil {
				fmt.Print("Failed to pull image ",err," nodeId: ",id)
			}
		}(id)
	}
}
func Update(serviceId string, nodeId string) {
	runningPods := datastore.GetAllServicesPods(serviceId)
	nodesToupdate := make(map[string]bool)
	connPool := make(map[string]*client.Client)
	if len(runningPods) > 0 {
		targetImage := config.ServicesDec.Def[serviceId].ImageUri
		tag := config.ServicesDec.Def[serviceId].Image
		if nodeId == "" {
			for _, servicePod := range runningPods {
				nodesToupdate[servicePod.Host] = true
			}
		} else {
			nodesToupdate[nodeId] = true
		}

		for id, _ := range config.Nodes.NodeList {
			if nodesToupdate[id] == true {
				func(id string) {
					newDockerClient, err := dockerClient.NewDockerClient(id)
					if err != nil {
					}
					connPool[id] = newDockerClient
				}(id)
			}
		}
		pullNewImagesInAllNodes(&connPool,targetImage,tag)
		for id , service := range runningPods {
			fmt.Println("Shutting Pod ",id)
			err := shutdownPod(serviceId,service.Id)
			if err != nil {
				fmt.Print("Failed to shutdown pod ",service.Id)
				continue
			}
			fmt.Println("scaling Pod ",id)
			err = ScaleUp(serviceId,nodeId)
			if err != nil {
				fmt.Print("Failed to scale up pod ",serviceId)
				continue
			}
		}
	}
}
