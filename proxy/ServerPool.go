// stores all the available backends
// and the current backend that its gonna start from to do the load balancer (move to the next server connection)

// im gonna start without load balancer , just implementing the roundrobbing here and test my code 

package proxy

import (
	"sync"
	"net/url"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current uint64 `json:"current"` // initially zero since not declared, current is a counter not a backend, counter for the round robbin
	mux sync.RWMutex
}


// we can add as much backend servers as we want
func (sPool *ServerPool) AddBackend(backend *Backend) {
	sPool.mux.Lock()
	defer sPool.mux.Unlock()
	sPool.Backends = append(sPool.Backends, backend)
}


func (sPool *ServerPool) SetBackendStatus(uri *url.URL, alive bool) {
	sPool.mux.Lock()
	defer sPool.mux.Unlock()

	for _,backend := range sPool.Backends {
		if backend.URL == uri {
			backend.SetAlive(alive)
			break 
		}
	}
}

