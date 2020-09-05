package datastore

import (
	"context"
	"github.com/rahultripathidev/docker-utility/config"
)

func GetNodeServicesCount() map[string]int64 {
	count := make(map[string]int64)
	for id, _ := range config.Nodes.NodeList {
		resp := RedisClient.LLen(context.Background(), XerxesNodes+id)
		count[id] = resp.Val()
	}
	return count
}
