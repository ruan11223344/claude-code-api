package claude

import (
	"sync"
	"time"
)

// MonitoredClient wraps Client with monitoring capabilities
type MonitoredClient struct {
	*Client
	metrics Metrics
	mu      sync.RWMutex
}

// Metrics holds monitoring metrics
type Metrics struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	TotalDuration   time.Duration
	AverageDuration time.Duration
}

// NewMonitoredClient creates a new monitored client
func NewMonitoredClient(opts ...Option) *MonitoredClient {
	return &MonitoredClient{
		Client: New(opts...),
	}
}

// Ask calls Claude and tracks metrics
func (c *MonitoredClient) Ask(question string) (string, error) {
	start := time.Now()
	
	response, err := c.Client.Ask(question)
	
	duration := time.Since(start)
	
	c.mu.Lock()
	c.metrics.TotalRequests++
	c.metrics.TotalDuration += duration
	
	if err == nil {
		c.metrics.SuccessRequests++
	} else {
		c.metrics.FailedRequests++
	}
	
	// Update average duration
	if c.metrics.TotalRequests > 0 {
		c.metrics.AverageDuration = c.metrics.TotalDuration / time.Duration(c.metrics.TotalRequests)
	}
	c.mu.Unlock()
	
	return response, err
}

// GetMetrics returns current metrics
func (c *MonitoredClient) GetMetrics() Metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.metrics
}

// ResetMetrics resets all metrics
func (c *MonitoredClient) ResetMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = Metrics{}
}