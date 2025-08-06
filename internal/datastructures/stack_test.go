package datastructures

import (
	"src/internal/models"
	"testing"
)

// Test helper function to create a test song for stack tests
func createStackTestSong(id, title, artist string) *models.Song {
	return models.NewSong(id, title, artist, "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
}

func TestNewPlaybackHistoryStack(t *testing.T) {
	tests := []struct {
		name        string
		maxSize     int
		expectedMax int
	}{
		{"Valid max size", 50, 50},
		{"Zero max size", 0, 50},       // Should default to 50
		{"Negative max size", -10, 50}, // Should default to 50
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewPlaybackHistoryStack(tt.maxSize)

			if stack.Top != nil {
				t.Errorf("NewPlaybackHistoryStack() Top should be nil")
			}
			if stack.Size != 0 {
				t.Errorf("NewPlaybackHistoryStack() Size = %v, want %v", stack.Size, 0)
			}
			if stack.MaxSize != tt.expectedMax {
				t.Errorf("NewPlaybackHistoryStack() MaxSize = %v, want %v", stack.MaxSize, tt.expectedMax)
			}
		})
	}
}

func TestPlaybackHistoryStack_Push(t *testing.T) {
	stack := NewPlaybackHistoryStack(3) // Small size for testing
	songs := []*models.Song{
		createStackTestSong("1", "Song 1", "Artist 1"),
		createStackTestSong("2", "Song 2", "Artist 2"),
		createStackTestSong("3", "Song 3", "Artist 3"),
		createStackTestSong("4", "Song 4", "Artist 4"),
	}

	// Test pushing first song
	stack.Push(songs[0])

	if stack.Size != 1 {
		t.Errorf("Push() Size = %v, want %v", stack.Size, 1)
	}
	if stack.Top == nil || stack.Top.Song.ID != "1" {
		t.Errorf("Push() Top song ID = %v, want %v", stack.Top.Song.ID, "1")
	}

	// Test pushing more songs
	stack.Push(songs[1])
	stack.Push(songs[2])

	if stack.Size != 3 {
		t.Errorf("Push() Size after 3 pushes = %v, want %v", stack.Size, 3)
	}
	if stack.Top.Song.ID != "3" {
		t.Errorf("Push() Top song ID = %v, want %v", stack.Top.Song.ID, "3")
	}

	// Test pushing beyond max size (should remove oldest)
	stack.Push(songs[3])

	if stack.Size != 3 {
		t.Errorf("Push() Size after exceeding max = %v, want %v", stack.Size, 3)
	}
	if stack.Top.Song.ID != "4" {
		t.Errorf("Push() Top song ID after overflow = %v, want %v", stack.Top.Song.ID, "4")
	}

	// Verify that the oldest song was removed
	allSongs := stack.ToSlice()
	songIDs := make([]string, len(allSongs))
	for i, song := range allSongs {
		songIDs[i] = song.ID
	}

	// Should contain songs 4, 3, 2 (newest to oldest)
	expected := []string{"4", "3", "2"}
	for i, expectedID := range expected {
		if songIDs[i] != expectedID {
			t.Errorf("Push() after overflow song[%d] = %v, want %v", i, songIDs[i], expectedID)
		}
	}
}

func TestPlaybackHistoryStack_Pop(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	// Test popping from empty stack
	_, err := stack.Pop()
	if err == nil {
		t.Errorf("Pop() on empty stack should return error")
	}

	// Add songs and test popping
	songs := []*models.Song{
		createStackTestSong("1", "Song 1", "Artist 1"),
		createStackTestSong("2", "Song 2", "Artist 2"),
		createStackTestSong("3", "Song 3", "Artist 3"),
	}

	for _, song := range songs {
		stack.Push(song)
	}

	// Pop songs and verify LIFO order
	for i := len(songs) - 1; i >= 0; i-- {
		song, err := stack.Pop()
		if err != nil {
			t.Errorf("Pop() error = %v, want nil", err)
		}
		if song.ID != songs[i].ID {
			t.Errorf("Pop() song ID = %v, want %v", song.ID, songs[i].ID)
		}
		if stack.Size != i {
			t.Errorf("Pop() Size = %v, want %v", stack.Size, i)
		}
	}

	// Verify stack is empty
	if !stack.IsEmpty() {
		t.Errorf("Pop() stack should be empty after popping all songs")
	}
}

func TestPlaybackHistoryStack_Peek(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	// Test peeking empty stack
	_, err := stack.Peek()
	if err == nil {
		t.Errorf("Peek() on empty stack should return error")
	}

	// Add song and test peeking
	song := createStackTestSong("1", "Song 1", "Artist 1")
	stack.Push(song)

	peekedSong, err := stack.Peek()
	if err != nil {
		t.Errorf("Peek() error = %v, want nil", err)
	}
	if peekedSong.ID != "1" {
		t.Errorf("Peek() song ID = %v, want %v", peekedSong.ID, "1")
	}
	if stack.Size != 1 {
		t.Errorf("Peek() should not change Size = %v, want %v", stack.Size, 1)
	}

	// Add another song and peek again
	song2 := createStackTestSong("2", "Song 2", "Artist 2")
	stack.Push(song2)

	peekedSong, err = stack.Peek()
	if err != nil {
		t.Errorf("Peek() error = %v, want nil", err)
	}
	if peekedSong.ID != "2" {
		t.Errorf("Peek() should return top song ID = %v, want %v", peekedSong.ID, "2")
	}
}

func TestPlaybackHistoryStack_UndoLastPlay(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	// Test undo on empty stack
	_, err := stack.UndoLastPlay()
	if err == nil {
		t.Errorf("UndoLastPlay() on empty stack should return error")
	}

	// Add songs and test undo
	songs := []*models.Song{
		createStackTestSong("1", "Song 1", "Artist 1"),
		createStackTestSong("2", "Song 2", "Artist 2"),
	}

	for _, song := range songs {
		stack.Push(song)
	}

	// Undo last play
	undoSong, err := stack.UndoLastPlay()
	if err != nil {
		t.Errorf("UndoLastPlay() error = %v, want nil", err)
	}
	if undoSong.ID != "2" {
		t.Errorf("UndoLastPlay() song ID = %v, want %v", undoSong.ID, "2")
	}
	if stack.Size != 1 {
		t.Errorf("UndoLastPlay() Size = %v, want %v", stack.Size, 1)
	}
}

func TestPlaybackHistoryStack_IsEmpty(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	if !stack.IsEmpty() {
		t.Errorf("IsEmpty() = %v, want %v", stack.IsEmpty(), true)
	}

	stack.Push(createStackTestSong("1", "Song 1", "Artist 1"))

	if stack.IsEmpty() {
		t.Errorf("IsEmpty() after Push = %v, want %v", stack.IsEmpty(), false)
	}
}

func TestPlaybackHistoryStack_GetSize(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	if stack.GetSize() != 0 {
		t.Errorf("GetSize() = %v, want %v", stack.GetSize(), 0)
	}

	// Add songs and verify size
	for i := 0; i < 3; i++ {
		stack.Push(createStackTestSong(string(rune('1'+i)), "Song", "Artist"))
		if stack.GetSize() != i+1 {
			t.Errorf("GetSize() after %d pushes = %v, want %v", i+1, stack.GetSize(), i+1)
		}
	}
}

func TestPlaybackHistoryStack_SetMaxSize(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	// Add 5 songs
	for i := 0; i < 5; i++ {
		stack.Push(createStackTestSong(string(rune('1'+i)), "Song", "Artist"))
	}

	// Reduce max size to 3
	stack.SetMaxSize(3)

	if stack.GetMaxSize() != 3 {
		t.Errorf("SetMaxSize() MaxSize = %v, want %v", stack.GetMaxSize(), 3)
	}
	if stack.GetSize() != 3 {
		t.Errorf("SetMaxSize() should reduce Size = %v, want %v", stack.GetSize(), 3)
	}

	// Verify remaining songs are the most recent ones
	songs := stack.ToSlice()
	expected := []string{"5", "4", "3"} // Most recent first
	for i, song := range songs {
		if song.ID != expected[i] {
			t.Errorf("SetMaxSize() remaining song[%d] = %v, want %v", i, song.ID, expected[i])
		}
	}

	// Test invalid max size
	stack.SetMaxSize(0)
	if stack.GetMaxSize() != 3 {
		t.Errorf("SetMaxSize(0) should not change MaxSize")
	}

	stack.SetMaxSize(-1)
	if stack.GetMaxSize() != 3 {
		t.Errorf("SetMaxSize(-1) should not change MaxSize")
	}
}

func TestPlaybackHistoryStack_Clear(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	// Add songs
	for i := 0; i < 3; i++ {
		stack.Push(createStackTestSong(string(rune('1'+i)), "Song", "Artist"))
	}

	stack.Clear()

	if !stack.IsEmpty() {
		t.Errorf("Clear() IsEmpty() = %v, want %v", stack.IsEmpty(), true)
	}
	if stack.GetSize() != 0 {
		t.Errorf("Clear() Size = %v, want %v", stack.GetSize(), 0)
	}
	if stack.Top != nil {
		t.Errorf("Clear() Top should be nil")
	}
}

func TestPlaybackHistoryStack_ToSlice(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	// Test empty stack
	slice := stack.ToSlice()
	if len(slice) != 0 {
		t.Errorf("ToSlice() empty stack length = %v, want %v", len(slice), 0)
	}

	// Add songs
	songs := []*models.Song{
		createStackTestSong("1", "Song 1", "Artist 1"),
		createStackTestSong("2", "Song 2", "Artist 2"),
		createStackTestSong("3", "Song 3", "Artist 3"),
	}

	for _, song := range songs {
		stack.Push(song)
	}

	slice = stack.ToSlice()
	if len(slice) != len(songs) {
		t.Errorf("ToSlice() length = %v, want %v", len(slice), len(songs))
	}

	// Verify order (should be newest first - LIFO order)
	expected := []string{"3", "2", "1"}
	for i, song := range slice {
		if song.ID != expected[i] {
			t.Errorf("ToSlice()[%d] ID = %v, want %v", i, song.ID, expected[i])
		}
	}
}

func TestPlaybackHistoryStack_GetRecentSongs(t *testing.T) {
	stack := NewPlaybackHistoryStack(10)

	// Test with empty stack
	recent := stack.GetRecentSongs(5)
	if len(recent) != 0 {
		t.Errorf("GetRecentSongs() empty stack length = %v, want %v", len(recent), 0)
	}

	// Add songs
	for i := 0; i < 7; i++ {
		stack.Push(createStackTestSong(string(rune('1'+i)), "Song", "Artist"))
	}

	// Test getting recent songs
	tests := []struct {
		name     string
		count    int
		expected int
	}{
		{"Get 3 recent", 3, 3},
		{"Get more than available", 10, 7},
		{"Get 0 songs", 0, 0},
		{"Get negative count", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recent := stack.GetRecentSongs(tt.count)
			if len(recent) != tt.expected {
				t.Errorf("GetRecentSongs(%d) length = %v, want %v", tt.count, len(recent), tt.expected)
			}

			// Verify order (newest first)
			if len(recent) > 0 {
				for i, song := range recent {
					expectedID := string(rune('7' - i)) // Most recent first
					if song.ID != expectedID {
						t.Errorf("GetRecentSongs(%d)[%d] ID = %v, want %v", tt.count, i, song.ID, expectedID)
					}
				}
			}
		})
	}
}

func TestPlaybackHistoryStack_ContainsSong(t *testing.T) {
	stack := NewPlaybackHistoryStack(5)

	// Test empty stack
	if stack.ContainsSong("1") {
		t.Errorf("ContainsSong() empty stack should return false")
	}

	// Add songs
	songs := []*models.Song{
		createStackTestSong("song1", "Song 1", "Artist 1"),
		createStackTestSong("song2", "Song 2", "Artist 2"),
		createStackTestSong("song3", "Song 3", "Artist 3"),
	}

	for _, song := range songs {
		stack.Push(song)
	}

	// Test existing songs
	for _, song := range songs {
		if !stack.ContainsSong(song.ID) {
			t.Errorf("ContainsSong(%s) = %v, want %v", song.ID, false, true)
		}
	}

	// Test non-existing song
	if stack.ContainsSong("nonexistent") {
		t.Errorf("ContainsSong(nonexistent) = %v, want %v", true, false)
	}
}

func TestPlaybackHistoryStack_GetPlaybackStats(t *testing.T) {
	stack := NewPlaybackHistoryStack(10)

	// Test empty stack
	stats := stack.GetPlaybackStats()
	expectedStats := map[string]interface{}{
		"total_songs":    0,
		"total_duration": 0,
		"unique_artists": 0,
		"unique_genres":  0,
	}

	for key, expected := range expectedStats {
		if stats[key] != expected {
			t.Errorf("GetPlaybackStats() empty stack %s = %v, want %v", key, stats[key], expected)
		}
	}

	// Add songs with different artists and genres
	songs := []*models.Song{
		models.NewSong("1", "Song 1", "Artist A", "Album", "Rock", "Alt", "Happy", 180, 120),
		models.NewSong("2", "Song 2", "Artist B", "Album", "Pop", "Dance", "Upbeat", 200, 130),
		models.NewSong("3", "Song 3", "Artist A", "Album", "Rock", "Hard", "Aggressive", 220, 140),
	}

	for _, song := range songs {
		stack.Push(song)
	}

	stats = stack.GetPlaybackStats()

	if stats["total_songs"] != 3 {
		t.Errorf("GetPlaybackStats() total_songs = %v, want %v", stats["total_songs"], 3)
	}
	if stats["total_duration"] != 600 {
		t.Errorf("GetPlaybackStats() total_duration = %v, want %v", stats["total_duration"], 600)
	}
	if stats["unique_artists"] != 2 {
		t.Errorf("GetPlaybackStats() unique_artists = %v, want %v", stats["unique_artists"], 2)
	}
	if stats["unique_genres"] != 2 {
		t.Errorf("GetPlaybackStats() unique_genres = %v, want %v", stats["unique_genres"], 2)
	}
}

// Benchmark tests
func BenchmarkPlaybackHistoryStack_Push(b *testing.B) {
	stack := NewPlaybackHistoryStack(1000)
	song := createStackTestSong("test", "Test Song", "Test Artist")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Push(song)
	}
}

func BenchmarkPlaybackHistoryStack_Pop(b *testing.B) {
	stack := NewPlaybackHistoryStack(b.N + 100)
	song := createStackTestSong("test", "Test Song", "Test Artist")

	// Pre-fill stack
	for i := 0; i < b.N+100; i++ {
		stack.Push(song)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Pop()
	}
}

func BenchmarkPlaybackHistoryStack_Peek(b *testing.B) {
	stack := NewPlaybackHistoryStack(100)
	song := createStackTestSong("test", "Test Song", "Test Artist")
	stack.Push(song)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.Peek()
	}
}

func BenchmarkPlaybackHistoryStack_ContainsSong(b *testing.B) {
	stack := NewPlaybackHistoryStack(1000)

	// Add 100 songs
	for i := 0; i < 100; i++ {
		stack.Push(createStackTestSong(string(rune(i)), "Song", "Artist"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.ContainsSong("50") // Search for song in middle
	}
}

func BenchmarkPlaybackHistoryStack_ToSlice(b *testing.B) {
	stack := NewPlaybackHistoryStack(1000)

	// Add 100 songs
	for i := 0; i < 100; i++ {
		stack.Push(createStackTestSong(string(rune(i)), "Song", "Artist"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stack.ToSlice()
	}
}

// Edge case and stress tests
func TestPlaybackHistoryStack_EdgeCases(t *testing.T) {
	stack := NewPlaybackHistoryStack(2)

	// Test with nil songs (implementation should handle gracefully)
	stack.Push(nil)
	if stack.GetSize() != 0 {
		t.Errorf("Push(nil) should not increase size")
	}

	// Test rapid push/pop operations
	song1 := createStackTestSong("1", "Song 1", "Artist 1")
	song2 := createStackTestSong("2", "Song 2", "Artist 2")

	stack.Push(song1)
	stack.Push(song2)

	popped, _ := stack.Pop()
	if popped.ID != "2" {
		t.Errorf("Rapid push/pop failed: expected 2, got %s", popped.ID)
	}

	stack.Push(song2)
	peeked, _ := stack.Peek()
	if peeked.ID != "2" {
		t.Errorf("Push after pop failed: expected 2, got %s", peeked.ID)
	}
}

func TestPlaybackHistoryStack_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	stack := NewPlaybackHistoryStack(1000)

	// Add many songs
	for i := 0; i < 10000; i++ {
		song := createStackTestSong(string(rune(i%1000)), "Song", "Artist")
		stack.Push(song)
	}

	// Stack should be bounded to max size
	if stack.GetSize() != 1000 {
		t.Errorf("Stress test: Size = %v, want %v", stack.GetSize(), 1000)
	}

	// Pop all songs
	count := 0
	for !stack.IsEmpty() {
		_, err := stack.Pop()
		if err != nil {
			t.Errorf("Stress test Pop() error: %v", err)
		}
		count++
	}

	if count != 1000 {
		t.Errorf("Stress test: Popped %v songs, want %v", count, 1000)
	}
}
