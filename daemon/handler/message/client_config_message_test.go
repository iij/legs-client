package message

import (
	"testing"
	"time"

	message "github.com/iij/legs-message"

	"github.com/iij/legs-client/config"
	"github.com/iij/legs-client/daemon/context"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHandleClientConfigMessage(t *testing.T) {
	var testClientConfig = message.ClientConfig{
		PingInterval: 20,
	}
	ctx := context.NewLegscContext(&config.Config{Viper: viper.New()})

	var tmpMessage = message.NewClientConfigure(testClientConfig)
	bytes, _ := message.Marshal(tmpMessage)

	testMessage := &message.ClientConfigure{}
	_ = message.Unmarshal(bytes, testMessage)

	ctx.PingInterval = 10 * time.Second

	go HandleClientConfigMessage(ctx, testMessage)
	time.Sleep(1 * time.Second)

	select {
	case <-ctx.Restart:
		t.Log("received restart message.")
	default:
		t.Fatal("it is need to restart but dis not receive restart message.")
	}
	assert.Equal(t, 20*time.Second, ctx.PingInterval)

	go HandleClientConfigMessage(ctx, testMessage)
	time.Sleep(1 * time.Second)

	select {
	case <-ctx.Restart:
		t.Fatal("it is no need to restart but received restart message.")
	default:
		t.Log("did not receive restart message.")
	}
	assert.Equal(t, 20*time.Second, ctx.PingInterval)
}
