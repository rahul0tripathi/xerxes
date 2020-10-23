package main

import (
	"github.com/gookit/color"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
	"github.com/rahultripathidev/docker-utility/kong"
	"github.com/rahultripathidev/docker-utility/service"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	err := config.LoadConfig.Bit()
	if err != nil {
		panic(err)
	}
	bitcask.InitClient()
	defer func() {
		err = bitcask.BitClient.Flock.Unlock()
		bitcask.GracefulClose()
	}()
	var commands []*cli.Command
	commands = append(commands, service.ServiceCommands...)
	commands = append(commands, kong.KongCommands)
	app := &cli.App{
		Name:     "Xerxes",
		Usage:    "🛳️ A \"simple\" container orchestration service ",
		Commands: commands,
	}

	err = app.Run(os.Args)
	if err != nil {
		color.Error.Println(err)
	}
}
