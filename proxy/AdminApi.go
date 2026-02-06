package proxy

import (
	"encoding/json"
	"log"
	"net/http"
)

type AdminApi struct {
	SPool *ServerPool
}

func NewAdminApi(sPool *ServerPool) *AdminApi {
	return &AdminApi{SPool: sPool}
}

type ResponseStructure struct {
	TotalBackends  int             `json:"total_backends"`
	ActiveBackends int             `json:"active_backends"`
	Backends       []BackendStatus `json:"backends"`
}
type BackendStatus struct {
	URL                string `json:"url"`
	Alive              bool   `json:"alive"`
	CurrentConnections int    `json:"current_connections"`
	Weight             int    `json:"weight"`
}

func (api *AdminApi) StatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	api.SPool.mux.Lock()
	defer api.SPool.mux.Unlock()

	active_backends := 0
	for _, backend := range api.SPool.Backends {
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
			Weight:             backend.Weight,
		})
	}

	response := ResponseStructure{
		TotalBackends:  len(api.SPool.Backends),
		ActiveBackends: active_backends,
		Backends:       backends,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type BackendReq struct {
	Url    string `json:"url"`
	Weight int    `json:"weight"`
}

type BackendReqDel struct {
	Url string `json:"url"`
}

// the r.Body contains the url to add and to delete (we will tranform it from json to the struct)
func (api *AdminApi) BackendsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

		var req BackendReq

		error := json.NewDecoder(r.Body).Decode(&req)
		if error != nil {
			http.Error(w, "Invalid Json", http.StatusBadRequest)
			return
		}

		backend, err := NewBackend(req.Url, req.Weight)
		if err != nil {
			http.Error(w, "invalid backend URL", http.StatusBadRequest)
			return
		}

		// another check other than the one in server pool
		api.SPool.mux.Lock()
		for _, existingBackend := range api.SPool.Backends {
			if existingBackend.URL.String() == backend.URL.String() {
				api.SPool.mux.Unlock()
				http.Error(w, "backend with this URL already exists", http.StatusConflict)
				return
			}
		}
		api.SPool.mux.Unlock()

		api.SPool.AddBackend(backend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "backend added successfully",
			"url":    backend.URL.String(),
		})

	case http.MethodDelete:
		var req BackendReqDel

		error := json.NewDecoder(r.Body).Decode(&req)
		if error != nil {
			http.Error(w, "Invalid Json", http.StatusBadRequest)
			return
		}

		api.SPool.mux.Lock()
		defer api.SPool.mux.Unlock()

		for idx, backend := range api.SPool.Backends {
			if backend.URL.String() == req.Url {
				api.SPool.Backends = append(api.SPool.Backends[:idx], api.SPool.Backends[idx+1:]...)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{
					"status": "backend deleted successfully",
					"url":    backend.URL.String(),
				})
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

