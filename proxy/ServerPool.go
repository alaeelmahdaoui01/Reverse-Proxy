// stores all the available backends
// and the current backend that its gonna start from to do the load balancer (move to the next server connection)
// current is the last used backend server, so we start from the next one in round robbin
// im gonna start without load balancer , just implementing the roundrobbing here and test my code 

package proxy

import (
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current uint64 `json:"current"` 
	mux sync.RWMutex
}


// we can add as much backend servers as we want
func (sPool *ServerPool) addBackend(backend *Backend) {
	sPool.mux.Lock()
	defer sPool.mux.Unlock()
	sPool.Backends = append(sPool.Backends, backend)
}

// we return one alive backend, we compute the current backend 
func (sPool *ServerPool) returnValidBackend() (*Backend) {
	sPool.mux.RLock()
	defer sPool.mux.RUnlock()

	n := len(sPool.Backends)
	if n == 0 {
		return nil
	}

	start := atomic.AddUint64(&sPool.Current, 1)

	for i := 0; i < n; i++ {
		idx := int((start + uint64(i)) % uint64(n))
		b := sPool.Backends[idx]
		if b.IsAlive() {
			return b
		}
	}
	return nil
	
}