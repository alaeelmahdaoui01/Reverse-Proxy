package proxy

// this is who the proxy talks to, not the server pool


import (
	"net/url"
)

type LoadBalancer interface {
	GetNextValidPeer() *Backend
	AddBackend(backend *Backend)
	SetBackendStatus(uri *url.URL, alive bool)
}