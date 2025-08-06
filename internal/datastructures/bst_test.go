package datastructures

import (
	"src/internal/models"
	"testing"
)

// Test helper function to create a test song for BST tests
func createBSTTestSong(id, title, artist string, rating int) *models.Song {
	song := models.NewSong(id, title, artist, "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
	song.SetRating(rating)
	return song
}

func TestNewSongRatingBST(t *testing.T) {
	bst := NewSongRatingBST()

	if bst.Root != nil {
		t.Errorf("NewSongRatingBST() Root should be nil")
	}
	if bst.NodeCount != 0 {
		t.Errorf("NewSongRatingBST() NodeCount = %v, want %v", bst.NodeCount, 0)
	}
}

func TestNewRatingBucket(t *testing.T) {
	rating := 4
	bucket := NewRatingBucket(rating)

	if bucket.Rating != rating {
		t.Errorf("NewRatingBucket() Rating = %v, want %v", bucket.Rating, rating)
	}
	if len(bucket.Songs) != 0 {
		t.Errorf("NewRatingBucket() Songs length = %v, want %v", len(bucket.Songs), 0)
	}
	if bucket.IsEmpty() != true {
		t.Errorf("NewRatingBucket() IsEmpty() = %v, want %v", bucket.IsEmpty(), true)
	}
}

func TestRatingBucket_AddSong(t *testing.T) {
	bucket := NewRatingBucket(4)
	song1 := createBSTTestSong("1", "Song 1", "Artist 1", 4)
	song2 := createBSTTestSong("2", "Song 2", "Artist 2", 4)

	bucket.AddSong(song1)
	if len(bucket.Songs) != 1 {
		t.Errorf("AddSong() Songs length = %v, want %v", len(bucket.Songs), 1)
	}
	if bucket.Songs[0].ID != "1" {
		t.Errorf("AddSong() song ID = %v, want %v", bucket.Songs[0].ID, "1")
	}

	bucket.AddSong(song2)
	if len(bucket.Songs) != 2 {
		t.Errorf("AddSong() Songs length after second add = %v, want %v", len(bucket.Songs), 2)
	}
}

func TestRatingBucket_RemoveSong(t *testing.T) {
	bucket := NewRatingBucket(4)
	songs := []*models.Song{
		createBSTTestSong("1", "Song 1", "Artist 1", 4),
		createBSTTestSong("2", "Song 2", "Artist 2", 4),
		createBSTTestSong("3", "Song 3", "Artist 3", 4),
	}

	// Add songs
	for _, song := range songs {
		bucket.AddSong(song)
	}

	// Test removing existing song
	removed := bucket.RemoveSong("2")
	if !removed {
		t.Errorf("RemoveSong() existing song should return true")
	}
	if len(bucket.Songs) != 2 {
		t.Errorf("RemoveSong() Songs length = %v, want %v", len(bucket.Songs), 2)
	}

	// Verify remaining songs
	foundSong2 := false
	for _, song := range bucket.Songs {
		if song.ID == "2" {
			foundSong2 = true
		}
	}
	if foundSong2 {
		t.Errorf("RemoveSong() song should be removed from bucket")
	}

	// Test removing non-existing song
	removed = bucket.RemoveSong("nonexistent")
	if removed {
		t.Errorf("RemoveSong() non-existing song should return false")
	}
	if len(bucket.Songs) != 2 {
		t.Errorf("RemoveSong() Songs length should not change = %v, want %v", len(bucket.Songs), 2)
	}
}

func TestSongRatingBST_InsertSong(t *testing.T) {
	bst := NewSongRatingBST()

	// Test inserting first song
	song1 := createBSTTestSong("1", "Song 1", "Artist 1", 3)
	bst.InsertSong(song1, 3)

	if bst.NodeCount != 1 {
		t.Errorf("InsertSong() NodeCount = %v, want %v", bst.NodeCount, 1)
	}
	if bst.Root == nil {
		t.Errorf("InsertSong() Root should not be nil")
	}
	if bst.Root.Bucket.Rating != 3 {
		t.Errorf("InsertSong() Root rating = %v, want %v", bst.Root.Bucket.Rating, 3)
	}

	// Test inserting song with same rating
	song2 := createBSTTestSong("2", "Song 2", "Artist 2", 3)
	bst.InsertSong(song2, 3)

	if bst.NodeCount != 1 {
		t.Errorf("InsertSong() same rating NodeCount = %v, want %v", bst.NodeCount, 1)
	}
	if len(bst.Root.Bucket.Songs) != 2 {
		t.Errorf("InsertSong() same rating bucket size = %v, want %v", len(bst.Root.Bucket.Songs), 2)
	}

	// Test inserting songs with different ratings
	song3 := createBSTTestSong("3", "Song 3", "Artist 3", 1)
	song4 := createBSTTestSong("4", "Song 4", "Artist 4", 5)
	bst.InsertSong(song3, 1)
	bst.InsertSong(song4, 5)

	if bst.NodeCount != 3 {
		t.Errorf("InsertSong() different ratings NodeCount = %v, want %v", bst.NodeCount, 3)
	}
	if bst.Root.Left == nil || bst.Root.Left.Bucket.Rating != 1 {
		t.Errorf("InsertSong() left child rating should be 1")
	}
	if bst.Root.Right == nil || bst.Root.Right.Bucket.Rating != 5 {
		t.Errorf("InsertSong() right child rating should be 5")
	}

	// Test invalid rating
	song5 := createBSTTestSong("5", "Song 5", "Artist 5", 6)
	initialNodeCount := bst.NodeCount
	bst.InsertSong(song5, 6)
	if bst.NodeCount != initialNodeCount {
		t.Errorf("InsertSong() invalid rating should not change NodeCount")
	}
}

func TestSongRatingBST_SearchByRating(t *testing.T) {
	bst := NewSongRatingBST()

	// Test searching empty BST
	songs := bst.SearchByRating(3)
	if len(songs) != 0 {
		t.Errorf("SearchByRating() empty BST should return empty slice")
	}

	// Add test songs
	testSongs := []*models.Song{
		createBSTTestSong("1", "Song 1", "Artist 1", 3),
		createBSTTestSong("2", "Song 2", "Artist 2", 3),
		createBSTTestSong("3", "Song 3", "Artist 3", 1),
		createBSTTestSong("4", "Song 4", "Artist 4", 5),
	}

	for _, song := range testSongs {
		bst.InsertSong(song, song.Rating)
	}

	// Test searching for rating 3 (should return 2 songs)
	songs = bst.SearchByRating(3)
	if len(songs) != 2 {
		t.Errorf("SearchByRating(3) length = %v, want %v", len(songs), 2)
	}

	// Test searching for rating 1 (should return 1 song)
	songs = bst.SearchByRating(1)
	if len(songs) != 1 {
		t.Errorf("SearchByRating(1) length = %v, want %v", len(songs), 1)
	}
	if songs[0].ID != "3" {
		t.Errorf("SearchByRating(1) song ID = %v, want %v", songs[0].ID, "3")
	}

	// Test searching for non-existing rating
	songs = bst.SearchByRating(2)
	if len(songs) != 0 {
		t.Errorf("SearchByRating(2) non-existing rating should return empty slice")
	}

	// Test invalid rating
	songs = bst.SearchByRating(6)
	if len(songs) != 0 {
		t.Errorf("SearchByRating(6) invalid rating should return empty slice")
	}
}

func TestSongRatingBST_DeleteSong(t *testing.T) {
	bst := NewSongRatingBST()

	// Test deleting from empty BST
	deleted := bst.DeleteSong("nonexistent")
	if deleted {
		t.Errorf("DeleteSong() from empty BST should return false")
	}

	// Add test songs
	testSongs := []*models.Song{
		createBSTTestSong("1", "Song 1", "Artist 1", 3),
		createBSTTestSong("2", "Song 2", "Artist 2", 3),
		createBSTTestSong("3", "Song 3", "Artist 3", 1),
		createBSTTestSong("4", "Song 4", "Artist 4", 5),
	}

	for _, song := range testSongs {
		bst.InsertSong(song, song.Rating)
	}

	// Test deleting song from bucket with multiple songs
	deleted = bst.DeleteSong("1")
	if !deleted {
		t.Errorf("DeleteSong() existing song should return true")
	}
	songs := bst.SearchByRating(3)
	if len(songs) != 1 {
		t.Errorf("DeleteSong() rating 3 bucket should have 1 song left")
	}
	if songs[0].ID == "1" {
		t.Errorf("DeleteSong() song should be removed from bucket")
	}

	// Test deleting last song from bucket (should remove node)
	deleted = bst.DeleteSong("2")
	if !deleted {
		t.Errorf("DeleteSong() last song in bucket should return true")
	}
	songs = bst.SearchByRating(3)
	if len(songs) != 0 {
		t.Errorf("DeleteSong() rating 3 bucket should be empty after removing last song")
	}
	if bst.NodeCount != 2 {
		t.Errorf("DeleteSong() should remove empty node, NodeCount = %v, want %v", bst.NodeCount, 2)
	}

	// Test deleting non-existing song
	deleted = bst.DeleteSong("nonexistent")
	if deleted {
		t.Errorf("DeleteSong() non-existing song should return false")
	}
}

func TestSongRatingBST_GetAllSongs(t *testing.T) {
	bst := NewSongRatingBST()

	// Test empty BST
	songs := bst.GetAllSongs()
	if len(songs) != 0 {
		t.Errorf("GetAllSongs() empty BST should return empty slice")
	}

	// Add test songs with different ratings
	testSongs := []*models.Song{
		createBSTTestSong("3", "Song 3", "Artist 3", 3),
		createBSTTestSong("1", "Song 1", "Artist 1", 1),
		createBSTTestSong("5", "Song 5", "Artist 5", 5),
		createBSTTestSong("2", "Song 2", "Artist 2", 2),
		createBSTTestSong("4", "Song 4", "Artist 4", 4),
	}

	for _, song := range testSongs {
		bst.InsertSong(song, song.Rating)
	}

	songs = bst.GetAllSongs()
	if len(songs) != 5 {
		t.Errorf("GetAllSongs() length = %v, want %v", len(songs), 5)
	}

	// Verify songs are returned in sorted order by rating (inorder traversal)
	expectedOrder := []string{"1", "2", "3", "4", "5"}
	for i, song := range songs {
		if song.ID != expectedOrder[i] {
			t.Errorf("GetAllSongs() song[%d] ID = %v, want %v", i, song.ID, expectedOrder[i])
		}
	}
}

func TestSongRatingBST_GetSongsByRatingRange(t *testing.T) {
	bst := NewSongRatingBST()

	// Add test songs
	testSongs := []*models.Song{
		createBSTTestSong("1", "Song 1", "Artist 1", 1),
		createBSTTestSong("2", "Song 2", "Artist 2", 2),
		createBSTTestSong("3", "Song 3", "Artist 3", 3),
		createBSTTestSong("4", "Song 4", "Artist 4", 4),
		createBSTTestSong("5", "Song 5", "Artist 5", 5),
	}

	for _, song := range testSongs {
		bst.InsertSong(song, song.Rating)
	}

	tests := []struct {
		name        string
		minRating   int
		maxRating   int
		expectedLen int
	}{
		{"Range 2-4", 2, 4, 3},
		{"Range 1-5", 1, 5, 5},
		{"Range 3-3", 3, 3, 1},
		{"Range 6-7", 6, 7, 0},
		{"Invalid range", 4, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			songs := bst.GetSongsByRatingRange(tt.minRating, tt.maxRating)
			if len(songs) != tt.expectedLen {
				t.Errorf("GetSongsByRatingRange(%d, %d) length = %v, want %v",
					tt.minRating, tt.maxRating, len(songs), tt.expectedLen)
			}

			// Verify all songs are within range
			for _, song := range songs {
				if song.Rating < tt.minRating || song.Rating > tt.maxRating {
					t.Errorf("GetSongsByRatingRange(%d, %d) song rating %d out of range",
						tt.minRating, tt.maxRating, song.Rating)
				}
			}
		})
	}
}

func TestSongRatingBST_GetRatingStats(t *testing.T) {
	bst := NewSongRatingBST()

	// Test empty BST
	stats := bst.GetRatingStats()
	if len(stats) != 0 {
		t.Errorf("GetRatingStats() empty BST should return empty map")
	}

	// Add test songs
	testSongs := []*models.Song{
		createBSTTestSong("1", "Song 1", "Artist 1", 3),
		createBSTTestSong("2", "Song 2", "Artist 2", 3),
		createBSTTestSong("3", "Song 3", "Artist 3", 1),
		createBSTTestSong("4", "Song 4", "Artist 4", 5),
		createBSTTestSong("5", "Song 5", "Artist 5", 5),
		createBSTTestSong("6", "Song 6", "Artist 6", 5),
	}

	for _, song := range testSongs {
		bst.InsertSong(song, song.Rating)
	}

	stats = bst.GetRatingStats()

	expected := map[int]int{
		1: 1,
		3: 2,
		5: 3,
	}

	for rating, count := range expected {
		if stats[rating] != count {
			t.Errorf("GetRatingStats() rating %d count = %v, want %v", rating, stats[rating], count)
		}
	}

	// Verify no unexpected ratings
	if len(stats) != len(expected) {
		t.Errorf("GetRatingStats() unexpected ratings in stats")
	}
}

func TestSongRatingBST_IsEmpty(t *testing.T) {
	bst := NewSongRatingBST()

	if !bst.IsEmpty() {
		t.Errorf("IsEmpty() = %v, want %v", bst.IsEmpty(), true)
	}

	song := createBSTTestSong("1", "Song 1", "Artist 1", 3)
	bst.InsertSong(song, 3)

	if bst.IsEmpty() {
		t.Errorf("IsEmpty() after insert = %v, want %v", bst.IsEmpty(), false)
	}
}

func TestSongRatingBST_GetNodeCount(t *testing.T) {
	bst := NewSongRatingBST()

	if bst.GetNodeCount() != 0 {
		t.Errorf("GetNodeCount() = %v, want %v", bst.GetNodeCount(), 0)
	}

	// Add songs with different ratings
	ratings := []int{3, 1, 5, 1, 3} // Should create 3 nodes
	for i, rating := range ratings {
		song := createBSTTestSong(string(rune('1'+i)), "Song", "Artist", rating)
		bst.InsertSong(song, rating)
	}

	if bst.GetNodeCount() != 3 {
		t.Errorf("GetNodeCount() after inserts = %v, want %v", bst.GetNodeCount(), 3)
	}
}

func TestSongRatingBST_GetTotalSongs(t *testing.T) {
	bst := NewSongRatingBST()

	if bst.GetTotalSongs() != 0 {
		t.Errorf("GetTotalSongs() empty = %v, want %v", bst.GetTotalSongs(), 0)
	}

	// Add 5 songs
	for i := 0; i < 5; i++ {
		song := createBSTTestSong(string(rune('1'+i)), "Song", "Artist", 3)
		bst.InsertSong(song, 3)
	}

	if bst.GetTotalSongs() != 5 {
		t.Errorf("GetTotalSongs() = %v, want %v", bst.GetTotalSongs(), 5)
	}
}

func TestSongRatingBST_Clear(t *testing.T) {
	bst := NewSongRatingBST()

	// Add some songs
	for i := 0; i < 3; i++ {
		song := createBSTTestSong(string(rune('1'+i)), "Song", "Artist", i+1)
		bst.InsertSong(song, i+1)
	}

	bst.Clear()

	if !bst.IsEmpty() {
		t.Errorf("Clear() IsEmpty() = %v, want %v", bst.IsEmpty(), true)
	}
	if bst.GetNodeCount() != 0 {
		t.Errorf("Clear() NodeCount = %v, want %v", bst.GetNodeCount(), 0)
	}
	if bst.Root != nil {
		t.Errorf("Clear() Root should be nil")
	}
}

// Benchmark tests
func BenchmarkSongRatingBST_InsertSong(b *testing.B) {
	bst := NewSongRatingBST()
	song := createBSTTestSong("test", "Test Song", "Test Artist", 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.InsertSong(song, 3)
	}
}

func BenchmarkSongRatingBST_SearchByRating(b *testing.B) {
	bst := NewSongRatingBST()

	// Pre-populate BST
	for i := 0; i < 1000; i++ {
		rating := (i % 5) + 1
		song := createBSTTestSong(string(rune(i)), "Song", "Artist", rating)
		bst.InsertSong(song, rating)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.SearchByRating(3)
	}
}

func BenchmarkSongRatingBST_DeleteSong(b *testing.B) {
	// Pre-populate multiple BSTs for deletion
	bsts := make([]*SongRatingBST, b.N)
	for i := 0; i < b.N; i++ {
		bst := NewSongRatingBST()
		for j := 0; j < 100; j++ {
			song := createBSTTestSong(string(rune(j)), "Song", "Artist", 3)
			bst.InsertSong(song, 3)
		}
		bsts[i] = bst
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bsts[i].DeleteSong("50")
	}
}

func BenchmarkSongRatingBST_GetAllSongs(b *testing.B) {
	bst := NewSongRatingBST()

	// Pre-populate BST with different ratings
	for i := 0; i < 1000; i++ {
		rating := (i % 5) + 1
		song := createBSTTestSong(string(rune(i)), "Song", "Artist", rating)
		bst.InsertSong(song, rating)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.GetAllSongs()
	}
}

func BenchmarkSongRatingBST_GetSongsByRatingRange(b *testing.B) {
	bst := NewSongRatingBST()

	// Pre-populate BST
	for i := 0; i < 1000; i++ {
		rating := (i % 5) + 1
		song := createBSTTestSong(string(rune(i)), "Song", "Artist", rating)
		bst.InsertSong(song, rating)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bst.GetSongsByRatingRange(2, 4)
	}
}

// Edge case tests
func TestSongRatingBST_EdgeCases(t *testing.T) {
	bst := NewSongRatingBST()

	// Test inserting song with nil
	bst.InsertSong(nil, 3) // Should handle gracefully
	if !bst.IsEmpty() {
		t.Errorf("Insert nil song should not add to BST")
	}

	// Test complex tree operations
	// Build a more complex tree
	ratings := []int{3, 1, 5, 2, 4}
	for i, rating := range ratings {
		song := createBSTTestSong(string(rune('a'+i)), "Song", "Artist", rating)
		bst.InsertSong(song, rating)
	}

	// Test deleting root node with two children
	bst.DeleteSong("a") // Song with rating 3 (root)

	// Verify tree structure is maintained
	songs := bst.GetAllSongs()
	if len(songs) != 4 {
		t.Errorf("After deleting root, should have 4 songs")
	}

	// Verify inorder traversal still works
	for i := 1; i < len(songs); i++ {
		if songs[i-1].Rating > songs[i].Rating {
			t.Errorf("Tree structure compromised after root deletion")
		}
	}
}

func TestSongRatingBST_TreeBalance(t *testing.T) {
	bst := NewSongRatingBST()

	// Insert songs in order that could create unbalanced tree
	for rating := 1; rating <= 5; rating++ {
		song := createBSTTestSong(string(rune('0'+rating)), "Song", "Artist", rating)
		bst.InsertSong(song, rating)
	}

	// Verify all operations still work correctly despite potential imbalance
	for rating := 1; rating <= 5; rating++ {
		songs := bst.SearchByRating(rating)
		if len(songs) != 1 {
			t.Errorf("SearchByRating(%d) in potentially unbalanced tree failed", rating)
		}
	}

	// Verify range query works
	songs := bst.GetSongsByRatingRange(2, 4)
	if len(songs) != 3 {
		t.Errorf("Range query in potentially unbalanced tree failed")
	}
}
