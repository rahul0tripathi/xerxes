package main

import (
	"fmt"
	"github.com/gookit/color"
	containerManager "github.com/rahultripathidev/docker-utility/proto"
	"github.com/rahultripathidev/docker-utility/scheduler/manager"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	PORT = "3333"
)

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", PORT))
	color.Style{color.FgCyan, color.OpBold}.Printf("[%s] Server Started \n", time.Now().String())
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	containerManager.RegisterContainerManagerServer(grpcServer, new(manager.ContainerManagerService))
	err = grpcServer.Serve(listen)
	if err != nil {
		panic(err)
	}
}
