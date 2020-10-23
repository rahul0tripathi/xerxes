package datastore

import (
	"errors"
	"fmt"
	"github.com/rahultripathidev/docker-utility/config"
	"net/http"
)

func NotifyChange() error {
	resp, err := http.Get(fmt.Sprintf("http://%s:8937/reload", config.XerxesHost.Host))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return errors.New("could not process request")
	}
	return nil
}
