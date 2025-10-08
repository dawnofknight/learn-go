#!/bin/bash

# Test script for Crawler REST API
API_BASE="http://localhost:8081/api/v1"
HEALTH_URL="http://localhost:8081/health"

echo "=== Crawler REST API Test Script ==="
echo

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s "$HEALTH_URL" | jq '.'
echo
echo

# Test 2: Submit a crawl job
echo "2. Submitting a crawl job..."
CRAWL_RESPONSE=$(curl -s -X POST "$API_BASE/crawl" \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["artificial intelligence", "machine learning", "technology"],
    "domains": ["iana.org", "example.com"],
    "start_date": "2024-01-01",
    "end_date": "2024-12-31",
    "max_depth": 2,
    "max_pages": 50
  }')

echo "$CRAWL_RESPONSE" | jq '.'
CRAWL_ID=$(echo "$CRAWL_RESPONSE" | jq -r '.crawl_id')
echo
echo "Crawl ID: $CRAWL_ID"
echo

# Test 3: Get crawl status
echo "3. Getting crawl status..."
curl -s "$API_BASE/crawl/$CRAWL_ID" | jq '.'
echo
echo

# Test 4: List all crawls
echo "4. Listing all crawls..."
curl -s "$API_BASE/crawl" | jq '.'
echo
echo

# Test 5: Get crawl results (may be empty initially)
echo "5. Getting crawl results..."
curl -s "$API_BASE/crawl/$CRAWL_ID/results?page=1&limit=10" | jq '.'
echo
echo

# Test 6: Submit another crawl job with different parameters
echo "6. Submitting another crawl job..."
CRAWL_RESPONSE_2=$(curl -s -X POST "$API_BASE/crawl" \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["web crawling", "data extraction"],
    "domains": ["github.com"],
    "max_depth": 1,
    "max_pages": 20
  }')

echo "$CRAWL_RESPONSE_2" | jq '.'
CRAWL_ID_2=$(echo "$CRAWL_RESPONSE_2" | jq -r '.crawl_id')
echo
echo "Second Crawl ID: $CRAWL_ID_2"
echo

# Test 7: List all crawls again
echo "7. Listing all crawls (should show 2 now)..."
curl -s "$API_BASE/crawl" | jq '.'
echo
echo

# Test 8: Test error handling - invalid request
echo "8. Testing error handling (missing required fields)..."
curl -s -X POST "$API_BASE/crawl" \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": [],
    "domains": []
  }' | jq '.'
echo
echo

# Test 9: Test error handling - invalid date format
echo "9. Testing error handling (invalid date format)..."
curl -s -X POST "$API_BASE/crawl" \
  -H "Content-Type: application/json" \
  -d '{
    "keywords": ["test"],
    "domains": ["example.com"],
    "start_date": "invalid-date",
    "end_date": "2024-12-31"
  }' | jq '.'
echo
echo

# Test 10: Test error handling - non-existent crawl ID
echo "10. Testing error handling (non-existent crawl ID)..."
curl -s "$API_BASE/crawl/non-existent-id" | jq '.'
echo
echo

echo "=== Test completed ==="
echo "Available crawl IDs for further testing:"
echo "- $CRAWL_ID"
echo "- $CRAWL_ID_2"
echo
echo "You can manually test these endpoints:"
echo "- Health: $HEALTH_URL"
echo "- List crawls: $API_BASE/crawl"
echo "- Get status: $API_BASE/crawl/{crawl_id}"
echo "- Get results: $API_BASE/crawl/{crawl_id}/results"
echo "- Cancel crawl: DELETE $API_BASE/crawl/{crawl_id}"