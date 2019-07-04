package cache

import (
	"os"
	"testing"
	"time"

	"github.com/iij/legs-client/daemon/model"

	"github.com/stretchr/testify/assert"
)

var testCache = routingCache{
	routings:   make(map[int64]model.Routing),
	updateChan: make(chan *model.Routing),
	deleteChan: make(chan *model.Routing),
	findChan:   make(chan string),
	outputChan: make(chan []model.Routing),
}

var testRouting = &model.Routing{
	ID:        1,
	AccountID: 1,
	Name:      "/test/update/:id",
	Params:    map[string]string{},
	Urls: []string{
		"https://hogege.jp",
	},
}

func TestMain(m *testing.M) {
	go testCache.run()
	os.Exit(m.Run())
}

func TestRoutingCache_Update_Delete(t *testing.T) {
	testCache.Update(testRouting)
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, len(testCache.FindRouting(testRouting.Name)))

	testCache.Delete(testRouting)
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 0, len(testCache.FindRouting(testRouting.Name)))
}

func TestRoutingCache_FindRouting(t *testing.T) {
	testCache.Update(testRouting)
	time.Sleep(500 * time.Millisecond)
	routings := testCache.FindRouting("/test/update/1")
	assert.Equal(t, 1, len(routings))
	assert.Equal(t, map[string]string{"id": "1"}, routings[0].Params)
}
