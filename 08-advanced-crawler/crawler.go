package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/google/uuid"
)

// CrawlRequest represents the request structure for crawling
type CrawlRequest struct {
	Domains   []string `json:"domains" binding:"required"`
	Keywords  []string `json:"keywords" binding:"required"`
	MaxPages  int      `json:"max_pages"`
	Depth     int      `json:"depth"`
	Parallel  int      `json:"parallel"`
	Delay     int      `json:"delay"` // delay in seconds
}

// CrawlResult represents a single crawl result
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

// CrawlJob represents a crawl job
type CrawlJob struct {
	ID           string        `json:"crawl_id"`
	Status       string        `json:"status"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      *time.Time    `json:"end_time,omitempty"`
	Progress     int           `json:"progress"`
	TotalResults int           `json:"total_results"`
	Results      []CrawlResult `json:"results"`
	mu           sync.RWMutex
}

// CrawlResponse represents the response structure
type CrawlResponse struct {
	CrawlID string `json:"crawl_id"`
	Message string `json:"message"`
}

// ResultsResponse represents the results response structure
type ResultsResponse struct {
	CrawlID      string        `json:"crawl_id"`
	GeneratedAt  time.Time     `json:"generated_at"`
	Status       string        `json:"status"`
	Progress     int           `json:"progress"`
	TotalResults int           `json:"total_results"`
	Results      []CrawlResult `json:"results"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      *time.Time    `json:"end_time,omitempty"`
}

// SummaryResult represents a summarized result
type SummaryResult struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	Domain     string    `json:"domain"`
	StatusCode int       `json:"status_code"`
	Timestamp  time.Time `json:"timestamp"`
}

// SummaryResponse represents the summary response structure
type SummaryResponse struct {
	CrawlID      string          `json:"crawl_id"`
	GeneratedAt  time.Time       `json:"generated_at"`
	Status       string          `json:"status"`
	TotalResults int             `json:"total_results"`
	Results      []SummaryResult `json:"results"`
}

// Global storage for crawl jobs
var crawlJobs = make(map[string]*CrawlJob)
var jobsMutex sync.RWMutex

// AdvancedCrawler represents the advanced crawler with Colly
type AdvancedCrawler struct {
	collector     *colly.Collector
	job           *CrawlJob
	keywords      []string
	maxPages      int
	pageCount     int
	mu            sync.Mutex
	allowedDomains []string
	visitedURLs   map[string]bool
}

// NewAdvancedCrawler creates a new advanced crawler instance
func NewAdvancedCrawler(domains []string, keywords []string, maxPages, depth, parallel, delay int) *AdvancedCrawler {
	// Expand domains to include www subdomains and vice versa
	expandedDomains := make([]string, 0, len(domains)*2)
	for _, domain := range domains {
		expandedDomains = append(expandedDomains, domain)
		if strings.HasPrefix(domain, "www.") {
			// If domain starts with www, add version without www
			expandedDomains = append(expandedDomains, domain[4:])
		} else {
			// If domain doesn't start with www, add www version
			expandedDomains = append(expandedDomains, "www."+domain)
		}
	}

	// Create collector with advanced configuration
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains(expandedDomains...),
	)

	// Set limits
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: parallel,
		Delay:       time.Duration(delay) * time.Second,
	})

	// Set user agent rotation
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}

	// Set random user agent
	c.UserAgent = userAgents[0]

	// Create crawl job
	job := &CrawlJob{
		ID:        uuid.New().String(),
		Status:    "running",
		StartTime: time.Now(),
		Progress:  0,
		Results:   make([]CrawlResult, 0),
	}

	crawler := &AdvancedCrawler{
		collector:      c,
		job:            job,
		keywords:       keywords,
		maxPages:       maxPages,
		pageCount:      0,
		allowedDomains: expandedDomains,
		visitedURLs:    make(map[string]bool),
	}

	// Store job globally
	jobsMutex.Lock()
	crawlJobs[job.ID] = job
	jobsMutex.Unlock()

	return crawler
}

// isAllowedDomain checks if a URL belongs to one of the allowed domains
func (ac *AdvancedCrawler) isAllowedDomain(urlStr string) bool {
	for _, domain := range ac.allowedDomains {
		if strings.Contains(urlStr, domain) {
			return true
		}
	}
	return false
}

// hasVisited checks if a URL has already been visited
func (ac *AdvancedCrawler) hasVisited(urlStr string) bool {
	return ac.visitedURLs[urlStr]
}

// markVisited marks a URL as visited
func (ac *AdvancedCrawler) markVisited(urlStr string) {
	ac.visitedURLs[urlStr] = true
}

// SetupCallbacks sets up the crawler callbacks
func (ac *AdvancedCrawler) SetupCallbacks() {
	// On HTML response
	ac.collector.OnHTML("html", func(e *colly.HTMLElement) {
		ac.mu.Lock()
		defer ac.mu.Unlock()

		// Mark this URL as visited first
		ac.markVisited(e.Request.URL.String())

		// Increment page count
		ac.pageCount++
		
		fmt.Printf("Processing page %d/%d: %s\n", ac.pageCount, ac.maxPages, e.Request.URL.String())

		if ac.pageCount > ac.maxPages {
			fmt.Printf("Reached max pages limit (%d), skipping: %s\n", ac.maxPages, e.Request.URL.String())
			return
		}

		title := e.ChildText("title")
		content := e.ChildText("body")
		
		// Check if content contains any of the keywords
		contentLower := strings.ToLower(content)
		titleLower := strings.ToLower(title)
		foundKeywords := make([]string, 0)
		
		for _, keyword := range ac.keywords {
			keywordLower := strings.ToLower(keyword)
			if strings.Contains(contentLower, keywordLower) || strings.Contains(titleLower, keywordLower) {
				foundKeywords = append(foundKeywords, keyword)
			}
		}

		// Store all results, but mark which ones contain keywords
		// This allows us to see what pages are being crawled
		result := CrawlResult{
			URL:        e.Request.URL.String(),
			Title:      title,
			Content:    content[:min(500, len(content))], // Limit content length
			Domain:     e.Request.URL.Host,
			Keywords:   foundKeywords, // Will be empty if no keywords found
			Timestamp:  time.Now(),
			StatusCode: 200,
			Metadata: map[string]string{
				"user_agent":      e.Request.Headers.Get("User-Agent"),
				"method":          "GET",
				"keywords_found":  fmt.Sprintf("%d", len(foundKeywords)),
				"content_length":  fmt.Sprintf("%d", len(content)),
			},
		}

		ac.job.mu.Lock()
		ac.job.Results = append(ac.job.Results, result)
		ac.job.TotalResults = len(ac.job.Results)
		ac.job.Progress = (ac.pageCount * 100) / ac.maxPages
		ac.job.mu.Unlock()

		fmt.Printf("Stored result #%d: %s (Title: %s, Keywords found: %d, Content length: %d)\n", 
			len(ac.job.Results), e.Request.URL.String(), title, len(foundKeywords), len(content))
	})

	// On every link found - comprehensive selector for news websites
	ac.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		ac.mu.Lock()
		defer ac.mu.Unlock()

		if ac.pageCount >= ac.maxPages {
			fmt.Printf("Max pages reached (%d), skipping link discovery\n", ac.maxPages)
			return
		}

		link := e.Attr("href")
		
		// Skip empty links, javascript links, and anchors
		if link == "" || strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") || strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") {
			return
		}
		
		// Convert relative URLs to absolute
		absoluteURL := e.Request.AbsoluteURL(link)
		
		// Debug: Print all found links for analysis
		fmt.Printf("Found link: %s -> %s\n", link, absoluteURL)
		
		// Check if the link is within allowed domains
		if !ac.isAllowedDomain(absoluteURL) {
			fmt.Printf("Skipping external link: %s\n", absoluteURL)
			return
		}
		
		// Check if we've already visited this URL
		if ac.hasVisited(absoluteURL) {
			fmt.Printf("Already visited: %s\n", absoluteURL)
			return
		}
		
		// Skip if it's the same as current URL
		if absoluteURL == e.Request.URL.String() {
			fmt.Printf("Skipping same URL: %s\n", absoluteURL)
			return
		}
		
		// Only follow links that look like article URLs (contain path segments)
		if strings.Count(absoluteURL, "/") > 3 {
			fmt.Printf("Following internal link: %s\n", absoluteURL)
			e.Request.Visit(absoluteURL)
		} else {
			fmt.Printf("Skipping homepage-like URL: %s\n", absoluteURL)
		}
	})

	// On request
	ac.collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL.String())
	})

	// On error
	ac.collector.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error visiting %s: %s\n", r.Request.URL.String(), err.Error())
	})

	// On response
	ac.collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response from %s: %d\n", r.Request.URL.String(), r.StatusCode)
	})
}

// Start begins the crawling process
func (ac *AdvancedCrawler) Start(domains []string) {
	ac.SetupCallbacks()

	// Start crawling from domain homepages
	for _, domain := range domains {
		if !strings.HasPrefix(domain, "http") {
			domain = "https://" + domain
		}
		ac.collector.Visit(domain)
	}

	// Wait for all requests to finish
	ac.collector.Wait()

	// Mark job as completed
	ac.job.mu.Lock()
	ac.job.Status = "completed"
	endTime := time.Now()
	ac.job.EndTime = &endTime
	ac.job.Progress = 100
	ac.job.mu.Unlock()
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// API Handlers

// submitCrawl handles POST /api/v1/crawl
func submitCrawl(c *gin.Context) {
	var req CrawlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.MaxPages == 0 {
		req.MaxPages = 10
	}
	if req.Depth == 0 {
		req.Depth = 2
	}
	if req.Parallel == 0 {
		req.Parallel = 2
	}
	if req.Delay == 0 {
		req.Delay = 1
	}

	// Create and start crawler in goroutine
	crawler := NewAdvancedCrawler(req.Domains, req.Keywords, req.MaxPages, req.Depth, req.Parallel, req.Delay)
	
	go crawler.Start(req.Domains)

	response := CrawlResponse{
		CrawlID: crawler.job.ID,
		Message: "Advanced crawl job submitted successfully",
	}

	c.JSON(http.StatusOK, response)
}

// getResults handles GET /api/v1/results/{crawl_id}
func getResults(c *gin.Context) {
	crawlID := c.Param("crawl_id")
	format := c.Query("format")

	jobsMutex.RLock()
	job, exists := crawlJobs[crawlID]
	jobsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Crawl job not found"})
		return
	}

	job.mu.RLock()
	defer job.mu.RUnlock()

	if format == "summary" {
		summaryResults := make([]SummaryResult, len(job.Results))
		for i, result := range job.Results {
			summaryResults[i] = SummaryResult{
				URL:        result.URL,
				Title:      result.Title,
				Domain:     result.Domain,
				StatusCode: result.StatusCode,
				Timestamp:  result.Timestamp,
			}
		}

		response := SummaryResponse{
			CrawlID:      job.ID,
			GeneratedAt:  time.Now(),
			Status:       job.Status,
			TotalResults: job.TotalResults,
			Results:      summaryResults,
		}

		c.JSON(http.StatusOK, response)
		return
	}

	response := ResultsResponse{
		CrawlID:      job.ID,
		GeneratedAt:  time.Now(),
		Status:       job.Status,
		Progress:     job.Progress,
		TotalResults: job.TotalResults,
		Results:      job.Results,
		StartTime:    job.StartTime,
		EndTime:      job.EndTime,
	}

	c.JSON(http.StatusOK, response)
}

// getStatus handles GET /api/v1/status/{crawl_id}
func getStatus(c *gin.Context) {
	crawlID := c.Param("crawl_id")

	jobsMutex.RLock()
	job, exists := crawlJobs[crawlID]
	jobsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Crawl job not found"})
		return
	}

	job.mu.RLock()
	defer job.mu.RUnlock()

	status := gin.H{
		"crawl_id":      job.ID,
		"status":        job.Status,
		"progress":      job.Progress,
		"total_results": job.TotalResults,
		"start_time":    job.StartTime,
	}

	if job.EndTime != nil {
		status["end_time"] = *job.EndTime
	}

	c.JSON(http.StatusOK, status)
}

func main() {
	// Create Gin router
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/crawl", submitCrawl)
		api.GET("/results/:crawl_id", getResults)
		api.GET("/status/:crawl_id", getStatus)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "advanced-crawler",
			"version": "1.0.0",
		})
	})

	fmt.Println("🚀 Advanced Crawler API starting on :8082")
	fmt.Println("📚 Endpoints:")
	fmt.Println("  POST /api/v1/crawl - Submit crawl job")
	fmt.Println("  GET  /api/v1/results/{crawl_id} - Get crawl results")
	fmt.Println("  GET  /api/v1/results/{crawl_id}?format=summary - Get summary results")
	fmt.Println("  GET  /api/v1/status/{crawl_id} - Get crawl status")
	fmt.Println("  GET  /health - Health check")

	log.Fatal(http.ListenAndServe(":8082", r))
}