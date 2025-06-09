package claude

import (
	"fmt"
	"sync"
)

// Session represents a conversation session with Claude
type Session struct {
	client      *Client
	history     []Message
	maxHistory  int
	mu          sync.Mutex
}

// Message represents a message in the conversation
type Message struct {
	Role    string // "user" or "assistant"
	Content string
}

// NewSession creates a new session with the given max history
func (c *Client) NewSession(maxHistory int) *Session {
	return &Session{
		client:     c,
		history:    make([]Message, 0),
		maxHistory: maxHistory,
	}
}

// Ask sends a question in the context of the session
func (s *Session) Ask(question string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add user message to history
	s.history = append(s.history, Message{
		Role:    "user",
		Content: question,
	})

	// Build conversation context
	context := s.buildContext()
	
	// Ask Claude
	response, err := s.client.Ask(context)
	if err != nil {
		return "", err
	}

	// Add assistant response to history
	s.history = append(s.history, Message{
		Role:    "assistant",
		Content: response,
	})

	// Trim history if needed
	s.trimHistory()

	return response, nil
}

// buildContext builds the conversation context
func (s *Session) buildContext() string {
	var context string
	
	for _, msg := range s.history {
		if msg.Role == "user" {
			context += fmt.Sprintf("Human: %s\n\n", msg.Content)
		} else {
			context += fmt.Sprintf("Assistant: %s\n\n", msg.Content)
		}
	}
	
	return context
}

// trimHistory trims the history to maintain max size
func (s *Session) trimHistory() {
	if len(s.history) > s.maxHistory {
		// Keep the most recent messages
		s.history = s.history[len(s.history)-s.maxHistory:]
	}
}

// Clear clears the session history
func (s *Session) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = make([]Message, 0)
}

// GetHistory returns the current session history
func (s *Session) GetHistory() []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Return a copy to prevent external modification
	historyCopy := make([]Message, len(s.history))
	copy(historyCopy, s.history)
	return historyCopy
}