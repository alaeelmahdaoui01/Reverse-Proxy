//  testing backend with no http, no proxy
//  one backend 

package main

import (
	"fmt"
	"log"

	"project.com/proxy"
)

func main() {
	b, err := proxy.NewBackend("http://localhost:8082")
	if err != nil {
		log.Fatal(err)
	}

	b.IncreaseConn()
	b.IncreaseConn()
	b.DecreaseConn()

	fmt.Println("Alive:", b.IsAlive())
	fmt.Println("Connections:", b.GetConnCount())
}
