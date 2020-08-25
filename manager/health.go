package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/kong"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	timeout  = 5 * time.Second

)

type (
	standalone struct {
		HostIp      string `json:"host_ip"`
		ContainerId string `json:"containerId"`
		BindingPort uint32 `json:"bindingPort"`
		Status      int    `json:"status"`
		Id          int    `json:"host_id"`
		ServiceId   string `json:"service_id"`
	}
	standaloneServers struct {
		Avaliable []standalone `json:"avaliable"`
	}
)

var (
	AvaliableServers   = standaloneServers{}
	modtime            = time.Time{}
	healthyStatusCodes = []int{201,
		202,
		203,
		204,
		205,
		206,
		207,
		208,
		226,
		300,
		301,
		302,
		303,
		304,
		305,
		306,
		307,
		308,
	}
	LockedInstances = []string{}
	ConnPool        = map[int]*client.Client{}
	Locked = false
	confFile = config.ConfigDir + "/discovery.json"
)

func init() {
	err := config.LoadHosts(config.ConfigDir)
	if err != nil {
		panic(err)
	}
	stat, err := os.Stat(confFile)
	if err != nil {
		panic(err)
	}
	modtime = stat.ModTime()
	discovery, _ := ioutil.ReadFile(confFile)
	err = json.Unmarshal(discovery, &AvaliableServers)
	if err != nil {
		panic(err)
	}
	ConnPool[0], _ = client.NewClient(config.Nodelist.Master.Host, config.Nodelist.Master.Version, nil, nil)
	for _, host := range config.Nodelist.Nodes {
		ConnPool[host.Id], err = client.NewClient(host.Host, host.Version, nil, nil)
		if err != nil {
			delete(ConnPool, host.Id)
			continue
		}
	}
}
func Emergency(server standalone) {
	Locked = true
	defer func() { Locked = false }()
	var DockerClient = &client.Client{}
	var err error
	if _, ok := ConnPool[server.Id]; !ok {
		return
	}
	DockerClient = ConnPool[server.Id]
	err = kong.AddUpstreamTarget(config.Config.Services[server.ServiceId].KongConf, kong.UpstreamTarget{
		Target: server.HostIp + ":" + strconv.FormatUint(uint64(server.BindingPort), 10),
		Weight: 0,
	})
	if err != nil {
		return
	}
	err = DockerClient.ContainerRestart(context.Background(), server.ContainerId, nil)
	if err != nil {
		return
	}
	fmt.Println("Container Restarted")
	err = kong.AddUpstreamTarget(config.Config.Services[server.ServiceId].KongConf, kong.UpstreamTarget{
		Target: server.HostIp + ":" + strconv.FormatUint(uint64(server.BindingPort), 10),
		Weight: 100,
	})
	func() {
		for i, id := range LockedInstances {
			if id == server.ContainerId {
				LockedInstances = append(LockedInstances[:i], LockedInstances[i+1:]...)
				break
			}
		}
	}()
	fmt.Println("Upstream Restarted ", server.Id)
}
func CheckHealth(endpoints standaloneServers) {
	for _, endpoint := range endpoints.Avaliable {
		exists := func() bool {
			for _, id := range LockedInstances {
				if id == endpoint.ContainerId {
					return true
				}
			}
			return false
		}()
		if exists != true {
			resp, err := http.Get("http://" + endpoint.HostIp + ":" + strconv.FormatUint(uint64(endpoint.BindingPort), 10) + "/")
			if err != nil {
				LockedInstances = append(LockedInstances, endpoint.ContainerId)
				go Emergency(endpoint)
				continue
			}
			codeHealthy := sort.SearchInts(healthyStatusCodes, resp.StatusCode)
			if codeHealthy > len(healthyStatusCodes) {
				go Emergency(endpoint)
				LockedInstances = append(LockedInstances, endpoint.ContainerId)
			} else {
			}
		} else {
		}

	}
}
func ReloadAndCheck(timestamp time.Time) {
	fmt.Println("Preforming check LOG: ", timestamp)
	newstat, err := os.Stat(confFile)
	if err != nil {
		panic(err)
	}
	if modtime != newstat.ModTime() {
		fmt.Println("File changed , Hot reloading")
		modtime = newstat.ModTime()
		discovery, _ := ioutil.ReadFile(confFile)
		json.Unmarshal(discovery, &AvaliableServers)
	}
	CheckHealth(AvaliableServers)

}

func HealthScheduler() {
	ticker := time.NewTicker(timeout)
	for timestamp := range ticker.C {
		if !Locked {
			ReloadAndCheck(timestamp)
		}else{
			Locked = false
		}
	}
}

func main() {
	exit := make(chan bool)
	HealthScheduler()
	<-exit
}
