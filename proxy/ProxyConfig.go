package proxy

import (
	"time"
	"errors"
)

type ProxyConfig struct {
	Backends []string `json:"backends"`
	Port int `json:"port"`
	Strategy string `json:"strategy"` // "round-robin" or "least-conn"
	HealthCheckFreq time.Duration `json:"health_check_frequency"`
}


// config.json contains the backends of the proxy, that we will load in the server pool
func (proxyConfig *ProxyConfig) BuildServerPool() (*ServerPool , error) {
	pool:=&ServerPool{}
	for _,backend := range(proxyConfig.Backends) {
		b, err := NewBackend(backend)
		if err != nil{
			return nil, err
		}
		pool.AddBackend(b)
	}
	return pool , nil 
}

// TO IMPLEMENT 
// load for the json file 
// func LoadConfig(path string)  (*ProxyConfig, error){

// }

// setting which strategy for the load balancer 
func (proxyConfig *ProxyConfig) SetStrategy(strategy string) {
	proxyConfig.Strategy = strategy 
}

// create the load balancer based on the strategy in the json file of the proxy
func (c *ProxyConfig) CreateLoadBalancer(pool *ServerPool) (LoadBalancer, error) {
	
	switch c.Strategy {
	case "round-robin" : 
		return NewRoundRobin(pool), nil
	case "least-conn" : 
		return NewLeastConn(pool), nil
	default:
		return nil, errors.New("Invalid strategy")
	}
	
}
