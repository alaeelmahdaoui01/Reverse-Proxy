package proxy

// this is who the proxy talks to, not the server pool


type LoadBalancer interface {
	GetNextValidPeer() *Backend
}