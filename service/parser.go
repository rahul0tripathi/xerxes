package service

import (
	"fmt"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"strconv"
)

func init() {
	err := config.LoadConfig.ServiceDef()
	if err != nil {
		panic(err)
	}
}

func GetServiceDef(service string) definitions.ServiceDeclaration {
	if _, ok := config.ServicesDec.Def[service]; ok {
		return config.ServicesDec.Def[service]
	} else {
		return definitions.ServiceDeclaration{}
	}
}

func GetAvailableConfig(def definitions.ServiceDeclaration,service []definitions.ServiceDef) (string, string) {
	available := false
	nodeWeights := datastore.GetNodeServicesCount()
	availableNode := ""
	availablePort := ""
	for id, weight := range nodeWeights {
		if _, ok := config.Nodes.NodeList[availableNode]; ok && nodeWeights[availableNode] < weight {
		} else {
			availableNode = id
		}
	}
	for {
		randomPort := definitions.GenRandPort(def.BasePort, def.MaxPort)
		for _, used := range service {
			_port, err := strconv.Atoi(used.Port)
			if err != nil {
				continue
			}
			if _port != randomPort && used.Host == availableNode {
				available = true
				break
			}
		}
		if available  || len(service) == 0{
			availablePort = fmt.Sprintf("%d",randomPort)
			break
		}
	}
	return availableNode, availablePort
}
