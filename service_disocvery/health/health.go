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
func PodUnhealthy(podDef definitions.ServiceDef) {
	color.Error.Println("["+ time.Now().String()+"] Pod "+podDef.Id+" Did not respond , trying to heal")
	serviceDef := config.ServicesDec.Def[podDef.Service]
	color.Info.Printf("[%s] setting target weight to 0", time.Now().String())
	err := kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
		Target: podDef.Ip + ":" + podDef.Port,
		Weight: 0,
	})
	if err == nil {
		dockerConn := DockerConnPool[podDef.Host]
		logReader , err := dockerConn.ContainerLogs(context.Background(), podDef.ContainerId,types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
		})
		if err == nil {
			func() {
				defer logReader.Close()
				color.Cyan.Println("Container Logs ======================>")
				io.Copy(os.Stdout, logReader)
				stdcopy.StdCopy(os.Stdout,os.Stderr,logReader)
				color.Cyan.Println("<======================")
			}()
		}
		err = dockerConn.ContainerRestart(context.Background(), podDef.ContainerId, nil)
		if err != nil {
			color.Info.Printf("[%s] unable to restart container", time.Now().String())
			service.ScaleDown(podDef.Service, podDef.Id, "")
			return
		} else {
			color.Info.Printf("[%s] container restarted , waiting for 30s for service to boot up\n", time.Now().String())
			<-time.After(30 * time.Second)
			endpoint := buildEndpoint(serviceDef, podDef)
			response, err := http.Get(endpoint)
			if err != nil {
				color.Error.Printf("[%s] pod failed to heal ,shutting it down", time.Now().String())
				service.ScaleDown(podDef.Service, podDef.Id, "")
				return
			}
			codeHealthy := sort.SearchInts(healthyStatusCodes, response.StatusCode)
			if codeHealthy > len(healthyStatusCodes) {
				color.Error.Printf("[%s] pod failed to heal ,shutting it down", time.Now().String())
				service.ScaleDown(podDef.Service, podDef.Id, "")
				return
			}
			color.Success.Printf("[%s] pod restarted , setting target to 100", time.Now().String())
			err = kong.AddUpstreamTarget(serviceDef.KongConf.Upstream, kong.UpstreamTarget{
				Target: podDef.Ip + ":" + podDef.Port,
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
func buildEndpoint(serviceDec definitions.ServiceDeclaration, serviceDef definitions.ServiceDef) string {
	var endpoint string
	if serviceDec.Health.Endpoint == "" {
		endpoint = fmt.Sprintf("http://%s:%s/", serviceDef.Ip, serviceDef.Port)
	} else {
		endpoint = fmt.Sprintf("http://%s:%s%s", serviceDef.Ip, serviceDef.Port, serviceDec.Health.Endpoint)
	}
	return endpoint
}
func CheckHealth() {
	isSchedulerRunning = true
	defer func() { isSchedulerRunning = false }()
	select {
	case <-time.After(1 * time.Millisecond):
		for _, service := range Services {
			func(serviceDef definitions.ServiceDef) {
				serviceDec := config.ServicesDec.Def[service.Service]
				endpoint := buildEndpoint(serviceDec, service)
				response, err := http.Get(endpoint)
				if err != nil {
					PodUnhealthy(service)
					return
				}
				codeHealthy := sort.SearchInts(healthyStatusCodes, response.StatusCode)
				if codeHealthy > len(healthyStatusCodes) {
					PodUnhealthy(service)
				}
			}(service)
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
