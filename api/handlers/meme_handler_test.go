package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"reddit-memeco/pkg/meme"
	"reddit-memeco/pkg/ratelimiter"

	"github.com/gorilla/mux"
)

func setupTestHandler() (*MemeHandler, *httptest.ResponseRecorder) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	memeService := meme.NewService(r)
	rateLimiter := ratelimiter.NewRateLimiter(time.Second, 3)
	handler := NewMemeHandler(memeService, rateLimiter)
	recorder := httptest.NewRecorder()
	return handler, recorder
}

func TestGetRandomMeme(t *testing.T) {
	handler, recorder := setupTestHandler()

	req := httptest.NewRequest("GET", "/memes/random", nil)
	handler.GetRandomMeme(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response meme.Meme
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID == "" || response.Title == "" || response.URL == "" {
		t.Error("Response missing required fields")
	}
}

func TestGetMemeByID(t *testing.T) {
	handler, recorder := setupTestHandler()

	// Test existing meme
	req := httptest.NewRequest("GET", "/memes/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	handler.GetMemeByID(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	var response meme.Meme
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != "1" {
		t.Errorf("Expected meme ID 1, got %s", response.ID)
	}

	// Test non-existent meme
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/memes/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	handler.GetMemeByID(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	handler, recorder := setupTestHandler()
	clientIP := "192.168.1.1"

	// Create a test handler function
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test within rate limit
	for i := 0; i < 3; i++ {
		recorder = httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Forwarded-For", clientIP)

		middleware := handler.RateLimitMiddleware(testHandler)
		middleware.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status code %d, got %d", i+1, http.StatusOK, recorder.Code)
		}
	}

	// Test exceeding rate limit
	recorder = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", clientIP)

	middleware := handler.RateLimitMiddleware(testHandler)
	middleware.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, recorder.Code)
	}
}

func TestGetClientIP(t *testing.T) {
	handler, _ := setupTestHandler()

	// Test X-Forwarded-For header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	ip := handler.getClientIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("Expected IP 192.168.1.1, got %s", ip)
	}

	// Test RemoteAddr fallback
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	ip = handler.getClientIP(req)
	if ip != "192.168.1.2:12345" {
		t.Errorf("Expected IP 192.168.1.2:12345, got %s", ip)
	}
}
