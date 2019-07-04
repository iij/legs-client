package message

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iij/legs-client/config"
	"github.com/iij/legs-client/daemon/context"
	"github.com/gorilla/websocket"
	message "github.com/iij/legs-message"
	"github.com/spf13/viper"
)

func TestHandleCommandMessage(t *testing.T) {
	testCmdMsg := &message.Command{}
	testCmdMsg.SetStatus("waiting")
	msgBytes, _ := message.Marshal(testCmdMsg)

	ctx := context.NewLegscContext(&config.Config{Viper: viper.New()})

	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	conn, _, _ := websocket.DefaultDialer.Dial(u, nil)
	defer conn.Close()

	ctx.Connection = conn

	HandleCommandMessage(ctx, msgBytes)
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, msg)
		if err != nil {
			break
		}
	}
}
