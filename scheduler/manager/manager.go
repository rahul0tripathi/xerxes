package manager

import (
	"context"
	"github.com/rahultripathidev/docker-utility/config"
	"github.com/rahultripathidev/docker-utility/datastore/bitcask"
	containerManager "github.com/rahultripathidev/docker-utility/proto"
	"github.com/rahultripathidev/docker-utility/scheduler/services"
	"github.com/rahultripathidev/docker-utility/service"
	"sync"
	"time"
)

type ContainerManagerService struct {
	containerManager.UnimplementedContainerManagerServer
}
type DatabaseLock struct {
	isOpen bool
	count  int
	timestamp time.Time
	lock sync.RWMutex
}
var (
	Database DatabaseLock
)
func (lock *DatabaseLock) Open() {
	lock.lock.Lock()
	defer lock.lock.Unlock()
	if !lock.isOpen {
		err := bitcask.InitClient()
		if err != nil && err == bitcask.ErrDatabaseLocked {
			lock.isOpen = true
			lock.count += 1
			lock.timestamp = time.Now()
			return
		}
		if err == nil {
			lock.isOpen = true
			lock.count += 1
			lock.timestamp = time.Now()
			return
		}
	}else{
		lock.count += 1
		lock.timestamp = time.Now()

	}
}
func (lock *DatabaseLock) Close() {
	lock.lock.Lock()
	defer lock.lock.Unlock()
	if lock.count > 0 {
		lock.count -= 1
	}
	if lock.count == 0 {
		//bitcask.BitClient.Flock.Unlock()
		bitcask.GracefulClose()
		lock.isOpen = false
		lock.count = 0
	}
	if lock.isOpen && time.Since(lock.timestamp).Seconds() > 30 {
		//bitcask.BitClient.Flock.Unlock()
		bitcask.GracefulClose()
		lock.isOpen = false
		lock.count = 0
	}
}
func init() {
	_ = config.LoadConfig.ServiceDef()
	_ = config.LoadConfig.Bit()
	Database = DatabaseLock{
		isOpen: false,
		count:  0,
	}
}

func (ContainerManagerService *ContainerManagerService) DeleteContainer(ctx context.Context, containerMeta *containerManager.DeleteContainerRequest) (response *containerManager.DeleteContainerResponse, err error) {
	go services.DeployScheduler(containerMeta.Meta.Timeout, containerMeta.Meta.XCid, containerMeta.Meta.XNodeId)
	return &containerManager.DeleteContainerResponse{Success: true}, nil
}
func (ContainerManagerService *ContainerManagerService) GetServices(ctx context.Context, serviceMetaRequest *containerManager.ServiceMetaRequest) (*containerManager.ServiceMetaResponse, error) {
	Database.Open()
	defer Database.Close()
	//err := bitcask.InitClient()
	//if err != nil {
	//	return nil, err
	//}
	//defer func() {
	//	bitcask.GracefulClose()
	//}()
	var metaData []*containerManager.ServiceMeta
	for serviceName, def := range config.ServicesDec.Def {
		metaData = append(metaData, &containerManager.ServiceMeta{
			Name:          serviceName,
			Image:         def.Image,
			ImageUri:      def.ImageUri,
			ContainerPort: def.ContainerPort,
			BasePort:      string(def.BasePort),
			MaxPort:       string(def.MaxPort),
			KongMeta: &containerManager.KongMeta{
				KongServiceMeta: &containerManager.KongServiceMeta{
					ServiceName: def.KongConf.Service.Name,
					Route:       def.KongConf.Service.Route,
					TargetPath:  def.KongConf.Service.TaregtPath,
				},
				KongUpstreamMeta: &containerManager.KongUpstreamMeta{
					UpstreamName: def.KongConf.Upstream.Name,
					Hashon:       def.KongConf.Upstream.Hashon,
				},
			},
		})
	}

	return &containerManager.ServiceMetaResponse{ServiceMeta: metaData}, nil
}

func (ContainerManagerService *ContainerManagerService) GetActiveServices(ctx context.Context, request *containerManager.FlakeMetaRequest) (*containerManager.FlakeMetaResponse, error) {
	Database.Open()
	defer Database.Close()
	var flakes []*containerManager.FlakeMeta
	for serviceName, _ := range config.ServicesDec.Def {
		flake, _ := bitcask.GetAllServiceFlakes(serviceName)
		for _, _flake := range flake {
			flakes = append(flakes, &containerManager.FlakeMeta{
				Id:          _flake.Id,
				ContainerId: _flake.ContainerId,
				HostId:      _flake.HostId,
				Service:     _flake.Service,
				Ip:          _flake.Ip,
				Port:        _flake.Port,
			})
		}
	}
	return &containerManager.FlakeMetaResponse{FlakeMeta: flakes}, nil
}
func (ContainerManagerService *ContainerManagerService) GetFlakeStats(ctx context.Context, request *containerManager.FlakeStatsRequest) (*containerManager.FlakeStatsResponse, error) {
	var flakeStats []*containerManager.FlakeStats
	Database.Open()
	defer Database.Close()
	for _, id := range request.Id {
		flake, err := bitcask.GetFlake(id)
		if err != nil {
			return nil, err
		}
		_stats := service.GetFlakeStatsOnce(flake)
		flakeStats = append(flakeStats, &containerManager.FlakeStats{
			Id:       id,
			MemUsage: _stats.MemUsage,
			CpuPer:   _stats.CpuPer,
			Network:  _stats.Network,
		})
	}
	return &containerManager.FlakeStatsResponse{
		FlakeStats: flakeStats,
	}, nil
}
func (ContainerManagerService *ContainerManagerService) GetFlakeLogs(ctx context.Context, request *containerManager.FlakeLogsRequest) (*containerManager.FlakeLogsResponse, error) {
	Database.Open()
	defer Database.Close()
	flake, err := bitcask.GetFlake(request.Id)
	if err != nil {
		return nil, err
	}
	logs, err := service.GetFlakeLogsOnce(flake, request.Tail)
	if err != nil {
		return nil, err
	}
	return &containerManager.FlakeLogsResponse{
		Logs: []byte(logs),
	}, nil
}
