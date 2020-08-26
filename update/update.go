package update

import (
	"encoding/json"
	"fmt"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/service"
	"io/ioutil"
)
type (
	updateHolder struct {
		Service map[string]updateDef `json:"services"`
	}
	updateDef struct {
		Image string `json:"image"`
		ImageUri string `json:"image_uri"`
		Type string `json:"type"`
		Factor int `json:"factor"`
	}
)
var (
	confFile string
	updateSchema = updateHolder{}
)
func init(){
	confFile =  config.ConfigDir + "/update.json"
	body , err := ioutil.ReadFile(confFile)
	if err != nil {
		fmt.Println("Error Reading update config ",err)
		return
	}
	err = json.Unmarshal(body,&updateSchema)
	if err != nil {
		fmt.Println("Error Reading update config ",err)
		return
	}
}

func Update(serviceId string , useNode bool) error {
	fmt.Print("\n Updating service [",serviceId,"]...")
	if _ , updateFound := updateSchema.Service[serviceId] ; !updateFound {
		fmt.Print("Update Definiton not found for service ",serviceId)
	}
	updateStruct := updateSchema.Service[serviceId]
	if _, serviceFound := config.Config.Services[serviceId] ; !serviceFound {
		fmt.Print("Service Definiton not found for service ",serviceId)
	}
	serviceStruct := config.Config.Services[serviceId]
	switch updateStruct.Type {
	case "rolling" : return  func() error{
		if updateStruct.Factor > 0{
			fmt.Printf("\nshutting down %d services having image %s ...",updateStruct.Factor,serviceStruct.Image)
			err := service.Shutdown(serviceId,serviceStruct.Type,updateStruct.Factor)
			if err != nil{
				fmt.Print("Unable to shutdown")
				return err
			}
			fmt.Print("OK\n")
			serviceStruct.ImageUri = updateStruct.ImageUri
			serviceStruct.Image = updateStruct.Image
			fmt.Print("\n Updating service with new image [",updateStruct.Image,"]...")
			service.ScalebyService(serviceId,serviceStruct,updateStruct.Factor,useNode)
			fmt.Print("OK\n")
			return nil
		}else{
			fmt.Print("factor Must be grater than 0")
			return nil
		}
	}()
	case "bluegreen":func(){
		//err := service.Shutdown(serviceId,serviceStruct.Type,updateStruct.Factor)
	}()
	default:
		return nil
	}
}
