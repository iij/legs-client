package golum

import (
	"net"
	"time"
)

// Client is a golum client.
type Client struct {
	// SocketName specifies the UNIX domain socket file name.
	SocketName string
	Timeout    time.Duration
}

// Do sends a golum request and returns a golum response.
func (c *Client) Do(req *Request) (*Response, error) {
	dialer := &net.Dialer{
		Timeout: c.Timeout,
	}
	conn, err := dialer.Dial("unix", c.SocketName)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	b, err := req.ToBinary()
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(b); err != nil {
		return nil, err
	}
	err = conn.(*net.UnixConn).CloseWrite()
	if err != nil {
		return nil, err
	}

	res, err := readResponse(conn)
	if err != nil {
		return nil, err
	}

	return res, nil
}
