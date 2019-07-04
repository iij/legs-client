package message

import (
	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/handler/socket"
	"github.com/iij/legs-client/daemon/log"
	message "github.com/iij/legs-message"
)

// HandleTransferMessage : TODO
func HandleTransferMessage(ctx *context.LegscContext, msgBytes []byte) {
	msg := message.Transfer{}
	err := message.Unmarshal(msgBytes, &msg)
	if err != nil {
		log.Error("failed in parsing message:", err)
		return
	}

	switch msg.GetModel() {
	case "http_transfer":
		handleTransferResponse(ctx, msg)
	}
}

func handleTransferResponse(ctx *context.LegscContext, msg message.Transfer) {
	socket.Response(msg)
}
