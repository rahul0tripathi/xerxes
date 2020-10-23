package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/datastore"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
	"github.com/rahultripathidev/docker-utility/kong"
	definitions "github.com/rahultripathidev/docker-utility/types"
)

func ScaleDown(serviceId string, flakeId string, instaceId string) error {
	if flakeId != "" && instaceId == "" {
		flake, err := bitcask.GetFlake(flakeId)
		if err != nil {
			return err
		}
		if flake.Ip != "" {
			err := shutdownPod(flake)
			if err != nil {
				return nil
			}
		} else {
			return errors.New("flake not found")
		}
		return nil
	} else if instaceId != "" && flakeId == "" {
		flakes, err := bitcask.GetAllInstanceFlakes(instaceId)
		if err != nil {
			return err
		}
		for _, flake := range flakes {
			err := shutdownPod(flake)
			if err != nil {
				return err
			}
		}
		return nil
	}
	flakes, err := bitcask.GetAllServiceFlakes(serviceId)
	if err != nil {
		return err
	}
	if len(flakes) > 0 {
		err = shutdownPod(flakes[0])
		if err != nil {
			return err
		}
	}
	return nil
}

func shutdownPod(flake definitions.FlakeDef) error {
	serviceDef := GetServiceDef(flake.Service)
	err := kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
		Target: flake.Ip + ":" + flake.Port,
		Weight: 0,
	})
	if err != nil {
		fmt.Print("Failed To shutdown service  [", flake.Id, "] error ", err)
	}
	Dockerclient, err := dockerClient.NewDockerClient(flake.HostId)
	if err != nil {
		return err
	}
	err = Dockerclient.ContainerRemove(context.Background(), flake.ContainerId, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		fmt.Print("Failed To shutdown service  [", flake.Id, "] error ", err)
		return err
	}
	if err = bitcask.DeleteFLake(flake.Id); err != nil {
		fmt.Print("Failed To Store  service config  [", flake.Id, "] error ", err)
		return err
	}
	err = datastore.NotifyChange()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
