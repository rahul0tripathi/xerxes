package services

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	"time"
)

func deleteContainer(_cid string, _hostId string) {
	Dockerclient, err := dockerClient.NewDockerClient(_hostId)
	if err != nil {
		fmt.Print(err)
		return
	}
	err = Dockerclient.ContainerRemove(context.Background(), _cid, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})
	if err != nil {
		fmt.Print(err)
		return
	}
}
func DeployScheduler(triggerTime int64, _cid string, _hostId string) {
	select {
	case <-time.After(time.Duration(triggerTime) * time.Second):
		deleteContainer(_cid, _hostId)
	}
}
