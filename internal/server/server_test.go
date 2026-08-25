package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handleHealth(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSolveEndpoint(t *testing.T) {
	body := `{"names":["N1","N2","N3"],"lambda":[0.1,0.05,0.0],"initial":[1000,0,0],"time":10}`
	req := httptest.NewRequest(http.MethodPost, "/api/solve", strings.NewReader(body))
	w := httptest.NewRecorder()
	handleSolve(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "solver") {
		t.Fatalf("expected solver in response, got %s", w.Body.String())
	}
}

func TestSolveInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/solve", strings.NewReader(`{bad`))
	w := httptest.NewRecorder()
	handleSolve(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSolveMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/solve", nil)
	w := httptest.NewRecorder()
	handleSolve(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
