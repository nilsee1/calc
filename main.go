package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

//go:embed static
var staticFiles embed.FS

// CalculateRequest represents the JSON request body for /api/calculate.
type CalculateRequest struct {
	Expression string `json:"expression"`
}

// CalculateResponse represents a successful JSON response.
type CalculateResponse struct {
	Result float64 `json:"result"`
}

// ErrorResponse represents an error JSON response.
type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Serve static files from the embedded filesystem
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Fehler beim Laden der statischen Dateien: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API endpoint for calculations
	mux.HandleFunc("/api/calculate", handleCalculate)

	// Health check endpoint
	mux.HandleFunc("/health", handleHealth)

	addr := ":" + port
	log.Printf("Server startet auf http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server-Fehler: %v", err)
	}
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Methode nicht erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Ungültiges JSON"})
		return
	}

	expr := strings.TrimSpace(req.Expression)
	if expr == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Leerer Ausdruck"})
		return
	}

	result, err := Calculate(expr)
	if err != nil {
		writeJSON(w, http.StatusOK, ErrorResponse{Error: err.Error()})
		return
	}

	// Format result: remove unnecessary trailing zeros
	writeJSON(w, http.StatusOK, CalculateResponse{Result: roundFloat(result, 10)})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// roundFloat rounds a float64 to n decimal places to avoid floating-point artifacts.
func roundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

