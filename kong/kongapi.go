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
	ExternalGatewayBasePath string
	InternalGatewayBasePath string
)

func init() {
	err := config.LoadConfig.KongConf()
	if err != nil {
		panic(err)
	}
	err = config.LoadConfig.InternalGateway()
	if err != nil {
		panic(err)
	}
	ExternalGatewayBasePath = config.KongConn.Host + ":" + config.KongConn.Admin
	InternalGatewayBasePath = config.InternalGateway.Host + ":" + config.InternalGateway.Admin
}
func CreateNewUpstream(service definitions.KongUpstream, id string) error {
	fmt.Print("\nCreating a new UpStream [", service.Name, "] .... ")
	_body, _ := json.Marshal(Upstream{
		service.Name,
		service.Hashon,
	})
	res, err := http.Post(ExternalGatewayBasePath+"/upstreams", "application/json", bytes.NewBuffer(_body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode <= 400 || res.StatusCode == 409 {
		fmt.Println("OK")
		fmt.Print("Creating a new UpStream in internal Gateway... ")
		_body, _ = json.Marshal(Upstream{
			fmt.Sprintf("%s_upstream", id),
			service.Hashon,
		})
		res, err = http.Post(InternalGatewayBasePath+"/upstreams", "application/json", bytes.NewBuffer(_body))
		if err != nil {
			return err
		}
		res.Body.Close()
		if res.StatusCode <= 400 || res.StatusCode == 409 {
			fmt.Println("OK")
			return nil
		} else {
			return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))
		}

	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func AddUpstreamTarget(service definitions.KongUpstream, target UpstreamTarget, id string) error {
	fmt.Print("\nUpdating Upstream Target [", target.Target, "] to [", service.Name, "] .... ")
	_body, err := json.Marshal(target)
	if err != nil {
		return err
	}
	res, err := http.Post(ExternalGatewayBasePath+"/upstreams/"+service.Name+"/targets", "application/json", bytes.NewBuffer(_body))
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode <= 400 {
		fmt.Println("OK")
		fmt.Print("\nUpdating Upstream Target In Internal Gateway [", target.Target, "] to [", fmt.Sprintf("%s_upstream", id), "] .... ")
		res, err = http.Post(InternalGatewayBasePath+"/upstreams/"+fmt.Sprintf("%s_upstream", id)+"/targets", "application/json", bytes.NewBuffer(_body))
		if err != nil {
			return err
		}
		res.Body.Close()
		if res.StatusCode <= 400 {
			fmt.Println("OK")
			return nil
		} else {
			return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))
		}
	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))
	}
}
func CreateService(service definitions.Kongdef, id string) error {
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

	res, err := http.Post(ExternalGatewayBasePath+"/services/", "application/json", bytes.NewBuffer(_body))
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode <= 400 || res.StatusCode == 409 {
		fmt.Println("OK")
		fmt.Print("\nCreating a new Service In Internal gateway [", id, "] .... ")
		_body, err = json.Marshal(Service{
			Name: id,
			Host: fmt.Sprintf("%s_upstream", id),
			Path: "/",
			Port: 80,
		})
		if err != nil {
			return err
		}
		res, err = http.Post(InternalGatewayBasePath+"/services/", "application/json", bytes.NewBuffer(_body))
		if err != nil {
			return err
		}
		if res.StatusCode <= 400 || res.StatusCode == 409 {
			fmt.Println("OK")
			return nil
		} else {
			return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))
		}
	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func CreateNewRoute(service definitions.KongService, id string) error {
	fmt.Print("\nCreating a new Route [", service.Route, "] [", service.Name, "] .... ")
	_body, err := json.Marshal(Route{
		Paths: []string{service.Route},
		Name:  service.Name + "-Route",
	})
	if err != nil {
		return err
	}
	res, err := http.Post(ExternalGatewayBasePath+"/services/"+service.Name+"/routes/", "application/json", bytes.NewBuffer(_body))
	if err != nil {
		return err
	}
	res.Body.Close()
	if res.StatusCode <= 400 || res.StatusCode == 409 {
		fmt.Println("OK")
		fmt.Print("\nCreating a new Route In Internal Gateway [", id, "] [", id, "] .... ")
		_body, err = json.Marshal(Route{
			Paths: []string{fmt.Sprintf("/%s", id)},
			Name:  id + "-Route",
		})
		res, err = http.Post(InternalGatewayBasePath+"/services/"+id+"/routes/", "application/json", bytes.NewBuffer(_body))
		if err != nil {
			return err
		}
		res.Body.Close()
		if res.StatusCode <= 400 || res.StatusCode == 409 {
			fmt.Println("OK")
			return nil
		} else {
			return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))
		}
	} else {
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d", res.StatusCode))

	}
}
func GetUpstreams() (UpstreamResp, error) {
	res, err := http.Get(ExternalGatewayBasePath + "/upstreams/")
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
	res, err := http.Get(ExternalGatewayBasePath + "/services/")
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
	res, err := http.Get(ExternalGatewayBasePath + "/services/" + serviceId + "/routes")
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
	res, err := http.Get(ExternalGatewayBasePath + "/upstreams/" + upstreamId + "/targets")
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
