package model

import "testing"
import "github.com/stretchr/testify/assert"

type params map[string]string

func TestComparePath(t *testing.T) {
	tests := []struct {
		input      string
		target     string
		wantOK     bool
		wantParams map[string]string
	}{
		{input: "", target: "hoge/path", wantOK: false, wantParams: params{}},
		{input: "/", target: "hoge/path", wantOK: false, wantParams: params{}},
		{input: "hoge", target: "hoge", wantOK: true, wantParams: params{}},
		{input: "hoge", target: "/hoge", wantOK: true, wantParams: params{}},
		{input: "/hoge", target: "hoge", wantOK: true, wantParams: params{}},
		{input: "test", target: ":hoge", wantOK: true, wantParams: params{"hoge": "test"}},
		{input: "test", target: "/:hoge", wantOK: true, wantParams: params{"hoge": "test"}},
		{input: "/hoge/path", target: "hoge/path", wantOK: true, wantParams: params{}},
		{input: "/hoge", target: "hoge/path", wantOK: false, wantParams: params{}},
		{input: "/hoge/path/to", target: "hoge/path", wantOK: false, wantParams: params{}},
		{input: "/hoge/create/path", target: "hoge/:action/path", wantOK: true, wantParams: params{"action": "create"}},
		{input: "/hoge/1/path", target: "hoge/:id/path", wantOK: true, wantParams: params{"id": "1"}},
		{input: "/hoge/1/huga", target: "hoge/:id/path", wantOK: false, wantParams: params{}},
		{input: "/hoge/update/1", target: "hoge/update/:id", wantOK: true, wantParams: params{"id": "1"}},
		{input: "/hoge/update/1", target: "hoge/update2/:id", wantOK: false, wantParams: params{}},
		{input: "/hoge/update/1", target: "hoge/update/:id-test", wantOK: true, wantParams: params{"id-test": "1"}},
		{input: "/hoge/update/1", target: "hoge/update/:id/huga", wantOK: false, wantParams: params{}},
		{input: "/hoge/update/1", target: "hoge/update/:id:test", wantOK: true, wantParams: params{"id:test": "1"}},
		{input: "/hoge/update/1/test", target: "hoge/:action/:id/test", wantOK: true, wantParams: params{"action": "update", "id": "1"}},
	}

	for _, test := range tests {
		routing := &Routing{
			Name: test.target,
		}
		ok, param := routing.ComparePath(test.input)
		assert.Equal(t, test.wantOK, ok)
		assert.Equal(t, test.wantParams, param)
	}
}
