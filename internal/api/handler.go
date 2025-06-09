package api

import (
	"claude-code-api/internal/claude"
	"claude-code-api/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Handler handles API requests
type Handler struct {
	claudeClient   *claude.Client
	fallbackClient *FallbackClient
}

// NewHandler creates a new API handler
func NewHandler() *Handler {
	h := &Handler{
		claudeClient:   claude.New(),
		fallbackClient: NewFallbackClient(),
	}

	// Log fallback providers
	if h.fallbackClient.HasFallbackProviders() {
		log.Printf("[HANDLER] Fallback providers configured: %v", h.fallbackClient.GetProviderNames())
	} else {
		log.Printf("[HANDLER] No fallback providers configured")
	}

	return h
}

// ChatCompletions handles the /v1/chat/completions endpoint
func (h *Handler) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "invalid_request_error")
		return
	}

	var req models.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body", "invalid_request_error")
		return
	}

	// Validate request
	if len(req.Messages) == 0 {
		h.sendError(w, http.StatusBadRequest, "Messages cannot be empty", "invalid_request_error")
		return
	}

	// Handle streaming response
	if req.Stream {
		h.handleStreamingResponse(w, req)
		return
	}

	// Handle non-streaming response
	h.handleNormalResponse(w, req)
}

// handleNormalResponse handles non-streaming chat completions
func (h *Handler) handleNormalResponse(w http.ResponseWriter, req models.ChatCompletionRequest) {
	// Prepare the conversation for Claude
	conversation := h.prepareConversation(req.Messages)

	log.Printf("[CLAUDE REQUEST] Model: %s, Messages: %d", req.Model, len(req.Messages))

	var response string
	var err error
	var usedFallback bool

	// Try Claude first
	response, err = h.claudeClient.Ask(conversation)
	// If Claude fails and we have fallback providers, try them
	if err != nil && h.fallbackClient.HasFallbackProviders() {
		log.Printf("[CLAUDE ERROR] %v, attempting fallback", err)

		// Convert messages to OpenAI format for fallback
		openAIMessages := make([]map[string]string, len(req.Messages))
		for i, msg := range req.Messages {
			openAIMessages[i] = map[string]string{
				"role":    msg.Role,
				"content": msg.Content,
			}
		}

		// Try each fallback provider
		for _, provider := range h.fallbackClient.providers {
			log.Printf("[FALLBACK] Trying %s", provider.Name)
			response, err = h.fallbackClient.CallOpenAICompatibleAPI(provider, openAIMessages, req.Model)
			if err == nil {
				usedFallback = true
				log.Printf("[FALLBACK] ✅ Success with %s", provider.Name)
				break
			}
			log.Printf("[FALLBACK] ❌ %s failed: %v", provider.Name, err)
		}
	}

	// If all attempts failed
	if err != nil {
		log.Printf("[API ERROR] All providers failed: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Internal server error", "api_error")
		return
	}

	if usedFallback {
		log.Printf("[FALLBACK RESPONSE] Length: %d characters", len(response))
	} else {
		log.Printf("[CLAUDE RESPONSE] Length: %d characters", len(response))
	}

	// Create response
	resp := models.ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%s", uuid.New().String()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []models.Choice{
			{
				Index: 0,
				Message: models.ChatMessage{
					Role:    "assistant",
					Content: response,
				},
				FinishReason: "stop",
			},
		},
		Usage: models.Usage{
			PromptTokens:     h.estimateTokens(conversation),
			CompletionTokens: h.estimateTokens(response),
			TotalTokens:      h.estimateTokens(conversation) + h.estimateTokens(response),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// handleStreamingResponse handles streaming chat completions
func (h *Handler) handleStreamingResponse(w http.ResponseWriter, req models.ChatCompletionRequest) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		h.sendError(w, http.StatusInternalServerError, "Streaming not supported", "api_error")
		return
	}

	// Prepare conversation
	conversation := h.prepareConversation(req.Messages)

	log.Printf("[CLAUDE STREAM REQUEST] Model: %s, Messages: %d", req.Model, len(req.Messages))

	// Generate a request ID
	requestID := fmt.Sprintf("chatcmpl-%s", uuid.New().String())

	// Send initial response
	initialChunk := models.StreamResponse{
		ID:      requestID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []models.StreamChoice{
			{
				Index: 0,
				Delta: models.DeltaContent{
					Role: "assistant",
				},
				FinishReason: nil,
			},
		},
	}

	h.sendStreamChunk(w, initialChunk)
	flusher.Flush()

	// Call Claude and get response
	response, err := h.claudeClient.Ask(conversation)
	if err != nil {
		log.Printf("[CLAUDE STREAM ERROR] %v", err)
		return
	}

	log.Printf("[CLAUDE STREAM RESPONSE] Length: %d characters", len(response))

	// Stream the response in chunks
	words := strings.Fields(response)
	for i, word := range words {
		chunk := models.StreamResponse{
			ID:      requestID,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []models.StreamChoice{
				{
					Index: 0,
					Delta: models.DeltaContent{
						Content: word + " ",
					},
					FinishReason: nil,
				},
			},
		}

		h.sendStreamChunk(w, chunk)
		flusher.Flush()

		// Simulate streaming delay
		if i < len(words)-1 {
			time.Sleep(20 * time.Millisecond)
		}
	}

	// Send final chunk
	finishReason := "stop"
	finalChunk := models.StreamResponse{
		ID:      requestID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []models.StreamChoice{
			{
				Index:        0,
				Delta:        models.DeltaContent{},
				FinishReason: &finishReason,
			},
		},
	}

	h.sendStreamChunk(w, finalChunk)
	flusher.Flush()

	// Send [DONE] marker
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// prepareConversation converts OpenAI messages to Claude format
func (h *Handler) prepareConversation(messages []models.ChatMessage) string {
	var parts []string

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			parts = append(parts, fmt.Sprintf("System: %s", msg.Content))
		case "user":
			parts = append(parts, fmt.Sprintf("Human: %s", msg.Content))
		case "assistant":
			parts = append(parts, fmt.Sprintf("Assistant: %s", msg.Content))
		}
	}

	// Add final Human prompt if needed
	if len(messages) > 0 && messages[len(messages)-1].Role != "user" {
		parts = append(parts, "Human: Please continue.")
	}

	return strings.Join(parts, "\n\n")
}

// sendStreamChunk sends a streaming chunk
func (h *Handler) sendStreamChunk(w io.Writer, chunk models.StreamResponse) {
	data, err := json.Marshal(chunk)
	if err != nil {
		log.Printf("Error marshaling chunk: %v", err)
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", string(data))
}

// estimateTokens estimates token count (rough approximation)
func (h *Handler) estimateTokens(text string) int {
	// Rough estimation: 1 token ≈ 4 characters
	return len(text) / 4
}

// sendError sends an error response
func (h *Handler) sendError(w http.ResponseWriter, statusCode int, message, errorType string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := models.ErrorResponse{
		Error: models.ErrorDetail{
			Message: message,
			Type:    errorType,
		},
	}

	if err := json.NewEncoder(w).Encode(errorResp); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// Models handles the /v1/models endpoint
func (h *Handler) Models(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "invalid_request_error")
		return
	}

	resp := models.ModelsResponse{
		Object: "list",
		Data: []models.ModelInfo{
			{
				ID:      "gpt-3.5-turbo",
				Object:  "model",
				Created: 1677610602,
				OwnedBy: "claude-code-api",
			},
			{
				ID:      "gpt-4",
				Object:  "model",
				Created: 1687882410,
				OwnedBy: "claude-code-api",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "claude-code-api",
	})
}
