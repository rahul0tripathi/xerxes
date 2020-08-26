package main

import (
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/kong"
	"github.com/rahultripathidev/docker-utility/service"
	"github.com/rahultripathidev/docker-utility/update"
	"github.com/urfave/cli/v2"
	"os"
	"log"
)

func main() {
	err := config.LoadServices(config.ConfigDir)
	if err != nil {
		panic(err)
	}
	app := &cli.App{
		Name: "Xerxes",
		Usage: "🛳️ A simple container orchestration service using Kong Gateway",
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
				Usage:   "scale a service to a desired number",
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
				Usage:   "shuts down all running instances of a service [not a swarm]",
				Action:  func(c *cli.Context) error {
					if _ ,ok :=config.Config.Services[c.String("service")] ; ok {
						service.Shutdown(c.String("service"),config.Config.Services[c.String("service")].Type,c.Int("number"))
						return nil
					}else{
						return  errors.New("Undefined service "+c.String("service"))
					}

				},
			},
			{
				Name:    "services",
				Aliases: []string{"serv"},
				Usage:   "List all running services",
				Action:  func(c *cli.Context) error {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Id", "Image", "createdAt","Endpoint","ports","replicas"})
					services ,_  := service.GetActiveServices()
					for _ , s := range services{
						port_config := ""
						for _ , port := range s.Spec.EndpointSpec.Ports {
							port_config += fmt.Sprintf("%d->%d/%s",port.PublishedPort,port.TargetPort,port.Protocol)
						}
						table.Append([]string{s.ID,s.Spec.TaskTemplate.ContainerSpec.Image,s.CreatedAt.String(),string(s.Spec.EndpointSpec.Mode),port_config , fmt.Sprintf("%d",*s.Spec.Mode.Replicated.Replicas)})
					}
					table.Render()
					return nil
				},
			},
			{
				Name:    "containers",
				Aliases: []string{"cont"},
				Usage:   "List all running containers",
				Action:  func(c *cli.Context) error {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Id", "serviceId", "Host","NodeId"})
					services ,_  := service.GetActiveContainers()
					for _ , s := range services.Avaliable{
						table.Append([]string{s.ContainerId[:10],s.ServiceId,fmt.Sprintf("%s:%d",s.HostIp,s.BindingPort),fmt.Sprintf("%d",s.Id)})
					}
					table.Render()
					return nil
				},
			},
			{
				Name:    "machines",
				Aliases: []string{"mach"},
				Usage:   "List all avaliable machines",
				Action:  func(c *cli.Context) error {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Id", "Docker_Host", "IP","Docker_API_VERSION"})
					table.Append([]string{fmt.Sprintf("%d",config.Nodelist.Master.Id),config.Nodelist.Master.Host,config.Nodelist.Master.Ip,config.Nodelist.Master.Version})
					for _ , node := range config.Nodelist.Nodes{
						table.Append([]string{fmt.Sprintf("%d",node.Id),node.Host,node.Ip,node.Version})
					}
					table.Render()
					return nil
				},
			},
			{
				Name:    "update",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "service",
						Usage: "Id of the service to update",
					},
					&cli.BoolFlag{
						Name: "userNode",
						Usage: "number of containers to shutdown to",
					},
				},
				Aliases: []string{"update"},
				Usage:   "Update a service machines",
				Action:  func(c *cli.Context) error {
					err := update.Update(c.String("service"),c.Bool("useNode"))
					return err
				},
			},
			{
				Name:    "kong",
				Aliases: []string{"k"},
				Usage:   "Interact with Kong api gateway",
				Subcommands: []*cli.Command{
					{
					Name:  "upstreams",
					Aliases: []string{"upstr"},
					Usage: "List all upstreams",
					Action: func(c *cli.Context) error {
						table := tablewriter.NewWriter(os.Stdout)
						table.SetHeader([]string{"Id", "Name", "HashOn"})
						upstreams, err := kong.GetUpstreams()
						if err != nil {
							return err
						}
						for _, stream := range upstreams.Data {
							table.Append([]string{stream.Id, stream.Name, stream.HashOn})
						}
						table.Render()
						return nil
					},
					},{
						Name:  "routes",
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
							if c.String("service") == ""{
								return errors.New("Empty service Id")	
							}
							routes, err := kong.GetRoutes(c.String("service"))
							if err != nil {
								return err
							}
							for _, route := range routes.Data {
								data := []string{route.Id, route.Name}
								data = append(data,route.Paths...)
								table.Append(data)
							}
							table.Render()
							return nil
						},
					},
					{
						Name:  "services",
						Aliases: []string{"serv"},
						Usage: "List all services",
						Action: func(c *cli.Context) error {
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Id", "Name", "Protocol","Host","port","path"})
							services, err := kong.GetServices()
							if err != nil {
								return err
							}
							for _, service := range services.Data {
								table.Append([]string{service.Id, service.Name, service.Protocol,service.Host,fmt.Sprintf("%d",service.Port),service.Path})
							}
							table.Render()
							return nil
						},
					},
					{
						Name:  "targets",
						Aliases: []string{"target"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "upstream",
								Usage: "Id of the upstream to Get targets for",
							},
						},
						Usage: "List all services",
						Action: func(c *cli.Context) error {
							if c.String("upstream") == ""{
								return errors.New("Empty Upstream Id")
							}
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Id", "Target", "Weight"})
							targets, err := kong.GetTargets(c.String("upstream"))
							if err != nil {
								return err
							}
							for _, target := range targets.Data {
								table.Append([]string{target.Id,target.Target,fmt.Sprintf("%d",target.Weight)})
							}
							table.Render()
							return nil
						},
					},
				},
			},
		},

	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
