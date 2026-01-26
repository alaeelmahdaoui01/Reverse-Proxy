// at the same time the health checker is gonna be running and chacking the sanity of the backends
//  we use goroutines

// Create a background goroutine that runs at a configurable interval (e.g., every 30 seconds).
// ● It should iterate through all backends in the ServerPool.
// ● Perform a "ping" (a simple GET request or TCP dial).
// ● Update the Alive status of the backend.
// ● Log changes in status (e.g., "Backend https://api1.host.com is DOWN")

// in the main we will use go HealthChecker(Spool,freq, client) (goroutine)

package proxy

import (
	"log"
	"net/http"
	"time"
)

// n addi timeout for context
type HealthChecker struct {
	SPool     *ServerPool
	Frequency time.Duration
	Client    *http.Client
}

func NewHealthChecker(pool *ServerPool, freq time.Duration) *HealthChecker {
	return &HealthChecker{
		SPool:     pool,
		Frequency: freq,
		Client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

// TO DO 
// adding context + timeout for health checks 
// NOW it runs forever, it should stop when proxy shuts down 
func (hc *HealthChecker) CheckBackend(b *Backend) bool {

	resp, err := hc.Client.Get(b.URL.String() + "/health")
	if err != nil {
		log.Println("Backend at ", b.URL, "down!")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Backend at ", b.URL, "down!", "status:", resp.StatusCode)
		return false
	}

	log.Println("Backend at ", b.URL, "up!")
	return true
}

// checks backends status and sets their status based on that
func (hc *HealthChecker) CheckBackends() {
	// hc.SPool.mux.RLock()
	// backends := hc.SPool.Backends
	// hc.SPool.mux.RUnlock()

	for _, backend := range hc.SPool.Backends {
		backend.SetAlive(hc.CheckBackend(backend))
	}
}

// call the check backends at every frequence (periodically)
func (hc *HealthChecker) Start() {

	ticker := time.NewTicker(hc.Frequency)
	defer ticker.Stop()
	for range ticker.C {
		hc.CheckBackends()
	}
}
