package health

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gookit/color"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/kong"
	"github.com/rahultripathidev/docker-utility/service"
	"github.com/rahultripathidev/docker-utility/service_disocvery/discovery"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"io"
	"net/http"
	"os"
	"sort"
	"time"
)

const (
	timeout = 5 * time.Second
)

var (
	healthyStatusCodes = []int{
		201,
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
	isSchedulerRunning = false
	StopScheduler      = make(chan bool)
)

func PodUnhealthy(flake definitions.FlakeDef) {
	color.Error.Println("[" + time.Now().String() + "] Pod " + flake.Id + " Did not respond , trying to heal")
	serviceDef := config.ServicesDec.Def[flake.Service]
	color.Info.Printf("[%s] setting target weight to 0", time.Now().String())
	err := kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
		Target: flake.Ip + ":" + flake.Port,
		Weight: 0,
	})
	if err == nil {
		dockerConn := DockerConnPool[flake.HostId]
		logReader, err := dockerConn.ContainerLogs(context.Background(), flake.ContainerId, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		})
		if err == nil {
			func() {
				defer logReader.Close()
				color.Cyan.Println("Container Logs ======================>")
				io.Copy(os.Stdout, logReader)
				stdcopy.StdCopy(os.Stdout, os.Stderr, logReader)
				color.Cyan.Println("<======================")
			}()
		}
		err = dockerConn.ContainerRestart(context.Background(), flake.ContainerId, nil)
		if err != nil {
			color.Info.Printf("[%s] unable to restart container", time.Now().String())
			service.ScaleDown(flake.Service, flake.Id, "")
			return
		} else {
			color.Info.Printf("[%s] container restarted , waiting for 30s for service to boot up\n", time.Now().String())
			<-time.After(30 * time.Second)
			endpoint := buildEndpoint(serviceDef, flake)
			response, err := http.Get(endpoint)
			if err != nil {
				color.Error.Printf("[%s] pod failed to heal ,shutting it down", time.Now().String())
				service.ScaleDown(flake.Service, flake.Id, "")
				return
			}
			codeHealthy := sort.SearchInts(healthyStatusCodes, response.StatusCode)
			if codeHealthy > len(healthyStatusCodes) {
				color.Error.Printf("[%s] pod failed to heal ,shutting it down", time.Now().String())
				service.ScaleDown(flake.Service, flake.Id, "")
				return
			}
			color.Success.Printf("[%s] pod restarted , setting target to 100", time.Now().String())
			err = kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
				Target: flake.Ip + ":" + flake.Port,
				Weight: 100,
			})
			if err == nil {
				color.Success.Printf("[%s] pod healing Completed", time.Now().String())
			}
		}
	} else {
		color.Error.Printf("[%s] failed to set target weight to 0", time.Now().String())
	}
}
func buildEndpoint(serviceDec definitions.ServiceDeclaration, flake definitions.FlakeDef) string {
	var endpoint string
	if serviceDec.Health.Endpoint == "" {
		endpoint = fmt.Sprintf("http://%s:%s/", flake.Ip, flake.Port)
	} else {
		endpoint = fmt.Sprintf("http://%s:%s%s", flake.Ip, flake.Port, serviceDec.Health.Endpoint)
	}
	return endpoint
}
func CheckHealth() {
	isSchedulerRunning = true
	defer func() { isSchedulerRunning = false }()
	select {
	case <-time.After(1 * time.Millisecond):
		for _, flake := range Flakes {
			func(flake definitions.FlakeDef) {
				serviceDec := config.ServicesDec.Def[flake.Service]
				endpoint := buildEndpoint(serviceDec, flake)
				response, err := http.Get(endpoint)
				if err != nil {
					PodUnhealthy(flake)
					return
				}
				codeHealthy := sort.SearchInts(healthyStatusCodes, response.StatusCode)
				if codeHealthy > len(healthyStatusCodes) {
					PodUnhealthy(flake)
				}
			}(flake)
		}
		return
	case <-StopScheduler:
		return

	}
}
func Scheduler() {
	ticker := time.NewTicker(timeout)
	for range ticker.C {
		if !isSchedulerRunning {
			go CheckHealth()
		}
	}
}

func GetTarget(service string) []byte {
	if _, ok := discovery.Targets[service]; ok {
		targetDef := discovery.Targets[service]
		targetResp := discovery.TargetResp{Host: targetDef.TargetList[targetDef.RoundRobin]}
		if targetDef.RoundRobin+1 > len(targetDef.TargetList)-1 {
			targetDef.RoundRobin = 0
		} else {
			targetDef.RoundRobin += 1
		}
		discovery.Targets[service] = targetDef
		_buf, _ := json.Marshal(targetResp)
		return _buf
	}
	return []byte{}

}
