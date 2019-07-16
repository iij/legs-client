package socket_test

import (
	"regexp"
	"testing"

	"github.com/iij/legs-client/daemon/handler/socket"
)

func TestGet_makeSavePath(t *testing.T) {
	cases := map[string]struct {
		url    string
		saveTo string
		want   string
		err    bool
	}{
		"saveTo is a directory": {
			url:    "http://example.com/test.txt",
			saveTo: "../../golum/testdata/",
			want:   "/.*/testdata/test.txt",
			err:    false,
		},
		"saveTo is a file": {
			url:    "http://example.com/test.txt",
			saveTo: "../../golum/testdata/test.txt",
			want:   "/.*/testdata/test.txt",
			err:    false,
		},
		"URL pathname is dir/filename": {
			url:    "http://example.com/test/test.txt",
			saveTo: "../../golum/testdata/",
			want:   "/.*/testdata/test.txt",
			err:    false,
		},
		"URL pathname is omitted": {
			url:    "http://example.com",
			saveTo: "../../golum/testdata/",
			want:   "/.*/testdata",
			err:    false,
		},
	}

	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			get, err := socket.GmakeSavePath(tt.url, tt.saveTo)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if matched, _ := regexp.MatchString(tt.want, get); !matched {
				t.Fatalf("\ngot : %v\nwant(regexp): %v", get, tt.want)
			}
		})
	}
}
