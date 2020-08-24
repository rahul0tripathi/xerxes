package main

import (
	"errors"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/service"
	"github.com/urfave/cli/v2"
	"os"
	"log"
)

func main() {
	config.LoadServices("./configuration")
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Flags: []cli.Flag {
					&cli.StringFlag{
						Name: "service",
						Usage: "Id of the service to scale",
					},
					&cli.Uint64Flag{
						Name: "number",
						Usage: "number of containers to scale to",
					},
					&cli.BoolFlag{
						Name : "useNodes",
						Usage : "force orchestrator to user Nodes",
					},
				},
				Name:    "scale",
				Aliases: []string{"scale"},
				Usage:   "complete a task on the list",
				Action:  func(c *cli.Context) error {
					if c.Bool("useNodes") {
						 service.Scale(c.String("service"),c.Uint64("number"),true)
					}else{
						service.Scale(c.String("service"),c.Uint64("number"),false)
					}
					return nil
				},
			},
			{
				Name:    "init",
				Aliases: []string{"init"},
				Usage:   "Initialises docker swarm , containers upto the required number and also creates corresponding Kong upstream,service,route",
				Action:  func(c *cli.Context) error {
					service.InitServices()
					return nil
				},
			},
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "service",
						Usage: "Id of the service to shutdown",
					},
					&cli.IntFlag{
						Name: "number",
						Usage: "number of containers to shutdown to",
					},
				},
				Name:    "shutdown",
				Aliases: []string{"shut"},
				Usage:   "Initialises docker swarm , containers upto the required number and also creates corresponding Kong upstream,service,route",
				Action:  func(c *cli.Context) error {
					if _ ,ok :=config.Config.Services[c.String("service")] ; ok {
						service.Shutdown(c.String("service"),config.Config.Services[c.String("service")].Type,c.Int("number"))
						return nil
					}else{
						return  errors.New("Undefined service "+c.String("service"))
					}

				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
