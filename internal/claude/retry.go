package claude

import (
	"fmt"
	"time"
)

// RetryableClient wraps Client with retry capabilities
type RetryableClient struct {
	*Client
	maxRetries int
	retryDelay time.Duration
}

// NewRetryableClient creates a new retryable client
func NewRetryableClient(maxRetries int, retryDelay time.Duration, opts ...Option) *RetryableClient {
	return &RetryableClient{
		Client:     New(opts...),
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// Ask calls Claude with retry logic
func (c *RetryableClient) Ask(question string) (string, error) {
	var lastErr error
	
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		response, err := c.Client.Ask(question)
		if err == nil {
			return response, nil
		}
		
		lastErr = err
		
		// Don't sleep after the last attempt
		if attempt < c.maxRetries {
			time.Sleep(c.retryDelay * time.Duration(attempt+1))
		}
	}
	
	return "", fmt.Errorf("failed after %d retries: %v", c.maxRetries+1, lastErr)
}