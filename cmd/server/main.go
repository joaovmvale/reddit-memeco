package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"reddit-memeco/api/handlers"
	"reddit-memeco/pkg/meme"
	"reddit-memeco/pkg/ratelimiter"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create services
	memeService := meme.NewService(r)
	rateLimiter := ratelimiter.NewRateLimiter(time.Second, 3) // 3 requests per second
	defer rateLimiter.Close()

	// Create handler
	memeHandler := handlers.NewMemeHandler(memeService, rateLimiter)

	// Create router
	router := mux.NewRouter()

	// Register routes with rate limiting middleware
	router.HandleFunc("/memes/random", memeHandler.RateLimitMiddleware(memeHandler.GetRandomMeme)).Methods("GET")
	router.HandleFunc("/memes/{id}", memeHandler.RateLimitMiddleware(memeHandler.GetMemeByID)).Methods("GET")

	// Add basic middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
