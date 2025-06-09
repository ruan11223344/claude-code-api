package models

// ChatCompletionRequest represents the OpenAI chat completion request
type ChatCompletionRequest struct {
	Model            string         `json:"model"`
	Messages         []ChatMessage  `json:"messages"`
	Temperature      float64        `json:"temperature,omitempty"`
	TopP             float64        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	PresencePenalty  float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
	
	// Claude-specific options (optional)
	ClaudeOptions    map[string]interface{} `json:"claude_options,omitempty"`
}

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ChatCompletionResponse represents the OpenAI chat completion response
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a single completion choice
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code,omitempty"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice represents a streaming choice
type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        DeltaContent `json:"delta"`
	FinishReason *string      `json:"finish_reason"`
}

// DeltaContent represents delta content in streaming
type DeltaContent struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// ModelInfo represents model information
type ModelInfo struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	OwnedBy string    `json:"owned_by"`
}

// ModelsResponse represents the models list response
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}