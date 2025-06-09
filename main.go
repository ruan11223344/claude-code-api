package main

import (
	"bytes"
	"claude-code-api/internal/api"
	"claude-code-api/internal/logger"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// Command line flags
	port := flag.String("port", "8082", "Port to run the server on")
	host := flag.String("host", "", "Host to bind to (default: all interfaces)")
	flag.Parse()

	// Override with environment variables if set
	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

	// Create API handler
	handler := api.NewHandler()

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", handler.ChatCompletions)
	mux.HandleFunc("/v1/models", handler.Models)
	mux.HandleFunc("/health", handler.HealthCheck)

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message":"Claude Code API - OpenAI Compatible API Server","version":"1.0.0","endpoints":["/v1/chat/completions","/v1/models","/health"]}`)
	})

	// Get API key from environment
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		logger.Log.Warn("API_KEY not set. API is publicly accessible!")
	} else {
		logger.Log.Info("API Key authentication enabled")
	}

	// Middleware
	wrappedMux := loggingMiddleware(corsMiddleware(authMiddleware(apiKey, mux)))

	// Start server
	addr := fmt.Sprintf("%s:%s", *host, *port)
	logger.Log.Info("=====================================")
	logger.Log.Info("🚀 Claude Code API Server")
	logger.Log.WithField("address", fmt.Sprintf("http://%s", addr)).Info("🔗 Listening on")
	logger.Log.Info("=====================================")
	logger.Log.Info("Available endpoints:")
	logger.Log.Infof("  POST   http://%s/v1/chat/completions", addr)
	logger.Log.Infof("  GET    http://%s/v1/models", addr)
	logger.Log.Infof("  GET    http://%s/health", addr)
	logger.Log.Infof("  GET    http://%s/", addr)
	logger.Log.Info("=====================================")
	logger.Log.Info("Ready to accept requests...")

	if err := http.ListenAndServe(addr, wrappedMux); err != nil {
		logger.Log.WithError(err).Fatal("Server failed to start")
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

		// Create log entry with request info
		entry := logger.Log.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"method":      r.Method,
			"path":        r.URL.Path,
		})

		entry.Info("Request received")

		// Log request body for POST requests
		if r.Method == "POST" && r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			if len(body) > 0 {
				// Sanitize the body before logging
				sanitized := logger.SanitizeRequest(string(body))
				entry.WithField("body", sanitized).Debug("Request body")
			}
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		entry.WithFields(logrus.Fields{
			"status":   wrapped.statusCode,
			"duration": duration,
		}).Info("Request completed")
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

// authMiddleware validates API key
func authMiddleware(apiKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth if no API key is configured
		if apiKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Skip auth for health check
		if r.URL.Path == "/health" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		// Get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.WithFields(logrus.Fields{
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
			}).Warn("Missing authorization header")
			respondWithAuthError(w)
			return
		}

		// Check Bearer token format
		const bearerPrefix = "Bearer "
		if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			logger.Log.WithFields(logrus.Fields{
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
			}).Warn("Invalid authorization format")
			respondWithAuthError(w)
			return
		}

		// Extract and validate token
		token := authHeader[len(bearerPrefix):]
		if token != apiKey {
			logger.Log.WithFields(logrus.Fields{
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
			}).Warn("Invalid API key")
			respondWithAuthError(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// respondWithAuthError returns OpenAI-style authentication error
func respondWithAuthError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, `{
		"error": {
			"message": "Incorrect API key provided. You can find your API key at https://platform.openai.com/account/api-keys.",
			"type": "invalid_request_error",
			"param": null,
			"code": "invalid_api_key"
		}
	}`)
}