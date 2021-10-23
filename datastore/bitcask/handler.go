package bitcask

import (
	"github.com/prologic/bitcask"
	"github.com/rahultripathidev/docker-utility/config"
)

var (
	BitClient *bitcask.Bitcask
	err       error
)

func InitClient() error {
	BitClient, err = bitcask.Open(config.BitConf.Dbpath)
	if err != nil {
		return err
	}
	return nil
}

func GracefulClose() error {
	err = BitClient.Close()
	if err != nil {
		return err
	}
	return nil
}
