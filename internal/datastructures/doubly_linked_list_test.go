package datastructures

import (
	"src/internal/models"
	"testing"
)

// Test helper function to create a test song
func createTestSong(id, title, artist string) *models.Song {
	return models.NewSong(id, title, artist, "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
}

func TestNewDoublyLinkedList(t *testing.T) {
	dll := NewDoublyLinkedList()

	if dll.Head != nil {
		t.Errorf("NewDoublyLinkedList() Head should be nil")
	}
	if dll.Tail != nil {
		t.Errorf("NewDoublyLinkedList() Tail should be nil")
	}
	if dll.Length != 0 {
		t.Errorf("NewDoublyLinkedList() Length = %v, want %v", dll.Length, 0)
	}
}

func TestDoublyLinkedList_AddSong(t *testing.T) {
	dll := NewDoublyLinkedList()
	song1 := createTestSong("1", "Song 1", "Artist 1")
	song2 := createTestSong("2", "Song 2", "Artist 2")

	// Test adding first song
	dll.AddSong(song1)

	if dll.Length != 1 {
		t.Errorf("AddSong() Length = %v, want %v", dll.Length, 1)
	}
	if dll.Head == nil || dll.Head.Song.ID != "1" {
		t.Errorf("AddSong() Head song ID = %v, want %v", dll.Head.Song.ID, "1")
	}
	if dll.Tail == nil || dll.Tail.Song.ID != "1" {
		t.Errorf("AddSong() Tail song ID = %v, want %v", dll.Tail.Song.ID, "1")
	}
	if dll.Head != dll.Tail {
		t.Errorf("AddSong() Head and Tail should be same for single element")
	}

	// Test adding second song
	dll.AddSong(song2)

	if dll.Length != 2 {
		t.Errorf("AddSong() Length = %v, want %v", dll.Length, 2)
	}
	if dll.Head.Song.ID != "1" {
		t.Errorf("AddSong() Head song ID = %v, want %v", dll.Head.Song.ID, "1")
	}
	if dll.Tail.Song.ID != "2" {
		t.Errorf("AddSong() Tail song ID = %v, want %v", dll.Tail.Song.ID, "2")
	}
	if dll.Head.Next != dll.Tail {
		t.Errorf("AddSong() Head.Next should point to Tail")
	}
	if dll.Tail.Prev != dll.Head {
		t.Errorf("AddSong() Tail.Prev should point to Head")
	}
}

func TestDoublyLinkedList_AddSongToBeginning(t *testing.T) {
	dll := NewDoublyLinkedList()
	song1 := createTestSong("1", "Song 1", "Artist 1")
	song2 := createTestSong("2", "Song 2", "Artist 2")

	// Add first song
	dll.AddSong(song1)

	// Add song to beginning
	dll.AddSongToBeginning(song2)

	if dll.Length != 2 {
		t.Errorf("AddSongToBeginning() Length = %v, want %v", dll.Length, 2)
	}
	if dll.Head.Song.ID != "2" {
		t.Errorf("AddSongToBeginning() Head song ID = %v, want %v", dll.Head.Song.ID, "2")
	}
	if dll.Tail.Song.ID != "1" {
		t.Errorf("AddSongToBeginning() Tail song ID = %v, want %v", dll.Tail.Song.ID, "1")
	}
	if dll.Head.Prev != nil {
		t.Errorf("AddSongToBeginning() Head.Prev should be nil")
	}
	if dll.Tail.Next != nil {
		t.Errorf("AddSongToBeginning() Tail.Next should be nil")
	}
}

func TestDoublyLinkedList_AddSongAtIndex(t *testing.T) {
	dll := NewDoublyLinkedList()
	songs := []*models.Song{
		createTestSong("1", "Song 1", "Artist 1"),
		createTestSong("2", "Song 2", "Artist 2"),
		createTestSong("3", "Song 3", "Artist 3"),
		createTestSong("4", "Song 4", "Artist 4"),
	}

	// Add songs to create initial list
	for _, song := range songs[:3] {
		dll.AddSong(song)
	}

	// Test adding at valid index
	err := dll.AddSongAtIndex(songs[3], 1)
	if err != nil {
		t.Errorf("AddSongAtIndex() error = %v, want nil", err)
	}
	if dll.Length != 4 {
		t.Errorf("AddSongAtIndex() Length = %v, want %v", dll.Length, 4)
	}

	// Verify order: Song 1, Song 4, Song 2, Song 3
	song, _ := dll.GetSong(0)
	if song.ID != "1" {
		t.Errorf("AddSongAtIndex() song at index 0 = %v, want %v", song.ID, "1")
	}
	song, _ = dll.GetSong(1)
	if song.ID != "4" {
		t.Errorf("AddSongAtIndex() song at index 1 = %v, want %v", song.ID, "4")
	}

	// Test invalid index
	err = dll.AddSongAtIndex(createTestSong("5", "Song 5", "Artist 5"), 10)
	if err == nil {
		t.Errorf("AddSongAtIndex() with invalid index should return error")
	}
}

func TestDoublyLinkedList_DeleteSong(t *testing.T) {
	dll := NewDoublyLinkedList()
	songs := []*models.Song{
		createTestSong("1", "Song 1", "Artist 1"),
		createTestSong("2", "Song 2", "Artist 2"),
		createTestSong("3", "Song 3", "Artist 3"),
	}

	// Add songs
	for _, song := range songs {
		dll.AddSong(song)
	}

	// Test deleting middle song
	deletedSong, err := dll.DeleteSong(1)
	if err != nil {
		t.Errorf("DeleteSong() error = %v, want nil", err)
	}
	if deletedSong.ID != "2" {
		t.Errorf("DeleteSong() deleted song ID = %v, want %v", deletedSong.ID, "2")
	}
	if dll.Length != 2 {
		t.Errorf("DeleteSong() Length = %v, want %v", dll.Length, 2)
	}

	// Verify list integrity
	if dll.Head.Next != dll.Tail {
		t.Errorf("DeleteSong() Head.Next should point to Tail")
	}
	if dll.Tail.Prev != dll.Head {
		t.Errorf("DeleteSong() Tail.Prev should point to Head")
	}

	// Test deleting first song
	_, err = dll.DeleteSong(0)
	if err != nil {
		t.Errorf("DeleteSong() error = %v, want nil", err)
	}
	if dll.Length != 1 {
		t.Errorf("DeleteSong() Length = %v, want %v", dll.Length, 1)
	}
	if dll.Head != dll.Tail {
		t.Errorf("DeleteSong() Head and Tail should be same for single element")
	}

	// Test deleting last song
	_, err = dll.DeleteSong(0)
	if err != nil {
		t.Errorf("DeleteSong() error = %v, want nil", err)
	}
	if dll.Length != 0 {
		t.Errorf("DeleteSong() Length = %v, want %v", dll.Length, 0)
	}
	if dll.Head != nil || dll.Tail != nil {
		t.Errorf("DeleteSong() Head and Tail should be nil for empty list")
	}

	// Test invalid index
	_, err = dll.DeleteSong(0)
	if err == nil {
		t.Errorf("DeleteSong() with invalid index should return error")
	}
}

func TestDoublyLinkedList_MoveSong(t *testing.T) {
	dll := NewDoublyLinkedList()
	songs := []*models.Song{
		createTestSong("1", "Song 1", "Artist 1"),
		createTestSong("2", "Song 2", "Artist 2"),
		createTestSong("3", "Song 3", "Artist 3"),
		createTestSong("4", "Song 4", "Artist 4"),
	}

	// Add songs
	for _, song := range songs {
		dll.AddSong(song)
	}

	// Test moving song from index 0 to index 2
	err := dll.MoveSong(0, 2)
	if err != nil {
		t.Errorf("MoveSong() error = %v, want nil", err)
	}

	// Verify new order: Song 2, Song 3, Song 1, Song 4
	expectedOrder := []string{"2", "3", "1", "4"}
	for i, expectedID := range expectedOrder {
		song, _ := dll.GetSong(i)
		if song.ID != expectedID {
			t.Errorf("MoveSong() song at index %d = %v, want %v", i, song.ID, expectedID)
		}
	}

	// Test moving to same position
	err = dll.MoveSong(1, 1)
	if err != nil {
		t.Errorf("MoveSong() same position error = %v, want nil", err)
	}

	// Test invalid indices
	err = dll.MoveSong(0, 10)
	if err == nil {
		t.Errorf("MoveSong() with invalid toIndex should return error")
	}
}

func TestDoublyLinkedList_ReversePlaylist(t *testing.T) {
	dll := NewDoublyLinkedList()
	songs := []*models.Song{
		createTestSong("1", "Song 1", "Artist 1"),
		createTestSong("2", "Song 2", "Artist 2"),
		createTestSong("3", "Song 3", "Artist 3"),
		createTestSong("4", "Song 4", "Artist 4"),
	}

	// Add songs
	for _, song := range songs {
		dll.AddSong(song)
	}

	// Reverse playlist
	dll.ReversePlaylist()

	// Verify reversed order: Song 4, Song 3, Song 2, Song 1
	expectedOrder := []string{"4", "3", "2", "1"}
	for i, expectedID := range expectedOrder {
		song, _ := dll.GetSong(i)
		if song.ID != expectedID {
			t.Errorf("ReversePlaylist() song at index %d = %v, want %v", i, song.ID, expectedID)
		}
	}

	// Test reversing empty list
	emptyDll := NewDoublyLinkedList()
	emptyDll.ReversePlaylist() // Should not panic

	// Test reversing single element
	singleDll := NewDoublyLinkedList()
	singleDll.AddSong(createTestSong("1", "Song 1", "Artist 1"))
	singleDll.ReversePlaylist()
	if singleDll.Head.Song.ID != "1" {
		t.Errorf("ReversePlaylist() single element should remain unchanged")
	}
}

func TestDoublyLinkedList_GetSong(t *testing.T) {
	dll := NewDoublyLinkedList()
	songs := []*models.Song{
		createTestSong("1", "Song 1", "Artist 1"),
		createTestSong("2", "Song 2", "Artist 2"),
		createTestSong("3", "Song 3", "Artist 3"),
	}

	// Test empty list
	_, err := dll.GetSong(0)
	if err == nil {
		t.Errorf("GetSong() on empty list should return error")
	}

	// Add songs
	for _, song := range songs {
		dll.AddSong(song)
	}

	// Test valid indices
	for i, expectedSong := range songs {
		song, err := dll.GetSong(i)
		if err != nil {
			t.Errorf("GetSong(%d) error = %v, want nil", i, err)
		}
		if song.ID != expectedSong.ID {
			t.Errorf("GetSong(%d) ID = %v, want %v", i, song.ID, expectedSong.ID)
		}
	}

	// Test invalid indices
	_, err = dll.GetSong(-1)
	if err == nil {
		t.Errorf("GetSong() with negative index should return error")
	}

	_, err = dll.GetSong(3)
	if err == nil {
		t.Errorf("GetSong() with out of bounds index should return error")
	}
}

func TestDoublyLinkedList_ToSlice(t *testing.T) {
	dll := NewDoublyLinkedList()

	// Test empty list
	slice := dll.ToSlice()
	if len(slice) != 0 {
		t.Errorf("ToSlice() empty list length = %v, want %v", len(slice), 0)
	}

	// Add songs
	songs := []*models.Song{
		createTestSong("1", "Song 1", "Artist 1"),
		createTestSong("2", "Song 2", "Artist 2"),
		createTestSong("3", "Song 3", "Artist 3"),
	}

	for _, song := range songs {
		dll.AddSong(song)
	}

	slice = dll.ToSlice()
	if len(slice) != len(songs) {
		t.Errorf("ToSlice() length = %v, want %v", len(slice), len(songs))
	}

	for i, song := range slice {
		if song.ID != songs[i].ID {
			t.Errorf("ToSlice()[%d] ID = %v, want %v", i, song.ID, songs[i].ID)
		}
	}
}

func TestDoublyLinkedList_Size(t *testing.T) {
	dll := NewDoublyLinkedList()

	if dll.Size() != 0 {
		t.Errorf("Size() = %v, want %v", dll.Size(), 0)
	}

	// Add songs and verify size
	for i := 0; i < 5; i++ {
		dll.AddSong(createTestSong(string(rune('1'+i)), "Song", "Artist"))
		if dll.Size() != i+1 {
			t.Errorf("Size() after adding %d songs = %v, want %v", i+1, dll.Size(), i+1)
		}
	}
}

func TestDoublyLinkedList_IsEmpty(t *testing.T) {
	dll := NewDoublyLinkedList()

	if !dll.IsEmpty() {
		t.Errorf("IsEmpty() = %v, want %v", dll.IsEmpty(), true)
	}

	dll.AddSong(createTestSong("1", "Song 1", "Artist 1"))

	if dll.IsEmpty() {
		t.Errorf("IsEmpty() after adding song = %v, want %v", dll.IsEmpty(), false)
	}
}

func TestDoublyLinkedList_Clear(t *testing.T) {
	dll := NewDoublyLinkedList()

	// Add some songs
	for i := 0; i < 3; i++ {
		dll.AddSong(createTestSong(string(rune('1'+i)), "Song", "Artist"))
	}

	dll.Clear()

	if !dll.IsEmpty() {
		t.Errorf("Clear() IsEmpty() = %v, want %v", dll.IsEmpty(), true)
	}
	if dll.Size() != 0 {
		t.Errorf("Clear() Size() = %v, want %v", dll.Size(), 0)
	}
	if dll.Head != nil || dll.Tail != nil {
		t.Errorf("Clear() Head and Tail should be nil")
	}
}

func TestDoublyLinkedList_GetTotalDuration(t *testing.T) {
	dll := NewDoublyLinkedList()

	// Test empty list
	if dll.GetTotalDuration() != 0 {
		t.Errorf("GetTotalDuration() empty list = %v, want %v", dll.GetTotalDuration(), 0)
	}

	// Add songs with different durations
	durations := []int{120, 180, 240}
	expectedTotal := 0

	for i, duration := range durations {
		song := models.NewSong(string(rune('1'+i)), "Song", "Artist", "Album", "Genre", "Sub", "Mood", duration, 120)
		dll.AddSong(song)
		expectedTotal += duration
	}

	if dll.GetTotalDuration() != expectedTotal {
		t.Errorf("GetTotalDuration() = %v, want %v", dll.GetTotalDuration(), expectedTotal)
	}
}

func TestDoublyLinkedList_FindSongByID(t *testing.T) {
	dll := NewDoublyLinkedList()
	songs := []*models.Song{
		createTestSong("song1", "Song 1", "Artist 1"),
		createTestSong("song2", "Song 2", "Artist 2"),
		createTestSong("song3", "Song 3", "Artist 3"),
	}

	// Test empty list
	_, err := dll.FindSongByID("song1")
	if err == nil {
		t.Errorf("FindSongByID() on empty list should return error")
	}

	// Add songs
	for _, song := range songs {
		dll.AddSong(song)
	}

	// Test finding existing songs
	for i, song := range songs {
		index, err := dll.FindSongByID(song.ID)
		if err != nil {
			t.Errorf("FindSongByID(%s) error = %v, want nil", song.ID, err)
		}
		if index != i {
			t.Errorf("FindSongByID(%s) index = %v, want %v", song.ID, index, i)
		}
	}

	// Test finding non-existing song
	_, err = dll.FindSongByID("nonexistent")
	if err == nil {
		t.Errorf("FindSongByID() non-existing song should return error")
	}
}

// Benchmark tests
func BenchmarkDoublyLinkedList_AddSong(b *testing.B) {
	dll := NewDoublyLinkedList()
	song := createTestSong("test", "Test Song", "Test Artist")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dll.AddSong(song)
	}
}

func BenchmarkDoublyLinkedList_AddSongToBeginning(b *testing.B) {
	dll := NewDoublyLinkedList()
	song := createTestSong("test", "Test Song", "Test Artist")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dll.AddSongToBeginning(song)
	}
}

func BenchmarkDoublyLinkedList_GetSong_Head(b *testing.B) {
	dll := NewDoublyLinkedList()
	// Add 1000 songs
	for i := 0; i < 1000; i++ {
		dll.AddSong(createTestSong(string(rune(i)), "Song", "Artist"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dll.GetSong(0) // Get from head
	}
}

func BenchmarkDoublyLinkedList_GetSong_Tail(b *testing.B) {
	dll := NewDoublyLinkedList()
	// Add 1000 songs
	for i := 0; i < 1000; i++ {
		dll.AddSong(createTestSong(string(rune(i)), "Song", "Artist"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dll.GetSong(999) // Get from tail
	}
}

func BenchmarkDoublyLinkedList_GetSong_Middle(b *testing.B) {
	dll := NewDoublyLinkedList()
	// Add 1000 songs
	for i := 0; i < 1000; i++ {
		dll.AddSong(createTestSong(string(rune(i)), "Song", "Artist"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dll.GetSong(500) // Get from middle
	}
}

func BenchmarkDoublyLinkedList_ReversePlaylist(b *testing.B) {
	dll := NewDoublyLinkedList()
	// Add 100 songs
	for i := 0; i < 100; i++ {
		dll.AddSong(createTestSong(string(rune(i)), "Song", "Artist"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dll.ReversePlaylist()
		// Reverse again to maintain state
		dll.ReversePlaylist()
	}
}

func BenchmarkDoublyLinkedList_ToSlice(b *testing.B) {
	dll := NewDoublyLinkedList()
	// Add 1000 songs
	for i := 0; i < 1000; i++ {
		dll.AddSong(createTestSong(string(rune(i)), "Song", "Artist"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dll.ToSlice()
	}
}

// Edge case tests
func TestDoublyLinkedList_EdgeCases(t *testing.T) {
	dll := NewDoublyLinkedList()

	// Test operations on empty list
	dll.ReversePlaylist() // Should not panic
	dll.Clear()           // Should not panic

	// Test with nil song (should handle gracefully)
	dll.AddSong(nil) // Implementation should handle this
	if dll.Size() != 0 {
		t.Errorf("Adding nil song should not increase size")
	}

	// Test bidirectional traversal optimization
	dll = NewDoublyLinkedList()
	for i := 0; i < 10; i++ {
		dll.AddSong(createTestSong(string(rune('0'+i)), "Song", "Artist"))
	}

	// Access near head (should traverse from head)
	song, _ := dll.GetSong(2)
	if song.ID != "2" {
		t.Errorf("GetSong(2) from head traversal failed")
	}

	// Access near tail (should traverse from tail)
	song, _ = dll.GetSong(8)
	if song.ID != "8" {
		t.Errorf("GetSong(8) from tail traversal failed")
	}
}
