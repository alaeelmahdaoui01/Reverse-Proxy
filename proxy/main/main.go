// mansach dak lblan dial when client disconnects it should also disconnect : 
// where am i gonna use context : Handle graceful shutdowns and request timeouts using the context package
// Master Context propagation to ensure backend requests are canceled if the client disconnects.

// In rpoxy handler :
// Context: Ensure the request context is passed through so that slow backend processing
// can be canceled.
// error 503 to add in case of no backend found 


// admin api A FAIRRREEE

// inside main function :
// Main Goroutine: Starts the HTTP Proxy server and the Admin API. (should i start the backends ??)
// Health Check Goroutine: Runs a time.Ticker loop to verify backend availability.
// Request Handling: When a request arrives, the ServerPool selects a backend,
// increments its connection count, forwards the request, and decrements the count
// upon completion.


// make a simpler backend server to try then test 
// then make a rest api 
// also should i make many backends or just one? hmmmm
// it should have a get /health to be used in the health checker to checkBackend 

// HealthChecker updates status, Admin API reads it and the proxyhandler (loadbalancer), client can see backend states using admin api get /health
// HealthChecker writes to ServerPool, Admin API reads from ServerPool
package main

import (
	"project.com/proxy"
	"log"
	"net/http"
	// "time"
	"fmt"
	// "flag"
)

// temporary main
// func main() {
// 	cfg, _ := proxy.LoadConfig("config.json")
// 	pool, _ := cfg.BuildServerPool()
// 	lb, _ := cfg.CreateLoadBalancer(pool)

// 	// start proxy handler
// 	handler := proxy.NewProxyHandler(lb, pool)
// 	go func() {
// 		log.Fatal(http.ListenAndServe(":8080", handler))
// 	}()

// 	// start health checker in background
// 	hc := &proxy.HealthChecker{
// 		SPool:      pool,
// 		Frequency: cfg.HealthCheckFreq,
// 		Client:    &http.Client{Timeout: 2 * time.Second},
// 	}
// 	go hc.Start()

// 	select {} // keep main alive

// }


// func main() {
//     cfg, _ := proxy.LoadConfig("config.json")

//     pool, _ := cfg.BuildServerPool()
//     lb, _ := cfg.CreateLoadBalancer(pool)

//     handler := proxy.NewProxyHandler(lb, pool)

//     server := &http.Server{
//         Addr:    fmt.Sprintf(":%d", cfg.Port),
//         Handler: handler,
//     }


//     log.Fatal(server.ListenAndServe())
// }


func main() {
	// configPath := flag.String("config", "../config.json", "path to config file")
	// flag.Parse()
	// cfg, err := proxy.LoadConfig(*configPath)

    cfg, err := proxy.LoadConfig("../../config.json" )
    if err != nil {
        log.Fatal(err)
    }

    pool, err := cfg.BuildServerPool()
    if err != nil {
        log.Fatal(err)
    }

    lb, err := cfg.CreateLoadBalancer(pool)
    if err != nil {
        log.Fatal(err)
    }

	healthChecker := proxy.NewHealthChecker(pool, cfg.HealthCheckFreq)
	go healthChecker.Start()

    handler := proxy.NewProxyHandler(lb, pool)

    log.Printf("Proxy listening on :%d\n", cfg.Port)
    log.Fatal(http.ListenAndServe(
        fmt.Sprintf(":%d", cfg.Port),
        handler,
    ))

	// healthChecker := proxy.NewHealthChecker(pool, cfg.HealthCheckFreq)
	// go healthChecker.Start()
}




// What this main.go does so far 
// Load config, Build server pool, Create load balancer, Create proxy handler, Start HTTP server