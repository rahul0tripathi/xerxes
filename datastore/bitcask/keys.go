package bitcask

import (
	"fmt"
	"github.com/prologic/bitcask"
)

var (
	ErrKeyNotFound = bitcask.ErrKeyNotFound
	keys = map[string]func(parmas ...string) (key []byte){
		"flake": func(params ...string) (key []byte) {
			return []byte(fmt.Sprintf("FLAKE:%s", params[0]))
		},
		"instance": func(params ...string) (key []byte) {
			if len(params) == 1 {
				return []byte(fmt.Sprintf("INSTANCE:%s", params[0]))
			} else {
				return []byte(fmt.Sprintf("INSTANCE:%s:%s", params[0], params[1]))
			}
		},
		"service": func(params ...string) (key []byte) {
			if len(params) == 1 {
				return []byte(fmt.Sprintf("SERVICE:%s", params[0]))
			} else {
				return []byte(fmt.Sprintf("SERVICE:%s:%s", params[0], params[1]))
			}

		},
	}
)
