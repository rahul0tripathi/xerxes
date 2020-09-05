package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gookit/color"
	"github.com/olekukonko/tablewriter"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore"
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
					&cli.StringFlag{
						Name:  "service",
						Usage: "Id of the service to scale",
					},
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
					return Scale(c.String("service"), c.String("node"), c.Int64("number"))
				},
			},
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "service",
						Usage: "Id of the service to scale",
					},
				},
				Name:    "definition",
				Aliases: []string{"def"},
				Usage:   "get definition of service",
				Action: func(c *cli.Context) error {
					if _, ok := config.ServicesDec.Def[c.String("service")]; ok {
						serv := config.ServicesDec.Def[c.String("service")]
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
						Name:  "service",
						Usage: "Id of the service to update",
					},
					&cli.StringFlag{
						Name:  "node",
						Usage: "specific node to update",
					},
				},
				Name:    "update",
				Aliases: []string{"upd"},
				Usage:   "update all running pods with new images",
				Action: func(c *cli.Context) error {
					Update(c.String("service"), c.String("node"))
					return nil
				},
			},
			{
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "service",
						Usage: "Id of the service to scale",
					},
				},
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "get running Pods of service",
				Action: func(c *cli.Context) error {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Id", "containerId", "host_node", "ip", "port"})
					if _, ok := config.ServicesDec.Def[c.String("service")]; ok {
						running := datastore.GetAllServicesPods(c.String("service"))
						for _, serv := range running {
							table.Append([]string{serv.Id, serv.ContainerId, serv.Host, serv.Ip, serv.Port})
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
		Name:  "pod",
		Usage: "Interact with running pods",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pod",
				Usage: "Id of the pod to inspect",
			},
			&cli.StringFlag{
				Name:  "service",
				Usage: "service under which the pod was created",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:  "inspect",
				Usage: "inspect a container",
				Action: func(c *cli.Context) error {
					podDef := datastore.GetpodById(c.String("service"), c.String("pod"))
					conn, err := dockerClient.NewDockerClient(podDef.Host)
					if err != nil {
						return err
					}
					podInspect, err := conn.ContainerInspect(context.Background(), podDef.ContainerId)
					if err != nil {
						return err
					}
					s, _ := json.MarshalIndent(podInspect, "", "\t")
					color.Style{color.FgCyan, color.OpBold}.Printf(string(s))
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"rm"},
				Usage:   "remove a container",
				Action: func(c *cli.Context) error {
					if c.String("pod") != "" && c.String("service") != "" {
						return shutdownPod(c.String("service"), c.String("pod"))
					} else {
						return errors.New("invalid podId / service")
					}
				},
			},
		},
	},
}
