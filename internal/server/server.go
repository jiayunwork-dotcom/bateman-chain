package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"bateman-chain/internal/eval"
)

type Config struct {
	Addr string
}

func ListenAndServe(cfg Config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/solve", handleSolve)
	mux.HandleFunc("/health", handleHealth)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

type solveRequest struct {
	Names   []string  `json:"names"`
	Lambda  []float64 `json:"lambda"`
	Initial []float64 `json:"initial"`
	Time    float64   `json:"time"`
}

func handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("read body: %v", err))
		return
	}
	var req solveRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %v", err))
		return
	}
	c := eval.Case{Names: req.Names, Lambda: req.Lambda, Initial: req.Initial}
	if err := c.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("validation: %v", err))
		return
	}
	t := req.Time
	if t <= 0 {
		writeError(w, http.StatusBadRequest, "time must be positive")
		return
	}
	request, err := eval.Single(t)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("request: %v", err))
		return
	}
	result, err := eval.Run(c, request)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("solve: %v", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"solver":     result.Solver,
		"times":      result.Times,
		"counts":     result.Counts,
		"activities": result.Activities,
	})
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
