// at the same time the health checker is gonna be running and chacking the sanity of the backends
//  we use goroutines

// Create a background goroutine that runs at a configurable interval (e.g., every 30 seconds).
// ● It should iterate through all backends in the ServerPool.
// ● Perform a "ping" (a simple GET request or TCP dial).
// ● Update the Alive status of the backend.
// ● Log changes in status (e.g., "Backend https://api1.host.com is DOWN")

package proxy

import (
	"net/http"
	"time"
)

type HealthChecker struct {
	SPool     *ServerPool
	Frequency time.Duration
	Client    *http.Client
}

func (hc *HealthChecker) Start() {
	
}
