package cmd

import (
	"fmt"
	"os"

	"github.com/iij/legs-client/config"
	"github.com/spf13/cobra"
)

var (
	// AppVer is application version by git tag
	AppVer string
	// Hash is commit hash when build
	Hash string
	// Builddate is date of build application
	Builddate string
	// Goversion is go version info in the building environment
	Goversion string
)

var rootCmd = &cobra.Command{
	Use:   "legsc",
	Short: "Legsc is a client for data transfer and remote command execution by websocket connection.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var configFile string
var conf *config.Config

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to config file")
}

// Execute execute the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	var err error
	conf, err = config.InitConfig(configFile)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
