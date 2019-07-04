package message

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	message "github.com/iij/legs-message"

	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/log"
)

// HandleProxyMessage handles proxy message
func HandleProxyMessage(ctx *context.LegscContext, msgBytes []byte) {
	msg := &message.Proxy{}
	err := message.Unmarshal(msgBytes, msg)
	if err != nil {
		log.Error("failed in parsing message:", err)
		return
	}

	switch msg.GetModel() {
	case "proxy_request":
		handleProxyRequest(ctx, msg)
	}
}

func handleProxyRequest(ctx *context.LegscContext, msg *message.Proxy) {
	proxyData := msg.Data
	log.Info("received proxy message. proxy id:", proxyData.ID)

	reader := bufio.NewReader(bytes.NewReader(proxyData.Request))
	r, err := http.ReadRequest(reader)

	request, _ := http.NewRequest(r.Method, proxyData.URL, r.Body)
	request.Header = r.Header

	if err != nil {
		log.Error("error in parsing request:", err.Error())
		return
	}

	log.Info("request to:", request.URL.String())

	client := new(http.Client)
	res, err := client.Do(request)
	if err != nil {
		log.Error("request was failed:", err.Error())
		res = &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		}
	}

	log.Info("response code:", res.StatusCode)

	dumpResponse, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Error("error in dump response:", err.Error())
		res = &http.Response{
			StatusCode: http.StatusBadGateway,
			Body:       ioutil.NopCloser(bytes.NewBufferString(err.Error())),
		}
		dumpResponse, err = httputil.DumpResponse(res, true)
		if err != nil {
			return
		}
	}

	responseMsg := message.NewProxyResponse(proxyData, dumpResponse)
	ctx.SendMessage(responseMsg)
}
