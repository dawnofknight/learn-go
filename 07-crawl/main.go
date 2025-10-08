package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// URLStatus represents the status of a URL during crawling
type URLStatus int

const (
	StatusPending URLStatus = iota
	StatusFetched
	StatusError
	StatusRedirect
)

// CrawlResult represents the result of crawling a URL
type CrawlResult struct {
	URL         string
	Status      URLStatus
	Content     string
	Links       []string
	Error       error
	StatusCode  int
	RedirectURL string
}

// URLFrontier manages the queue of URLs to be crawled
type URLFrontier struct {
	urls     chan string
	visited  map[string]bool
	mu       sync.RWMutex
	maxDepth int
	depth    map[string]int
}

// NewURLFrontier creates a new URL frontier
func NewURLFrontier(maxDepth int) *URLFrontier {
	return &URLFrontier{
		urls:     make(chan string, 1000),
		visited:  make(map[string]bool),
		maxDepth: maxDepth,
		depth:    make(map[string]int),
	}
}

// AddURL adds a URL to the frontier if not already visited
func (uf *URLFrontier) AddURL(rawURL string, currentDepth int) {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	// Normalize URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	normalizedURL := parsedURL.String()

	// Check if already visited or max depth exceeded
	if uf.visited[normalizedURL] || currentDepth >= uf.maxDepth {
		return
	}

	uf.visited[normalizedURL] = true
	uf.depth[normalizedURL] = currentDepth

	select {
	case uf.urls <- normalizedURL:
	default:
		// Channel is full, skip this URL
	}
}

// GetURL retrieves the next URL to crawl
func (uf *URLFrontier) GetURL() (string, int, bool) {
	select {
	case url := <-uf.urls:
		uf.mu.RLock()
		depth := uf.depth[url]
		uf.mu.RUnlock()
		return url, depth, true
	default:
		return "", 0, false
	}
}

// Close closes the URL frontier
func (uf *URLFrontier) Close() {
	close(uf.urls)
}

// Fetcher handles HTTP requests with politeness and rate limiting
type Fetcher struct {
	client      *http.Client
	userAgent   string
	rateLimiter map[string]time.Time
	mu          sync.Mutex
	delay       time.Duration
}

// NewFetcher creates a new fetcher with rate limiting
func NewFetcher(delay time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent:   "GoCrawler/1.0 (+https://example.com/bot)",
		rateLimiter: make(map[string]time.Time),
		delay:       delay,
	}
}

// Fetch retrieves content from a URL with politeness
func (f *Fetcher) Fetch(rawURL string) *CrawlResult {
	result := &CrawlResult{
		URL:    rawURL,
		Status: StatusPending,
	}

	// Parse URL to get hostname for rate limiting
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		result.Status = StatusError
		result.Error = err
		return result
	}

	hostname := parsedURL.Hostname()

	// Apply rate limiting per hostname
	f.mu.Lock()
	if lastRequest, exists := f.rateLimiter[hostname]; exists {
		if time.Since(lastRequest) < f.delay {
			time.Sleep(f.delay - time.Since(lastRequest))
		}
	}
	f.rateLimiter[hostname] = time.Now()
	f.mu.Unlock()

	// Create request
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		result.Status = StatusError
		result.Error = err
		return result
	}

	req.Header.Set("User-Agent", f.userAgent)

	// Perform request
	resp, err := f.client.Do(req)
	if err != nil {
		result.Status = StatusError
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// Handle redirects
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		result.Status = StatusRedirect
		result.RedirectURL = resp.Header.Get("Location")
		return result
	}

	// Read content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = StatusError
		result.Error = err
		return result
	}

	result.Content = string(body)
	result.Status = StatusFetched
	return result
}

// Parser extracts links and content from HTML
type Parser struct {
	baseURL *url.URL
}

// NewParser creates a new HTML parser
func NewParser(baseURL string) (*Parser, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &Parser{baseURL: parsedURL}, nil
}

// Parse extracts links from HTML content
func (p *Parser) Parse(content string, currentURL string) []string {
	var links []string

	// Parse current URL for resolving relative links
	currentParsedURL, err := url.Parse(currentURL)
	if err != nil {
		return links
	}

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return links
	}

	// Extract links recursively
	p.extractLinks(doc, currentParsedURL, &links)

	return links
}

// extractLinks recursively extracts links from HTML nodes
func (p *Parser) extractLinks(n *html.Node, baseURL *url.URL, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				// Resolve relative URLs
				if resolvedURL, err := baseURL.Parse(attr.Val); err == nil {
					// Only include HTTP/HTTPS URLs
					if resolvedURL.Scheme == "http" || resolvedURL.Scheme == "https" {
						*links = append(*links, resolvedURL.String())
					}
				}
			}
		}
	}

	// Recursively process child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.extractLinks(c, baseURL, links)
	}
}

// Indexer handles the indexing/output of crawled content
type Indexer struct {
	output io.Writer
}

// NewIndexer creates a new indexer
func NewIndexer(output io.Writer) *Indexer {
	return &Indexer{output: output}
}

// Index processes and outputs the crawled content
func (i *Indexer) Index(result *CrawlResult) {
	switch result.Status {
	case StatusFetched:
		// Extract text content (simplified)
		text := i.extractText(result.Content)
		fmt.Fprintf(i.output, "=== CRAWLED: %s ===\n", result.URL)
		fmt.Fprintf(i.output, "Status Code: %d\n", result.StatusCode)
		fmt.Fprintf(i.output, "Content Length: %d bytes\n", len(result.Content))
		fmt.Fprintf(i.output, "Links Found: %d\n", len(result.Links))
		fmt.Fprintf(i.output, "Text Preview: %s\n", i.truncate(text, 200))
		fmt.Fprintf(i.output, "Links: %v\n", result.Links[:min(len(result.Links), 5)])
		fmt.Fprintln(i.output, "")
	case StatusError:
		fmt.Fprintf(i.output, "ERROR crawling %s: %v\n", result.URL, result.Error)
	case StatusRedirect:
		fmt.Fprintf(i.output, "REDIRECT %s -> %s\n", result.URL, result.RedirectURL)
	}
}

// extractText extracts plain text from HTML (simplified)
func (i *Indexer) extractText(htmlContent string) string {
	// Remove HTML tags using regex (simplified approach)
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(htmlContent, " ")
	
	// Clean up whitespace
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	
	return strings.TrimSpace(text)
}

// truncate truncates text to specified length
func (i *Indexer) truncate(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[:length] + "..."
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Crawler orchestrates the crawling process
type Crawler struct {
	frontier *URLFrontier
	fetcher  *Fetcher
	parser   *Parser
	indexer  *Indexer
	workers  int
}

// NewCrawler creates a new crawler
func NewCrawler(maxDepth, workers int, delay time.Duration) *Crawler {
	return &Crawler{
		frontier: NewURLFrontier(maxDepth),
		fetcher:  NewFetcher(delay),
		indexer:  NewIndexer(os.Stdout),
		workers:  workers,
	}
}

// Crawl starts the crawling process
func (c *Crawler) Crawl(startURL string) error {
	// Initialize parser with base URL
	parser, err := NewParser(startURL)
	if err != nil {
		return err
	}
	c.parser = parser

	// Add initial URL
	c.frontier.AddURL(startURL, 0)

	// Start worker goroutines
	var wg sync.WaitGroup
	results := make(chan *CrawlResult, 100)

	// Start workers
	for i := 0; i < c.workers; i++ {
		wg.Add(1)
		go c.worker(&wg, results)
	}

	// Start result processor
	go c.processResults(results)

	// Wait for all workers to complete
	wg.Wait()
	close(results)

	return nil
}

// worker processes URLs from the frontier
func (c *Crawler) worker(wg *sync.WaitGroup, results chan<- *CrawlResult) {
	defer wg.Done()

	for {
		url, depth, ok := c.frontier.GetURL()
		if !ok {
			// No more URLs, wait a bit and try again
			time.Sleep(100 * time.Millisecond)
			if _, _, stillOk := c.frontier.GetURL(); !stillOk {
				break
			}
			continue
		}

		// Fetch the URL
		result := c.fetcher.Fetch(url)

		// Parse links if successful
		if result.Status == StatusFetched {
			links := c.parser.Parse(result.Content, url)
			result.Links = links

			// Add new URLs to frontier
			for _, link := range links {
				c.frontier.AddURL(link, depth+1)
			}
		}

		// Send result for processing
		select {
		case results <- result:
		default:
			// Results channel is full, skip
		}
	}
}

// processResults processes crawl results
func (c *Crawler) processResults(results <-chan *CrawlResult) {
	for result := range results {
		c.indexer.Index(result)
	}
}

func main() {
	fmt.Println("🕷️  Go Web Crawler (inspired by StormCrawler)")
	fmt.Println("============================================")

	// Get URL from user input
	fmt.Print("Enter the URL to crawl: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	startURL := strings.TrimSpace(scanner.Text())

	if startURL == "" {
		fmt.Println("❌ No URL provided. Exiting.")
		return
	}

	// Validate URL
	if _, err := url.Parse(startURL); err != nil {
		fmt.Printf("❌ Invalid URL: %v\n", err)
		return
	}

	// Add scheme if missing
	if !strings.HasPrefix(startURL, "http://") && !strings.HasPrefix(startURL, "https://") {
		startURL = "https://" + startURL
	}

	fmt.Printf("🚀 Starting crawl of: %s\n", startURL)
	fmt.Println("📊 Configuration:")
	fmt.Println("   - Max Depth: 2")
	fmt.Println("   - Workers: 3")
	fmt.Println("   - Delay: 1s between requests per host")
	fmt.Println()

	// Create and start crawler
	crawler := NewCrawler(2, 3, 1*time.Second)
	
	start := time.Now()
	if err := crawler.Crawl(startURL); err != nil {
		fmt.Printf("❌ Crawl failed: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Crawl completed in %v\n", time.Since(start))
}