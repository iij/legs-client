package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newConfigureSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret [secret key to configure]",
		Short: "Add 'secret' parameter to config file. 'secret' is a parameter for authorization to server.",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			err := conf.Save("secret", args[0])
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newConfigureSecretCmd())
}
