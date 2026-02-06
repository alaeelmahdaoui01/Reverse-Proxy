package proxy

import (
	"time"
	"errors"
	"os"
	"encoding/json"
)

type ProxyConfig struct {
	Backends []BackendConfig `json:"backends"`
	Port int `json:"port"`
	Strategy string `json:"strategy"` // "round-robin" or "least-conn" or "weighted"
	HealthCheckFreqRaw string `json:"health_check_frequency"`
	HealthCheckFreq    time.Duration `json:"-"`  // to not read from json
}

type BackendConfig struct {
	URL    string `json:"url"`
	Weight int    `json:"weight"` 
}


// config.json contains the backends of the proxy, that we will load in the server pool
func (proxyConfig *ProxyConfig) BuildServerPool() (*ServerPool , error) {
	pool:=&ServerPool{}
	for _,backend := range(proxyConfig.Backends) {
		b, err := NewBackend(backend.URL,backend.Weight)
		if err != nil{
			return nil, err
		}
		pool.AddBackend(b)
	}
	return pool , nil 
}

func LoadConfig(path string) (*ProxyConfig, error) {

	var config ProxyConfig

	if config.HealthCheckFreqRaw == "" {
		config.HealthCheckFreq = 30 * time.Second 
	} else {
		d, err := time.ParseDuration(config.HealthCheckFreqRaw)
		if err != nil {
			return nil, errors.New("invalid health_check_frequency (example: \"30s\")")
		}
		config.HealthCheckFreq = d
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	if config.Port == 0 {
		return nil, errors.New("proxy port must be specified")
	}
	if len(config.Backends) == 0 {
		return nil, errors.New("at least one backend is required")
	}
	if config.Strategy == "" {
		config.Strategy = "round-robin" 
	}

	return &config, nil
}



// create the load balancer based on the strategy in the json file of the proxy
func (c *ProxyConfig) CreateLoadBalancer(pool *ServerPool) (LoadBalancer, error) {
	
	switch c.Strategy {
	case "round-robin" : 
		return NewRoundRobin(pool), nil
	case "least-conn" : 
		return NewLeastConn(pool), nil
	case "weightedRB":
		return NewWeightedRoundRobin(pool), nil
	default:
		return nil, errors.New("Invalid strategy")
	}
	
}
