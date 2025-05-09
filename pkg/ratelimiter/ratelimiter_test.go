package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	windowSize := time.Second
	maxRequests := 3

	rl := NewRateLimiter(windowSize, maxRequests)
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	if rl.windowSize != windowSize {
		t.Errorf("Expected window size %v, got %v", windowSize, rl.windowSize)
	}

	if rl.maxRequests != maxRequests {
		t.Errorf("Expected max requests %d, got %d", maxRequests, rl.maxRequests)
	}

	// Test cleanup goroutine is running
	rl.Close()
}

func TestRateLimit_SingleClient(t *testing.T) {
	rl := NewRateLimiter(time.Second, 3)
	defer rl.Close()

	client := "test_client"

	// First 3 requests should be allowed
	for i := 0; i < 3; i++ {
		if !rl.RateLimit(client) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be rejected
	if rl.RateLimit(client) {
		t.Error("Request should be rejected after exceeding rate limit")
	}

	// Wait for window to reset
	time.Sleep(time.Second)

	// Should be allowed again after window reset
	if !rl.RateLimit(client) {
		t.Error("Request should be allowed after window reset")
	}
}

func TestRateLimit_MultipleClients(t *testing.T) {
	rl := NewRateLimiter(time.Second, 3)
	defer rl.Close()

	clients := []string{"client1", "client2", "client3"}

	// Each client should be able to make 3 requests
	for _, client := range clients {
		for i := 0; i < 3; i++ {
			if !rl.RateLimit(client) {
				t.Errorf("Client %s: Request %d should be allowed", client, i+1)
			}
		}
	}

	// Each client should be rejected after their limit
	for _, client := range clients {
		if rl.RateLimit(client) {
			t.Errorf("Client %s: Request should be rejected after exceeding rate limit", client)
		}
	}
}

func TestRateLimit_ConcurrentRequests(t *testing.T) {
	rl := NewRateLimiter(time.Second, 3)
	defer rl.Close()

	client := "concurrent_client"
	var wg sync.WaitGroup
	results := make(chan bool, 5)

	// Launch 5 concurrent requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- rl.RateLimit(client)
		}()
	}

	wg.Wait()
	close(results)

	// Count allowed requests
	allowed := 0
	for result := range results {
		if result {
			allowed++
		}
	}

	if allowed != 3 {
		t.Errorf("Expected 3 allowed requests, got %d", allowed)
	}
}

func TestRateLimit_BurstRequests(t *testing.T) {
	rl := NewRateLimiter(time.Second, 3)
	defer rl.Close()

	client := "burst_client"
	allowed := 0

	// Make 5 requests in rapid succession
	for i := 0; i < 5; i++ {
		if rl.RateLimit(client) {
			allowed++
		}
	}

	if allowed != 3 {
		t.Errorf("Expected 3 allowed requests in burst, got %d", allowed)
	}
}

func TestRateLimit_WindowSliding(t *testing.T) {
	rl := NewRateLimiter(time.Second, 3)
	defer rl.Close()

	client := "sliding_client"

	// Make 3 requests
	for i := 0; i < 3; i++ {
		if !rl.RateLimit(client) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Wait for half the window
	time.Sleep(500 * time.Millisecond)

	// Should still be rejected
	if rl.RateLimit(client) {
		t.Error("Request should be rejected within the same window")
	}

	// Wait for the full window
	time.Sleep(500 * time.Millisecond)

	// Should be allowed again
	if !rl.RateLimit(client) {
		t.Error("Request should be allowed after window reset")
	}
}

func TestRateLimit_Cleanup(t *testing.T) {
	rl := NewRateLimiter(time.Second, 3)
	defer rl.Close()

	client := "cleanup_client"

	// Make some requests
	for i := 0; i < 3; i++ {
		rl.RateLimit(client)
	}

	// Wait for cleanup (wait for 2 seconds to ensure the cleanup ticker has run)
	time.Sleep(2 * time.Second)

	// Check if the client's data was cleaned up
	rl.mu.RLock()
	_, exists := rl.requests[client]
	rl.mu.RUnlock()

	if exists {
		t.Error("Client data should be cleaned up after window expiration")
	}
}
