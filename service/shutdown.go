package service

import (
	"errors"
	"fmt"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
	"github.com/rahultripathidev/docker-utility/kong"
	"github.com/rahultripathidev/docker-utility/scheduler_client"
	definitions "github.com/rahultripathidev/docker-utility/types"
)

func ScaleDown(serviceId string, flakeId string, instaceId string, timeout int64) error {
	if flakeId != "" {
		flake, err := bitcask.GetFlake(flakeId)
		if err != nil {
			return err
		}
		if flake.Ip != "" {
			err := shutdownPod(flake, timeout)
			if err != nil {
				return nil
			}
		} else {
			return errors.New("flake not found")
		}
		return nil
	} else if instaceId != "" {
		flakes, err := bitcask.GetAllInstanceFlakes(instaceId)
		if err != nil {
			return err
		}
		for _, flake := range flakes {
			if flake.HostId == instaceId {
				err := shutdownPod(flake, timeout)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}
	flakes, err := bitcask.GetAllServiceFlakes(serviceId)
	if err != nil {
		return err
	}
	if len(flakes) > 0 {
		err = shutdownPod(flakes[0], timeout)
		if err != nil {
			return err
		}
	}
	return nil
}

func shutdownPod(flake definitions.FlakeDef, timeout int64) error {
	serviceDef := GetServiceDef(flake.Service)
	err := kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
		Target: flake.Ip + ":" + flake.Port,
		Weight: 0,
	}, flake.Service)
	if err != nil {
		fmt.Print("Failed To shutdown service  [", flake.Id, "] error ", err)
	}
	if err = bitcask.DeleteFLake(flake.Id); err != nil {
		fmt.Print("Failed To Store  service config  [", flake.Id, "] error ", err)
		return err
	}
	if err != nil {
		fmt.Println(err)
	}
	scheduler_client.ScheduleContainerForDeletion(flake.ContainerId, flake.HostId, true, timeout)
	return nil
}
