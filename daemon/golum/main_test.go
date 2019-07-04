package golum_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/iij/legs-client/daemon/golum"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

var testServer = map[string]func(){
	"echo": func() {
		s := &golum.Server{SocketName: filepath.Join("testdata", "test.sock")}
		golum.HandleFunc("hoge", func(w io.Writer, req *golum.Request) {
			golum.RenderResponse(w, req.Body)
		})
		s.ListenAndServe()
	},
}

var testClient = map[string]*golum.Client{
	"base":         &golum.Client{SocketName: filepath.Join("testdata", "test.sock")},
	"timeout":      &golum.Client{SocketName: filepath.Join("testdata", "test.sock"), Timeout: 1 * time.Nanosecond},
	"empty socket": &golum.Client{},
}

var testRawRequest = map[string][]byte{}
