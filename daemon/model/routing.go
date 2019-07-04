package model

import (
	"strings"
)

// Routing has routing configuration cache configured by server.
type Routing struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"account_id"`
	Name      string `json:"name"`
	Params    map[string]string
	Urls      []string `json:"urls"`
}

// ComparePath distribute routing path and extraction the matched parameters.
// For example, routing name is a `/sample/:id` and take a path to `/sample/1`,
// method returns match=true and params={id:1}.
func (r *Routing) ComparePath(path string) (match bool, params map[string]string) {
	match = false
	params = map[string]string{}
	targetName := r.Name

	if len(path) == 0 || path[0:1] != "/" {
		path = "/" + path
	}
	if targetName[0:1] != "/" {
		targetName = "/" + targetName
	}

	pathNodes := strings.Split(path, "/")
	if len(pathNodes) == 0 {
		return
	}

	nameNodes := strings.Split(targetName, "/")
	if len(pathNodes) != len(nameNodes) {
		return
	}

	for i, node := range nameNodes {
		if len(node) > 0 && node[0:1] == ":" {
			params[node[1:]] = pathNodes[i]
			match = true
		} else if node == pathNodes[i] {
			match = true
		} else {
			match = false
			break
		}
	}

	if !match {
		params = map[string]string{}
	}

	return
}
