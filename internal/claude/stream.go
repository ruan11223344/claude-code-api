package claude

import (
	"context"
	"fmt"
)

// StreamOption defines options for streaming responses
type StreamOption struct {
	OnChunk    func(chunk string)
	OnComplete func()
	OnError    func(error)
}

// AskStream streams the response from Claude
func (c *Client) AskStream(ctx context.Context, question string, opt StreamOption) error {
	// For now, we'll simulate streaming by calling the regular Ask method
	// In a real implementation, this would use Claude's streaming API
	
	response, err := c.Ask(question)
	if err != nil {
		if opt.OnError != nil {
			opt.OnError(err)
		}
		return err
	}
	
	// Simulate streaming by breaking response into chunks
	words := splitIntoWords(response)
	for _, word := range words {
		select {
		case <-ctx.Done():
			if opt.OnError != nil {
				opt.OnError(fmt.Errorf("streaming cancelled"))
			}
			return ctx.Err()
		default:
			if opt.OnChunk != nil {
				opt.OnChunk(word + " ")
			}
		}
	}
	
	if opt.OnComplete != nil {
		opt.OnComplete()
	}
	
	return nil
}

// splitIntoWords splits text into words for streaming simulation
func splitIntoWords(text string) []string {
	var words []string
	var word string
	
	for _, char := range text {
		if char == ' ' || char == '\n' || char == '\t' {
			if word != "" {
				words = append(words, word)
				word = ""
			}
			if char == '\n' {
				words = append(words, "\n")
			}
		} else {
			word += string(char)
		}
	}
	
	if word != "" {
		words = append(words, word)
	}
	
	return words
}