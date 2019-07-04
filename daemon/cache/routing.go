package cache

import (
	"github.com/iij/legs-client/daemon/log"
	"github.com/iij/legs-client/daemon/model"
)

// RoutingCache is a interface for cache of routing data.
type RoutingCache interface {
	Update(*model.Routing)
	Delete(*model.Routing)
	FindRouting(key string) []model.Routing
	run()
}

type routingCache struct {
	routings   map[int64]model.Routing
	updateChan chan *model.Routing
	deleteChan chan *model.Routing
	findChan   chan string
	outputChan chan []model.Routing
}

// NewRoutingCache make a RoutingCache instance, and start goroutine for select loop.
func NewRoutingCache() RoutingCache {
	cache := routingCache{
		routings:   make(map[int64]model.Routing),
		updateChan: make(chan *model.Routing),
		deleteChan: make(chan *model.Routing),
		findChan:   make(chan string),
		outputChan: make(chan []model.Routing),
	}
	go cache.run()
	return &cache
}

func (r *routingCache) Update(routing *model.Routing) {
	r.updateChan <- routing
}

func (r *routingCache) Delete(routing *model.Routing) {
	r.deleteChan <- routing
}

func (r *routingCache) FindRouting(path string) []model.Routing {
	r.findChan <- path
	return <-r.outputChan
}

func (r *routingCache) run() {
	for {
		select {
		case routing := <-r.updateChan:
			r.routings[routing.ID] = *routing
			log.Info("update routing cache: ", routing.Name)
		case routing := <-r.deleteChan:
			delete(r.routings, routing.ID)
			log.Info("delete routing cache: ", routing.Name)
		case path := <-r.findChan:
			r.outputChan <- findRouting(r.routings, path)
		}
	}
}

func findRouting(routings map[int64]model.Routing, path string) (result []model.Routing) {
	for _, routing := range routings {
		match, params := routing.ComparePath(path)
		if match {
			routing.Params = params
			result = append(result, routing)
		}
	}

	return
}
