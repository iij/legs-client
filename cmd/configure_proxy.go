package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newConfigureProxyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy [proxy url with protocol like 'http']",
		Short: "Add 'proxy' parameter to config file.",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			err := conf.Save("proxy", args[0])
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newConfigureProxyCmd())
}
