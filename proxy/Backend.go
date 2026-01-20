package reverse_proxy

import (
	"net/url"
	"sync"
)

type Backend struct {
	URL *url.URL `json:"url"`
	Alive bool `json:"alive"`
	CurrentConns int64 `json:"current_connections"`
	mux sync.RWMutex
}

// TO IMPLEMENT 


// pointer receiver (we're modifying the value)
func (backend *Backend) setAlive(alive bool) {

}


func (backend *Backend) isAlive() bool{
	return true 
}