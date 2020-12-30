package scheduler_client

import (
	containerManager "github.com/rahultripathidev/docker-utility/proto"
)

func deleteCommand(cid string, nodeId string, force bool, timeout int64) (command *containerManager.DeleteContainerRequest) {
	command = &containerManager.DeleteContainerRequest{
		Meta: &containerManager.ContainerMeta{
			XCid:    cid,
			XNodeId: nodeId,
			Force:   force,
			Timeout: timeout,
		},
	}
	return
}
