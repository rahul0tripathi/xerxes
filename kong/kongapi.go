package kong

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rahultripathidev/docker-utility/config"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"io/ioutil"
	"net/http"
)

var (
	basePath string
)
func init() {
	err := config.LoadConfig.KongConf()
	if err != nil {
		panic(err)
	}
	basePath = config.KongConn.Host + ":" + config.KongConn.Admin
}
func CreateNewUpstream(service definitions.KongUpstream) error {
	fmt.Print("\nCreating a new UpStream [", service.Name, "] .... ")
	_body, _ := json.Marshal(Upstream{
		service.Name,
		service.Hashon,
	})
	res, err := http.Post(basePath+"/upstreams", "application/json", bytes.NewBuffer(_body))
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
func AddUpstreamTarget(service definitions.KongUpstream, target UpstreamTarget) error {
	fmt.Print("\nUpdating Upstream Target [", target.Target, "] to [", service.Name, "] .... ")
	_body, err := json.Marshal(target)
	if err != nil {
		return err
	}
	res, err := http.Post(basePath+"/upstreams/"+service.Name+"/targets", "application/json", bytes.NewBuffer(_body))
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
func CreateService(service definitions.Kongdef) error {
	fmt.Print("\nCreating a new Service [", service.Service.Name, "] .... ")
	var _body []byte
	var err error
	_body, err = json.Marshal(Service{
		Name: service.Service.Name,
		Host: service.Upstream.Name,
		Path: service.Service.TaregtPath,
		Port: 80,
	})
	if err != nil {
		return err
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
func CreateNewRoute(service definitions.KongService) error {
	fmt.Print("\nCreating a new Route [", service.Route, "] [", service.Name, "] .... ")
	_body, err := json.Marshal(Route{
		Paths: []string{service.Route},
		Name:  service.Name + "-Route",
	})
	if err != nil {
		return err
	}
	res, err := http.Post(basePath+"/services/"+service.Name+"/routes/", "application/json", bytes.NewBuffer(_body))
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
