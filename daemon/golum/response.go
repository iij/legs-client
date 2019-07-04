package golum

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
)

// RenderError : TODO
func RenderError(w io.Writer, err error) {
	res := Response{
		Body: []byte(err.Error()),
		Code: StatusError,
	}
	b, err := res.ToBinary()
	if err != nil {
		return
	}
	w.Write(b)
}

// RenderResponse : TODO
func RenderResponse(w io.Writer, body []byte) {
	res := Response{
		Body: body,
		Code: StatusSuccess,
	}
	b, err := res.ToBinary()
	if err != nil {
		return
	}
	w.Write(b)
}

// Response : TODO
type Response struct {
	Body []byte `json:"body"`
	Code int    `json:"code"`
}

// ReadResponse parses []byte(JSON) data by json.Decoder, and transforms to the Response.
func readResponse(conn net.Conn) (*Response, error) {
	b := make([]byte, 0)
	for {
		buf := make([]byte, 512)
		nr, err := conn.Read(buf)
		if err != nil {
			break
		}
		buf = buf[:nr]
		b = append(b, buf...)
	}
	res := &Response{}
	decoder := json.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(res)
	return res, err
}

// ToBinary encodes the Response to []byte.
func (r *Response) ToBinary() ([]byte, error) {
	var buff bytes.Buffer
	encoder := json.NewEncoder(&buff)

	err := encoder.Encode(*r)
	bytes := buff.Bytes()
	return bytes, err
}
