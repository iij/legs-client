package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/iij/legs-client/daemon/model"
	"github.com/iij/legs-client/daemon/model/status"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get legsc status.",
		Run: func(c *cobra.Command, args []string) {
			stat := getStatus(conf.GetString("status_file"))
			fmt.Println("device name:", conf.GetString("device_name"))
			fmt.Println("device id:", stat.DeviceID)
			fmt.Println("daemon status:", stat.Daemon)
			fmt.Println("server connection:", stat.Conn)
			if stat.Conn == status.Disconnected || stat.Daemon == status.Stopped {
				os.Exit(1)
			}
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newStatusCmd())
}

func getStatus(statusFile string) model.Status {
	bytes, err := ioutil.ReadFile(statusFile)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error", err)
		os.Exit(1)
	}

	var status model.Status

	if err = json.Unmarshal(bytes, &status); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error", err)
		os.Exit(1)
	}

	return status
}
