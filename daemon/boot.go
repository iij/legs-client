package daemon

import (
	"flag"
	"fmt"
	"github.com/iij/legs-client/util"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/iij/legs-client/config"
	"github.com/iij/legs-client/daemon/context"
	signalHandler "github.com/iij/legs-client/daemon/handler/signal"
	legscLog "github.com/iij/legs-client/daemon/log"
	"github.com/iij/legs-client/daemon/model/status"
	"github.com/kr/pty"
	"github.com/sevlyar/go-daemon"
)

var (
	sig = flag.String("s", "", `send sig to the daemon
		stop — shutdown daemon
		reload — reloading the configuration file`)
)

// NonDaemonMarkValue marks process is not daemon
const NonDaemonMarkValue = "0"

// Boot boot legsc application.
// When foreground is false, boot as daemon program.
func Boot(configFile string, foreground bool) (pid int) {
	conf, _ := config.InitConfig(configFile)
	ctx := context.NewLegscContext(conf)
	err := createDirs(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if foreground {
		return bootForeground(ctx)
	}

	return bootDaemon(configFile, ctx)
}

func bootForeground(ctx *context.LegscContext) int {
	legscLog.InitLogger(true)

	legscLog.Info("legsc daemon stated.")
	ctx.Status.Daemon = status.Startted
	ctx.UpdateStatus()

	signal.Notify(ctx.Interrupt, os.Interrupt)

	initSocket(ctx)
	initConnection(ctx)

	legscLog.Info("legsc terminated.")
	ctx.Status.Daemon = status.Stopped
	ctx.UpdateStatus()

	return 0
}

func bootDaemon(configFile string, ctx *context.LegscContext) (pid int) {
	daemon.AddCommand(daemon.StringFlag(sig, "stop"), syscall.SIGKILL|syscall.SIGTERM|syscall.SIGQUIT, func(sig os.Signal) (err error) {
		signalHandler.HandleTerm(ctx.Interrupt, sig)
		return daemon.ErrStop
	})

	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	daemonContext := &daemon.Context{
		PidFileName: ctx.PidFileName,
		PidFilePerm: 0644,
		Umask:       027,
		WorkDir:     "./",
		Args:        []string{exe, "--daemon", "-c", configFile},
	}

	d, err := daemonContext.Reborn()
	if err != nil {
		log.Fatal("Unable to start daemon: ", err)
	}
	if d != nil {
		return d.Pid
	}
	defer daemonContext.Release()

	legscLog.InitLogger(false)

	legscLog.Info("--------------------")
	legscLog.Info("legsc daemon was started.")
	ctx.Status.Daemon = status.Startted
	ctx.UpdateStatus()

	// Responding to problems receiving SIGHUP when opening pty while running on OSX.
	// When SIGHUP invalidated, a permission error occurs when pty is opened for the first time.
	if runtime.GOOS == "darwin" {
		err := tryPtyOpen()
		if err != nil {
			legscLog.Error(err)
		}
	}

	serveSignal(ctx)
	initSocket(ctx)
	initConnection(ctx)

	legscLog.Info("legsc terminated.")
	ctx.Status.Daemon = status.Stopped
	ctx.UpdateStatus()

	return
}

func createDirs(ctx *context.LegscContext) (err error) {
	// create pid file directory
	if err = util.CreateDir(ctx.PidFileName); err != nil {
		return err
	}

	// create sock file directory
	if err = util.CreateDir(ctx.SockFileName); err != nil {
		return err
	}

	// create status file directory
	if err = util.CreateDir(ctx.StatusFileName); err != nil {
		return err
	}

	return
}

func tryPtyOpen() error {
	sigch := make(chan os.Signal, 1)
	defer close(sigch)
	signal.Notify(sigch, syscall.SIGHUP)
	defer signal.Reset(syscall.SIGHUP)

	p, t, err := pty.Open()
	if err != nil {
		return fmt.Errorf("failed to open pty: %s", err)
	}
	if err := p.Close(); err != nil {
		return fmt.Errorf("failed to close pty: %s", err)
	}
	if err := t.Close(); err != nil {
		return fmt.Errorf("failed to close tty: %s", err)
	}

	return nil
}
