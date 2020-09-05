package health

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gookit/color"
	"github.com/rahultripathidev/docker-utility/datastore"
	"time"
)

var (
	Notifier <-chan *redis.Message
)

func InitSub() {
	ctx := context.Background()
	sub := datastore.RedisPubSub.PSubscribe(ctx, datastore.XerxesPubSub)
	_, err := sub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	Notifier = sub.Channel()
}

func StopTheWholeWorldAndReload() {
	color.Style{color.FgCyan, color.OpBold}.Printf("[%s] Reloading \n",time.Now().String())
	defer func() {
		<-time.After(30 * time.Second)
		isSchedulerRunning = false
	}()
	func() {
		select {
		case StopScheduler <- true:
			return
		default:
			return
		}
	}()
	isSchedulerRunning = true
	LoadMetaData()
}

func Reloader() {
	for range Notifier {
		go StopTheWholeWorldAndReload()
	}
}
