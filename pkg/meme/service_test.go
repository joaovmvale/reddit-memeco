package meme

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	service := NewService(r)

	if service == nil {
		t.Fatal("NewService returned nil")
	}

	if len(service.memes) != 5 {
		t.Errorf("Expected 5 memes, got %d", len(service.memes))
	}

	// Verify meme IDs are unique
	ids := make(map[string]bool)
	for _, meme := range service.memes {
		if ids[meme.ID] {
			t.Errorf("Duplicate meme ID found: %s", meme.ID)
		}
		ids[meme.ID] = true
	}
}

func TestGetRandomMeme(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	service := NewService(r)

	// Test multiple random selections
	seenMemes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		meme, err := service.GetRandomMeme()
		if err != nil {
			t.Errorf("GetRandomMeme failed: %v", err)
		}
		if meme == nil {
			t.Error("GetRandomMeme returned nil meme")
		}
		seenMemes[meme.ID] = true
	}

	// Verify we've seen all memes at least once
	if len(seenMemes) != 5 {
		t.Errorf("Expected to see all 5 memes, got %d unique memes", len(seenMemes))
	}
}

func TestGetMemeByID(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	service := NewService(r)

	// Test existing memes
	for i := 1; i <= 5; i++ {
		id := string(rune('0' + i))
		meme, err := service.GetMemeByID(id)
		if err != nil {
			t.Errorf("GetMemeByID failed for ID %s: %v", id, err)
		}
		if meme == nil {
			t.Errorf("GetMemeByID returned nil for ID %s", id)
		}
		if meme.ID != id {
			t.Errorf("Expected meme ID %s, got %s", id, meme.ID)
		}
	}

	// Test non-existent meme
	_, err := service.GetMemeByID("999")
	if err == nil {
		t.Error("Expected error for non-existent meme ID")
	}
}

func TestGetRandomMemeEmpty(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	service := &Service{
		memes: []Meme{},
		r:     r,
	}

	_, err := service.GetRandomMeme()
	if err == nil {
		t.Error("Expected error for empty meme collection")
	}
}
