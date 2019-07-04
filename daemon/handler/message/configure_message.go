package message

import (
	"fmt"

	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/log"
	message "github.com/iij/legs-message"
)

// HandleConfigureMessage handle a configure message receiving from server.
func HandleConfigureMessage(ctx *context.LegscContext, messageBytes []byte) {
	msg := &message.ClientConfigure{}
	err := message.Unmarshal(messageBytes, msg)
	if err != nil {
		log.Error("error in parsing message: ", err)
		return
	}
	log.Info(fmt.Sprintf("configure message: %+v", msg))

	switch msg.Model {
	case "client-config":
		HandleClientConfigMessage(ctx, msg)
	}
}
