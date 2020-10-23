package service

import (
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
)

func Scale(serviceId string, nodeId string, factor int) error {
	runningServices, _ := bitcask.GetAllServiceFlakes(serviceId)
	var i int
	if len(runningServices) < factor {
		for i = 0; i < factor-len(runningServices); i++ {
			err := ScaleUp(serviceId, nodeId)
			if err != nil {
				return err
			}
		}
	} else if len(runningServices) > factor {
		for i = 0; i < len(runningServices)-factor; i++ {
			err := ScaleDown(serviceId, "", nodeId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
