package config

import (
	"fmt"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"github.com/spf13/viper"
	"os"
)

// Variables Of all the decelerations that are loaded
// Includes Caches connection config , Kong Admin Api , List of available nodes , Services definitions

var (
	Config      config
	ServicesDec struct {
		Def map[string]definitions.ServiceDeclaration `mapstructure:"services"`
	}
	KongConn    definitions.KongConn
	Nodes       definitions.Nodes
	ConfigDir   string
	CacheConfig definitions.CacheClient
	ServiceDiscovery map[string]string
	XerxesHost struct{
		Host string `mapstructure:"host"`
	}
)

// func ReadAndUnmarshal reads a config field / or its sub and unmarshal its json to the corresponding variable
// and returns error
func ReadAndUnmarshal(configname string, format string, object interface{}, sub interface{}) error {
	conf := viper.New()
	conf.SetConfigName(configname)
	conf.SetConfigType(format)
	conf.AddConfigPath(ConfigDir)
	if err := conf.ReadInConfig(); err != nil {
		fmt.Printf("error reading Sub config %v", err)
		return err
	}
	if sub != nil {
		subItem := conf.Sub(sub.(string))
		return subItem.Unmarshal(object)
	}
	return conf.Unmarshal(object)
}
func setEnv(object map[string]string) {
	for key, value := range object {
		err := os.Setenv(key, value)
		if err != nil {
			fmt.Println("Unable to set env ", key)
		}
	}
}

// init initializes the config directory , the default is $HOME/.orchestrator/configuration
func init() {
	//ConfigDir = "./configuration"
	ConfigDir = func() string { HOME , _ := os.UserHomeDir()
		return HOME }() + "/.orchestrator/configuration"
	ReadAndUnmarshal("config", "json", &XerxesHost, "config.xerxes_host")
}

// Load config contains the functions to Load Configs into their respective variables
var (
	LoadConfig = struct {
		Cache        func() error
		Nodes        func() error
		ServiceDef   func() error
		RegistryAuth func() error
		KongConf     func() error
		ServiceDiscovery func() error
	}{
		Cache: func() error {
			return ReadAndUnmarshal("config", "json", &CacheConfig, "config.cache")
		},
		Nodes: func() error {
			return ReadAndUnmarshal("host", "json", &Nodes, nil)
		},
		ServiceDef: func() error {
			return ReadAndUnmarshal("service", "json", &ServicesDec, nil)
		},
		RegistryAuth: func() error {
			credentials := make(map[string]string)
			err := ReadAndUnmarshal("config", "json", &credentials, "config.registry.AWS")
			if err != nil {
				return err
			}
			fmt.Println(credentials)
			setEnv(credentials)
			return nil
		},
		KongConf: func() error {
			return ReadAndUnmarshal("config", "json", &KongConn, "config.kong")
		},
		ServiceDiscovery: func() error {
			return ReadAndUnmarshal("config", "json", &ServiceDiscovery, "config.discovery")
		},
	}
)
