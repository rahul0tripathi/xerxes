package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
	"github.com/rahultripathidev/docker-utility/kong"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"github.com/rs/xid"
)

const (
	XERXES_HOST = "service.xerxes.discovery"
)

func ScaleUp(serviceId string, nodeId string) error {
	ctx := context.Background()
	serviceDef := GetServiceDef(serviceId)
	runningPods, err := bitcask.GetAllServiceFlakes(serviceId)
	if err != nil && err != bitcask.ErrKeyNotFound{
		return err
	}
	availableNode, availablePort := GetAvailableConfig(serviceDef, runningPods)
	if nodeId == "" {
		nodeId = availableNode
	}
	dockerDaemon, err := dockerClient.NewDockerClient(nodeId)
	if err != nil {
		return err
	}
	if availablePort == "" {
		return errors.New("undefined hostport")
	}
	portset := make(map[nat.Port]struct{})
	port := nat.Port(fmt.Sprintf("%s/tcp", serviceDef.ContainerPort))
	FlakeId := xid.New()
	portset[port] = struct{}{}
	containerConfig := container.Config{
		Image: serviceDef.Image,
		Env:   []string{fmt.Sprintf("INSTANCEID=%s", FlakeId.String())},
	}
	hostConfig := container.HostConfig{
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: availablePort,
				},
			},
		},
		ExtraHosts: []string{fmt.Sprintf("%s:%s", XERXES_HOST, config.XerxesHost.Host)},
	}
	containerResp, err := dockerDaemon.ContainerCreate(ctx, &containerConfig, &hostConfig, nil, FlakeId.String())
	if err != nil {
		fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
		return err
	}
	if err := dockerDaemon.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
		return err
	}
	func() {
		if len(runningPods) == 0 {
			err := kong.CreateService(serviceDef.KongConf)
			if err != nil {
				fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
				return
			}
			err = kong.CreateNewRoute(serviceDef.KongConf.Service)
			if err != nil {
				fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
				return
			}
			err = kong.CreateNewUpstream(serviceDef.KongConf.Upstream)
			if err != nil {
				fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
				return
			}
		}
		err = kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
			Target: config.Nodes.NodeList[nodeId].Ip + ":" + availablePort,
			Weight: 100,
		})
	}()
	err = bitcask.NewFlake(definitions.FlakeDef{
		Id:          FlakeId.String(),
		HostId:      nodeId,
		Ip:          config.Nodes.NodeList[nodeId].Ip,
		Port:        availablePort,
		ContainerId: containerResp.ID,
		Service:     serviceId,
	})
	if err != nil {
		return err
	}
	err = datastore.NotifyChange()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
