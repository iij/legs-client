package message

import (
	"os"
	"os/exec"

	message "github.com/iij/legs-message"

	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/log"
	"github.com/kr/pty"
)

// HandleConsoleMessage handles console message from server
func HandleConsoleMessage(ctx *context.LegscContext, msgBytes []byte) {
	msg := &message.Console{}
	err := message.Unmarshal(msgBytes, msg)
	if err != nil {
		log.Error("failed in parsing message: ", err)
		sendConsoleError(err, msg, ctx)
		return
	}

	switch msg.State {
	case message.ConsoleStartState:
		handleStartMessage(ctx, msg)
	case message.ConsoleCloseState:
		handleCloseMessage(ctx, msg)
	case message.ConsoleInputState:
		handleInputMessage(ctx, msg)
	}
}

func handleStartMessage(ctx *context.LegscContext, msg *message.Console) {
	log.Info("start console...")

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	log.Info("target shell:", shell)

	c := exec.Command(shell)
	ptmx, err := pty.Start(c)
	if err != nil {
		log.Error("can not open pty", err)
		sendConsoleError(err, msg, ctx)
		return
	}
	session := context.ConsoleSession{
		Fd:     ptmx,
		Closer: make(chan struct{}),
		Closed: make(chan struct{}),
	}
	ctx.ConsoleSessions[msg.SessionID] = session

	go readConsole(ctx, msg.SessionID)

	_ = c.Wait()
}

func handleCloseMessage(ctx *context.LegscContext, msg *message.Console) {
	closeConsole(ctx, msg.SessionID)
}

func handleInputMessage(ctx *context.LegscContext, msg *message.Console) {
	sessionID := msg.SessionID

	_, err := ctx.ConsoleSessions[sessionID].Fd.Write(msg.Data)
	if err != nil {
		log.Error("can not write console:", err)
		sendConsoleError(err, msg, ctx)
		closeConsole(ctx, sessionID)
	}
}

func sendConsoleError(err error, msg *message.Console, ctx *context.LegscContext) {
	emsg := message.NewConsoleClose(msg.SessionID)
	sendConsoleMessage(emsg, ctx)
}

func sendConsoleMessage(msg message.Message, ctx *context.LegscContext) {
	ctx.SendMessage(msg)
}

func closeConsole(ctx *context.LegscContext, sessionID string) {
	session, ok := ctx.ConsoleSessions[sessionID]
	if !ok {
		return
	}

	close(session.Closer)

	_, err := session.Fd.Write([]byte("closing..."))
	if err != nil {
		log.Error("err in write to fd:", err)
	}

	<-session.Closed

	delete(ctx.ConsoleSessions, sessionID)
	log.Info("console closed...")
}

func readConsole(ctx *context.LegscContext, sessionID string) {
	errOccurred := false
	session := ctx.ConsoleSessions[sessionID]

	defer func() {
		if err := session.Fd.Close(); err != nil {
			log.Error("failed in closing fd:", err)
		}

		close(session.Closed)

		if errOccurred {
			closeConsole(ctx, sessionID)
		}
	}()

	for {
		select {
		case <-session.Closer:
			return
		default:
			buf := make([]byte, 512)
			nr, err := session.Fd.Read(buf)
			if err != nil {
				log.Error("error in read console", err)
				errOccurred = true
				return
			}
			data := buf[0:nr]
			msg := message.NewConsoleOutput(sessionID, data)
			sendConsoleMessage(msg, ctx)
		}
	}
}
