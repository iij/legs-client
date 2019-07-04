package context

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"

	"github.com/iij/legs-client/config"
	"github.com/iij/legs-client/daemon/log"
	"github.com/iij/legs-client/daemon/model"
	"github.com/gorilla/websocket"
	message "github.com/iij/legs-message"
)

// LegscContext is global variables for through all over application.
type LegscContext struct {
	PidFileName     string
	SockFileName    string
	StatusFileName  string
	Server          string
	Proxy           string
	DeviceName      string
	Secret          string
	UseWss          bool
	Status          *model.Status
	PingInterval    time.Duration
	Interrupt       chan os.Signal
	Restart         chan interface{}
	SocketListener  net.Listener
	Connection      *websocket.Conn
	ConsoleSessions map[string]ConsoleSession
}

// NewLegscContext return LegscContext instance made by configs.
func NewLegscContext(conf *config.Config) *LegscContext {
	return &LegscContext{
		PidFileName:     conf.GetString("pid_file"),
		SockFileName:    conf.GetString("sock_file"),
		StatusFileName:  conf.GetString("status_file"),
		Server:          conf.GetString("server"),
		Proxy:           conf.GetString("proxy"),
		DeviceName:      conf.GetString("device_name"),
		Secret:          conf.GetString("secret"),
		UseWss:          conf.GetBool("use_wss"),
		Status:          &model.Status{},
		PingInterval:    30 * time.Second,
		Interrupt:       make(chan os.Signal),
		Restart:         make(chan interface{}),
		ConsoleSessions: map[string]ConsoleSession{},
	}
}

var connMutex = sync.Mutex{}

// SendMessage sends message to server.
// It is not good to have this method here.
// But there is no other place to write this.
func (ctx *LegscContext) SendMessage(msg message.Message) {
	msgBytes, err := message.Marshal(msg)
	if err != nil {
		log.Error("error in encode message: ", err)
		return
	}

	if msg.GetMessageType() != "console" {
		log.Info(fmt.Sprintf("send message to server. msg_type: %s, msg_model: %s", msg.GetMessageType(), msg.GetModel()))
	}

	connMutex.Lock()
	defer connMutex.Unlock()
	_ = ctx.Connection.WriteMessage(websocket.BinaryMessage, msgBytes)
}

var mutex sync.Mutex

// UpdateStatus update legsc status file
func (ctx *LegscContext) UpdateStatus() (err error) {
	mutex.Lock()
	bytes, err := json.Marshal(ctx.Status)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(ctx.StatusFileName, bytes, 0600); err != nil {
		return err
	}
	mutex.Unlock()

	return nil
}

// ConsoleSession is container of console session.
type ConsoleSession struct {
	Fd     *os.File
	Closer chan struct{}
	Closed chan struct{}
}
