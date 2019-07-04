package main

import (
	"os"

	goDaemon "github.com/sevlyar/go-daemon"

	"github.com/iij/legs-client/cmd"
	"github.com/iij/legs-client/daemon"
)

var (
	version   = "dev"
	hash      = "none"
	builddate = "none"
	goversion = "none"
)

func main() {
	cmd.AppVer = version
	cmd.Hash = hash
	cmd.Builddate = builddate
	cmd.Goversion = goversion

	if isDaemon() {
		_ = os.Setenv(goDaemon.MARK_NAME, goDaemon.MARK_VALUE)
		daemon.Boot(os.Args[3], false)
	} else {
		cmd.Execute()
	}
}

func isDaemon() bool {
	if len(os.Args) < 4 {
		return false
	}

	return os.Args[1] == "--daemon"
}
