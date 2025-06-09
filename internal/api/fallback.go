package api

import (
	"bytes"
	"claude-code-api/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// FallbackProvider represents a fallback API provider
type FallbackProvider struct {
	Name    string
	BaseURL string
	APIKey  string
	Model   string
}

// FallbackClient handles fallback API calls
type FallbackClient struct {
	providers []FallbackProvider
	client    *http.Client
}

// NewFallbackClient creates a new fallback client
func NewFallbackClient() *FallbackClient {
	providers := []FallbackProvider{}
	
	// Load primary fallback API
	if apiKey := os.Getenv("FALLBACK_API_KEY_1"); apiKey != "" {
		provider := FallbackProvider{
			Name:    getEnvOrDefault("FALLBACK_API_NAME_1", "Fallback-1"),
			BaseURL: getEnvOrDefault("FALLBACK_API_URL_1", "https://api.openai.com/v1"),
			APIKey:  apiKey,
			Model:   getEnvOrDefault("FALLBACK_API_MODEL_1", "gpt-3.5-turbo"),
		}
		providers = append(providers, provider)
		logger.Log.Infof("[FALLBACK] Provider 1 '%s' configured at %s", provider.Name, provider.BaseURL)
	}
	
	// Load secondary fallback API
	if apiKey := os.Getenv("FALLBACK_API_KEY_2"); apiKey != "" {
		provider := FallbackProvider{
			Name:    getEnvOrDefault("FALLBACK_API_NAME_2", "Fallback-2"),
			BaseURL: getEnvOrDefault("FALLBACK_API_URL_2", "https://api.openai.com/v1"),
			APIKey:  apiKey,
			Model:   getEnvOrDefault("FALLBACK_API_MODEL_2", "gpt-3.5-turbo"),
		}
		providers = append(providers, provider)
		logger.Log.Infof("[FALLBACK] Provider 2 '%s' configured at %s", provider.Name, provider.BaseURL)
	}
	
	// Load additional fallback APIs (3-5)
	for i := 3; i <= 5; i++ {
		keyEnv := fmt.Sprintf("FALLBACK_API_KEY_%d", i)
		if apiKey := os.Getenv(keyEnv); apiKey != "" {
			provider := FallbackProvider{
				Name:    getEnvOrDefault(fmt.Sprintf("FALLBACK_API_NAME_%d", i), fmt.Sprintf("Fallback-%d", i)),
				BaseURL: getEnvOrDefault(fmt.Sprintf("FALLBACK_API_URL_%d", i), "https://api.openai.com/v1"),
				APIKey:  apiKey,
				Model:   getEnvOrDefault(fmt.Sprintf("FALLBACK_API_MODEL_%d", i), "gpt-3.5-turbo"),
			}
			providers = append(providers, provider)
			logger.Log.Infof("[FALLBACK] Provider %d '%s' configured at %s", i, provider.Name, provider.BaseURL)
		}
	}

	return &FallbackClient{
		providers: providers,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CallWithFallback attempts to call the primary API and falls back if it fails
func (fc *FallbackClient) CallWithFallback(primaryFunc func() (string, error)) (string, error) {
	// Try primary Claude API first
	response, err := primaryFunc()
	if err == nil {
		return response, nil
	}

	logger.Log.Infof("[FALLBACK] Primary API failed: %v", err)

	// Try fallback providers
	for i, provider := range fc.providers {
		logger.Log.Infof("[FALLBACK] Trying fallback provider %d/%d: %s", i+1, len(fc.providers), provider.Name)
		
		response, err := fc.callProvider(provider)
		if err == nil {
			logger.Log.Infof("[FALLBACK] ✅ Success with %s", provider.Name)
			return response, nil
		}
		
		logger.Log.Infof("[FALLBACK] ❌ Provider %s failed: %v", provider.Name, err)
	}

	return "", fmt.Errorf("all providers failed, last error: %v", err)
}

// callProvider calls a specific fallback provider
func (fc *FallbackClient) callProvider(provider FallbackProvider) (string, error) {
	// All providers use OpenAI-compatible API
	return fc.callOpenAICompatible(provider)
}

// callOpenAICompatible calls any OpenAI-compatible API
func (fc *FallbackClient) callOpenAICompatible(provider FallbackProvider) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", provider.BaseURL)
	
	// Create request body
	reqBody := map[string]interface{}{
		"model": provider.Model,
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
		"temperature": 0.7,
		"max_tokens":  500,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.APIKey))
	
	resp, err := fc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	
	// Extract content from response
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}
	
	return "", fmt.Errorf("unexpected response format")
}

// CallOpenAICompatibleAPI calls an OpenAI-compatible API with full request
func (fc *FallbackClient) CallOpenAICompatibleAPI(provider FallbackProvider, messages []map[string]string, model string) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", provider.BaseURL)
	
	// Use provided model or default to provider's model
	if model == "" {
		model = provider.Model
	}
	
	// Create request body
	reqBody := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": 0.7,
		"max_tokens":  2000,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.APIKey))
	
	startTime := time.Now()
	resp, err := fc.client.Do(req)
	duration := time.Since(startTime)
	
	if err != nil {
		return "", fmt.Errorf("request failed after %v: %v", duration, err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}
	
	// Extract content from response
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					logger.Log.Infof("[FALLBACK] %s responded in %v", provider.Name, duration)
					return content, nil
				}
			}
		}
	}
	
	return "", fmt.Errorf("unexpected response format from %s", provider.Name)
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// HasFallbackProviders checks if any fallback providers are configured
func (fc *FallbackClient) HasFallbackProviders() bool {
	return len(fc.providers) > 0
}

// GetProviderNames returns the names of configured providers
func (fc *FallbackClient) GetProviderNames() []string {
	names := make([]string, len(fc.providers))
	for i, p := range fc.providers {
		names[i] = p.Name
	}
	return names
}

// GetProviders returns all configured providers
func (fc *FallbackClient) GetProviders() []FallbackProvider {
	return fc.providers
}