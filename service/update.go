package service

import (
	"fmt"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
)

func Update(serviceId string, nodeId string) {
	runningFlakes, err := bitcask.GetAllServiceFlakes(serviceId)
	if err != nil {
		return
	}
		for _, flake := range runningFlakes {
			fmt.Println("Shutting flake ", flake.Id)
			err := shutdownPod(flake,0)
			if err != nil {
				fmt.Print("Failed to shutdown flake ", flake.Id)
				continue
			}
			fmt.Println("scaling service back ", serviceId)
			err = ScaleUp(flake.Service, flake.HostId)
			if err != nil {
				fmt.Print("Failed to scale up service ", serviceId)
				continue
			}
		}
}
