#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Testing MemeCo API..."
echo "====================="

# Test 1: Get a random meme
echo -e "\n${GREEN}Test 1: Get a random meme${NC}"
curl -s http://localhost:8080/memes/random

# Test 2: Get a specific meme
echo -e "\n${GREEN}Test 2: Get a specific meme (ID: 1)${NC}"
curl -s http://localhost:8080/memes/1

# Test 3: Get a non-existent meme
echo -e "\n${GREEN}Test 3: Get a non-existent meme (ID: 999)${NC}"
curl -s http://localhost:8080/memes/999

# Test 4: Rate limiting test
echo -e "\n${GREEN}Test 4: Rate limiting test (making 10 requests in quick succession)${NC}"
echo "Making 10 requests to /memes/random..."
for i in {1..10}; do
    response=$(curl -s -w "\n%{http_code}" http://localhost:8080/memes/random)
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 429 ]; then
        echo -e "${RED}Request $i: Rate limited (429)${NC}"
    else
        echo -e "${GREEN}Request $i: Success ($status_code)${NC}"
        echo "$body"
    fi
done

# Test 5: Multiple clients test
echo -e "\n${GREEN}Test 5: Multiple clients test${NC}"
echo "Making requests from different IPs..."

# Simulate different clients using X-Forwarded-For header
for i in {1..5}; do
    ip="192.168.1.$i"
    echo -e "\nTesting with IP: $ip"
    for j in {1..3}; do
        response=$(curl -s -w "\n%{http_code}" -H "X-Forwarded-For: $ip" http://localhost:8080/memes/random)
        status_code=$(echo "$response" | tail -n1)
        body=$(echo "$response" | sed '$d')
        
        if [ "$status_code" -eq 429 ]; then
            echo -e "${RED}Request $j: Rate limited (429)${NC}"
        else
            echo -e "${GREEN}Request $j: Success ($status_code)${NC}"
            echo "$body"
        fi
    done
done 