package datastore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"github.com/rs/xid"
)

func RedisResponse(resp *redis.IntCmd) bool {
	if statusCode, _ := resp.Uint64(); statusCode == 1 {
		return true
	} else {
		return false
	}
}
func GetAllServicesPods(service string) []definitions.ServiceDef {
	var services []definitions.ServiceDef
	ctx := context.Background()
	val := RedisClient.HGetAll(ctx, XerxesService+service)
	for id, value := range val.Val() {
		func(id string, serviceValue string) {
			newService := definitions.ServiceDef{}
			err := json.Unmarshal([]byte(serviceValue), &newService)
			if err != nil {
				fmt.Println("error occurred while parsing json ", err)
			}
			newService.Id = id
			services = append(services, newService)
		}(id, value)
	}
	return services
}
func GetpodById(service string, pod string) definitions.ServiceDef {
	ctx := context.Background()
	podInfo := definitions.ServiceDef{}
	val := RedisClient.HGet(ctx, XerxesService+service, pod)
	err := json.Unmarshal([]byte(val.Val()), &podInfo)
	if err != nil {
		fmt.Println("error occurred while parsing json", err)
	}
	podInfo.Id = pod
	return podInfo
}
func GetAllpodsInNode(nodeId string) []string {
	ctx := context.Background()
	val := RedisClient.LRange(ctx, XerxesNodes+nodeId, 0, -1)
	return val.Val()
}
func addNewPodToNode(podId string, nodeId string) {
	ctx := context.Background()
	RedisClient.LPush(ctx, XerxesNodes+nodeId, podId)
}
func removePodFromNode(podId string, nodeId string) {
	ctx := context.Background()
	RedisClient.LRem(ctx, XerxesNodes+nodeId, 1, podId)
}
func RegisterPod(service string, pod definitions.ServiceDef) error {
	ctx := context.Background()
	data, err := json.Marshal(pod)
	if err != nil {
		return err
	}
	_id := xid.New()
	resp := RedisClient.HSet(ctx, XerxesService+service, _id.String(), string(data))
	if RedisResponse(resp) {
		addNewPodToNode(_id.String(), pod.Host)
		return nil
	} else {
		return errors.New(resp.String())
	}
}
func DeregisterPod(service string, podId string, nodeId string) error {
	ctx := context.Background()
	resp := RedisClient.HDel(ctx, XerxesService+service, podId)
	if RedisResponse(resp) {
		removePodFromNode(podId, nodeId)
		return nil
	} else {
		return errors.New(resp.String())
	}
}
func GetTotalServices(serviceId string) int64 {
	val := RedisClient.HLen(context.Background(), XerxesService+serviceId);
	return val.Val()
}