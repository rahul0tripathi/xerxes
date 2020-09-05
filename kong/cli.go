package kong

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"os"
)

var KongCommands = &cli.Command{
	Name:    "kong",
	Aliases: []string{"k"},
	Usage:   "Interact with Kong api gateway",
	Subcommands: []*cli.Command{
		{
			Name:    "upstreams",
			Aliases: []string{"upstr"},
			Usage:   "List all upstreams",
			Action: func(c *cli.Context) error {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "HashOn"})
				upstreams, err := GetUpstreams()
				if err != nil {
					return err
				}
				for _, stream := range upstreams.Data {
					table.Append([]string{stream.Id, stream.Name, stream.HashOn})
				}
				table.Render()
				return nil
			},
		}, {
			Name:    "routes",
			Aliases: []string{"route"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "service",
					Usage: "Id of the service to Get routes for",
				},
			},
			Usage: "Routes of particular service",
			Action: func(c *cli.Context) error {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Paths"})
				if c.String("service") == "" {
					return errors.New("Empty service Id")
				}
				routes, err := GetRoutes(c.String("service"))
				if err != nil {
					return err
				}
				for _, route := range routes.Data {
					data := []string{route.Id, route.Name}
					data = append(data, route.Paths...)
					table.Append(data)
				}
				table.Render()
				return nil
			},
		},
		{
			Name:    "services",
			Aliases: []string{"serv"},
			Usage:   "List all services",
			Action: func(c *cli.Context) error {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Name", "Protocol", "Host", "port", "path"})
				services, err := GetServices()
				if err != nil {
					return err
				}
				for _, service := range services.Data {
					table.Append([]string{service.Id, service.Name, service.Protocol, service.Host, fmt.Sprintf("%d", service.Port), service.Path})
				}
				table.Render()
				return nil
			},
		},
		{
			Name:    "targets",
			Aliases: []string{"target"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "upstream",
					Usage: "Id of the upstream to Get targets for",
				},
			},
			Usage: "List all services",
			Action: func(c *cli.Context) error {
				if c.String("upstream") == "" {
					return errors.New("Empty Upstream Id")
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Id", "Target", "Weight"})
				targets, err := GetTargets(c.String("upstream"))
				if err != nil {
					return err
				}
				for _, target := range targets.Data {
					table.Append([]string{target.Id, target.Target, fmt.Sprintf("%d", target.Weight)})
				}
				table.Render()
				return nil
			},
		},
	},
}
