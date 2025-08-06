package models

import (
	"testing"
	"time"
)

func TestNewSong(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		title      string
		artist     string
		album      string
		genre      string
		subgenre   string
		mood       string
		duration   int
		bpm        int
		wantTitle  string
		wantArtist string
	}{
		{
			name:       "Valid song creation",
			id:         "test-song-1",
			title:      "Test Song",
			artist:     "Test Artist",
			album:      "Test Album",
			genre:      "Rock",
			subgenre:   "Alternative",
			mood:       "Energetic",
			duration:   180,
			bpm:        120,
			wantTitle:  "Test Song",
			wantArtist: "Test Artist",
		},
		{
			name:       "Empty fields handled",
			id:         "",
			title:      "",
			artist:     "",
			album:      "",
			genre:      "",
			subgenre:   "",
			mood:       "",
			duration:   0,
			bpm:        0,
			wantTitle:  "",
			wantArtist: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			song := NewSong(tt.id, tt.title, tt.artist, tt.album, tt.genre, tt.subgenre, tt.mood, tt.duration, tt.bpm)

			if song.ID != tt.id {
				t.Errorf("NewSong() ID = %v, want %v", song.ID, tt.id)
			}
			if song.Title != tt.wantTitle {
				t.Errorf("NewSong() Title = %v, want %v", song.Title, tt.wantTitle)
			}
			if song.Artist != tt.wantArtist {
				t.Errorf("NewSong() Artist = %v, want %v", song.Artist, tt.wantArtist)
			}
			if song.Duration != tt.duration {
				t.Errorf("NewSong() Duration = %v, want %v", song.Duration, tt.duration)
			}
			if song.Rating != 0 {
				t.Errorf("NewSong() Rating = %v, want %v", song.Rating, 0)
			}
			if song.PlayCount != 0 {
				t.Errorf("NewSong() PlayCount = %v, want %v", song.PlayCount, 0)
			}
			if time.Since(song.AddedAt) > time.Second {
				t.Errorf("NewSong() AddedAt should be recent")
			}
		})
	}
}

func TestSong_Play(t *testing.T) {
	song := NewSong("test-1", "Test Song", "Test Artist", "Test Album", "Rock", "Alt", "Happy", 180, 120)

	initialPlayCount := song.PlayCount
	if song.LastPlayed != nil {
		t.Errorf("Song.LastPlayed should be nil initially")
	}

	// Play the song
	song.Play()

	if song.PlayCount != initialPlayCount+1 {
		t.Errorf("Song.Play() PlayCount = %v, want %v", song.PlayCount, initialPlayCount+1)
	}
	if song.LastPlayed == nil {
		t.Errorf("Song.LastPlayed should not be nil after playing")
	}
	if time.Since(*song.LastPlayed) > time.Second {
		t.Errorf("Song.LastPlayed should be recent")
	}

	// Play again to test increment
	song.Play()
	if song.PlayCount != initialPlayCount+2 {
		t.Errorf("Song.Play() second call PlayCount = %v, want %v", song.PlayCount, initialPlayCount+2)
	}
}

func TestSong_SetRating(t *testing.T) {
	song := NewSong("test-1", "Test Song", "Test Artist", "Test Album", "Rock", "Alt", "Happy", 180, 120)

	tests := []struct {
		name           string
		rating         int
		expectedRating int
	}{
		{"Valid rating 1", 1, 1},
		{"Valid rating 5", 5, 5},
		{"Valid rating 3", 3, 3},
		{"Invalid rating 0", 0, 3},   // Should not change from previous
		{"Invalid rating 6", 6, 3},   // Should not change from previous
		{"Invalid rating -1", -1, 3}, // Should not change from previous
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			song.SetRating(tt.rating)
			if song.Rating != tt.expectedRating {
				t.Errorf("Song.SetRating(%v) Rating = %v, want %v", tt.rating, song.Rating, tt.expectedRating)
			}
		})
	}
}

func TestSong_IsSimilar(t *testing.T) {
	baseSong := NewSong("base", "Base Song", "Base Artist", "Base Album", "Rock", "Alternative", "Energetic", 180, 120)

	tests := []struct {
		name     string
		other    *Song
		expected bool
	}{
		{
			name:     "Same genre and mood, similar duration",
			other:    NewSong("other1", "Other Song", "Other Artist", "Other Album", "Rock", "Progressive", "Energetic", 190, 130),
			expected: true,
		},
		{
			name:     "Same genre and mood, exact duration",
			other:    NewSong("other2", "Other Song", "Other Artist", "Other Album", "Rock", "Grunge", "Energetic", 180, 100),
			expected: true,
		},
		{
			name:     "Same genre and mood, duration difference at boundary (30s)",
			other:    NewSong("other3", "Other Song", "Other Artist", "Other Album", "Rock", "Hard Rock", "Energetic", 210, 140),
			expected: true,
		},
		{
			name:     "Same genre and mood, duration difference too large",
			other:    NewSong("other4", "Other Song", "Other Artist", "Other Album", "Rock", "Punk", "Energetic", 220, 150),
			expected: false,
		},
		{
			name:     "Different genre, same mood",
			other:    NewSong("other5", "Other Song", "Other Artist", "Other Album", "Pop", "Dance Pop", "Energetic", 180, 128),
			expected: false,
		},
		{
			name:     "Same genre, different mood",
			other:    NewSong("other6", "Other Song", "Other Artist", "Other Album", "Rock", "Alternative", "Sad", 180, 120),
			expected: false,
		},
		{
			name:     "Different genre and mood",
			other:    NewSong("other7", "Other Song", "Other Artist", "Other Album", "Jazz", "Smooth Jazz", "Relaxing", 240, 80),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := baseSong.IsSimilar(tt.other)
			if result != tt.expected {
				t.Errorf("Song.IsSimilar() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSong_DurationString(t *testing.T) {
	tests := []struct {
		name     string
		duration int
		expected string
	}{
		{"0 seconds", 0, "00:00"},
		{"30 seconds", 30, "00:30"},
		{"1 minute", 60, "01:00"},
		{"3 minutes 30 seconds", 210, "03:30"},
		{"10 minutes 5 seconds", 605, "10:05"},
		{"59 minutes 59 seconds", 3599, "59:59"},
		{"1 hour (3600 seconds)", 3600, "60:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			song := NewSong("test", "Test", "Test", "Test", "Test", "Test", "Test", tt.duration, 120)
			result := song.DurationString()
			if result != tt.expected {
				t.Errorf("Song.DurationString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSong_GetMetadata(t *testing.T) {
	song := NewSong("test-id", "Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
	song.SetRating(4)
	song.Play()

	metadata := song.GetMetadata()

	expectedKeys := []string{"id", "title", "artist", "album", "duration", "genre", "subgenre", "mood", "bpm", "rating", "playcount", "added_at", "last_played"}

	for _, key := range expectedKeys {
		if _, exists := metadata[key]; !exists {
			t.Errorf("Song.GetMetadata() missing key: %s", key)
		}
	}

	if metadata["id"] != "test-id" {
		t.Errorf("Song.GetMetadata() id = %v, want %v", metadata["id"], "test-id")
	}
	if metadata["title"] != "Test Song" {
		t.Errorf("Song.GetMetadata() title = %v, want %v", metadata["title"], "Test Song")
	}
	if metadata["rating"] != 4 {
		t.Errorf("Song.GetMetadata() rating = %v, want %v", metadata["rating"], 4)
	}
	if metadata["playcount"] != 1 {
		t.Errorf("Song.GetMetadata() playcount = %v, want %v", metadata["playcount"], 1)
	}
}

// Benchmark tests for performance analysis
func BenchmarkNewSong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewSong("test-song", "Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
	}
}

func BenchmarkSong_Play(b *testing.B) {
	song := NewSong("test-song", "Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		song.Play()
	}
}

func BenchmarkSong_IsSimilar(b *testing.B) {
	song1 := NewSong("song1", "Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 180, 120)
	song2 := NewSong("song2", "Song 2", "Artist 2", "Album 2", "Rock", "Grunge", "Energetic", 190, 130)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		song1.IsSimilar(song2)
	}
}

func BenchmarkSong_DurationString(b *testing.B) {
	song := NewSong("test-song", "Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		song.DurationString()
	}
}

func BenchmarkSong_GetMetadata(b *testing.B) {
	song := NewSong("test-song", "Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
	song.SetRating(4)
	song.Play()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		song.GetMetadata()
	}
}

// Test helpers
func createTestSong(id, title, artist string) *Song {
	return NewSong(id, title, artist, "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
}

func TestSongFieldsInitialization(t *testing.T) {
	song := NewSong("test", "Title", "Artist", "Album", "Genre", "Subgenre", "Mood", 240, 130)

	// Test all fields are properly initialized
	if song.ID != "test" {
		t.Errorf("ID not initialized correctly")
	}
	if song.Title != "Title" {
		t.Errorf("Title not initialized correctly")
	}
	if song.Artist != "Artist" {
		t.Errorf("Artist not initialized correctly")
	}
	if song.Album != "Album" {
		t.Errorf("Album not initialized correctly")
	}
	if song.Genre != "Genre" {
		t.Errorf("Genre not initialized correctly")
	}
	if song.SubGenre != "Subgenre" {
		t.Errorf("SubGenre not initialized correctly")
	}
	if song.Mood != "Mood" {
		t.Errorf("Mood not initialized correctly")
	}
	if song.Duration != 240 {
		t.Errorf("Duration not initialized correctly")
	}
	if song.BPM != 130 {
		t.Errorf("BPM not initialized correctly")
	}
	if song.Rating != 0 {
		t.Errorf("Rating should default to 0")
	}
	if song.PlayCount != 0 {
		t.Errorf("PlayCount should default to 0")
	}
	if song.LastPlayed != nil {
		t.Errorf("LastPlayed should default to nil")
	}
	if song.AddedAt.IsZero() {
		t.Errorf("AddedAt should be set to current time")
	}
}

func TestSongRatingBoundaryValues(t *testing.T) {
	song := NewSong("test", "Test", "Test", "Test", "Test", "Test", "Test", 180, 120)

	// Test boundary values for rating
	boundaryTests := []struct {
		rating   int
		expected int
		name     string
	}{
		{1, 1, "minimum valid rating"},
		{5, 5, "maximum valid rating"},
		{0, 0, "below minimum (should not change)"},
		{6, 5, "above maximum (should not change from previous)"},
		{-10, 5, "negative rating (should not change)"},
		{100, 5, "very high rating (should not change)"},
	}

	for _, test := range boundaryTests {
		t.Run(test.name, func(t *testing.T) {
			song.SetRating(test.rating)
			if song.Rating != test.expected {
				t.Errorf("SetRating(%d) = %d, want %d", test.rating, song.Rating, test.expected)
			}
		})
	}
}
