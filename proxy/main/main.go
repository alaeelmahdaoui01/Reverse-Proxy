// TEST 2 : Test ServerPool Round-Robin (NO HTTP)
// testing round robbin logic

package main

import (
	"fmt"

	"project.com/proxy"
)

func main() {
	pool := &proxy.ServerPool{}

	b1, _ := proxy.NewBackend("http://backend1")
	b2, _ := proxy.NewBackend("http://backend2")

	// server pool has 2 backends
	pool.AddBackend(b1)
	pool.AddBackend(b2)

	for i := 0; i < 6; i++ {
		b := pool.ReturnValidBackend()
		fmt.Println(b.URL)
	}
}
