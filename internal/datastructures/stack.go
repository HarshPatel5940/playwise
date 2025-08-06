package datastructures

import (
	"fmt"
	"src/internal/models"
)

// PlaybackHistoryNode represents a node in the stack for playback history
// Each node contains a song and pointer to the next node below it
// Time Complexity: O(1) for all field operations
// Space Complexity: O(1) per node
type PlaybackHistoryNode struct {
	Song *models.Song
	Next *PlaybackHistoryNode
}

// PlaybackHistoryStack represents a LIFO stack for managing playback history
// Supports undo functionality by maintaining recently played songs
// Time Complexity: O(1) for push, pop, peek operations
// Space Complexity: O(n) where n is the number of songs in history
type PlaybackHistoryStack struct {
	Top     *PlaybackHistoryNode
	Size    int
	MaxSize int // Maximum number of songs to keep in history
}

// NewPlaybackHistoryStack creates a new playback history stack
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewPlaybackHistoryStack(maxSize int) *PlaybackHistoryStack {
	if maxSize <= 0 {
		maxSize = 50 // Default maximum history size
	}
	return &PlaybackHistoryStack{
		Top:     nil,
		Size:    0,
		MaxSize: maxSize,
	}
}

// Push adds a song to the top of the history stack
// If stack exceeds maxSize, removes the oldest entry from bottom
// Time Complexity: O(1) amortized, O(n) worst case when removing old entries
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) Push(song *models.Song) {
	newNode := &PlaybackHistoryNode{
		Song: song,
		Next: phs.Top,
	}

	phs.Top = newNode
	phs.Size++

	// If we exceed max size, remove the oldest entry (bottom of stack)
	if phs.Size > phs.MaxSize {
		phs.removeBottom()
	}
}

// Pop removes and returns the most recently played song from history
// This is the main undo functionality - removes last played song
// Time Complexity: O(1)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) Pop() (*models.Song, error) {
	if phs.IsEmpty() {
		return nil, fmt.Errorf("playback history is empty")
	}

	song := phs.Top.Song
	phs.Top = phs.Top.Next
	phs.Size--

	return song, nil
}

// Peek returns the most recently played song without removing it
// Time Complexity: O(1)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) Peek() (*models.Song, error) {
	if phs.IsEmpty() {
		return nil, fmt.Errorf("playback history is empty")
	}

	return phs.Top.Song, nil
}

// UndoLastPlay removes the last played song from history and returns it
// This allows re-queueing the song back to current playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) UndoLastPlay() (*models.Song, error) {
	return phs.Pop()
}

// IsEmpty checks if the playback history is empty
// Time Complexity: O(1)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) IsEmpty() bool {
	return phs.Size == 0
}

// GetSize returns the current number of songs in history
// Time Complexity: O(1)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) GetSize() int {
	return phs.Size
}

// GetMaxSize returns the maximum capacity of the history stack
// Time Complexity: O(1)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) GetMaxSize() int {
	return phs.MaxSize
}

// SetMaxSize updates the maximum size of the history stack
// If new size is smaller, removes oldest entries
// Time Complexity: O(k) where k is number of entries to remove
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) SetMaxSize(newMaxSize int) {
	if newMaxSize <= 0 {
		return
	}

	phs.MaxSize = newMaxSize

	// Remove excess entries if current size exceeds new max size
	for phs.Size > phs.MaxSize {
		phs.removeBottom()
	}
}

// Clear removes all songs from playback history
// Time Complexity: O(1)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) Clear() {
	phs.Top = nil
	phs.Size = 0
}

// ToSlice returns all songs in history as a slice (top to bottom)
// Time Complexity: O(n)
// Space Complexity: O(n)
func (phs *PlaybackHistoryStack) ToSlice() []*models.Song {
	songs := make([]*models.Song, 0, phs.Size)
	current := phs.Top

	for current != nil {
		songs = append(songs, current.Song)
		current = current.Next
	}

	return songs
}

// GetRecentSongs returns the n most recently played songs
// Time Complexity: O(min(n, size))
// Space Complexity: O(min(n, size))
func (phs *PlaybackHistoryStack) GetRecentSongs(n int) []*models.Song {
	if n <= 0 {
		return []*models.Song{}
	}

	songs := make([]*models.Song, 0, min(n, phs.Size))
	current := phs.Top
	count := 0

	for current != nil && count < n {
		songs = append(songs, current.Song)
		current = current.Next
		count++
	}

	return songs
}

// ContainsSong checks if a specific song is in the playback history
// Time Complexity: O(n)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) ContainsSong(songID string) bool {
	current := phs.Top

	for current != nil {
		if current.Song.ID == songID {
			return true
		}
		current = current.Next
	}

	return false
}

// GetPlaybackStats returns statistics about the playback history
// Time Complexity: O(n)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) GetPlaybackStats() map[string]interface{} {
	if phs.IsEmpty() {
		return map[string]interface{}{
			"total_songs":    0,
			"total_duration": 0,
			"unique_artists": 0,
			"unique_genres":  0,
		}
	}

	totalDuration := 0
	artistSet := make(map[string]bool)
	genreSet := make(map[string]bool)

	current := phs.Top
	for current != nil {
		totalDuration += current.Song.Duration
		artistSet[current.Song.Artist] = true
		genreSet[current.Song.Genre] = true
		current = current.Next
	}

	return map[string]interface{}{
		"total_songs":    phs.Size,
		"total_duration": totalDuration,
		"unique_artists": len(artistSet),
		"unique_genres":  len(genreSet),
	}
}

// removeBottom is a helper method to remove the bottom (oldest) entry
// Used when stack exceeds maximum size
// Time Complexity: O(n)
// Space Complexity: O(1)
func (phs *PlaybackHistoryStack) removeBottom() {
	if phs.Size <= 1 {
		phs.Clear()
		return
	}

	// Traverse to second-to-last node
	current := phs.Top
	for current.Next.Next != nil {
		current = current.Next
	}

	// Remove the last node
	current.Next = nil
	phs.Size--
}

// String returns a string representation of the playback history
// Time Complexity: O(n)
// Space Complexity: O(n)
func (phs *PlaybackHistoryStack) String() string {
	if phs.IsEmpty() {
		return "No playback history"
	}

	result := "Playback History (Most Recent First):\n"
	current := phs.Top
	index := 1

	for current != nil {
		result += fmt.Sprintf("%d. %s - %s (%s)\n",
			index, current.Song.Title, current.Song.Artist, current.Song.DurationString())
		current = current.Next
		index++
	}

	return result
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
