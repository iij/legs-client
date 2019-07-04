package daemon

import (
	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/golum"
	"github.com/iij/legs-client/daemon/handler/socket"
)

func initSocket(ctx *context.LegscContext) {
	s := golum.Server{SocketName: ctx.SockFileName}

	sh := socket.NewSendHandler(ctx)
	golum.HandleFunc("send", sh.Send)
	gh := socket.NewGetHandler(ctx)
	golum.HandleFunc("get", gh.Get)
	go s.ListenAndServe()
}
