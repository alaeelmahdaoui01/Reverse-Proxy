package proxy


import (
	"encoding/json"
	"net/http"
	"log"
) 


type AdminApi struct{
	SPool *ServerPool
}

func NewAdminApi(sPool *ServerPool) *AdminApi {
	return &AdminApi{SPool: sPool}
}

type ResponseStructure struct {
	TotalBackends    int   `json:"total_backends"`
	ActiveBackends     int      `json:"active_backends"`
	Backends []BackendStatus   `json:"backends"`
}
type BackendStatus struct {
    URL                string `json:"url"`
    Alive              bool   `json:"alive"`
    CurrentConnections int    `json:"current_connections"`
}

func (api *AdminApi) StatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return 
	} 
	api.SPool.mux.Lock()
	defer api.SPool.mux.Unlock()

	active_backends := 0 
	for _ ,backend := range api.SPool.Backends {
		if backend.IsAlive() {
			active_backends++
		}
	}

	backends := make([]BackendStatus, 0)

	for _, backend := range api.SPool.Backends {
		backends = append(backends, BackendStatus{
			URL:                backend.URL.String(),
			Alive:              backend.IsAlive(),
			CurrentConnections: int(backend.GetConnCount()),
		})
	}


	response := ResponseStructure{
		TotalBackends : len(api.SPool.Backends),
		ActiveBackends: active_backends,
		Backends: backends,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type BackendReq struct {
	Url string `json:"url"`
}

// the r.Body contains the url to add and to delete (we will tranform it from json to the struct)
func (api *AdminApi) BackendsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case http.MethodPost:

			var req BackendReq
			// api.SPool.mux.Lock()
			// defer api.SPool.mux.Unlock()
			error := json.NewDecoder(r.Body).Decode(&req)
			if error != nil {
				http.Error(w, "Invalid Json" , http.StatusBadRequest)
				return
			}

			backend, err := NewBackend(req.Url)
			if err!=nil {
				http.Error(w, "invalid backend URL", http.StatusBadRequest)
				return
			}

			api.SPool.AddBackend(backend)
			w.WriteHeader(http.StatusCreated)



		case http.MethodDelete:
			var req BackendReq

			

			error := json.NewDecoder(r.Body).Decode(&req)
			if error != nil {
				http.Error(w, "Invalid Json" , http.StatusBadRequest)
				return
			}

			api.SPool.mux.Lock()
			defer api.SPool.mux.Unlock()

			for idx, backend := range api.SPool.Backends {
				if backend.URL.String() == req.Url{
					api.SPool.Backends = append(api.SPool.Backends[:idx],api.SPool.Backends[idx+1:]...,)
					w.WriteHeader(http.StatusOK)
					return 
				}
				
			}

			http.Error(w, "backend not found", http.StatusNotFound)

		
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (api *AdminApi) StartAdminServer() {
	http.HandleFunc("/backends", api.BackendsHandler)
	http.HandleFunc("/status", api.StatusHandler)
	log.Println("Admin server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}


