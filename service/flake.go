package service

import (
	"context"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"io"
	"os"
	"time"
)

func GetFlakeLogs(flake definitions.FlakeDef) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dockerDaemon, err := dockerClient.NewDockerClient(flake.HostId)
	if err != nil {
		return err
	}
	reader, err := dockerDaemon.ContainerLogs(ctx, flake.ContainerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
