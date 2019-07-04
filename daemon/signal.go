package daemon

import (
	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/log"
	daemon "github.com/sevlyar/go-daemon"
)

func serveSignal(ctx *context.LegscContext) {
	go func() {
		err := daemon.ServeSignals()
		if err != nil {
			log.Error("Error:", err)
		}
	}()
}
