package golum_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/iij/legs-client/daemon/golum"
)

func TestResponse_readResponse(t *testing.T) {
	cases := map[string]struct {
		input []byte
		want  *golum.Response
		err   bool
	}{
		"valid code 0": {
			input: []byte(`{ "body": [97], "code": 0 }`),
			want:  &golum.Response{Body: []byte("a"), Code: golum.StatusSuccess},
			err:   false,
		},
		"valid code 1": {
			input: []byte(`{ "body": [97], "code": 1 }`),
			want:  &golum.Response{Body: []byte("a"), Code: golum.StatusError},
			err:   false,
		},
		"empty code": {
			input: []byte(`{ "body": [97] }`),
			want:  &golum.Response{Body: []byte("a")},
			err:   false,
		},
		"empty body": {
			input: []byte(`{ "code": 0}`),
			want:  &golum.Response{Code: golum.StatusSuccess},
			err:   false,
		},
		"empty": {
			input: []byte(""),
			want:  &golum.Response{},
			err:   true,
		},
		"nonsense": {
			input: []byte("hogehoge"),
			want:  &golum.Response{},
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
			get, err := golum.ReadResponse(client)
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

func TestResponse_ToBinary(t *testing.T) {
	cases := map[string]struct {
		input *golum.Response
		want  []byte
		err   bool
	}{
		"valid code 0": {
			input: &golum.Response{Body: []byte("a"), Code: golum.StatusSuccess},
			want:  append([]byte(`{"body":"YQ==","code":0}`), 10),
			err:   false,
		},
		"valid code 1": {
			input: &golum.Response{Body: []byte("a"), Code: golum.StatusError},
			want:  append([]byte(`{"body":"YQ==","code":1}`), 10),
			err:   false,
		},
		"empty code": {
			input: &golum.Response{Body: []byte("a")},
			want:  append([]byte(`{"body":"YQ==","code":0}`), 10),
			err:   false,
		},
		"empty body": {
			input: &golum.Response{Code: golum.StatusSuccess},
			want:  append([]byte(`{"body":null,"code":0}`), 10),
			err:   false,
		},
		"empty": {
			input: &golum.Response{},
			want:  append([]byte(`{"body":null,"code":0}`), 10),
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
