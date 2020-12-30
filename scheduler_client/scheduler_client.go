package scheduler_client

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	containerManager "github.com/rahultripathidev/docker-utility/proto"
	"google.golang.org/grpc"
)

func forceDelete(cid string, nodeId string) error {
	Dockerclient, err := dockerClient.NewDockerClient(nodeId)
	if err != nil {
		return err
	}
	err = Dockerclient.ContainerRemove(context.Background(), cid, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		fmt.Print("Failed To shutdown service  [", cid, "] error ", err)
		return err
	}
	return nil
}
func ScheduleContainerForDeletion(cid string, nodeId string, force bool, timeout int64) {
	if timeout == 0 {
		timeout = 30
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 3333), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to schedule deletion , forcing delete")
		forceDelete(cid, nodeId)
		return
	}
	client := containerManager.NewContainerManagerClient(conn)
	req := deleteCommand(cid, nodeId, force, timeout)
	var resp *containerManager.DeleteContainerResponse
	resp, err = client.DeleteContainer(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to schedule deletion , forcing delete")
		forceDelete(cid, nodeId)
		return
	}
	fmt.Println(fmt.Sprintf("Scheduled to shutdown container in %d Seconds : %t", timeout, resp.Success))
	return
}
