package meme

import (
	"fmt"
	"math/rand"
)

// Meme represents a meme with its metadata
type Meme struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Service handles meme operations
type Service struct {
	memes []Meme
	r     *rand.Rand
}

// NewService creates a new meme service with some sample memes
func NewService(r *rand.Rand) *Service {
	return &Service{
		memes: []Meme{
			{ID: "1", Title: "First Meme", URL: "https://example.com/meme1.gif"},
			{ID: "2", Title: "Second Meme", URL: "https://example.com/meme2.gif"},
			{ID: "3", Title: "Third Meme", URL: "https://example.com/meme3.gif"},
			{ID: "4", Title: "Fourth Meme", URL: "https://example.com/meme4.gif"},
			{ID: "5", Title: "Fifth Meme", URL: "https://example.com/meme5.gif"},
		},
		r: r,
	}
}

// GetRandomMeme returns a random meme from the collection
func (s *Service) GetRandomMeme() (*Meme, error) {
	if len(s.memes) == 0 {
		return nil, fmt.Errorf("no memes available")
	}

	randomIndex := s.r.Intn(len(s.memes))
	return &s.memes[randomIndex], nil
}

// GetMemeByID returns a meme by its ID
func (s *Service) GetMemeByID(id string) (*Meme, error) {
	for _, meme := range s.memes {
		if meme.ID == id {
			return &meme, nil
		}
	}
	return nil, fmt.Errorf("meme not found with ID: %s", id)
}
