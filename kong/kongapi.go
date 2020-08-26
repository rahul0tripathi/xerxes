package kong

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rahultripathidev/docker-utility/config"
	"io/ioutil"
	"net/http"
)

var (
	basePath string
)

func init() {
	config.LoadServices(config.ConfigDir)
	err := config.LoadKongConf(config.ConfigDir)
	if err != nil {
		panic(err)
	}
	basePath = config.KongConf.Host + ":" + config.KongConf.Admin
}
func CreateNewUpstream(service config.KongDesc) error {
	fmt.Print("\nCreating a new UpStream [", service.UpstreamService, "] .... ")
	_body, _ := json.Marshal(Upstream{
		service.UpstreamService,
		"ip",
	})
	res, err := http.Post(basePath+"/upstreams", "application/json", bytes.NewBuffer(_body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode == 200 || res.StatusCode == 409 {
		fmt.Println("OK")
		return nil
	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func AddUpstreamTarget(service config.KongDesc, target UpstreamTarget) error {
	fmt.Print("\nAdding Upstream Target [", target.Target, "] to [", service.UpstreamService, "] .... ")
	_body, err := json.Marshal(target)
	if err != nil {
		return err
	}
	res, err := http.Post(basePath+"/upstreams/"+service.UpstreamService+"/targets", "application/json", bytes.NewBuffer(_body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode <= 400 {
		fmt.Println("OK")
		return nil
	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func CreateService(serviceId string, service config.Services, serviceEndpointHost string, serviceEndpointPort int) error {
	fmt.Print("\nCreating a new Service [", serviceId, "] .... ")
	var _body []byte
	var err error
	if service.KongConf.Upstream == true {
		_body, err = json.Marshal(Service{
			Name: serviceId,
			Host: service.KongConf.UpstreamService,
			Path: service.KongConf.ServicePath,
			Port: 80,
		})
		if err != nil {
			return err
		}
	} else {
		_body, err = json.Marshal(Service{
			Name: serviceId,
			Host: serviceEndpointHost,
			Path: service.KongConf.ServicePath,
			Port: serviceEndpointPort,
		})
		if err != nil {
			return err
		}
	}
	res, err := http.Post(basePath+"/services/", "application/json", bytes.NewBuffer(_body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode <= 400 || res.StatusCode == 409 {
		fmt.Println("OK")
		return nil
	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func CreateNewRoute(serviceId string, service config.Services) error {
	fmt.Print("\nCreating a new Route [", service.KongConf.Route, "] [", serviceId, "] .... ")
	_body, err := json.Marshal(Route{
		Paths: []string{service.KongConf.Route},
		Name:  serviceId + "-Route",
	})
	if err != nil {
		return err
	}
	res, err := http.Post(basePath+"/services/"+serviceId+"/routes/", "application/json", bytes.NewBuffer(_body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode <= 400 || res.StatusCode == 409 {
		fmt.Println("OK")
		return nil
	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func GetUpstreams() (UpstreamResp, error) {
	res, err := http.Get(basePath + "/upstreams/")
	body, _ := ioutil.ReadAll(res.Body)
	upstreams := UpstreamResp{}
	json.Unmarshal(body, &upstreams)
	defer res.Body.Close()
	if err != nil {
		return upstreams, err
	}
	if res.StatusCode <= 400 {
		return upstreams, nil
	} else {
		return upstreams, errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func GetServices() (ServiceResp, error) {
	res, err := http.Get(basePath + "/services/")
	body, _ := ioutil.ReadAll(res.Body)
	services := ServiceResp{}
	json.Unmarshal(body, &services)
	defer res.Body.Close()
	if err != nil {
		return services, err
	}
	if res.StatusCode <= 400 {
		return services, nil
	} else {
		return services, errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func GetRoutes(serviceId string) (RouteResp, error) {
	res, err := http.Get(basePath + "/services/" + serviceId + "/routes")
	body, _ := ioutil.ReadAll(res.Body)
	routes := RouteResp{}
	json.Unmarshal(body, &routes)
	defer res.Body.Close()
	if err != nil {
		return routes, err
	}
	if res.StatusCode <= 400 {
		return routes, nil
	} else {
		return routes, errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func GetTargets(upstreamId string) (TargetResp, error) {
	res, err := http.Get(basePath + "/upstreams/" + upstreamId + "/targets")
	body, _ := ioutil.ReadAll(res.Body)
	targets := TargetResp{}
	json.Unmarshal(body, &targets)
	defer res.Body.Close()
	if err != nil {
		return targets, err
	}
	if res.StatusCode <= 400 {
		return targets, nil
	} else {
		return targets, errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
