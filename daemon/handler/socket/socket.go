package socket

import (
	"fmt"
	"sync"

	message "github.com/iij/legs-message"
)

var resChs = map[string]chan []message.TransferResponseData{}
var resChLock = sync.RWMutex{}

// Response : TODO
func Response(msg message.Transfer) {
	resCh, ok := getCh(msg.Data.SessionID)
	if !ok {
		return
	}
	resCh <- msg.Data.Responses
}

func setCh(id string, c chan []message.TransferResponseData) {
	resChLock.Lock()
	defer resChLock.Unlock()

	resChs[id] = c
}

func getCh(id string) (chan []message.TransferResponseData, bool) {
	resChLock.RLock()
	defer resChLock.RUnlock()

	c, ok := resChs[id]
	return c, ok
}

func deleteCh(id string) {
	resChLock.Lock()
	defer resChLock.Unlock()

	close(resChs[id])
	delete(resChs, id)
}

// TimeoutError : TODO
type TimeoutError struct {
	msg string
}

func (e *TimeoutError) Error() string {
	return e.msg
}

// NewTimeoutError : TODO
func NewTimeoutError() *TimeoutError {
	return &TimeoutError{
		msg: "timeout",
	}
}

// RoutingNotFoundError : TODO
type RoutingNotFoundError struct {
	msg string
}

func (e *RoutingNotFoundError) Error() string {
	return e.msg
}

// NewRoutingNotFoundError : TODO
func NewRoutingNotFoundError(routingName string) *RoutingNotFoundError {
	return &RoutingNotFoundError{
		msg: fmt.Sprintf("routing: %s is not found", routingName),
	}
}
