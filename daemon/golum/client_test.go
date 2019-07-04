package golum_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/iij/legs-client/daemon/golum"
)

func TestClient_Do(t *testing.T) {
	cases := map[string]struct {
		server func()
		client *golum.Client
		input  *golum.Request
		want   *golum.Response
		err    bool
	}{
		"valid": {
			server: testServer["echo"],
			client: testClient["base"],
			input:  golum.NewRequest("hoge", []byte("hoge")),
			want:   &golum.Response{Body: []byte("hoge"), Code: golum.StatusSuccess},
			err:    false,
		},
		"timeout": {
			server: testServer["echo"],
			client: testClient["timeout"],
			input:  golum.NewRequest("hoge", []byte("hoge")),
			want:   nil,
			err:    true,
		},
		"empty socket": {
			server: testServer["echo"],
			client: testClient["empty socket"],
			input:  golum.NewRequest("hoge", []byte("hoge")),
			want:   nil,
			err:    true,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			go tt.server()
			for {
				if _, err := os.Stat("testdata/test.sock"); err == nil {
					time.Sleep(1 * time.Millisecond)
					break
				}
			}
			get, err := tt.client.Do(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if !reflect.DeepEqual(get, tt.want) {
				t.Fatalf("\n\tgot: %v\n\twant: %v", get, tt.want)
			}
		})
	}
}
