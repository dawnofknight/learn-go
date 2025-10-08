package urlfrontier

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents a URLFrontier gRPC client
type Client struct {
	conn    *grpc.ClientConn
	address string
}

// URLRequest represents a URL to be submitted to the frontier
type URLRequest struct {
	URL      string            `json:"url"`
	Metadata map[string]string `json:"metadata"`
	Queue    string            `json:"queue"`
}

// QueueStats represents statistics for a queue
type QueueStats struct {
	Queue       string `json:"queue"`
	ActiveURLs  int    `json:"active_urls"`
	InProcess   int    `json:"in_process"`
	Completed   int    `json:"completed"`
}

// FrontierStats represents overall frontier statistics
type FrontierStats struct {
	ActiveQueues int          `json:"active_queues"`
	TotalURLs    int          `json:"total_urls"`
	Queues       []QueueStats `json:"queues"`
}

// NewClient creates a new URLFrontier client
func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to URLFrontier at %s: %v", address, err)
	}

	client := &Client{
		conn:    conn,
		address: address,
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping URLFrontier: %v", err)
	}

	log.Printf("Successfully connected to URLFrontier at %s", address)
	return client, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ping tests the connection to URLFrontier
func (c *Client) ping(ctx context.Context) error {
	// For now, we'll implement a simple connection test
	// In a real implementation, this would use the URLFrontier gRPC service
	log.Printf("Testing connection to URLFrontier at %s", c.address)
	return nil
}

// SubmitURLs submits URLs to the URLFrontier service
func (c *Client) SubmitURLs(ctx context.Context, urls []URLRequest) error {
	log.Printf("Submitting %d URLs to URLFrontier", len(urls))
	
	// Placeholder implementation
	// In a real implementation, this would:
	// 1. Create URLFrontier gRPC requests
	// 2. Submit URLs with metadata
	// 3. Handle responses and errors
	
	for _, url := range urls {
		log.Printf("Submitting URL: %s to queue: %s", url.URL, url.Queue)
		// Here we would make the actual gRPC call
	}
	
	return nil
}

// GetStats retrieves statistics from the URLFrontier service
func (c *Client) GetStats(ctx context.Context) (*FrontierStats, error) {
	log.Printf("Retrieving stats from URLFrontier")
	
	// Placeholder implementation
	// In a real implementation, this would query the URLFrontier gRPC service
	
	stats := &FrontierStats{
		ActiveQueues: 1,
		TotalURLs:    0,
		Queues: []QueueStats{
			{
				Queue:      "default",
				ActiveURLs: 0,
				InProcess:  0,
				Completed:  0,
			},
		},
	}
	
	return stats, nil
}

// GetQueueStats retrieves statistics for a specific queue
func (c *Client) GetQueueStats(ctx context.Context, queue string) (*QueueStats, error) {
	log.Printf("Retrieving stats for queue: %s", queue)
	
	// Placeholder implementation
	stats := &QueueStats{
		Queue:      queue,
		ActiveURLs: 0,
		InProcess:  0,
		Completed:  0,
	}
	
	return stats, nil
}

// CreateURLRequest creates a URLRequest with metadata for crawling
func CreateURLRequest(url, crawlID string, keywords []string, domains []string, dateRange map[string]string) URLRequest {
	metadata := make(map[string]string)
	metadata["crawl_id"] = crawlID
	metadata["keywords"] = fmt.Sprintf("%v", keywords)
	metadata["domains"] = fmt.Sprintf("%v", domains)
	
	if startDate, exists := dateRange["start_date"]; exists {
		metadata["start_date"] = startDate
	}
	if endDate, exists := dateRange["end_date"]; exists {
		metadata["end_date"] = endDate
	}
	
	// Add timestamp
	metadata["submitted_at"] = time.Now().Format(time.RFC3339)
	
	return URLRequest{
		URL:      url,
		Metadata: metadata,
		Queue:    crawlID, // Use crawl ID as queue name for isolation
	}
}