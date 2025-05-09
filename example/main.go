package main

import (
	"fmt"
	"sync"
	"time"

	"reddit-memeco/pkg/ratelimiter"
)

func main() {
	// Create a rate limiter that allows 3 requests per second
	limiter := ratelimiter.NewRateLimiter(time.Second, 3)
	defer limiter.Close()

	// Example 1: Single client making requests
	fmt.Println("Example 1: Single client making requests")
	client := "client1"
	for i := 0; i < 5; i++ {
		allowed := limiter.RateLimit(client)
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(200 * time.Millisecond)
	}

	// Wait a bit to reset the window
	time.Sleep(time.Second)
	fmt.Println("\nExample 2: Multiple clients making concurrent requests")

	// Example 2: Multiple clients making concurrent requests
	var wg sync.WaitGroup
	clients := []string{"client1", "client2", "client3"}

	for _, client := range clients {
		wg.Add(1)
		go func(client string) {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				allowed := limiter.RateLimit(client)
				fmt.Printf("Client %s, Request %d: %v\n", client, i+1, allowed)
				time.Sleep(100 * time.Millisecond)
			}
		}(client)
	}

	wg.Wait()

	// Example 3: Burst of requests
	fmt.Println("\nExample 3: Burst of requests")
	burstClient := "burst_client"
	for i := 0; i < 5; i++ {
		allowed := limiter.RateLimit(burstClient)
		fmt.Printf("Burst Request %d: %v\n", i+1, allowed)
		// No sleep between requests to simulate a burst
	}
}
