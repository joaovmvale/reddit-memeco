package handlers

import (
	"encoding/json"
	"net/http"

	"reddit-memeco/pkg/meme"
	"reddit-memeco/pkg/ratelimiter"

	"github.com/gorilla/mux"
)

// MemeHandler handles HTTP requests for memes
type MemeHandler struct {
	memeService *meme.Service
	rateLimiter *ratelimiter.RateLimiter
}

// NewMemeHandler creates a new meme handler
func NewMemeHandler(memeService *meme.Service, rateLimiter *ratelimiter.RateLimiter) *MemeHandler {
	return &MemeHandler{
		memeService: memeService,
		rateLimiter: rateLimiter,
	}
}

// getClientIP extracts the client IP from the request
func (h *MemeHandler) getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		return forwardedFor
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// responseWithError sends an error response
func responseWithError(w http.ResponseWriter, code int, message string) {
	responseWithJSON(w, code, map[string]string{"error": message})
}

// responseWithJSON sends a JSON response
func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RateLimitMiddleware checks if the request should be rate limited
func (h *MemeHandler) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := h.getClientIP(r)

		if !h.rateLimiter.RateLimit(clientIP) {
			responseWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		next(w, r)
	}
}

// GetRandomMeme handles requests for random memes
func (h *MemeHandler) GetRandomMeme(w http.ResponseWriter, r *http.Request) {
	meme, err := h.memeService.GetRandomMeme()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Failed to get random meme")
		return
	}

	responseWithJSON(w, http.StatusOK, meme)
}

// GetMemeByID handles requests for specific memes
func (h *MemeHandler) GetMemeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	meme, err := h.memeService.GetMemeByID(id)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "Meme not found")
		return
	}

	responseWithJSON(w, http.StatusOK, meme)
}
