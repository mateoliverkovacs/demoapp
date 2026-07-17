package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const VERSION = "1.0.0"

var configStore = struct {
	sync.RWMutex
	data map[string]string
}{data: make(map[string]string)}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": VERSION})
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "undefined"
	}
	writeJSON(w, http.StatusOK, map[string]string{"environment": env})
}

func postConfigHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	configStore.Lock()
	configStore.data[body.Name] = body.Value
	configStore.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{
		"name":  body.Name,
		"value": body.Value,
	})
}

func getConfigHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/config/")

	configStore.RLock()
	value, exists := configStore.data[name]
	configStore.RUnlock()

	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "config not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"name":  name,
		"value": value,
	})
}

func deleteConfigHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/config/")

	configStore.Lock()
	_, exists := configStore.data[name]
	if exists {
		delete(configStore.data, name)
	}
	configStore.Unlock()

	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "config not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/env", envHandler)

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postConfigHandler(w, r)
			return
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/config/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getConfigHandler(w, r)
		case http.MethodDelete:
			deleteConfigHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	srv := &http.Server{
		Addr:         ":8000",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Server running on http://localhost:8000")
	log.Fatal(srv.ListenAndServe())
}
