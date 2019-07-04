package golum

import (
	"bufio"
	"bytes"
	"encoding/json"
)

// Request represents a golum request.
type Request struct {
	Type string `json:"type"`
	Body []byte `json:"body"`
}

// NewRequest returns a new Request given a type and body.
func NewRequest(typ string, data []byte) *Request {
	return &Request{
		Type: typ,
		Body: data,
	}
}

// ReadRequest parses []byte data by json.Decoder, and transforms to the Request.
func readRequest(b *bufio.Reader) (*Request, error) {
	req := &Request{}
	decoder := json.NewDecoder(b)
	err := decoder.Decode(req)
	return req, err
}

// ToBinary encodes the Request to []byte.
func (r *Request) ToBinary() ([]byte, error) {
	var buff bytes.Buffer
	encoder := json.NewEncoder(&buff)

	err := encoder.Encode(*r)
	bytes := buff.Bytes()
	return bytes, err
}
