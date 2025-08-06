package datastructures

import (
	"src/internal/models"
	"testing"
)

// Test helper function to create a test song for HashMap tests
func createHashMapTestSong(id, title, artist string) *models.Song {
	return models.NewSong(id, title, artist, "Test Album", "Rock", "Alternative", "Energetic", 180, 120)
}

func TestNewSongHashMap(t *testing.T) {
	tests := []struct {
		name             string
		capacity         int
		expectedCapacity int
	}{
		{"Valid capacity", 32, 32},
		{"Zero capacity", 0, 16},      // Should default to 16
		{"Negative capacity", -5, 16}, // Should default to 16
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashMap := NewSongHashMap(tt.capacity)

			if hashMap.Capacity != tt.expectedCapacity {
				t.Errorf("NewSongHashMap() Capacity = %v, want %v", hashMap.Capacity, tt.expectedCapacity)
			}
			if hashMap.Size != 0 {
				t.Errorf("NewSongHashMap() Size = %v, want %v", hashMap.Size, 0)
			}
			if len(hashMap.Buckets) != tt.expectedCapacity {
				t.Errorf("NewSongHashMap() Buckets length = %v, want %v", len(hashMap.Buckets), tt.expectedCapacity)
			}
		})
	}
}

func TestSongHashMap_hash(t *testing.T) {
	hashMap := NewSongHashMap(16)

	tests := []struct {
		name string
		key  string
	}{
		{"Simple key", "test"},
		{"Empty key", ""},
		{"Long key", "this-is-a-very-long-song-id-for-testing"},
		{"Special characters", "test-song-123!@#"},
		{"Unicode key", "song-with-üñíçødé"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1 := hashMap.hash(tt.key)
			hash2 := hashMap.hash(tt.key)

			// Hash should be consistent
			if hash1 != hash2 {
				t.Errorf("hash(%s) inconsistent: %d != %d", tt.key, hash1, hash2)
			}

			// Hash should be within bounds
			if hash1 < 0 || hash1 >= hashMap.Capacity {
				t.Errorf("hash(%s) = %d, should be within [0, %d)", tt.key, hash1, hashMap.Capacity)
			}
		})
	}

	// Test that different keys produce different hashes (mostly)
	key1 := "song1"
	key2 := "song2"
	hash1 := hashMap.hash(key1)
	hash2 := hashMap.hash(key2)

	// While not guaranteed, different keys should usually produce different hashes
	// This is a probabilistic test
	if hash1 == hash2 {
		t.Logf("hash collision detected: %s and %s both hash to %d", key1, key2, hash1)
	}
}

func TestSongHashMap_Put(t *testing.T) {
	hashMap := NewSongHashMap(4) // Small capacity for testing
	song1 := createHashMapTestSong("song1", "Song 1", "Artist 1")
	song2 := createHashMapTestSong("song2", "Song 2", "Artist 2")

	// Test inserting first song
	hashMap.Put(song1)

	if hashMap.Size != 1 {
		t.Errorf("Put() Size = %v, want %v", hashMap.Size, 1)
	}

	retrievedSong, err := hashMap.Get("song1")
	if err != nil {
		t.Errorf("Put() song not retrievable: %v", err)
	}
	if retrievedSong.ID != "song1" {
		t.Errorf("Put() retrieved song ID = %v, want %v", retrievedSong.ID, "song1")
	}

	// Test inserting second song
	hashMap.Put(song2)

	if hashMap.Size != 2 {
		t.Errorf("Put() second song Size = %v, want %v", hashMap.Size, 2)
	}

	// Test updating existing song
	updatedSong := createHashMapTestSong("song1", "Updated Song", "Updated Artist")
	hashMap.Put(updatedSong)

	if hashMap.Size != 2 {
		t.Errorf("Put() update should not increase Size = %v, want %v", hashMap.Size, 2)
	}

	retrievedSong, err = hashMap.Get("song1")
	if err != nil {
		t.Errorf("Put() updated song not retrievable: %v", err)
	}
	if retrievedSong.Title != "Updated Song" {
		t.Errorf("Put() updated song title = %v, want %v", retrievedSong.Title, "Updated Song")
	}

	// Test with nil song
	hashMap.Put(nil) // Should handle gracefully
	if hashMap.Size != 2 {
		t.Errorf("Put(nil) should not change Size")
	}

	// Test with empty ID
	emptySong := createHashMapTestSong("", "Empty ID", "Artist")
	hashMap.Put(emptySong)
	if hashMap.Size != 2 {
		t.Errorf("Put() with empty ID should not change Size")
	}
}

func TestSongHashMap_Get(t *testing.T) {
	hashMap := NewSongHashMap(16)
	songs := []*models.Song{
		createHashMapTestSong("song1", "Song 1", "Artist 1"),
		createHashMapTestSong("song2", "Song 2", "Artist 2"),
		createHashMapTestSong("song3", "Song 3", "Artist 3"),
	}

	// Test getting from empty map
	_, err := hashMap.Get("song1")
	if err == nil {
		t.Errorf("Get() from empty map should return error")
	}

	// Add songs
	for _, song := range songs {
		hashMap.Put(song)
	}

	// Test getting existing songs
	for _, expectedSong := range songs {
		retrievedSong, err := hashMap.Get(expectedSong.ID)
		if err != nil {
			t.Errorf("Get(%s) error = %v, want nil", expectedSong.ID, err)
		}
		if retrievedSong.ID != expectedSong.ID {
			t.Errorf("Get(%s) ID = %v, want %v", expectedSong.ID, retrievedSong.ID, expectedSong.ID)
		}
		if retrievedSong.Title != expectedSong.Title {
			t.Errorf("Get(%s) Title = %v, want %v", expectedSong.ID, retrievedSong.Title, expectedSong.Title)
		}
	}

	// Test getting non-existing song
	_, err = hashMap.Get("nonexistent")
	if err == nil {
		t.Errorf("Get() non-existing song should return error")
	}

	// Test with empty ID
	_, err = hashMap.Get("")
	if err == nil {
		t.Errorf("Get() with empty ID should return error")
	}
}

func TestSongHashMap_PutByTitle(t *testing.T) {
	hashMap := NewSongHashMap(16)
	song1 := createHashMapTestSong("song1", "Unique Title", "Artist 1")
	song2 := createHashMapTestSong("song2", "Another Title", "Artist 2")

	// Test putting by title
	hashMap.PutByTitle(song1)
	hashMap.PutByTitle(song2)

	if hashMap.Size != 2 {
		t.Errorf("PutByTitle() Size = %v, want %v", hashMap.Size, 2)
	}

	// Test retrieving by title
	retrievedSong, err := hashMap.GetByTitle("Unique Title")
	if err != nil {
		t.Errorf("PutByTitle() song not retrievable by title: %v", err)
	}
	if retrievedSong.ID != "song1" {
		t.Errorf("PutByTitle() retrieved song ID = %v, want %v", retrievedSong.ID, "song1")
	}

	// Test with nil song
	hashMap.PutByTitle(nil)
	if hashMap.Size != 2 {
		t.Errorf("PutByTitle(nil) should not change Size")
	}

	// Test with empty title
	emptySong := createHashMapTestSong("song3", "", "Artist 3")
	hashMap.PutByTitle(emptySong)
	if hashMap.Size != 2 {
		t.Errorf("PutByTitle() with empty title should not change Size")
	}
}

func TestSongHashMap_GetByTitle(t *testing.T) {
	hashMap := NewSongHashMap(16)
	songs := []*models.Song{
		createHashMapTestSong("song1", "Title One", "Artist 1"),
		createHashMapTestSong("song2", "Title Two", "Artist 2"),
	}

	// Add songs by title
	for _, song := range songs {
		hashMap.PutByTitle(song)
	}

	// Test getting by title
	for _, expectedSong := range songs {
		retrievedSong, err := hashMap.GetByTitle(expectedSong.Title)
		if err != nil {
			t.Errorf("GetByTitle(%s) error = %v, want nil", expectedSong.Title, err)
		}
		if retrievedSong.Title != expectedSong.Title {
			t.Errorf("GetByTitle(%s) Title = %v, want %v", expectedSong.Title, retrievedSong.Title, expectedSong.Title)
		}
	}

	// Test non-existing title
	_, err := hashMap.GetByTitle("Non-existing Title")
	if err == nil {
		t.Errorf("GetByTitle() non-existing title should return error")
	}

	// Test empty title
	_, err = hashMap.GetByTitle("")
	if err == nil {
		t.Errorf("GetByTitle() empty title should return error")
	}
}

func TestSongHashMap_Delete(t *testing.T) {
	hashMap := NewSongHashMap(16)
	songs := []*models.Song{
		createHashMapTestSong("song1", "Song 1", "Artist 1"),
		createHashMapTestSong("song2", "Song 2", "Artist 2"),
		createHashMapTestSong("song3", "Song 3", "Artist 3"),
	}

	// Test deleting from empty map
	_, err := hashMap.Delete("song1")
	if err == nil {
		t.Errorf("Delete() from empty map should return error")
	}

	// Add songs
	for _, song := range songs {
		hashMap.Put(song)
	}

	// Test deleting existing song
	deletedSong, err := hashMap.Delete("song2")
	if err != nil {
		t.Errorf("Delete() error = %v, want nil", err)
	}
	if deletedSong.ID != "song2" {
		t.Errorf("Delete() returned song ID = %v, want %v", deletedSong.ID, "song2")
	}
	if hashMap.Size != 2 {
		t.Errorf("Delete() Size = %v, want %v", hashMap.Size, 2)
	}

	// Verify song is actually deleted
	_, err = hashMap.Get("song2")
	if err == nil {
		t.Errorf("Delete() song should not be retrievable after deletion")
	}

	// Test deleting non-existing song
	_, err = hashMap.Delete("nonexistent")
	if err == nil {
		t.Errorf("Delete() non-existing song should return error")
	}

	// Test with empty ID
	_, err = hashMap.Delete("")
	if err == nil {
		t.Errorf("Delete() with empty ID should return error")
	}
}

func TestSongHashMap_Contains(t *testing.T) {
	hashMap := NewSongHashMap(16)
	song := createHashMapTestSong("song1", "Song 1", "Artist 1")

	// Test contains on empty map
	if hashMap.Contains("song1") {
		t.Errorf("Contains() empty map should return false")
	}

	// Add song and test
	hashMap.Put(song)

	if !hashMap.Contains("song1") {
		t.Errorf("Contains() existing song should return true")
	}

	if hashMap.Contains("nonexistent") {
		t.Errorf("Contains() non-existing song should return false")
	}
}

func TestSongHashMap_ContainsByTitle(t *testing.T) {
	hashMap := NewSongHashMap(16)
	song := createHashMapTestSong("song1", "Unique Song Title", "Artist 1")

	// Test contains by title on empty map
	if hashMap.ContainsByTitle("Unique Song Title") {
		t.Errorf("ContainsByTitle() empty map should return false")
	}

	// Add song by title and test
	hashMap.PutByTitle(song)

	if !hashMap.ContainsByTitle("Unique Song Title") {
		t.Errorf("ContainsByTitle() existing title should return true")
	}

	if hashMap.ContainsByTitle("Non-existing Title") {
		t.Errorf("ContainsByTitle() non-existing title should return false")
	}
}

func TestSongHashMap_GetAllSongs(t *testing.T) {
	hashMap := NewSongHashMap(16)

	// Test empty map
	songs := hashMap.GetAllSongs()
	if len(songs) != 0 {
		t.Errorf("GetAllSongs() empty map length = %v, want %v", len(songs), 0)
	}

	// Add songs
	testSongs := []*models.Song{
		createHashMapTestSong("song1", "Song 1", "Artist 1"),
		createHashMapTestSong("song2", "Song 2", "Artist 2"),
		createHashMapTestSong("song3", "Song 3", "Artist 3"),
	}

	for _, song := range testSongs {
		hashMap.Put(song)
	}

	songs = hashMap.GetAllSongs()
	if len(songs) != len(testSongs) {
		t.Errorf("GetAllSongs() length = %v, want %v", len(songs), len(testSongs))
	}

	// Verify all songs are present (order not guaranteed in hash map)
	songIDs := make(map[string]bool)
	for _, song := range songs {
		songIDs[song.ID] = true
	}

	for _, expectedSong := range testSongs {
		if !songIDs[expectedSong.ID] {
			t.Errorf("GetAllSongs() missing song ID: %s", expectedSong.ID)
		}
	}
}

func TestSongHashMap_GetAllKeys(t *testing.T) {
	hashMap := NewSongHashMap(16)

	// Test empty map
	keys := hashMap.GetAllKeys()
	if len(keys) != 0 {
		t.Errorf("GetAllKeys() empty map length = %v, want %v", len(keys), 0)
	}

	// Add songs
	expectedKeys := []string{"song1", "song2", "song3"}
	for _, key := range expectedKeys {
		song := createHashMapTestSong(key, "Song", "Artist")
		hashMap.Put(song)
	}

	keys = hashMap.GetAllKeys()
	if len(keys) != len(expectedKeys) {
		t.Errorf("GetAllKeys() length = %v, want %v", len(keys), len(expectedKeys))
	}

	// Verify all keys are present
	keySet := make(map[string]bool)
	for _, key := range keys {
		keySet[key] = true
	}

	for _, expectedKey := range expectedKeys {
		if !keySet[expectedKey] {
			t.Errorf("GetAllKeys() missing key: %s", expectedKey)
		}
	}
}

func TestSongHashMap_UpdateSong(t *testing.T) {
	hashMap := NewSongHashMap(16)
	originalSong := createHashMapTestSong("song1", "Original Title", "Original Artist")
	updatedSong := createHashMapTestSong("song1", "Updated Title", "Updated Artist")

	// Test updating non-existing song
	err := hashMap.UpdateSong(updatedSong)
	if err == nil {
		t.Errorf("UpdateSong() non-existing song should return error")
	}

	// Add original song
	hashMap.Put(originalSong)

	// Update song
	err = hashMap.UpdateSong(updatedSong)
	if err != nil {
		t.Errorf("UpdateSong() error = %v, want nil", err)
	}

	// Verify update
	retrievedSong, err := hashMap.Get("song1")
	if err != nil {
		t.Errorf("UpdateSong() updated song not retrievable: %v", err)
	}
	if retrievedSong.Title != "Updated Title" {
		t.Errorf("UpdateSong() title = %v, want %v", retrievedSong.Title, "Updated Title")
	}

	// Test with nil song
	err = hashMap.UpdateSong(nil)
	if err == nil {
		t.Errorf("UpdateSong(nil) should return error")
	}

	// Test with empty ID
	emptySong := createHashMapTestSong("", "Empty", "Artist")
	err = hashMap.UpdateSong(emptySong)
	if err == nil {
		t.Errorf("UpdateSong() with empty ID should return error")
	}
}

func TestSongHashMap_GetSize(t *testing.T) {
	hashMap := NewSongHashMap(16)

	if hashMap.GetSize() != 0 {
		t.Errorf("GetSize() empty map = %v, want %v", hashMap.GetSize(), 0)
	}

	// Add songs
	for i := 0; i < 5; i++ {
		song := createHashMapTestSong(string(rune('1'+i)), "Song", "Artist")
		hashMap.Put(song)
	}

	if hashMap.GetSize() != 5 {
		t.Errorf("GetSize() after 5 inserts = %v, want %v", hashMap.GetSize(), 5)
	}
}

func TestSongHashMap_GetCapacity(t *testing.T) {
	hashMap := NewSongHashMap(32)

	if hashMap.GetCapacity() != 32 {
		t.Errorf("GetCapacity() = %v, want %v", hashMap.GetCapacity(), 32)
	}
}

func TestSongHashMap_GetLoadFactor(t *testing.T) {
	hashMap := NewSongHashMap(10)

	// Empty map
	if hashMap.GetLoadFactor() != 0.0 {
		t.Errorf("GetLoadFactor() empty map = %v, want %v", hashMap.GetLoadFactor(), 0.0)
	}

	// Add 3 songs
	for i := 0; i < 3; i++ {
		song := createHashMapTestSong(string(rune('1'+i)), "Song", "Artist")
		hashMap.Put(song)
	}

	expected := 3.0 / 10.0
	if hashMap.GetLoadFactor() != expected {
		t.Errorf("GetLoadFactor() = %v, want %v", hashMap.GetLoadFactor(), expected)
	}
}

func TestSongHashMap_IsEmpty(t *testing.T) {
	hashMap := NewSongHashMap(16)

	if !hashMap.IsEmpty() {
		t.Errorf("IsEmpty() empty map = %v, want %v", hashMap.IsEmpty(), true)
	}

	song := createHashMapTestSong("song1", "Song 1", "Artist 1")
	hashMap.Put(song)

	if hashMap.IsEmpty() {
		t.Errorf("IsEmpty() after insert = %v, want %v", hashMap.IsEmpty(), false)
	}
}

func TestSongHashMap_Clear(t *testing.T) {
	hashMap := NewSongHashMap(16)

	// Add songs
	for i := 0; i < 3; i++ {
		song := createHashMapTestSong(string(rune('1'+i)), "Song", "Artist")
		hashMap.Put(song)
	}

	hashMap.Clear()

	if !hashMap.IsEmpty() {
		t.Errorf("Clear() IsEmpty() = %v, want %v", hashMap.IsEmpty(), true)
	}
	if hashMap.GetSize() != 0 {
		t.Errorf("Clear() Size = %v, want %v", hashMap.GetSize(), 0)
	}
}

func TestSongHashMap_resize(t *testing.T) {
	hashMap := NewSongHashMap(4) // Small capacity to trigger resize

	// Add enough songs to trigger resize (threshold is 2 * capacity = 8)
	for i := 0; i < 9; i++ {
		song := createHashMapTestSong(string(rune('1'+i)), "Song", "Artist")
		hashMap.Put(song)
	}

	// Capacity should have doubled
	if hashMap.GetCapacity() != 8 {
		t.Errorf("resize() Capacity = %v, want %v", hashMap.GetCapacity(), 8)
	}

	// All songs should still be retrievable
	for i := 0; i < 9; i++ {
		songID := string(rune('1' + i))
		_, err := hashMap.Get(songID)
		if err != nil {
			t.Errorf("resize() song %s not retrievable after resize: %v", songID, err)
		}
	}
}

func TestSongHashMap_GetBucketDistribution(t *testing.T) {
	hashMap := NewSongHashMap(8)

	// Test empty map
	distribution := hashMap.GetBucketDistribution()
	if distribution["capacity"] != 8 {
		t.Errorf("GetBucketDistribution() capacity = %v, want %v", distribution["capacity"], 8)
	}
	if distribution["size"] != 0 {
		t.Errorf("GetBucketDistribution() size = %v, want %v", distribution["size"], 0)
	}
	if distribution["empty_buckets"] != 8 {
		t.Errorf("GetBucketDistribution() empty_buckets = %v, want %v", distribution["empty_buckets"], 8)
	}

	// Add songs
	for i := 0; i < 5; i++ {
		song := createHashMapTestSong(string(rune('1'+i)), "Song", "Artist")
		hashMap.Put(song)
	}

	distribution = hashMap.GetBucketDistribution()
	if distribution["size"] != 5 {
		t.Errorf("GetBucketDistribution() size after inserts = %v, want %v", distribution["size"], 5)
	}

	// Verify distribution array length
	distArray, ok := distribution["distribution"].([]int)
	if !ok {
		t.Errorf("GetBucketDistribution() distribution should be []int")
	}
	if len(distArray) != 8 {
		t.Errorf("GetBucketDistribution() distribution length = %v, want %v", len(distArray), 8)
	}
}

func TestSongHashMap_GetMetadata(t *testing.T) {
	hashMap := NewSongHashMap(16)
	song := createHashMapTestSong("song1", "Test Song", "Test Artist")
	song.SetRating(4)

	hashMap.Put(song)

	metadata, err := hashMap.GetMetadata("song1")
	if err != nil {
		t.Errorf("GetMetadata() error = %v, want nil", err)
	}

	if metadata["id"] != "song1" {
		t.Errorf("GetMetadata() id = %v, want %v", metadata["id"], "song1")
	}
	if metadata["title"] != "Test Song" {
		t.Errorf("GetMetadata() title = %v, want %v", metadata["title"], "Test Song")
	}
	if metadata["rating"] != 4 {
		t.Errorf("GetMetadata() rating = %v, want %v", metadata["rating"], 4)
	}

	// Test non-existing song
	_, err = hashMap.GetMetadata("nonexistent")
	if err == nil {
		t.Errorf("GetMetadata() non-existing song should return error")
	}
}

// Benchmark tests
func BenchmarkSongHashMap_Put(b *testing.B) {
	hashMap := NewSongHashMap(1000)
	_ = createHashMapTestSong("test", "Test Song", "Test Artist")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testSong := createHashMapTestSong(string(rune(i)), "Song", "Artist")
		hashMap.Put(testSong)
	}
}

func BenchmarkSongHashMap_Get(b *testing.B) {
	hashMap := NewSongHashMap(10000)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		song := createHashMapTestSong(string(rune(i)), "Song", "Artist")
		hashMap.Put(song)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashMap.Get(string(rune(i % 1000)))
	}
}

func BenchmarkSongHashMap_Delete(b *testing.B) {
	// Pre-populate multiple maps for deletion
	maps := make([]*SongHashMap, b.N)
	for i := 0; i < b.N; i++ {
		hashMap := NewSongHashMap(1000)
		for j := 0; j < 100; j++ {
			song := createHashMapTestSong(string(rune(j)), "Song", "Artist")
			hashMap.Put(song)
		}
		maps[i] = hashMap
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maps[i].Delete(string(rune(50)))
	}
}

func BenchmarkSongHashMap_hash(b *testing.B) {
	hashMap := NewSongHashMap(1000)
	key := "test-song-id-for-benchmark"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashMap.hash(key)
	}
}

// Collision and stress tests
func TestSongHashMap_CollisionHandling(b *testing.T) {
	if testing.Short() {
		b.Skip("Skipping collision test in short mode")
	}

	// Use small capacity to increase collision probability
	hashMap := NewSongHashMap(4)

	// Add many songs to test collision handling
	numSongs := 20
	for i := 0; i < numSongs; i++ {
		song := createHashMapTestSong(string(rune('a'+i)), "Song", "Artist")
		hashMap.Put(song)
	}

	// All songs should be retrievable despite collisions
	for i := 0; i < numSongs; i++ {
		songID := string(rune('a' + i))
		_, err := hashMap.Get(songID)
		if err != nil {
			b.Errorf("CollisionHandling: song %s not retrievable: %v", songID, err)
		}
	}

	// Test bucket distribution
	distribution := hashMap.GetBucketDistribution()
	maxChainLength := distribution["max_chain_length"].(int)
	if maxChainLength <= 1 {
		b.Logf("CollisionHandling: max_chain_length = %d (may indicate poor collision distribution)", maxChainLength)
	}
}

func TestSongHashMap_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	hashMap := NewSongHashMap(100)

	// Add many songs
	numSongs := 10000
	for i := 0; i < numSongs; i++ {
		song := createHashMapTestSong(string(rune(i)), "Song", "Artist")
		hashMap.Put(song)
	}

	// Verify all songs are retrievable
	for i := 0; i < numSongs; i++ {
		songID := string(rune(i))
		_, err := hashMap.Get(songID)
		if err != nil {
			t.Errorf("StressTest: song %s not retrievable: %v", songID, err)
		}
	}

	// Test random deletions
	for i := 0; i < numSongs/2; i += 2 {
		songID := string(rune(i))
		_, err := hashMap.Delete(songID)
		if err != nil {
			t.Errorf("StressTest: failed to delete song %s: %v", songID, err)
		}
	}

	// Verify deleted songs are gone and remaining songs are still there
	for i := 0; i < numSongs; i++ {
		songID := string(rune(i))
		_, err := hashMap.Get(songID)

		if i%2 == 0 {
			// Should be deleted
			if err == nil {
				t.Errorf("StressTest: song %s should have been deleted", songID)
			}
		} else {
			// Should still exist
			if err != nil {
				t.Errorf("StressTest: song %s should still exist: %v", songID, err)
			}
		}
	}
}

func TestSongHashMap_EdgeCases(t *testing.T) {
	hashMap := NewSongHashMap(16)

	// Test operations with special characters in IDs
	specialSong := createHashMapTestSong("song-with-special!@#$%^&*()chars", "Special Song", "Special Artist")
	hashMap.Put(specialSong)

	retrievedSong, err := hashMap.Get("song-with-special!@#$%^&*()chars")
	if err != nil {
		t.Errorf("EdgeCases: special characters in ID failed: %v", err)
	}
	if retrievedSong.Title != "Special Song" {
		t.Errorf("EdgeCases: special characters song not retrieved correctly")
	}

	// Test very long IDs
	longID := string(make([]byte, 1000))
	for i := range longID {
		longID = string(append([]byte(longID[:i]), byte('a'+i%26)))
	}
	longSong := createHashMapTestSong(longID, "Long ID Song", "Artist")
	hashMap.Put(longSong)

	_, err = hashMap.Get(longID)
	if err != nil {
		t.Errorf("EdgeCases: very long ID failed: %v", err)
	}

	// Test unicode characters
	unicodeSong := createHashMapTestSong("söng-wíth-üñíçødé", "Unicode Song", "Unicode Artist")
	hashMap.Put(unicodeSong)

	_, err = hashMap.Get("söng-wíth-üñíçødé")
	if err != nil {
		t.Errorf("EdgeCases: unicode characters failed: %v", err)
	}
}
