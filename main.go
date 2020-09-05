package main

import (
	 "github.com/gookit/color"
	"github.com/rahultripathidev/docker-utility/kong"
	"github.com/rahultripathidev/docker-utility/service"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	var commands []*cli.Command
	commands = append(commands, service.ServiceCommands...)
	commands = append(commands, kong.KongCommands)
	app := &cli.App{
		Name:     "Xerxes",
		Usage:    "🛳️ A simple container orchestration service using Kong Gateway",
		Commands: commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		color.Error.Println(err)
	}
}
