package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
	"github.com/urfave/cli/v2"
	"os"
)

var ServiceCommands = []*cli.Command{
	{
		Name:  "service",
		Usage: "Interact with service module",
		Subcommands: []*cli.Command{
			{
				Flags: []cli.Flag{
					&cli.Uint64Flag{
						Name:  "number",
						Usage: "number of containers to scale to",
					},
					&cli.StringFlag{
						Name:  "node",
						Usage: "force orchestrator to use Node",
					},
				},
				Name:    "scale",
				Aliases: []string{"scale"},
				Usage:   "scale a service to a desired number",
				Action: func(c *cli.Context) error {
					return Scale(c.Args().Get(0), c.String("node"), c.Int("number"))
				},
			},
			{

				Name:    "definition",
				Aliases: []string{"def"},
				Usage:   "get definition of service",
				Action: func(c *cli.Context) error {
					if _, ok := config.ServicesDec.Def[c.Args().Get(0)]; ok {
						serv := config.ServicesDec.Def[c.Args().Get(0)]
						s, _ := json.MarshalIndent(serv, "", "\t")
						color.Style{color.FgCyan, color.OpBold}.Println(string(s))
						return nil
					} else {
						return errors.New("service not found")
					}
				},
			},
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "instance",
						Usage: "specific node to update",
					},
				},
				Name:    "update",
				Aliases: []string{"upd"},
				Usage:   "update all running pods with new images",
				Action: func(c *cli.Context) error {
					Update(c.Args().Get(0), c.String("instance"))
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "get running Pods of service",
				Action: func(c *cli.Context) error {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Id", "containerId", "host_node", "ip", "port"})
					if _, ok := config.ServicesDec.Def[c.Args().Get(0)]; ok {
						running, err := bitcask.GetAllServiceFlakes(c.Args().Get(0))
						if err == nil {
							for _, flake := range running {
								table.Append([]string{flake.Id, flake.ContainerId, flake.HostId, flake.Ip, flake.Port})
							}
						}
						table.Render()
						return nil
					} else {
						return errors.New("service not found")
					}
				},
			},
		},
	}, {
		Name:  "flake",
		Usage: "Interact with running pods",
		Subcommands: []*cli.Command{
			{
				Name:  "inspect",
				Usage: "inspect a container",
				Action: func(c *cli.Context) error {
					flakeDef, err := bitcask.GetFlake(c.Args().Get(0))
					if err != nil {
						return err
					}
					conn, err := dockerClient.NewDockerClient(flakeDef.HostId)
					if err != nil {
						return err
					}
					flakeInspect, err := conn.ContainerInspect(context.Background(), flakeDef.ContainerId)
					if err != nil {
						return err
					}
					s, _ := json.MarshalIndent(flakeInspect, "", "\t")
					color.Style{color.FgCyan, color.OpBold}.Printf(string(s))
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "remove a container",
				Action: func(c *cli.Context) error {
					if c.Args().Get(0) != "" {
						flakeDef, err := bitcask.GetFlake(c.Args().Get(0))
						if err != nil {
							return err
						}
						return shutdownPod(flakeDef)
					} else {
						return errors.New("invalid podId / service")
					}
				},
			},
			{
				Name:  "logs",
				Usage: "Get Container logs",
				Action: func(c *cli.Context) error {
					flakeDef, err := bitcask.GetFlake(c.Args().Get(0))
					if err != nil {
						return err
					}
					err = GetFlakeLogs(flakeDef)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	},
}
