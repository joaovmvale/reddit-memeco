package ratelimiter

import (
	"log"
	"sync"
	"time"
)

// RateLimiter implements a sliding window rate limiter
type RateLimiter struct {
	mu            sync.RWMutex
	requests      map[string][]time.Time // client -> slice of request timestamps
	windowSize    time.Duration          // size of the sliding window
	maxRequests   int                    // maximum number of requests allowed in the window
	cleanupTicker *time.Ticker
	done          chan struct{}
}

// NewRateLimiter creates a new rate limiter with the specified window size and maximum requests
func NewRateLimiter(windowSize time.Duration, maxRequests int) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(map[string][]time.Time),
		windowSize:  windowSize,
		maxRequests: maxRequests,
		done:        make(chan struct{}),
	}

	// Start cleanup goroutine to remove old entries
	rl.cleanupTicker = time.NewTicker(windowSize)
	go rl.cleanup()

	return rl
}

// RateLimit checks if a request from the given client should be allowed
func (rl *RateLimiter) RateLimit(clientName string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.windowSize)

	// Get or create the client's request history
	requests, exists := rl.requests[clientName]
	if !exists {
		requests = make([]time.Time, 0)
	}

	// Remove requests outside the window
	validRequests := make([]time.Time, 0)
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	// Debug logging
	log.Printf("Client: %s, Valid requests in window: %d/%d", clientName, len(validRequests), rl.maxRequests)

	// Check if we're under the limit
	if len(validRequests) < rl.maxRequests {
		validRequests = append(validRequests, now)
		rl.requests[clientName] = validRequests
		return true
	}

	return false
}

// cleanup periodically removes old entries from the requests map
func (rl *RateLimiter) cleanup() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.mu.Lock()
			now := time.Now()
			windowStart := now.Add(-rl.windowSize)

			for client, requests := range rl.requests {
				validRequests := make([]time.Time, 0)
				for _, t := range requests {
					if t.After(windowStart) {
						validRequests = append(validRequests, t)
					}
				}
				if len(validRequests) == 0 {
					delete(rl.requests, client)
					log.Printf("Cleaned up client: %s", client)
				} else {
					rl.requests[client] = validRequests
				}
			}
			rl.mu.Unlock()
		case <-rl.done:
			rl.cleanupTicker.Stop()
			return
		}
	}
}

// Close stops the cleanup goroutine
func (rl *RateLimiter) Close() {
	close(rl.done)
}
