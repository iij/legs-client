package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version info",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("legsc version:", AppVer)
			fmt.Println("commit hash:", Hash)
			fmt.Println("build date:", Builddate)
			fmt.Println("go version:", Goversion)
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(newVersionCmd())
}
