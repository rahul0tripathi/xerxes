package kong

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rahultripathidev/docker-utility/config"
	"net/http"
)
var (
	basePath string
)
func init(){
	config.LoadServices("$HOME/.orchestrator/configuration")
	err := config.LoadKongConf("$HOME/.orchestrator/configuration")
	if err != nil {
		panic(err)
	}
	basePath = config.KongConf.Host +":"+config.KongConf.Admin
}
func CreateNewUpstream(service config.KongDesc) error {
	_body , _:= json.Marshal(Upstream{
		service.UpstreamService,
		"ip",
	})
	res , err := http.Post(basePath+"/upstreams","application/json",bytes.NewBuffer(_body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode == 200 || res.StatusCode == 409{
		return nil
	}else{
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d",res.StatusCode))

	}
}
func AddUpstreamTarget(service config.KongDesc , target UpstreamTarget) error {
	_body , err := json.Marshal(target)
	if err != nil {
		return err
	}
	res , err := http.Post(basePath+"/upstreams/"+service.UpstreamService+"/targets","application/json",bytes.NewBuffer(_body))
	//body ,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode <= 400{
		return nil
	}else{
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d",res.StatusCode))

	}
}
func CreateService(serviceId string,service config.Services , serviceEndpointHost string , serviceEndpointPort int) error {
	var _body []byte
	var err error
	if service.KongConf.Upstream == true{
		_body , err = json.Marshal(Service{
			Name : serviceId,
			Host: service.KongConf.UpstreamService,
			Path: service.KongConf.ServicePath,
			Port: 80,
		})
		if err != nil {
			return err
		}
	}else{
		_body , err = json.Marshal(Service{
			Name : serviceId,
			Host:  serviceEndpointHost,
			Path:  service.KongConf.ServicePath,
			Port: serviceEndpointPort,
		})
		if err != nil {
			return err
		}
	}
	res , err := http.Post(basePath+"/services/","application/json",bytes.NewBuffer(_body))
	//body ,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode <= 400 || res.StatusCode == 409{
		return nil
	}else{
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d",res.StatusCode))

	}
}
func CreateNewRoute(serviceId string ,service config.Services) error{
	_body , err := json.Marshal(Route{
		Paths: []string{service.KongConf.Route},
		Name: serviceId+"-Route",
	})
	if err != nil {
		return err
	}
	res , err := http.Post(basePath+"/services/"+serviceId+"/routes/","application/json",bytes.NewBuffer(_body))
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode <= 400 || res.StatusCode == 409{
		return nil
	}else{
		return errors.New("Request Failed with a status of " + fmt.Sprintf("%d",res.StatusCode))

	}
}
