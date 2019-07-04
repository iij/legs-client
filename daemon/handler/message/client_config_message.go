package message

import (
	"errors"
	"time"

	"github.com/iij/legs-client/daemon/log"

	message "github.com/iij/legs-message"

	"github.com/iij/legs-client/daemon/context"
)

// HandleClientConfigMessage handling message of client-config from server
func HandleClientConfigMessage(ctx *context.LegscContext, msg *message.ClientConfigure) {
	if msg == nil {
		return
	}

	pingInterval, err := getInt(msg.Data.PingInterval)
	if err != nil {
		log.Error(err)
		return
	}

	deviceID := msg.Data.DeviceID
	log.Info("assigned device-id:", deviceID)
	log.Info("set ping interval:", pingInterval)

	ctx.Status.DeviceID = deviceID
	if err = ctx.UpdateStatus(); err != nil {
		log.Error("failed to write status file:", err.Error())
	}

	needRestart := setPingInterval(ctx, pingInterval)

	if needRestart {
		ctx.Restart <- struct{}{}
	}
}

func setPingInterval(ctx *context.LegscContext, pingInterval int) (needRestart bool) {
	interval := time.Duration(pingInterval) * time.Second
	if ctx.PingInterval != interval {
		ctx.PingInterval = interval
		needRestart = true
	}

	return
}

func getInt(val interface{}) (int, error) {
	switch i := val.(type) {
	case int:
		return i, nil
	case int8:
		return int(i), nil
	case int16:
		return int(i), nil
	case int32:
		return int(i), nil
	case int64:
		return int(i), nil
	default:
		return 0, errors.New("unknown type")
	}
}
