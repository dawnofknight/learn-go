#!/bin/bash

# ScyllaDB CRUD REST API Test Script
# This script tests all the API endpoints

API_BASE="http://localhost:8080/api/v1"

echo "🧪 Testing ScyllaDB CRUD REST API"
echo "================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print test results
print_test() {
    echo -e "\n${YELLOW}$1${NC}"
    echo "-----------------------------------"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Test 1: Health Check
print_test "1. Health Check"
response=$(curl -s "$API_BASE/health")
if [[ $? -eq 0 ]]; then
    print_success "Health check passed"
    echo "Response: $response"
else
    print_error "Health check failed"
    echo "Make sure the server is running on port 8080"
    exit 1
fi

# Test 2: Create User
print_test "2. Create User"
user_response=$(curl -s -X POST "$API_BASE/users" \
    -H "Content-Type: application/json" \
    -d '{"name": "Test User", "email": "test@example.com"}')

if [[ $? -eq 0 ]]; then
    print_success "User created successfully"
    echo "Response: $user_response"
    
    # Extract user ID from response
    user_id=$(echo "$user_response" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "User ID: $user_id"
else
    print_error "Failed to create user"
    exit 1
fi

# Test 3: Get All Users
print_test "3. Get All Users"
response=$(curl -s "$API_BASE/users")
if [[ $? -eq 0 ]]; then
    print_success "Retrieved all users"
    echo "Response: $response"
else
    print_error "Failed to get users"
fi

# Test 4: Get User by ID
if [[ -n "$user_id" ]]; then
    print_test "4. Get User by ID"
    response=$(curl -s "$API_BASE/users/$user_id")
    if [[ $? -eq 0 ]]; then
        print_success "Retrieved user by ID"
        echo "Response: $response"
    else
        print_error "Failed to get user by ID"
    fi
fi

# Test 5: Update User
if [[ -n "$user_id" ]]; then
    print_test "5. Update User"
    response=$(curl -s -X PUT "$API_BASE/users/$user_id" \
        -H "Content-Type: application/json" \
        -d '{"name": "Updated Test User", "email": "updated@example.com"}')
    
    if [[ $? -eq 0 ]]; then
        print_success "User updated successfully"
        echo "Response: $response"
    else
        print_error "Failed to update user"
    fi
fi

# Test 6: Get Updated User
if [[ -n "$user_id" ]]; then
    print_test "6. Verify Update"
    response=$(curl -s "$API_BASE/users/$user_id")
    if [[ $? -eq 0 ]]; then
        print_success "Verified user update"
        echo "Response: $response"
    else
        print_error "Failed to verify update"
    fi
fi

# Test 7: Delete User
if [[ -n "$user_id" ]]; then
    print_test "7. Delete User"
    response=$(curl -s -X DELETE "$API_BASE/users/$user_id")
    if [[ $? -eq 0 ]]; then
        print_success "User deleted successfully"
        echo "Response: $response"
    else
        print_error "Failed to delete user"
    fi
fi

# Test 8: Verify Deletion
if [[ -n "$user_id" ]]; then
    print_test "8. Verify Deletion"
    response=$(curl -s "$API_BASE/users/$user_id")
    if [[ $response == *"user not found"* ]]; then
        print_success "Confirmed user deletion"
        echo "Response: $response"
    else
        print_error "User still exists after deletion"
        echo "Response: $response"
    fi
fi

echo -e "\n${GREEN}🎉 API Testing Complete!${NC}"
echo "================================="