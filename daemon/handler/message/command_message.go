package message

import (
	"errors"
	"os/exec"
	"time"

	message "github.com/iij/legs-message"

	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/log"
	shellwords "github.com/mattn/go-shellwords"
)

// HandleCommandMessage handle a command message receiving from server.
func HandleCommandMessage(ctx *context.LegscContext, msgBytes []byte) {
	msg := &message.Command{}
	err := message.Unmarshal(msgBytes, msg)
	if err != nil {
		log.Error("failed in parsing message", err)
		sendError(err, msg, ctx)
		return
	}
	execData := msg.Data

	log.Info("start command execution. cmd:", execData.Command, "exec_id:", execData.ID)

	cmd, err := genCmd(execData.Command)
	if err != nil {
		log.Error("faild in parsing command", err)
		sendError(err, msg, ctx)
		return
	}

	cmdResult := make(chan string, 1)
	go func() {
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error("cmd execution error:", err)
			cmdResult <- err.Error()
		}
		cmdResult <- string(out)
	}()

	select {
	case outStr := <-cmdResult:
		log.Info("command executed. cmd:", execData.Command, "exec_id:", execData.ID)
		msg.SetStatus("executed")
		msg.SetResult(outStr)

		ctx.SendMessage(msg)
	case <-time.After(60 * time.Second):
		log.Info("timeout in command execution")
		msg.SetStatus("error")
		msg.SetResult("command timeout")

		ctx.SendMessage(msg)
	}
}

func sendError(err error, msg *message.Command, ctx *context.LegscContext) {
	msg.SetStatus("executed")
	msg.SetResult(err.Error())

	ctx.SendMessage(msg)
}

func genCmd(cmdStr string) (cmd *exec.Cmd, err error) {
	c, err := shellwords.Parse(cmdStr)
	if err != nil {
		err = errors.New("empty command")
		return
	}

	switch len(c) {
	case 0:
		err = errors.New("empty command")
		return
	case 1:
		cmd = exec.Command(c[0])
	default:
		cmd = exec.Command(c[0], c[1:]...)
	}

	return
}
