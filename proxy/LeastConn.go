package proxy


type LeastConn struct {
	pool    *ServerPool
}

func NewLeastConn(pool *ServerPool) *LeastConn {
	return &LeastConn{pool: pool}
}


func (lc *LeastConn) GetNextValidPeer() *Backend {
	lc.pool.mux.RLock()
	defer lc.pool.mux.RUnlock()

	var best *Backend
	for _, b := range lc.pool.Backends {
		if !b.IsAlive() {
			continue
		}
		if best == nil || b.GetConnCount() < best.GetConnCount() {
			best = b
		}
	}
	return best
}

