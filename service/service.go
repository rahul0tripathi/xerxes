package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/swarm"
	DockerClient "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/kong"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

type (
	standalone struct {
		HostIp      string `json:"host_ip"`
		ContainerId string `json:"containerId"`
		BindingPort uint32 `json:"bindingPort"`
		Status      int    `json:"status"`
		Id int `json:"host_id"`
		ServiceId string `json:"service_id"`
	}
	standaloneServers struct {
		Avaliable []standalone `json:"avaliable"`
	}
)

var (
	nodeRoundRobin = 0
	DiscoverWriteLock = sync.RWMutex{}
	confFile = func() string { HOME , _ := os.UserHomeDir()
		return HOME }() + "/.orchestrator/configuration/discovery.json"
)

func GetActiveServices() ([]swarm.Service, error) {
	ctx := context.Background()
	return client.DockerClient.ServiceList(ctx, types.ServiceListOptions{})
}

func InitServices() {
	alreadyRunning := make(map[string]swarm.Service)
	avaliable, err := GetActiveServices()
	if err != nil {
		return
	}
	for _, existingService := range avaliable {
		alreadyRunning[existingService.Spec.Annotations.Labels["serviceId"]] = existingService
	}
	discovery, _ := ioutil.ReadFile(confFile)
	avaliableServers := standaloneServers{}
	json.Unmarshal(discovery, &avaliableServers)
	for id, service := range config.Config.Services {
		switch service.Type {
		case "swarm":
			if _, ok := alreadyRunning[id]; !ok {
				registerService(service, id,config.Nodelist.Master,true)
			} else {
				fmt.Println("Swarm Service Already Exists")
			}
		case "standalone":
			running := false
			for _, conf := range avaliableServers.Avaliable {
				for _ , ports := range service.HostPortRange{
					if conf.BindingPort == ports {
						running = true
						break
					}
				}
			}
			if !running {
				func() {
					if len(config.Nodelist.Nodes) > 0 {
						if nodeRoundRobin == 0 {
							registerService(service, id, config.Nodelist.Master, false)
							nodeRoundRobin += 1
						} else {
							node_ip := config.Nodelist.Nodes[nodeRoundRobin-1]
							nodeRoundRobin = func() int {
								if nodeRoundRobin  < len(config.Nodelist.Nodes) {
									return nodeRoundRobin + 1
								} else {
									return 0
								}
							}()
							registerService(service, id, node_ip,false)
						}
					}else{
						registerService(service, id,config.Nodelist.Master,false)
					}
				}()
			} else {
				fmt.Println("Standalone Service Already Exists ")
			}
		}

	}
}
func registerService(service config.Services, id string, node config.Daemon ,init bool) {
	ctx := context.Background()
	switch service.Type {
	case "swarm":
		func() {
			servicetask := swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Name: id,
					Labels: map[string]string{
						"serviceId": id,
					},
				},
				TaskTemplate: swarm.TaskSpec{

					ContainerSpec: swarm.ContainerSpec{

						Image: service.Image,
					},
					RestartPolicy: &swarm.RestartPolicy{
						Condition: "any",
					}},
				Mode: swarm.ServiceMode{
					Replicated: &swarm.ReplicatedService{
						Replicas: &service.Init,
					},
				},
				EndpointSpec: &swarm.EndpointSpec{
					Ports: []swarm.PortConfig{
						{
							Name:          "PORT",
							Protocol:      "tcp",
							TargetPort:    service.Containerport,
							PublishedPort: service.Hostport,
						},
					},
				},
			}
			_, err := client.DockerClient.ServiceCreate(ctx, servicetask, types.ServiceCreateOptions{})
			if err != nil {
				fmt.Println("Create service error ", err)
				return
			} else {
				err := kong.CreateService(id,service,config.Nodelist.Master.Ip,int(service.Hostport))
				if err != nil{
					fmt.Println(err)
					return
				}
				err = kong.CreateNewRoute(id,service)
				if err != nil{
					fmt.Println(err)
					return
				}
			}
		}()
	case "standalone":
		func() {
			var dockerClient *DockerClient.Client
			fmt.Printf("%+v",node)
			if node.Ip != config.Nodelist.Master.Ip {
				dockerClient, _ = DockerClient.NewClient(node.Host, node.Version, nil, nil)
			} else {
				dockerClient = client.DockerClient
			}
			discovery, _ := ioutil.ReadFile(confFile)
			avaliableServers := standaloneServers{}
			json.Unmarshal(discovery, &avaliableServers)
			var Hostport uint32
			if len(avaliableServers.Avaliable) > 0 {
				for _, freePort := range service.HostPortRange {
					free := true
					for _, conf := range avaliableServers.Avaliable {
						if freePort == conf.BindingPort {
							free = false
							break
						}
					}
					if free {
						Hostport = freePort
					}
				}
			} else {
				Hostport = service.HostPortRange[0]
			}
			if Hostport == 0 {
				return
			}
			portset := make(map[nat.Port]struct{})
			port := nat.Port(strconv.FormatUint(uint64(service.Containerport), 10) + "/tcp")
			portset[port] = struct{}{}
			containerConfig := container.Config{
				Image: service.Image,
			}
			HostIp := node.Ip
			hostConfig := container.HostConfig{
				PortBindings: nat.PortMap{
					port: []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: strconv.FormatUint(uint64(Hostport), 10),
						},
					},
				},
			}
			containerResp, err := dockerClient.ContainerCreate(ctx, &containerConfig, &hostConfig, nil, id+"_"+strconv.FormatUint(uint64(Hostport), 10))
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := dockerClient.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); err != nil {
				fmt.Println(err)
				return
			}
			avaliableServers.Avaliable = append(avaliableServers.Avaliable, standalone{ContainerId: containerResp.ID, BindingPort: Hostport , HostIp: HostIp , Status: 1,Id: node.Id,ServiceId: id})
			str, _ := json.Marshal(avaliableServers)
			DiscoverWriteLock.Lock()
			err = ioutil.WriteFile(confFile, str, os.ModePerm)
			DiscoverWriteLock.Unlock()
			func() {
				if init == true {
					err := kong.CreateService(id, service, "",0)
					if err != nil {
						fmt.Println(err)
						return
					}
					err =  kong.CreateNewRoute(id, service)
					if err != nil {
						fmt.Println(err)
						return
					}
					err =  kong.CreateNewUpstream(service.KongConf)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				err = kong.AddUpstreamTarget(service.KongConf,kong.UpstreamTarget{
					Target: HostIp+":"+strconv.FormatUint(uint64(Hostport), 10),
					Weight: 100,
				})
				if err != nil {
					fmt.Println(err)
					return
				}
			}()
			fmt.Println("Service scaled by 1")
		}()
	}
}
func Scale(id string, factor uint64 , useNode bool) {
	ctx := context.Background()
	avaliable, err := GetActiveServices()
	if err != nil {
		return
	}
	service := swarm.Service{}
	for _, existingService := range avaliable {
		if existingService.Spec.Annotations.Labels["serviceId"] == id {
			service = existingService
		}
	}
	if service.ID != "" {
		service, _, err := client.DockerClient.ServiceInspectWithRaw(ctx, service.ID)
		if err != nil {
			return
		}

		serviceMode := &service.Spec.Mode
		if serviceMode.Replicated != nil {
			serviceMode.Replicated.Replicas = &factor
		} else {
			fmt.Println(errors.Errorf("scale can only be used with replicated mode"))
			return
		}

		response, err := client.DockerClient.ServiceUpdate(ctx, service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
		if err != nil {
			return
		}

		for _, warning := range response.Warnings {
			fmt.Println(warning)
		}
		fmt.Println("Scaled to ", factor)
	} else {
		if config.Config.Services[id].Type == "standalone" {
			discovery, _ := ioutil.ReadFile(confFile)
			avaliableServers := standaloneServers{}
			json.Unmarshal(discovery, &avaliableServers)
			currentCount := 0
			for _, server := range avaliableServers.Avaliable {
				if server.ServiceId == id {
					currentCount ++
				}
			}
			fmt.Println(currentCount)
			if currentCount > int(factor) {
				err = Shutdown(id ,"standalone", currentCount-int(factor) )
			}else{
				for i := 0; i < int(factor)-currentCount; i++ {
					func() {
						fmt.Println(nodeRoundRobin,useNode)
						if len(config.Nodelist.Nodes) > 0 {
							if nodeRoundRobin == 0 && useNode == false{
								registerService(config.Config.Services[id], id, config.Nodelist.Master, false)
								nodeRoundRobin += 1
							} else if useNode == true || nodeRoundRobin >0 {
								node_ip := config.Daemon{}
								if useNode == true && nodeRoundRobin-1 < 0{
									node_ip = config.Nodelist.Nodes[nodeRoundRobin]
								}else{
									node_ip = config.Nodelist.Nodes[nodeRoundRobin-1]
								}
								useNode = false
								registerService(config.Config.Services[id], id, node_ip,false)
								nodeRoundRobin = func() int {
									if nodeRoundRobin  < len(config.Nodelist.Nodes) {
										return nodeRoundRobin + 1
									} else {
										return 0
									}
								}()
							}
						}else{
							registerService(config.Config.Services[id], id,config.Nodelist.Master,false)
						}
					}()

				}
			}

		}
	}
}
func Shutdown(service string , serviceType string , containers int) error {
	switch serviceType {
	case "swarm":return func() error {
	return nil
	}()
	case "standalone":return func() error {
		fmt.Println("Shutting down ", service)
		DiscoverWriteLock.Lock()
		defer DiscoverWriteLock.Unlock()
		discovery, _ := ioutil.ReadFile(confFile)
		avaliableServers := standaloneServers{}
		NewavaliableServer := standaloneServers{
			Avaliable: []standalone{},
		}
		servicesCount := 0
		json.Unmarshal(discovery, &avaliableServers)
		for _, server := range avaliableServers.Avaliable {
			if containers >0 {
				if servicesCount == containers {
					break
				}
			}
			if server.ServiceId == service {
				err := kong.AddUpstreamTarget(config.Config.Services[server.ServiceId].KongConf, kong.UpstreamTarget{
					Target: server.HostIp + ":" + strconv.FormatUint(uint64(server.BindingPort), 10),
					Weight: 0,
				})
				if err != nil {
					fmt.Println(err)
				}
				var Dockerclient = &DockerClient.Client{}
				if server.Id == 0 {
					host := config.Nodelist.Master
					Dockerclient, err = DockerClient.NewClient(host.Host, host.Version, nil, nil)
				} else {
					for _, host := range config.Nodelist.Nodes {
						if host.Id == server.Id {
							Dockerclient, err = DockerClient.NewClient(host.Host, host.Version, nil, nil)
							break
						}
					}
				}
				err = Dockerclient.ContainerRemove(context.Background(), server.ContainerId, types.ContainerRemoveOptions{
					Force: true,
				})
				if err != nil {
					fmt.Println(err)
				}
				servicesCount ++
			} else {
				NewavaliableServer.Avaliable = append(NewavaliableServer.Avaliable, server)
			}
		}
		str, _ := json.Marshal(NewavaliableServer)
		err := ioutil.WriteFile(confFile, str, os.ModePerm)
		if err != nil {
			return err
		}
		return nil }()
	default: return errors.New("invalid type "+serviceType)
	}
}
