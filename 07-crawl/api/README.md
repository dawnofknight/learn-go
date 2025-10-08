# Crawler REST API

A Go REST API application that communicates with StormCrawler to crawl data based on keywords, specific domains, and optional date range filters.

## Features

- **Keyword-based crawling**: Submit multiple keywords for targeted crawling
- **Domain filtering**: Specify multiple domains to crawl
- **Date range filtering**: Optional start and end date filters
- **Crawl job management**: Submit, monitor, and cancel crawl jobs
- **Real-time status tracking**: Monitor crawl progress and results
- **URLFrontier integration**: Communicates with URLFrontier service for distributed crawling

## API Endpoints

### Health Check
```
GET /health
```

### Submit Crawl Job
```
POST /api/v1/crawl
```

**Request Body:**
```json
{
  "keywords": ["artificial intelligence", "machine learning"],
  "domains": ["example.com", "test.org"],
  "start_date": "2024-01-01",
  "end_date": "2024-12-31",
  "max_depth": 3,
  "max_pages": 100
}
```

**Response:**
```json
{
  "crawl_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "submitted",
  "message": "Crawl job submitted successfully with 6 seed URLs",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Get Crawl Status
```
GET /api/v1/crawl/{crawl_id}
```

**Response:**
```json
{
  "crawl_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "running",
  "progress": 45,
  "total_urls": 100,
  "processed_urls": 45,
  "start_time": "2024-01-15T10:30:00Z",
  "end_time": null,
  "results": []
}
```

### Get Crawl Results
```
GET /api/v1/crawl/{crawl_id}/results?page=1&limit=50
```

**Response:**
```json
{
  "crawl_id": "550e8400-e29b-41d4-a716-446655440000",
  "results": [
    {
      "url": "https://example.com/page1",
      "title": "Example Page",
      "content": "Page content...",
      "domain": "example.com",
      "keywords": ["artificial intelligence"],
      "timestamp": "2024-01-15T10:35:00Z",
      "status_code": 200,
      "metadata": {
        "content_type": "text/html",
        "content_length": "1024"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 45,
    "pages": 1
  }
}
```

### List All Crawls
```
GET /api/v1/crawl
```

### Cancel Crawl Job
```
DELETE /api/v1/crawl/{crawl_id}
```

## Request Parameters

### Required Parameters
- `keywords`: Array of keywords to search for (minimum 1)
- `domains`: Array of domains to crawl (minimum 1)

### Optional Parameters
- `start_date`: Start date for filtering (format: YYYY-MM-DD)
- `end_date`: End date for filtering (format: YYYY-MM-DD)
- `max_depth`: Maximum crawl depth (default: 3)
- `max_pages`: Maximum pages to crawl (default: 100)

## Running the API

### Prerequisites
- Go 1.24 or later
- StormCrawler cluster running
- URLFrontier service accessible

### Installation
```bash
cd api
go mod tidy
```

### Configuration
The API connects to URLFrontier at `host.docker.internal:7071` by default. This can be configured in the `main()` function.

### Start the Server
```bash
go run main.go
```

The API will start on port 8080:
- Health check: http://localhost:8080/health
- API base: http://localhost:8080/api/v1

## Integration with StormCrawler

The API integrates with your existing StormCrawler setup by:

1. **URLFrontier Communication**: Submits crawl requests to the URLFrontier service
2. **Metadata Enrichment**: Adds keywords, domains, and date filters as metadata
3. **Queue Management**: Uses crawl IDs as queue names for job isolation
4. **Status Monitoring**: Queries URLFrontier for real-time crawl statistics

## Example Usage

### Submit a crawl job
```bash
curl -X POST http://localhost:8080/api/v1/crawl \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["technology", "innovation"],
    "domains": ["techcrunch.com", "wired.com"],
    "start_date": "2024-01-01",
    "max_depth": 2,
    "max_pages": 50
  }'
```

### Check crawl status
```bash
curl http://localhost:8080/api/v1/crawl/550e8400-e29b-41d4-a716-446655440000
```

### Get crawl results
```bash
curl http://localhost:8080/api/v1/crawl/550e8400-e29b-41d4-a716-446655440000/results?page=1&limit=10
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid request format or parameters
- `404 Not Found`: Crawl job not found
- `500 Internal Server Error`: Server or URLFrontier communication errors

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   REST API      │    │   URLFrontier    │    │  StormCrawler   │
│   (Go/Gin)      │◄──►│   Service        │◄──►│   Topology      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Crawl Jobs    │    │   URL Queues     │    │   Fetcher/      │
│   Management    │    │   & Metadata     │    │   Parser/       │
│                 │    │                  │    │   Indexer       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

The API serves as a user-friendly interface to the StormCrawler system, allowing easy submission and monitoring of keyword-based crawl jobs with domain and date filtering capabilities.