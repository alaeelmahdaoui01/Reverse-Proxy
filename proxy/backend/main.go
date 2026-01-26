// from powershell to run : 

// $env:PORT="8001"
// go run main.go

// Then in another terminal:
// $env:PORT="8002"
// go run main.go

// And another:
// $env:PORT="8003"
// go run main.go

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

type Student struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Address string   `json:"address"`
	Courses []string `json:"courses"`
}

var (
	students []Student
	mu       sync.Mutex
	port     = "8000"
)

func main() {
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/students", studentsHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Backend running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(students)

	case http.MethodPost:
		var s Student
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		mu.Lock()
		students = append(students, s)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("student added\n"))

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}


//  for the backend server i dont need client.go, since im using just running the backend servers and the response/client is equivalent to running the proxy and curl 