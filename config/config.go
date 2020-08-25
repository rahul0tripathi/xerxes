package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var (
	Config   config
	KongConf Kong
	Nodelist NodeList
	ConfigDir string
)

func LoadServices(file string) error {
	services := viper.New()
	services.SetConfigName("service")
	services.SetConfigType("json")
	services.AddConfigPath(file)
	if err := services.ReadInConfig(); err != nil {
		fmt.Printf("ERROR READING CONFIG %v", err)
		return err
	}
	return services.Unmarshal(&Config)
}

func LoadKongConf(file string) error {
	services := viper.New()
	services.SetConfigName("kong")
	services.SetConfigType("json")
	services.AddConfigPath(file)
	if err := services.ReadInConfig(); err != nil {
		fmt.Printf("ERROR READING CONFIG %v", err)
		return err
	}
	return services.Unmarshal(&KongConf)
}

func LoadHosts(file string) error {
	services := viper.New()
	services.SetConfigName("host")
	services.SetConfigType("json")
	services.AddConfigPath(file)
	if err := services.ReadInConfig(); err != nil {
		fmt.Printf("ERROR READING CONFIG %v", err)
		return err
	}
	return services.Unmarshal(&Nodelist)
}
func init(){
	//ConfigDir = "./configuration"
	ConfigDir = func() string { HOME , _ := os.UserHomeDir()
		return HOME }() + "/.orchestrator/configuration"
}
