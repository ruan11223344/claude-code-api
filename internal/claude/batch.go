package claude

import (
	"sync"
)

// BatchClient handles batch requests
type BatchClient struct {
	client        *Client
	concurrentNum int
}

// BatchResult represents a single batch result
type BatchResult struct {
	Question string
	Response string
	Error    error
}

// NewBatchClient creates a new batch client
func NewBatchClient(concurrentNum int) *BatchClient {
	return &BatchClient{
		client:        New(),
		concurrentNum: concurrentNum,
	}
}

// AskBatch processes multiple questions in parallel
func (b *BatchClient) AskBatch(questions []string) []BatchResult {
	results := make([]BatchResult, len(questions))
	
	// Create a channel for work distribution
	workCh := make(chan int, len(questions))
	for i := range questions {
		workCh <- i
	}
	close(workCh)
	
	// Create wait group for workers
	var wg sync.WaitGroup
	
	// Start workers
	for i := 0; i < b.concurrentNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for idx := range workCh {
				response, err := b.client.Ask(questions[idx])
				results[idx] = BatchResult{
					Question: questions[idx],
					Response: response,
					Error:    err,
				}
			}
		}()
	}
	
	// Wait for all workers to complete
	wg.Wait()
	
	return results
}