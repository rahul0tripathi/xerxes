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
	if err != nil && err != bitcask.ErrKeyNotFound {
		return err
	}
	var availableNode, availablePort string
	if nodeId != "" {
		availablePort = fmt.Sprintf("%d", definitions.GenRandPort(serviceDef.BasePort, serviceDef.MaxPort))
	} else {
		availableNode, availablePort = GetAvailableConfig(serviceDef, runningPods)
		nodeId = availableNode
	}
	dockerDaemon, err := dockerClient.NewDockerClient(nodeId)
	defer dockerDaemon.Close()
	if err != nil {
		return err
	}
	if availablePort == "" {
		return errors.New("undefined hostport")
	}
	portset := make(map[nat.Port]struct{})
	port := nat.Port(fmt.Sprintf("%s/tcp", serviceDef.ContainerPort))
	FlakeId := xid.New()
	fmt.Println("Service Name:",serviceId,"Host Name:",nodeId,"Service Id:",FlakeId.String())
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
		RestartPolicy: container.RestartPolicy{
			Name:              "on-failure",
			MaximumRetryCount: 1,
		},
		LogConfig: container.LogConfig{
			Type: "local",
			Config: map[string]string{
				"max-size": "1m",
				"max-file": "1",
				"compress":"false",
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
			err = kong.CreateNewUpstream(serviceDef.KongConf.Upstream, serviceId)
			if err != nil {
				fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
				return
			}
			err = kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
				Target: config.Nodes.NodeList[nodeId].Ip + ":" + availablePort,
				Weight: 100,
			}, serviceId)
			err := kong.CreateService(serviceDef.KongConf, serviceId)
			if err != nil {
				fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
				return
			}
			err = kong.CreateNewRoute(serviceDef.KongConf.Service, serviceId)
			if err != nil {
				fmt.Print("Failed To Create service  [", serviceId, "] error ", err)
				return
			}
		} else {
			err = kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
				Target: config.Nodes.NodeList[nodeId].Ip + ":" + availablePort,
				Weight: 100,
			}, serviceId)
		}
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
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
