#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Testing Rate Limiter..."
echo "====================="

# Function to make a request and print the result
make_request() {
    local request_num=$1
    local ip=$2
    
    if [ -z "$ip" ]; then
        response=$(curl -s -w "\n%{http_code}" http://localhost:8080/memes/random)
    else
        response=$(curl -s -w "\n%{http_code}" -H "X-Forwarded-For: $ip" http://localhost:8080/memes/random)
    fi
    
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq 429 ]; then
        echo -e "${RED}Request $request_num: Rate limited (429)${NC}"
    else
        echo -e "${GREEN}Request $request_num: Success ($status_code)${NC}"
    fi
}

# Test 1: Single client burst
echo -e "\n${GREEN}Test 1: Single client burst (20 requests)${NC}"
for i in {1..20}; do
    make_request $i &
done
wait

# Wait for rate limit window to reset
echo -e "\nWaiting for rate limit window to reset..."
sleep 2

# Test 2: Multiple clients burst
echo -e "\n${GREEN}Test 2: Multiple clients burst${NC}"
for i in {1..5}; do
    ip="192.168.1.$i"
    echo -e "\nTesting with IP: $ip"
    for j in {1..10}; do
        make_request $j $ip &
    done
    wait
    sleep 0.1
done 