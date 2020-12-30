package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-units"
	"github.com/olekukonko/tablewriter"
	dockerClient "github.com/rahultripathidev/docker-utility/client"
	definitions "github.com/rahultripathidev/docker-utility/types"
	"io"
	"os"
	"os/exec"
	"time"
)

var (
	_previousCPU    uint64 = 0.0
	_previousSystem uint64 = 0.0
)

func GetFlakeLogs(flake definitions.FlakeDef,tail string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dockerDaemon, err := dockerClient.NewDockerClient(flake.HostId)
	if err != nil {
		return err
	}
	reader, err := dockerDaemon.ContainerLogs(ctx, flake.ContainerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail: tail,
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
		onlineCPUs  = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	_previousCPU = v.CPUStats.CPUUsage.TotalUsage
	_previousSystem = v.CPUStats.SystemUsage
	return cpuPercent
}
func calculateNetwork(network map[string]types.NetworkStats) string {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return fmt.Sprintf("%s / %s", units.HumanSizeWithPrecision(rx, 3), units.HumanSizeWithPrecision(tx, 3))
}
func GetFlakeStats(flake definitions.FlakeDef) (err error) {
	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()
	dockerDaemon, err := dockerClient.NewDockerClient(flake.HostId)
	if err != nil {
		return err
	}
	stats, err := dockerDaemon.ContainerStats(ctx, flake.ContainerId, true)
	if err != nil {
		return err
	}
	var p = make([]byte, 5024)
	var n int
	var _stats types.StatsJSON
	for {
		n, err = stats.Body.Read(p)
		if err == io.EOF {
			break
		}
		err = json.Unmarshal(p[:n], &_stats)
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		err = c.Run()
		if err != nil {
			break
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Id", "MEM USAGE / LIMIT ", "CPU%", "NET I/O", "PID"})
		table.Append([]string{_stats.Name, fmt.Sprintf("%s / %s  ", units.BytesSize(float64(_stats.MemoryStats.Usage)), units.BytesSize(float64(_stats.MemoryStats.Limit))), fmt.Sprintf("%.2f", calculateCPUPercentUnix(_previousCPU, _previousSystem, &_stats)), calculateNetwork(_stats.Networks), fmt.Sprintf("%d", _stats.PidsStats.Current)})
		table.Render()

	}
	return nil
}
