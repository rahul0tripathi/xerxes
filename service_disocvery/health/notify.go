package health

import (
	"github.com/gookit/color"
	"time"
)


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

