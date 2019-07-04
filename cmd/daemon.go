package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"

	goDaemon "github.com/sevlyar/go-daemon"

	"github.com/iij/legs-client/daemon"
	"github.com/spf13/cobra"
)

type daemonOption struct {
	foreground bool
}

func newDaemonStartCmd() *cobra.Command {
	o := daemonOption{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start legsc daemon.",
		Run:   o.startDaemon,
	}
	cmd.Flags().BoolVarP(&o.foreground, "foreground", "f", false, "start legsc in foreground for debug")
	return cmd
}

func newDaemonStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop legsc daemon.",
		Run:   stopDaemon,
	}
	return cmd
}

func newDaemonRestartCmd() *cobra.Command {
	o := daemonOption{}
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart legsc daemon. This command only uses the stop/start commands.",
		Run: func(c *cobra.Command, args []string) {
			stopDaemon(c, args)
			time.Sleep(3 * time.Second)
			o.startDaemon(c, args)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newDaemonStartCmd())
	rootCmd.AddCommand(newDaemonStopCmd())
	rootCmd.AddCommand(newDaemonRestartCmd())

}

func (o *daemonOption) startDaemon(c *cobra.Command, args []string) {
	_ = os.Setenv(goDaemon.MARK_NAME, daemon.NonDaemonMarkValue)
	pid := daemon.Boot(configFile, o.foreground)
	fmt.Println("pid:", pid)
}

func stopDaemon(c *cobra.Command, args []string) {
	value, err := ioutil.ReadFile(conf.GetString("pid_file"))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "cannot open pid file:", err)
		return
	}
	pid, _ := strconv.ParseInt(string(value), 10, 32)

	p, err := os.FindProcess(int(pid))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "cannot fined process:", err)
		return
	}
	p.Signal(syscall.SIGTERM)
	fmt.Println("send term signal to legsc daemon.")
}
