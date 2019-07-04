package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newExportConfCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export current configuration to file",
		Run: func(c *cobra.Command, args []string) {
			err := conf.Write()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println("export to", conf.ConfigFileUsed())
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newExportConfCmd())
}
