
// Master Context propagation to ensure backend requests are canceled if the client disconnects.

// In proxy handler :.
// error 503 to add in case of no backend found 


// admin api A FAIRRREEE



// HealthChecker updates status, Admin API reads it and the proxyhandler (loadbalancer), client can see backend states using admin api get /health
// HealthChecker writes to ServerPool, Admin API reads from ServerPool


package main

import (
	"project.com/proxy"
	"log"
	"net/http"
	"time"
	"fmt"
	// "flag"
	"context"
	"os"
	"os/signal"
	"syscall"
)

// with : srv.Shutdown(ctx)
// Go does ALL of this automatically: Stop accepting new connections, Wait for active ServeHTTP calls, 
// Propagate cancellation to their contexts, Let them finish cleanly, Exit only after timeout

// proxy handler relies on main.go to trigger cancellation globally

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

	srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Port),
        Handler:      handler,
    }

	// Run server in goroutine
    go func() {
        log.Printf("Proxy listening on :%d\n", cfg.Port)
		// http.listenandserve(addr,handler) to run server of proxy
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Listen error: %v", err)
        }
    }()

	// Wait for CTRL+C or kill
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    <-stop
    log.Println("Shutting down proxy...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Shutdown failed: %v", err)
    }

    log.Println("Proxy stopped gracefully")

}

// before running this main.go 

// RUNNING THE BACKENDS 
// $env:PORT="8001"
// go run main.go
// $env:PORT="8002"
// go run main.go
// $env:PORT="8003"
// go run main.go

// then run the main.go 

//  then check proxy response with curl : curl.exe http://localhost:9000/students




// What this main.go does so far 
// Load config, Build server pool, Create load balancer, Create proxy handler, Start HTTP server