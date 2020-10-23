package bitcask

import (
	"github.com/prologic/bitcask"
	"github.com/rahultripathidev/docker-utility/config"
)

var (
	BitClient *bitcask.Bitcask
	err       error
)

func InitClient() {
	BitClient, err = bitcask.Open(config.BitConf.Dbpath)
	if err != nil {
		panic(err)
	}
}

func GracefulClose() {
	err = BitClient.Close()
	if err != nil {
		panic(err)
	}
}
