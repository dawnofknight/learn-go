# Advanced Web Crawler with Colly

This is an advanced web crawler implementation using the [Colly](https://go-colly.org/) framework, featuring parallel crawling, rate limiting, user agent rotation, and comprehensive API endpoints.

## Features

### Advanced Crawling Capabilities
- **Parallel Crawling**: Configurable parallelism for faster crawling
- **Rate Limiting**: Built-in delay mechanisms to respect server resources
- **User Agent Rotation**: Multiple user agents to avoid detection
- **Domain Filtering**: Restrict crawling to specific domains
- **Keyword Filtering**: Only collect pages containing specified keywords
- **Depth Control**: Limit crawling depth
- **Page Limits**: Set maximum pages to crawl

### API Endpoints
- `POST /api/v1/crawl` - Submit a new crawl job
- `GET /api/v1/results/{crawl_id}` - Get detailed crawl results
- `GET /api/v1/results/{crawl_id}?format=summary` - Get summarized results
- `GET /api/v1/status/{crawl_id}` - Get crawl job status
- `GET /health` - Health check endpoint

## Installation

1. Install dependencies:
```bash
go mod init advanced-crawler
go get -u github.com/gocolly/colly/...
go get github.com/gin-gonic/gin
go get github.com/google/uuid
```

2. Run the crawler:
```bash
go run crawler.go
```

The API will be available at `http://localhost:8082`

## Usage Examples

### Submit a Crawl Job
```bash
curl -X POST http://localhost:8082/api/v1/crawl \
  -H "Content-Type: application/json" \
  -d '{
    "domains": ["kompas.com", "detik.com"],
    "keywords": ["teknologi", "programming"],
    "max_pages": 20,
    "depth": 3,
    "parallel": 4,
    "delay": 2
  }'
```

### Get Crawl Results
```bash
# Detailed results
curl http://localhost:8082/api/v1/results/{crawl_id}

# Summary results
curl http://localhost:8082/api/v1/results/{crawl_id}?format=summary
```

### Check Crawl Status
```bash
curl http://localhost:8082/api/v1/status/{crawl_id}
```

## Configuration Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `domains` | List of allowed domains | Required |
| `keywords` | Keywords to search for | Required |
| `max_pages` | Maximum pages to crawl | 10 |
| `depth` | Maximum crawling depth | 2 |
| `parallel` | Number of parallel workers | 2 |
| `delay` | Delay between requests (seconds) | 1 |

## Response Format

### Crawl Result
```json
{
  "url": "https://example.com/page",
  "title": "Page Title",
  "content": "Page content...",
  "domain": "example.com",
  "keywords": ["keyword1", "keyword2"],
  "timestamp": "2024-01-01T12:00:00Z",
  "status_code": 200,
  "metadata": {
    "user_agent": "Mozilla/5.0...",
    "method": "GET"
  }
}
```

### Job Status
```json
{
  "crawl_id": "uuid-string",
  "status": "running|completed",
  "progress": 75,
  "total_results": 15,
  "start_time": "2024-01-01T12:00:00Z",
  "end_time": "2024-01-01T12:05:00Z"
}
```

## Advanced Features

### User Agent Rotation
The crawler automatically rotates between different user agents to avoid detection:
- Chrome on Windows
- Chrome on macOS
- Chrome on Linux

### Rate Limiting
Built-in rate limiting prevents overwhelming target servers:
- Configurable delay between requests
- Per-domain parallelism control
- Respectful crawling practices

### Error Handling
Comprehensive error handling for:
- Network timeouts
- HTTP errors
- Invalid responses
- Domain restrictions

### Debugging
Built-in debug logging shows:
- URLs being visited
- Response status codes
- Error messages
- Request details

## Performance Considerations

- **Memory Usage**: Results are stored in memory; consider implementing persistent storage for large crawls
- **Concurrency**: Higher parallelism increases speed but may overwhelm servers
- **Rate Limiting**: Balance between speed and server respect
- **Content Filtering**: Only pages with matching keywords are stored to save memory

## Comparison with Basic Crawler

| Feature | Basic Crawler | Advanced Crawler |
|---------|---------------|------------------|
| Framework | Custom HTTP | Colly |
| Parallelism | Limited | Full support |
| Rate Limiting | Basic | Advanced |
| User Agents | Single | Rotation |
| Error Handling | Basic | Comprehensive |
| Debugging | Limited | Full logging |
| Performance | Good | Excellent |
| Scalability | Limited | High |