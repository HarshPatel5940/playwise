package models

import (
	"fmt"
	"time"
)

// Song represents a music track with metadata
// Time Complexity: O(1) for all field access operations
// Space Complexity: O(1) per song instance
type Song struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Artist     string     `json:"artist"`
	Album      string     `json:"album"`
	Duration   int        `json:"duration"` // in seconds
	Genre      string     `json:"genre"`
	SubGenre   string     `json:"subgenre"`
	Mood       string     `json:"mood"`
	BPM        int        `json:"bpm"`
	Rating     int        `json:"rating"` // 1-5 stars
	PlayCount  int        `json:"playcount"`
	AddedAt    time.Time  `json:"added_at"`
	LastPlayed *time.Time `json:"last_played,omitempty"`
}

// NewSong creates a new song instance
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewSong(id, title, artist, album, genre, subgenre, mood string, duration, bpm int) *Song {
	return &Song{
		ID:        id,
		Title:     title,
		Artist:    artist,
		Album:     album,
		Duration:  duration,
		Genre:     genre,
		SubGenre:  subgenre,
		Mood:      mood,
		BPM:       bpm,
		Rating:    0,
		PlayCount: 0,
		AddedAt:   time.Now(),
	}
}

// Play increments play count and updates last played time
// Time Complexity: O(1)
// Space Complexity: O(1)
func (s *Song) Play() {
	s.PlayCount++
	now := time.Now()
	s.LastPlayed = &now
}

// SetRating sets the song rating (1-5)
// Time Complexity: O(1)
// Space Complexity: O(1)
func (s *Song) SetRating(rating int) {
	if rating >= 1 && rating <= 5 {
		s.Rating = rating
	}
}

// IsSimilar checks if two songs are similar based on genre, mood, and duration
// Time Complexity: O(1)
// Space Complexity: O(1)
func (s *Song) IsSimilar(other *Song) bool {
	if s.Genre == other.Genre && s.Mood == other.Mood {
		durationDiff := s.Duration - other.Duration
		if durationDiff < 0 {
			durationDiff = -durationDiff
		}
		// Consider similar if duration difference is less than 30 seconds
		return durationDiff <= 30
	}
	return false
}

// DurationString returns formatted duration as MM:SS
// Time Complexity: O(1)
// Space Complexity: O(1)
func (s *Song) DurationString() string {
	minutes := s.Duration / 60
	seconds := s.Duration % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// GetMetadata returns song metadata as a map for quick lookup
// Time Complexity: O(1)
// Space Complexity: O(1)
func (s *Song) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"id":          s.ID,
		"title":       s.Title,
		"artist":      s.Artist,
		"album":       s.Album,
		"duration":    s.Duration,
		"genre":       s.Genre,
		"subgenre":    s.SubGenre,
		"mood":        s.Mood,
		"bpm":         s.BPM,
		"rating":      s.Rating,
		"playcount":   s.PlayCount,
		"added_at":    s.AddedAt,
		"last_played": s.LastPlayed,
	}
}
