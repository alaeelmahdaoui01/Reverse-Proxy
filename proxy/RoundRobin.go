package proxy

import (
	"sync/atomic"
	// "log"
)

type RoundRobin struct {
	pool    *ServerPool
}

func NewRoundRobin(pool *ServerPool) *RoundRobin {
	return &RoundRobin{pool: pool}
}


func (rr *RoundRobin) GetNextValidPeer() *Backend {
	rr.pool.mux.RLock()
	defer rr.pool.mux.RUnlock()

	n := len(rr.pool.Backends)
	if n == 0 {
		return nil
	}

	start := atomic.AddUint64(&rr.pool.Current, 1)

	for i := 0; i < n; i++ {
		idx := int((start + uint64(i)) % uint64(n))
		b := rr.pool.Backends[idx]
		// log.Printf("Checking backend %s alive=%v", b.URL, b.IsAlive())
		if b.IsAlive() {
			return b
		}
	}

	
	return nil
}



