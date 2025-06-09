package main

import (
	"bytes"
	"claude-code-api/internal/api"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Command line flags
	port := flag.String("port", "8082", "Port to run the server on")
	host := flag.String("host", "0.0.0.0", "Host to bind the server to")
	flag.Parse()

	// Create handler
	handler := api.NewHandler()

	// Setup routes
	mux := http.NewServeMux()

	// OpenAI compatible endpoints
	mux.HandleFunc("/v1/chat/completions", handler.ChatCompletions)
	mux.HandleFunc("/v1/models", handler.Models)
	mux.HandleFunc("/health", handler.HealthCheck)

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Claude Code API - OpenAI Compatible API Server","version":"1.0.0","endpoints":["/v1/chat/completions","/v1/models","/health"]}`)
	})

	// Middleware
	wrappedMux := loggingMiddleware(corsMiddleware(mux))

	// Start server
	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Println("=====================================")
	log.Printf("🚀 Claude Code API Server")
	log.Printf("🔗 Listening on: http://%s", addr)
	log.Println("=====================================")
	log.Println("Available endpoints:")
	log.Printf("  POST   http://%s/v1/chat/completions", addr)
	log.Printf("  GET    http://%s/v1/models", addr)
	log.Printf("  GET    http://%s/health", addr)
	log.Printf("  GET    http://%s/", addr)
	log.Println("=====================================")
	log.Println("Ready to accept requests...")

	if err := http.ListenAndServe(addr, wrappedMux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs all requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[REQUEST] %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)

		// Log request body for POST requests
		if r.Method == "POST" && r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			if len(body) > 0 {
				log.Printf("[REQUEST BODY] %s", string(body))
			}
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("[RESPONSE] %s %s %s - Status: %d - Duration: %v",
			r.RemoteAddr, r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
