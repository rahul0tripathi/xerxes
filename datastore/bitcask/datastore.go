package bitcask

import (
	"encoding/json"
	"github.com/prologic/bitcask"
	definitions "github.com/rahultripathidev/docker-utility/types"
)

type FlakeStores interface {
	GetFlakes() []string
}

type serviceStore struct {
	Flakes []string `json:"flakes"`
	Id     string   `json:"id"`
}

func (s *serviceStore) GetFlakes() []string {
	return s.Flakes
}
func (s *serviceStore) AddFlake(id string) {
	s.Flakes = append(s.Flakes, id)
}
func (s *serviceStore) RemoveFlake(id string) {
	s.Flakes = append(s.Flakes, id)
}
func (s *serviceStore) Save() (err error) {
	var buf []byte
	buf, err = json.Marshal(s)
	err = BitClient.Put(keys["service"](s.Id), buf)
	return
}

type instanceStore struct {
	Id     string   `json:"id"`
	Flakes []string `json:"flakes"`
}

func (i *instanceStore) GetFlakes() []string {
	return i.Flakes
}
func (i *instanceStore) AddFlake(id string) {
	i.Flakes = append(i.Flakes, id)
}
func (i *instanceStore) RemoveFlake(id string) {
	i.Flakes = append(i.Flakes, id)
}
func (i *instanceStore) Save() (err error) {
	var buf []byte
	buf, err = json.Marshal(i)
	err = BitClient.Put(keys["instance"](i.Id), buf)
	return
}
func GetInstanceFlakeCount(instanceId string) (count int) {
	data, err := get(keys["instance"](instanceId))
	if err != nil {
		return
	}
	tempInstance := &instanceStore{}
	err = json.Unmarshal(data, tempInstance)
	if err == nil {
		count = len(tempInstance.Flakes)
	}
	return
}
func get(key []byte) (val []byte, err error) {
	return BitClient.Get(key)
}

func put(key []byte, val []byte) (err error) {
	return BitClient.Put(key, val)
}
func GetFlakeDef(stores FlakeStores) (flakes []definitions.FlakeDef) {
	var tmp []byte
	var err error
	for _, flake := range stores.GetFlakes() {
		tmp, err = get(keys["flake"](flake))
		if err == nil {
			var tempFlake definitions.FlakeDef
			err = json.Unmarshal(tmp, &tempFlake)
			if err != nil {
				continue
			}
			flakes = append(flakes, tempFlake)
		}
	}
	return
}
func GetAllServiceFlakes(id string) (flake []definitions.FlakeDef, err error) {
	var tmp []byte
	tmp, err = BitClient.Get(keys["service"](id))
	if err != nil {
		return
	}
	serviceMeta := serviceStore{}
	err = json.Unmarshal(tmp, &serviceMeta)
	if err != nil {
		return
	}
	flake = GetFlakeDef(&serviceMeta)
	return
}

func GetFlake(id string) (flake definitions.FlakeDef, err error) {
	var tmp []byte
	tmp, err = get(keys["flake"](id))
	if err != nil {
		return definitions.FlakeDef{}, err
	}
	err = json.Unmarshal(tmp, &flake)
	return
}

func GetAllInstanceFlakes(id string) (flake []definitions.FlakeDef, err error) {
	var tmp []byte
	keys["instance"](id)
	tmp, err = BitClient.Get(keys["instance"](id))
	if err != nil {
		return
	}
	instanceMeta := instanceStore{}
	err = json.Unmarshal(tmp, &instanceMeta)
	if err != nil {
		return
	}
	flake = GetFlakeDef(&instanceMeta)
	return
}
func NewFlake(flake definitions.FlakeDef) (err error) {
	var data []byte
	data, err = json.Marshal(flake)
	if err != nil {
		return
	}
	err = put(keys["flake"](flake.Id), data)
	if err != nil {
		return
	}
	data, err = get(keys["service"](flake.Service))
	tempService := &serviceStore{}
	if err != nil {
		if err == bitcask.ErrKeyNotFound {
			tempService.Id = flake.Service
		} else {
			return
		}
	} else {
		err = json.Unmarshal(data, tempService)
		if err != nil {
			return
		}
	}
	data, err = get(keys["instance"](flake.HostId))
	tempInstance := &instanceStore{}
	if err != nil {
		if err == bitcask.ErrKeyNotFound {
			tempInstance.Id = flake.HostId
		} else {
			return
		}
	} else {
		err = json.Unmarshal(data, tempInstance)
		if err != nil {
			return
		}
	}
	tempInstance.AddFlake(flake.Id)
	err = tempInstance.Save()
	if err != nil {
		return
	}
	tempService.AddFlake(flake.Id)
	err = tempService.Save()
	return
}
func DeleteFLake(flakeId string) (err error) {
	var flake definitions.FlakeDef
	flake, err = GetFlake(flakeId)
	if err != nil {
		return
	}
	var data []byte
	tempService := &serviceStore{}
	data, err = get(keys["service"](flake.Service))
	if err != nil {
		return
	}
	err = json.Unmarshal(data, tempService)
	if err != nil {
		return
	}
	data, err = get(keys["instance"](flake.HostId))
	if err != nil {
		return
	}
	tempInstance := &instanceStore{}
	err = json.Unmarshal(data, tempInstance)
	if err != nil {
		return
	}
	tempInstance.RemoveFlake(flake.Id)
	err = tempInstance.Save()
	if err != nil {
		return
	}
	tempService.RemoveFlake(flake.Id)
	err = tempService.Save()
	if err != nil {
		return
	}
	err = BitClient.Delete(keys["flake"](flakeId))
	return

}
