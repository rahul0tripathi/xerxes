package service

import "github.com/rahultripathidev/docker-utility/datastore"

func Scale(serviceId string , nodeId string , factor int64 ) error {
	runningServices := datastore.GetTotalServices(serviceId)
	var i int64
	if runningServices < factor {
		for i = 0 ; i < factor-runningServices; i++ {
			err := ScaleUp(serviceId , nodeId)
			if err != nil {
				return err
			}
		}
	}else if runningServices > factor {
		for i = 0 ; i < runningServices-factor; i++ {
			err := ScaleDown(serviceId , "",nodeId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}