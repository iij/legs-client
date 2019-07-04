package message

import (
	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/log"
	message "github.com/iij/legs-message"
)

// HandleMessage handle a message from server.
// This method routing to message handlers by message types.
func HandleMessage(ctx *context.LegscContext, messageBytes []byte) {
	msg := &message.BaseMessage{}
	err := message.Unmarshal(messageBytes, msg)
	if err != nil {
		log.Error("error in parsing message: ", err)
		return
	}

	switch msg.GetMessageType() {
	case "configure":
		HandleConfigureMessage(ctx, messageBytes)
	case "execute":
		switch msg.GetModel() {
		case "command":
			HandleCommandMessage(ctx, messageBytes)
		}
	case "console":
		HandleConsoleMessage(ctx, messageBytes)
	case "proxy":
		HandleProxyMessage(ctx, messageBytes)
	case "transfer":
		HandleTransferMessage(ctx, messageBytes)
	}
}
