package handlers

import (
	"encoding/json"
	"net/http"
)

type resp map[string]any

func Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/hello", hello)
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, resp{"status": "ok"})
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	writeJSON(w, http.StatusOK, resp{"message": "hello " + name})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
