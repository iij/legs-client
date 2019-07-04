package signal

import (
	"os"

	"github.com/iij/legs-client/daemon/log"
)

// HandleTerm handle the interrupt signal like Ctrl-c.
func HandleTerm(interrupt chan os.Signal, sig os.Signal) {
	log.Info("terminationg legsc daemon...")
	interrupt <- sig
}
