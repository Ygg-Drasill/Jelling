package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ygg-Drasill/Jelling/backend/model"
)

func main() {
	h := model.Health{Status: 0}
	hBytes, _ := json.Marshal(h)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) { w.Write(hBytes) })
	http.ListenAndServe("localhost:30420", mux)
}
