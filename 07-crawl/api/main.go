package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"sync"
	"math/rand"

	"crawler-api/urlfrontier"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CrawlRequest represents a crawl request from the API
type CrawlRequest struct {
	Keywords    []string  `json:"keywords" binding:"required"`
	Domains     []string  `json:"domains" binding:"required"`
	StartDate   *string   `json:"start_date,omitempty"`
	EndDate     *string   `json:"end_date,omitempty"`
	MaxDepth    int       `json:"max_depth,omitempty"`
	MaxPages    int       `json:"max_pages,omitempty"`
}

// CrawlResponse represents the response after submitting a crawl request
type CrawlResponse struct {
	CrawlID   string `json:"crawl_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// CrawlStatus represents the status of a crawl job
type CrawlStatus struct {
	CrawlID     string    `json:"crawl_id"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
	TotalURLs   int       `json:"total_urls"`
	ProcessedURLs int     `json:"processed_urls"`
	StartTime   time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Results     []CrawlResult `json:"results,omitempty"`
}

// CrawlResult represents a single crawled page result
type CrawlResult struct {
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Domain      string            `json:"domain"`
	Keywords    []string          `json:"keywords"`
	Timestamp   time.Time         `json:"timestamp"`
	StatusCode  int               `json:"status_code"`
	Metadata    map[string]string `json:"metadata"`
}

// URLFrontierClient handles communication with URLFrontier service
type URLFrontierClient struct {
	client *urlfrontier.Client
}

// CrawlManager manages crawl jobs and their status
type CrawlManager struct {
	jobs           map[string]*CrawlStatus
	urlFrontier    *URLFrontierClient
	resultStore    *ResultStore
	mutex          sync.RWMutex
}

// ResultStore handles storage and retrieval of crawl results
type ResultStore struct {
	results map[string][]CrawlResult
	mutex   sync.RWMutex
}

// NewResultStore creates a new result store
func NewResultStore() *ResultStore {
	return &ResultStore{
		results: make(map[string][]CrawlResult),
	}
}

// AddResult adds a crawl result to the store
func (rs *ResultStore) AddResult(crawlID string, result CrawlResult) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	
	if rs.results[crawlID] == nil {
		rs.results[crawlID] = make([]CrawlResult, 0)
	}
	rs.results[crawlID] = append(rs.results[crawlID], result)
}

// GetResults retrieves results for a crawl ID with pagination
func (rs *ResultStore) GetResults(crawlID string, page, limit int) ([]CrawlResult, int) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	results, exists := rs.results[crawlID]
	if !exists {
		return []CrawlResult{}, 0
	}
	
	total := len(results)
	start := (page - 1) * limit
	end := start + limit
	
	if start >= total {
		return []CrawlResult{}, total
	}
	
	if end > total {
		end = total
	}
	
	return results[start:end], total
}

// GetAllResults returns all results for a crawl ID
func (rs *ResultStore) GetAllResults(crawlID string) []CrawlResult {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	results, exists := rs.results[crawlID]
	if !exists {
		return []CrawlResult{}
	}
	
	return results
}

// NewCrawlManager creates a new crawl manager
func NewCrawlManager() *CrawlManager {
	return &CrawlManager{
		jobs:        make(map[string]*CrawlStatus),
		resultStore: NewResultStore(),
	}
}

// InitURLFrontierClient initializes connection to URLFrontier service
func (cm *CrawlManager) InitURLFrontierClient(address string) error {
	client, err := urlfrontier.NewClient(address)
	if err != nil {
		return fmt.Errorf("failed to connect to URLFrontier: %v", err)
	}
	
	cm.urlFrontier = &URLFrontierClient{
		client: client,
	}
	
	return nil
}

// SubmitCrawlJob submits a new crawl job
func (cm *CrawlManager) SubmitCrawlJob(req *CrawlRequest) (*CrawlResponse, error) {
	crawlID := uuid.New().String()
	
	// Create crawl status
	status := &CrawlStatus{
		CrawlID:       crawlID,
		Status:        "submitted",
		Progress:      0,
		TotalURLs:     0,
		ProcessedURLs: 0,
		StartTime:     time.Now(),
		Results:       []CrawlResult{},
	}
	
	cm.mutex.Lock()
	cm.jobs[crawlID] = status
	cm.mutex.Unlock()
	
	// Generate seed URLs based on domains and keywords
	seedURLs := cm.generateSeedURLs(req.Domains, req.Keywords)
	
	// Submit URLs to URLFrontier (if available)
	if cm.urlFrontier != nil {
		err := cm.submitURLsToFrontier(crawlID, seedURLs, req)
		if err != nil {
			status.Status = "failed"
			return nil, fmt.Errorf("failed to submit URLs to frontier: %v", err)
		}
	}
	
	status.Status = "running"
	status.TotalURLs = len(seedURLs)
	
	// Start simulating crawl results for demonstration
	cm.SimulateCrawlResults(crawlID, req.Domains, req.Keywords)

	return &CrawlResponse{
		CrawlID:   crawlID,
		Status:    "submitted",
		Message:   fmt.Sprintf("Crawl job submitted successfully with %d seed URLs", len(seedURLs)),
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

// GetCrawlStatus retrieves the status of a crawl job
func (cm *CrawlManager) GetCrawlStatus(crawlID string) (*CrawlStatus, error) {
	cm.mutex.RLock()
	status, exists := cm.jobs[crawlID]
	cm.mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("crawl job not found")
	}
	
	// Update status from URLFrontier if available
	if cm.urlFrontier != nil {
		cm.updateCrawlStatusFromFrontier(status)
	}
	
	return status, nil
}

// updateCrawlStatusFromFrontier updates crawl status from URLFrontier
func (cm *CrawlManager) updateCrawlStatusFromFrontier(status *CrawlStatus) {
	if cm.urlFrontier == nil || cm.urlFrontier.client == nil {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Get queue statistics for this crawl
	queueStats, err := cm.urlFrontier.client.GetQueueStats(ctx, status.CrawlID)
	if err != nil {
		log.Printf("Failed to get queue stats for crawl %s: %v", status.CrawlID, err)
		return
	}
	
	// Update status based on queue statistics
	status.ProcessedURLs = queueStats.Completed
	if status.TotalURLs > 0 {
		status.Progress = (status.ProcessedURLs * 100) / status.TotalURLs
	}
	
	// Update status based on progress
	if queueStats.ActiveURLs == 0 && queueStats.InProcess == 0 && queueStats.Completed > 0 {
		status.Status = "completed"
		if status.EndTime == nil {
			now := time.Now()
			status.EndTime = &now
		}
	} else if queueStats.ActiveURLs > 0 || queueStats.InProcess > 0 {
		status.Status = "running"
	}
}

// API Handlers

func setupRoutes(cm *CrawlManager) *gin.Engine {
	r := gin.Default()
	
	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	api := r.Group("/api/v1")
	{
		api.POST("/crawl", handleSubmitCrawl(cm))
		api.GET("/crawl/:crawl_id", handleGetCrawlStatus(cm))
		api.GET("/crawl/:crawl_id/results", handleGetCrawlResults(cm))
		api.GET("/crawl", handleListCrawls(cm))
		api.DELETE("/crawl/:crawl_id", handleCancelCrawl(cm))
		
		// New endpoint for getting all crawl results in JSON format
		api.GET("/results/:crawl_id", handleGetAllCrawlResults(cm))
	}
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	
	return r
}

func handleSubmitCrawl(cm *CrawlManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CrawlRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
				"details": err.Error(),
			})
			return
		}
		
		// Validate request
		if len(req.Keywords) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "At least one keyword is required",
			})
			return
		}
		
		if len(req.Domains) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "At least one domain is required",
			})
			return
		}
		
		// Set defaults
		if req.MaxDepth == 0 {
			req.MaxDepth = 3
		}
		if req.MaxPages == 0 {
			req.MaxPages = 100
		}
		
		// Validate date range if provided
		if req.StartDate != nil && req.EndDate != nil {
			startDate, err1 := time.Parse("2006-01-02", *req.StartDate)
			endDate, err2 := time.Parse("2006-01-02", *req.EndDate)
			
			if err1 != nil || err2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid date format. Use YYYY-MM-DD",
				})
				return
			}
			
			if startDate.After(endDate) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Start date must be before end date",
				})
				return
			}
		}
		
		response, err := cm.SubmitCrawlJob(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to submit crawl job",
				"details": err.Error(),
			})
			return
		}
		
		c.JSON(http.StatusCreated, response)
	}
}

func handleGetCrawlStatus(cm *CrawlManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		crawlID := c.Param("crawl_id")
		
		status, err := cm.GetCrawlStatus(crawlID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Crawl job not found",
				"crawl_id": crawlID,
			})
			return
		}
		
		c.JSON(http.StatusOK, status)
	}
}

func handleGetCrawlResults(cm *CrawlManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		crawlID := c.Param("crawl_id")
		
		status, err := cm.GetCrawlStatus(crawlID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Crawl job not found",
				"crawl_id": crawlID,
			})
			return
		}
		
		// Parse query parameters for pagination
		page := 1
		limit := 50
		
		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
				limit = l
			}
		}
		
		// Calculate pagination
		start := (page - 1) * limit
		end := start + limit
		
		results := status.Results
		total := len(results)
		
		if start >= total {
			results = []CrawlResult{}
		} else if end > total {
			results = results[start:]
		} else {
			results = results[start:end]
		}
		
		c.JSON(http.StatusOK, gin.H{
			"crawl_id": crawlID,
			"results": results,
			"pagination": gin.H{
				"page": page,
				"limit": limit,
				"total": total,
				"pages": (total + limit - 1) / limit,
			},
		})
	}
}

func handleListCrawls(cm *CrawlManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var crawls []map[string]interface{}
		
		for crawlID, status := range cm.jobs {
			crawls = append(crawls, map[string]interface{}{
				"crawl_id": crawlID,
				"status": status.Status,
				"progress": status.Progress,
				"total_urls": status.TotalURLs,
				"processed_urls": status.ProcessedURLs,
				"start_time": status.StartTime,
				"end_time": status.EndTime,
			})
		}
		
		c.JSON(http.StatusOK, gin.H{
			"crawls": crawls,
			"total": len(crawls),
		})
	}
}

func handleCancelCrawl(cm *CrawlManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		crawlID := c.Param("crawl_id")
		
		status, exists := cm.jobs[crawlID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Crawl job not found",
				"crawl_id": crawlID,
			})
			return
		}
		
		if status.Status == "completed" || status.Status == "failed" || status.Status == "cancelled" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot cancel crawl job in current status",
				"status": status.Status,
			})
			return
		}
		
		// Cancel the crawl job (placeholder implementation)
		status.Status = "cancelled"
		now := time.Now()
		status.EndTime = &now
		
		c.JSON(http.StatusOK, gin.H{
			"message": "Crawl job cancelled successfully",
			"crawl_id": crawlID,
		})
	}
}

func main() {
	// Initialize crawl manager
	cm := NewCrawlManager()
	
	// Initialize URLFrontier client
	frontierAddress := "host.docker.internal:7071"
	if err := cm.InitURLFrontierClient(frontierAddress); err != nil {
		log.Printf("Warning: Failed to connect to URLFrontier: %v", err)
		log.Println("API will start but crawl functionality may be limited")
	}
	
	// Setup routes
	r := setupRoutes(cm)
	
	// Start server
	port := ":8081"
	log.Printf("Starting Crawler API server on port %s", port)
	log.Printf("Health check: http://localhost%s/health", port)
	log.Printf("API documentation: http://localhost%s/api/v1", port)
	
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// handleGetAllCrawlResults returns all crawl results in JSON format
func handleGetAllCrawlResults(cm *CrawlManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		crawlID := c.Param("crawl_id")
		
		// Check if crawl exists
		cm.mutex.RLock()
		status, exists := cm.jobs[crawlID]
		cm.mutex.RUnlock()
		
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Crawl job not found",
				"crawl_id": crawlID,
			})
			return
		}
		
		// Get all results
		results := cm.resultStore.GetAllResults(crawlID)
		
		// Parse query parameters for filtering
		format := c.DefaultQuery("format", "detailed") // detailed or summary
		
		if format == "summary" {
			// Return summary format
			summaryResults := make([]gin.H, len(results))
			for i, result := range results {
				summaryResults[i] = gin.H{
					"url":         result.URL,
					"title":       result.Title,
					"domain":      result.Domain,
					"status_code": result.StatusCode,
					"timestamp":   result.Timestamp.Format(time.RFC3339),
				}
			}
			
			c.JSON(http.StatusOK, gin.H{
				"crawl_id": crawlID,
				"status":   status.Status,
				"total_results": len(results),
				"results":  summaryResults,
				"generated_at": time.Now().Format(time.RFC3339),
			})
		} else {
			// Return detailed format
			c.JSON(http.StatusOK, gin.H{
				"crawl_id": crawlID,
				"status":   status.Status,
				"progress": status.Progress,
				"total_urls": status.TotalURLs,
				"processed_urls": status.ProcessedURLs,
				"start_time": status.StartTime.Format(time.RFC3339),
				"end_time": func() *string {
					if status.EndTime != nil {
						t := status.EndTime.Format(time.RFC3339)
						return &t
					}
					return nil
				}(),
				"total_results": len(results),
				"results": results,
				"generated_at": time.Now().Format(time.RFC3339),
			})
		}
	}
}

// generateSampleResults creates sample crawl results for demonstration
func (cm *CrawlManager) generateSampleResults(domains []string, keywords []string) []CrawlResult {
	results := make([]CrawlResult, 0)
	
	samplePages := []string{
		"/", "/about", "/products", "/services", "/contact", 
		"/blog", "/news", "/support", "/pricing", "/features",
	}
	
	sampleTitles := []string{
		"Home Page", "About Us", "Our Products", "Services", "Contact Us",
		"Blog", "Latest News", "Support Center", "Pricing Plans", "Features",
	}
	
	for i, domain := range domains {
		for j, page := range samplePages {
			if len(results) >= 20 { // Limit to 20 results for demo
				break
			}
			
			url := fmt.Sprintf("https://%s%s", domain, page)
			title := fmt.Sprintf("%s - %s", sampleTitles[j%len(sampleTitles)], domain)
			content := fmt.Sprintf("This is sample content from %s containing keywords: %s. Lorem ipsum dolor sit amet, consectetur adipiscing elit.", 
				url, strings.Join(keywords, ", "))
			
			result := CrawlResult{
				URL:        url,
				Title:      title,
				Content:    content,
				Domain:     domain,
				Keywords:   keywords,
				Timestamp:  time.Now().Add(time.Duration(i*j) * time.Second),
				StatusCode: 200,
				Metadata: map[string]string{
					"content_type":   "text/html",
					"content_length": fmt.Sprintf("%d", len(content)),
					"crawl_depth":    fmt.Sprintf("%d", rand.Intn(3)+1),
				},
			}
			
			results = append(results, result)
		}
	}
	
	return results
}

// SimulateCrawlResults simulates crawl results for demonstration
func (cm *CrawlManager) SimulateCrawlResults(crawlID string, domains []string, keywords []string) {
	go func() {
		// Wait a bit before starting to simulate processing
		time.Sleep(2 * time.Second)
		
		// Generate some sample results
		sampleResults := cm.generateSampleResults(domains, keywords)
		
		for i, result := range sampleResults {
			// Add delay between results to simulate real crawling
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
			
			// Add result to store
			cm.resultStore.AddResult(crawlID, result)
			
			// Update crawl status
			cm.mutex.Lock()
			if status, exists := cm.jobs[crawlID]; exists {
				status.ProcessedURLs = i + 1
				if status.TotalURLs > 0 {
					status.Progress = (status.ProcessedURLs * 100) / status.TotalURLs
				}
				status.Results = cm.resultStore.GetAllResults(crawlID)
			}
			cm.mutex.Unlock()
		}
		
		// Mark as completed
		cm.mutex.Lock()
		if status, exists := cm.jobs[crawlID]; exists {
			status.Status = "completed"
			now := time.Now()
			status.EndTime = &now
		}
		cm.mutex.Unlock()
	}()
}

// generateSeedURLs creates seed URLs from domains and keywords
func (cm *CrawlManager) generateSeedURLs(domains []string, keywords []string) []string {
	var seedURLs []string
	
	for _, domain := range domains {
		// Add base domain
		if !strings.HasPrefix(domain, "http") {
			domain = "https://" + domain
		}
		seedURLs = append(seedURLs, domain)
		
		// Add search URLs with keywords (example patterns)
		for _, keyword := range keywords {
			searchURL := fmt.Sprintf("%s/search?q=%s", domain, strings.ReplaceAll(keyword, " ", "+"))
			seedURLs = append(seedURLs, searchURL)
		}
	}
	
	return seedURLs
}

// submitURLsToFrontier submits URLs to the URLFrontier service
func (cm *CrawlManager) submitURLsToFrontier(crawlID string, urls []string, req *CrawlRequest) error {
	if cm.urlFrontier == nil || cm.urlFrontier.client == nil {
		log.Printf("URLFrontier client not available, simulating submission for %d URLs", len(urls))
		return nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Prepare date range metadata
	dateRange := make(map[string]string)
	if req.StartDate != nil {
		dateRange["start_date"] = *req.StartDate
	}
	if req.EndDate != nil {
		dateRange["end_date"] = *req.EndDate
	}
	
	// Create URL requests with metadata
	var urlRequests []urlfrontier.URLRequest
	for _, url := range urls {
		urlReq := urlfrontier.CreateURLRequest(url, crawlID, req.Keywords, req.Domains, dateRange)
		urlRequests = append(urlRequests, urlReq)
	}
	
	// Submit URLs to URLFrontier
	err := cm.urlFrontier.client.SubmitURLs(ctx, urlRequests)
	if err != nil {
		return fmt.Errorf("failed to submit URLs to URLFrontier: %v", err)
	}
	
	log.Printf("Successfully submitted %d URLs to URLFrontier for crawl %s", len(urls), crawlID)
	return nil
}