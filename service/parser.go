package service

import (
	"fmt"
	"github.com/rahultripathidev/docker-utility/config"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"math/rand"
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

func GetAvailableConfig(def definitions.ServiceDeclaration, service []definitions.FlakeDef) (string, string) {
	availableNode := ""
	availablePort := ""
	flatMap := make([]string,0,len(config.Nodes.NodeList))
	for k := range config.Nodes.NodeList {
		flatMap = append(flatMap,k)
	}
	availableNode = flatMap[rand.Intn(len(flatMap))]
	availablePort = fmt.Sprintf("%d",definitions.GenRandPort(def.BasePort, def.MaxPort))
	return availableNode, availablePort
}
