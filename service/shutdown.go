package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/datastore"
	"github.com/rahultripathidev/docker-utility/kong"
)

func ScaleDown(serviceId string, podId string, nodeId string) error {
	if podId != "" && nodeId == "" {
		podDef := datastore.GetpodById(serviceId, podId)
		if podDef.Ip != "" {
			err := shutdownPod(serviceId, podDef.Id)
			if err != nil {
				return nil
			}
		} else {
			return errors.New("Pod not found")
		}
		return nil
	} else if nodeId != "" && podId == "" {
		pods := datastore.GetAllpodsInNode(nodeId)
		for _, pod := range pods {
			err := shutdownPod(serviceId, pod)
			if err != nil {
				return err
			}
		}
		return nil
	}
	pods := datastore.GetAllServicesPods(serviceId)
	err :=  shutdownPod(serviceId, pods[0].Id)
	return err
}

func shutdownPod(serviceId string, podId string) error {
	podDef := datastore.GetpodById(serviceId, podId)
	serviceDef := GetServiceDef(serviceId)
	err := kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
		Target: podDef.Ip + ":" + podDef.Port,
		Weight: 0,
	})
	if err != nil {
		fmt.Print("Failed To shutdown service  [", serviceId, "] error ", err)
	}
	Dockerclient, err := dockerClient.NewDockerClient(podDef.Host)
	if err != nil {
		return err
	}
	err = Dockerclient.ContainerRemove(context.Background(), podDef.ContainerId, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		fmt.Print("Failed To shutdown service  [", serviceId, "] error ", err)
		return err
	}
	if err = datastore.DeregisterPod(serviceId, podDef.Id, podDef.Host); err != nil {
		fmt.Print("Failed To shutdown service  [", serviceId, "] error ", err)
		return err
	}
	key := datastore.XerxesEvents["reloadAll"]()
	datastore.NotifyChange(key, "")
	return nil
}
