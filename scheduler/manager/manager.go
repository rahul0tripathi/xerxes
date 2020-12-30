package manager

import (
	"context"
	containerManager "github.com/rahultripathidev/docker-utility/proto"
	"github.com/rahultripathidev/docker-utility/scheduler/services"
)

type ContainerManagerService struct {
	containerManager.UnimplementedContainerManagerServer
}

func (containerManagerService *ContainerManagerService) DeleteContainer(ctx context.Context, containerMeta *containerManager.DeleteContainerRequest) (response *containerManager.DeleteContainerResponse, err error) {
	go services.DeployScheduler(containerMeta.Meta.Timeout, containerMeta.Meta.XCid, containerMeta.Meta.XNodeId)
	return &containerManager.DeleteContainerResponse{Success: true}, nil
}
