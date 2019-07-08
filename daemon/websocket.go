package daemon

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/handler/message"
	"github.com/iij/legs-client/daemon/log"
	"github.com/iij/legs-client/daemon/model/status"
	"github.com/iij/legs-client/util"
)

func initConnection(ctx *context.LegscContext) {
	prot := "ws"
	if ctx.UseWss {
		prot = "wss"
	}

	serverURL := url.URL{Scheme: prot, Host: ctx.Server, Path: "/hello"}
	header := createHeader(ctx.DeviceName, ctx.Secret)

reconnect:
	startTime := time.Now()
	timer := time.NewTimer(30 * time.Second)
	for {
		log.Info("connecting to server:", serverURL.String())

		err := connectServer(ctx, serverURL.String(), header)
		if err == nil {
			ctx.Status.Conn = status.Connected
			ctx.UpdateStatus()
			break
		}

		reconnectDuration := getReconnectDuration(startTime, time.Now())
		timer.Reset(reconnectDuration)
		if reconnectDuration == 0 {
			timer.Stop()
		}
		select {
		case <-ctx.Interrupt:
			log.Info("stop connecting server")
			timer.Stop()
			return
		case <-timer.C:
			continue
		}
	}
	defer ctx.Connection.Close()

	done := make(chan struct{})
	go readMessage(ctx, ctx.Connection, done)

restart:
	pingTicker := time.NewTicker(ctx.PingInterval)

	for {
		select {
		case <-pingTicker.C:
			err := ctx.Connection.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Info("error in sending ping:", err)
			}
		case <-done:
			ctx.Connection.Close()
			ctx.Status.Conn = status.Disconnected
			ctx.UpdateStatus()
			// Wait for a random time(5 - 15 seconds).
			rand.Seed(time.Now().UnixNano())
			sleep := rand.Intn(10) + 5
			log.Info(fmt.Sprintf("wait for %d sec...", sleep))
			time.Sleep(time.Duration(sleep) * time.Second)
			goto reconnect
		case <-ctx.Interrupt:
			log.Info("close connection")

			err := ctx.Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Info("error in sending close:", err)
			}
			ctx.Status.Conn = status.Disconnected
			ctx.UpdateStatus()
			return
		case <-ctx.Restart:
			log.Info("restart listening")
			goto restart
		}
	}
}

func createHeader(deviceName string, secret string) http.Header {
	header := http.Header{}
	hash := util.ToHash(deviceName + ":" + secret)

	header.Set("AuthUser", deviceName)
	header.Set("AuthHash", hash)

	return header
}

func connectServer(ctx *context.LegscContext, u string, header http.Header) error {
	dialer := websocket.DefaultDialer

	if ctx.Proxy != "" {
		proxyURL, err := url.Parse(ctx.Proxy)
		if err != nil {
			log.Error("invalid proxy url:", err)
			os.Exit(1)
		}
		dialer.Proxy = func(request *http.Request) (*url.URL, error) {
			return proxyURL, nil
		}
	}

	conn, res, err := dialer.Dial(u, header)
	if res != nil {
		body, _ := ioutil.ReadAll(res.Body)
		log.Info(fmt.Sprintf("status: %d, message: %s", res.StatusCode, string(body)))
		res.Body.Close()
	}
	if err != nil {
		log.Error("failed connect server:", err)
		return err
	}
	log.Info("connected to server!")
	ctx.Connection = conn
	return nil
}

func readMessage(ctx *context.LegscContext, conn *websocket.Conn, done chan struct{}) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error("closed by server:", err)
			close(done)
			return
		}
		go message.HandleMessage(ctx, msg)
	}
}

// getReconnectDuration return reconnect duration.
// In the first 10 minite, legsc reconnects every 30 second.
// In the day, legsc reconnects every 5 minute.
// In the week, legsc reconnects every 30 minute.
// After that, legsc never reconnect. In this case, it return 0.
func getReconnectDuration(start, now time.Time) time.Duration {
	duration := now.Sub(start)
	switch {
	case duration < 10*time.Minute:
		return 30 * time.Second
	case duration < 24*time.Hour:
		return 5 * time.Minute
	case duration < 168*time.Hour:
		return 30 * time.Minute
	default:
		return 0
	}
}
