package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iij/legs-client/util"
	message "github.com/iij/legs-message"
)

// SendRequest represents send command request.
type SendRequest struct {
	SessionID string `json:"session_id"`
	Target    string `json:"target"`
	Body      string `json:"body"`
}

// SendResponse represents send command response.
type SendResponse struct {
	Responses []message.TransferResponseData
}

// NewSendRequest returns sendData instance.
func NewSendRequest(target string, body []string) *SendRequest {
	data := &SendRequest{
		SessionID: util.GenID(8),
		Target:    target,
	}
	if body != nil {
		data.Body = strings.Join(body, " ")
	}
	return data
}

// ParseSendRequest parses []byte(JSON) data, and create DTO.
func ParseSendRequest(b []byte) (*SendRequest, error) {
	dto := &SendRequest{}
	decoder := json.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(dto)
	return dto, err
}

// ParseSendResponse : TODO
func ParseSendResponse(b []byte) (*SendResponse, error) {
	dto := &SendResponse{}
	decoder := json.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(dto)
	return dto, err

}

// ToBinary compiles the commandMessage object to messagepack.
func (s *SendRequest) ToBinary() ([]byte, error) {
	var buff bytes.Buffer
	encoder := json.NewEncoder(&buff)

	err := encoder.Encode(*s)
	bytes := buff.Bytes()

	return bytes, err
}

// ToBinary : TODO
func (s *SendResponse) ToBinary() ([]byte, error) {
	var buff bytes.Buffer
	encoder := json.NewEncoder(&buff)

	err := encoder.Encode(*s)
	bytes := buff.Bytes()

	return bytes, err
}

func (s *SendResponse) String() string {
	msg := ""
	for i, res := range s.Responses {
		msg += fmt.Sprintf("---Response[%d]---\n", i)
		msg += fmt.Sprintf("URL: %s\n", res.URL)
		msg += fmt.Sprintf("Status code: %d\n", res.StatusCode)
		msg += fmt.Sprintf("body: \n%s", res.Body)
	}
	return msg
}
