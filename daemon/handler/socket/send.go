package socket

import (
	"errors"
	"io"
	"time"

	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/golum"
	"github.com/iij/legs-client/daemon/handler/socket/dto"
	"github.com/iij/legs-client/daemon/log"
	message "github.com/iij/legs-message"
)

// SendHandler : TODo
type SendHandler struct {
	ctx *context.LegscContext
}

// NewSendHandler : TODO
func NewSendHandler(ctx *context.LegscContext) *SendHandler {
	return &SendHandler{
		ctx: ctx,
	}
}

// Send : TODO
func (s *SendHandler) Send(w io.Writer, r *golum.Request) {
	data, err := dto.ParseSendRequest(r.Body)
	if err != nil {
		log.Error("failed to parse request")
		golum.RenderError(w, err)
		return
	}

	action := message.TransferActionPost

	msgData := message.TransferRequestData{
		Action: action,
		Target: data.Target,
		Value:  data.Body,
	}
	msg := message.NewTransferRequest(data.SessionID, msgData)
	s.ctx.SendMessage(msg)

	s.handleResponse(data.SessionID, w)
}

func (s *SendHandler) handleResponse(resID string, w io.Writer) {
	timer := time.NewTimer(20 * time.Second)
	resCh := make(chan []message.TransferResponseData, 1)
	setCh(resID, resCh)
	for {
		select {
		case res := <-resCh:
			success := true
			for _, r := range res {
				if r.StatusCode < 200 || r.StatusCode >= 300 {
					success = false
				}
			}

			data := dto.SendResponse{Responses: res}
			b, err := data.ToBinary()
			if err != nil {
				log.Error("failed to parse request")
				golum.RenderError(w, errors.New("failed to parse request"))
				return
			}

			if success {
				golum.RenderResponse(w, b)
			} else {
				golum.RenderError(w, errors.New(string(b)))
			}
			return
		case <-timer.C:
			log.Error("request timeout")
			golum.RenderError(w, NewTimeoutError())
			return
		}
	}
}
