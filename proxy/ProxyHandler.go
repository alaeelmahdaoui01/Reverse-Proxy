package proxy 


import (
	"net/http"
	"net/http/httputil"
	"context"
	"time"
)

// // an HTTP handler: used forward HTTP traffic properly

// Receive request
// Pick backend
// Forward request to backend server 
// Return the backend's response to the client 
// increment connection, once done after request finishes 

// proxy config give us the info in the port of the proxy and the backends, which we store in server pool 
// when we will run the main, we will use proxy config for startup
// here we just need serverPool
// here loadbalancer.getNextPeer will behave in main depending on proxyconfig and the strategy will already be defined ig?? 


type ProxyHandler struct {
	lb LoadBalancer  // to choose the backend at request
	sPool *ServerPool // kill the backend in case of error (mark it dead instead of alive)
}



func NewProxyHandler(lb LoadBalancer, pool *ServerPool) *ProxyHandler {
	return &ProxyHandler{
		lb:   lb,
		sPool: pool,
	}
}


// A ajouter 
// If no backends are available, the proxy should return an appropriate HTTP error (e.g., 503 Service Unavailable)
// using context : when client diconnects, cancel the client backend request, OR if backend takes too long (5s) the request is cancelled 
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	backend := p.lb.GetNextValidPeer()
	if backend == nil {
		http.Error(w, "No available backend", http.StatusServiceUnavailable)
		return
	}

	backend.IncreaseConn()
	defer backend.DecreaseConn()

	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	// error handler of when backend fails, set as non alive 
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		p.sPool.SetBackendStatus(backend.URL, false)
		http.Error(rw, "Backend unavailable", http.StatusBadGateway)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) 
	defer cancel()

	rq := r.WithContext(ctx)

	proxy.ServeHTTP(w, rq)
}
