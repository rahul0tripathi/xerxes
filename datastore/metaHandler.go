package datastore

import (
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
)

func GetInstanceFlakeCounts() map[string]int {
	count := make(map[string]int)
	for id, _ := range config.Nodes.NodeList {
		count[id] = bitcask.GetInstanceFlakeCount(id)
	}
	return count
}
