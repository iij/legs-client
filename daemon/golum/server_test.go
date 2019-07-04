package golum_test

import (
	"io"
	"net"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/iij/legs-client/daemon/golum"
)

var testServeMux = map[string]*golum.ServeMux{
	"default": golum.DefaultServeMux,
	"empty":   &golum.ServeMux{},
}

func TestServeMux_ServeSocket(t *testing.T) {
	cases := map[string]struct {
		mux   *golum.ServeMux
		input *golum.Request
		want  *golum.Response
	}{
		"empty": {
			mux:   testServeMux["empty"],
			input: &golum.Request{Type: "hoge"},
			want:  &golum.Response{Body: []byte("not found"), Code: 1},
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			server, client := net.Pipe()
			go func() {
				tt.mux.ServeSocket(server, tt.input)
				server.Close()
			}()
			get, _ := golum.ReadResponse(client)
			if !reflect.DeepEqual(get, tt.want) {
				t.Fatalf("\n\tgot: %v\n\twant: %v", get, tt.want)
			}
		})
	}
}

var testHandler = map[string]golum.HandlerFunc{
	"echo": golum.HandlerFunc(func(w io.Writer, r *golum.Request) {
		golum.RenderResponse(w, r.Body)
	}),
}

func TestServeMux_HandleAndHandler(t *testing.T) {
	type input struct {
		typ     string
		handler golum.Handler
	}

	cases := map[string]struct {
		input      input
		isNotFound bool
	}{
		"valid": {
			input: input{
				typ:     "hoge",
				handler: testHandler["echo"],
			},
			isNotFound: false,
		},
		"empty type": {
			input: input{
				typ:     "",
				handler: testHandler["echo"],
			},
			isNotFound: true,
		},
		"empty handler": {
			input: input{
				typ:     "hoge",
				handler: nil,
			},
			isNotFound: true,
		},
	}

	for caseName, tt := range cases {
		mux := golum.ServeMux{}
		t.Run(caseName, func(t *testing.T) {
			mux.Handle(tt.input.typ, tt.input.handler)
			h, _ := mux.Handler(&golum.Request{Type: tt.input.typ})

			// mux.Handler returns the NotFound handler if no handler is registered for the specified type.
			// Here we get the name of the handler function.
			hv := reflect.ValueOf(h)
			get := strings.TrimPrefix(path.Base(runtime.FuncForPC(hv.Pointer()).Name()), "golum.")

			if (get == "NotFound") != tt.isNotFound {
				t.Fatalf("\n\tgot: %v\n\twant: NotFound", get)
			}
		})
	}
}
