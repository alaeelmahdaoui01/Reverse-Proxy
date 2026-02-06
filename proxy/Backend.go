
package proxy

import (
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	URL *url.URL `json:"url"`
	Alive bool `json:"alive"`
	CurrentConns int64 `json:"current_connections"`
	mux sync.RWMutex
}


// constructor for backends
func NewBackend(rawURL string) (*Backend, error) {

	parsedURL, err := url.Parse(rawURL) 
	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:          parsedURL,
		Alive:        true,
		CurrentConns: 0,
	}, nil
}


func (backend *Backend) SetAlive(alive bool) {
	// Lock instead of Rlock bcs here we're writing and the lock should be exclusive, accessed by only one goroutine 
	backend.mux.Lock()
	defer backend.mux.Unlock()
	backend.Alive = alive  
	
}


func (backend *Backend) IsAlive() bool{
	// locking the current state of backend (to avoid collision with the health checker goroutine access at the same time)
	// Rlock for the lock on reading only 
	backend.mux.RLock()
	defer backend.mux.RUnlock()
	return backend.Alive
}


// manipulating the connections to the backend 
// must use atomic to avoid race conditions 
// mutex would work but since this could be running multiple times at a time
// bcs the code would be executing for every request, it would be too expensive for a high load
// atomic allows to do the operation as well, doesnt need mutex 

func (backend *Backend) IncreaseConn() {
	atomic.AddInt64(&backend.CurrentConns, 1)
}

func (backend *Backend) DecreaseConn() {
	atomic.AddInt64(&backend.CurrentConns, -1)
}

// to use for least connections strategy
func (backend *Backend) GetConnCount() int64 {
	return atomic.LoadInt64(&backend.CurrentConns)
}


// to make logs cleaner 
func (backend *Backend) String() string {
	return backend.URL.String()
}


// A backend becomes non-alive when your proxy decides it cannot be trusted to receive traffic.
// There are exactly two moments when this happens 
// Backend becomes non-alive during health checking 
// Backend becomes non-alive during proxying errors : httputil.ReverseProxy calls  ErrorHandler