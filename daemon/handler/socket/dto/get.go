package dto

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/iij/legs-client/util"
	message "github.com/iij/legs-message"
)

// GetRequest represents get command request.
type GetRequest struct {
	SessionID string `json:"session_id"`
	Target    string `json:"target"`

	// SaveTo replesents the save file path.
	// If "/path/to/file" is specified, the target file is saved to /path/to/file.
	// If "/path/to/directory" is specified, the target file is saved to /path/to/directory + target URL filename.
	// (ex.) If SaveTo is specified ~/dir and target URL is http://example.com/test/file.txt, the target file is saved to ~/dir/file.txt.
	SaveTo  string     `json:"save_to"`
	Options GetOptions `json:"options"`
}

// GetOptions represents get command options.
type GetOptions struct {
	IsDryRun bool `json:"is_dry_run"`
}

// GetResponse represents get command response.
type GetResponse struct {
	Responses []message.TransferResponseData
}

// NewGetRequest returns GetRequest instance.
func NewGetRequest(target string, saveTo string, options GetOptions) *GetRequest {
	data := &GetRequest{
		SessionID: util.GenID(8),
		Target:    target,
		SaveTo:    saveTo,
		Options:   options,
	}
	return data
}

// ParseGetRequest parses []byte(JSON) data, and create DTO.
func ParseGetRequest(b []byte) (*GetRequest, error) {
	dto := &GetRequest{}
	decoder := json.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(dto)
	return dto, err
}

// ParseGetResponse : TODO
func ParseGetResponse(b []byte) (*GetResponse, error) {
	dto := &GetResponse{}
	decoder := json.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(dto)
	return dto, err
}

// ToBinary compiles the commandMessage object to messagepack.
func (s *GetRequest) ToBinary() ([]byte, error) {
	var buff bytes.Buffer
	encoder := json.NewEncoder(&buff)

	err := encoder.Encode(*s)
	bytes := buff.Bytes()

	return bytes, err
}

// ToBinary : TODO
func (s *GetResponse) ToBinary() ([]byte, error) {
	var buff bytes.Buffer
	encoder := json.NewEncoder(&buff)

	err := encoder.Encode(*s)
	bytes := buff.Bytes()

	return bytes, err
}

func (s *GetResponse) String() string {
	msg := ""
	for i, res := range s.Responses {
		msg += fmt.Sprintf("---Response[%d]---\n", i)
		msg += fmt.Sprintf("URL: %s\n", res.URL)
		msg += fmt.Sprintf("Status code: %d\n", res.StatusCode)
		msg += fmt.Sprintf("body: \n%s", res.Body)
	}
	return msg
}
