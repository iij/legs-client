package golum_test

import (
	"bufio"
	"net"
	"reflect"
	"testing"

	"github.com/iij/legs-client/daemon/golum"
)

func TestRequest_ReadRequest(t *testing.T) {
	cases := map[string]struct {
		input []byte
		want  *golum.Request
		err   bool
	}{
		"valid": {
			input: []byte(`{ "type": "hoge", "body": [97]}`),
			want:  &golum.Request{Type: "hoge", Body: []byte("a")},
			err:   false,
		},
		"empty body": {
			input: []byte(`{ "type": "hoge"}`),
			want:  &golum.Request{Type: "hoge"},
			err:   false,
		},
		"empty type": {
			input: []byte(`{ "body": [97]}`),
			want:  &golum.Request{Body: []byte("a")},
			err:   false,
		},
		"empty": {
			input: []byte(""),
			want:  &golum.Request{},
			err:   true,
		},
		"nonsense": {
			input: []byte("hogehoge"),
			want:  &golum.Request{},
			err:   true,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			server, client := net.Pipe()
			go func() {
				server.Write(tt.input)
				server.Close()
			}()
			get, err := golum.ReadRequest(bufio.NewReader(client))
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

func TestRequest_ToBinary(t *testing.T) {
	cases := map[string]struct {
		input *golum.Request
		want  []byte
		err   bool
	}{
		"valid": {
			input: &golum.Request{Type: "hoge", Body: []byte("a")},
			want:  append([]byte(`{"type":"hoge","body":"YQ=="}`), 10),
			err:   false,
		},
		"empty body": {
			input: &golum.Request{Type: "hoge"},
			want:  append([]byte(`{"type":"hoge","body":null}`), 10),
			err:   false,
		},
		"empty type": {
			input: &golum.Request{Body: []byte("a")},
			want:  append([]byte(`{"type":"","body":"YQ=="}`), 10),
			err:   false,
		},
		"empty": {
			input: &golum.Request{},
			want:  append([]byte(`{"type":"","body":null}`), 10),
			err:   false,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			get, err := tt.input.ToBinary()
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
