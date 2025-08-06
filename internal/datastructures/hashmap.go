package datastructures

import (
	"fmt"
	"src/internal/models"
)

// HashMapEntry represents a key-value pair in the hash map
// Contains song ID as key and song metadata as value
// Time Complexity: O(1) for field access
// Space Complexity: O(1) per entry
type HashMapEntry struct {
	Key   string
	Value *models.Song
	Next  *HashMapEntry // For handling collisions using chaining
}

// SongHashMap represents a hash map for instant song lookup
// Provides O(1) average time complexity for insert, search, and delete operations
// Uses separate chaining for collision resolution
// Time Complexity: O(1) average, O(n) worst case for operations
// Space Complexity: O(n) where n is the number of songs
type SongHashMap struct {
	Buckets  []*HashMapEntry
	Size     int // Number of songs stored
	Capacity int // Number of buckets
}

// NewSongHashMap creates a new song hash map with specified initial capacity
// Time Complexity: O(capacity)
// Space Complexity: O(capacity)
func NewSongHashMap(capacity int) *SongHashMap {
	if capacity <= 0 {
		capacity = 16 // Default capacity
	}
	return &SongHashMap{
		Buckets:  make([]*HashMapEntry, capacity),
		Size:     0,
		Capacity: capacity,
	}
}

// hash function using djb2 algorithm for string hashing
// Time Complexity: O(k) where k is the length of the key
// Space Complexity: O(1)
func (shm *SongHashMap) hash(key string) int {
	hash := 5381
	for _, c := range key {
		hash = ((hash << 5) + hash) + int(c) // hash * 33 + c
	}
	if hash < 0 {
		hash = -hash
	}
	return hash % shm.Capacity
}

// Put inserts or updates a song in the hash map
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) Put(song *models.Song) {
	if song == nil || song.ID == "" {
		return
	}

	index := shm.hash(song.ID)
	entry := shm.Buckets[index]

	// If bucket is empty, create new entry
	if entry == nil {
		shm.Buckets[index] = &HashMapEntry{
			Key:   song.ID,
			Value: song,
			Next:  nil,
		}
		shm.Size++
		return
	}

	// Traverse the chain to find if key exists or add to end
	for {
		if entry.Key == song.ID {
			// Update existing entry
			entry.Value = song
			return
		}
		if entry.Next == nil {
			break
		}
		entry = entry.Next
	}

	// Add new entry at the end of the chain
	entry.Next = &HashMapEntry{
		Key:   song.ID,
		Value: song,
		Next:  nil,
	}
	shm.Size++

	// Check if we need to resize
	if shm.Size > shm.Capacity*2 {
		shm.resize()
	}
}

// PutByTitle adds a song indexed by title for title-based lookup
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) PutByTitle(song *models.Song) {
	if song == nil || song.Title == "" {
		return
	}

	// Use title as key instead of ID
	index := shm.hash(song.Title)
	entry := shm.Buckets[index]

	if entry == nil {
		shm.Buckets[index] = &HashMapEntry{
			Key:   song.Title,
			Value: song,
			Next:  nil,
		}
		shm.Size++
		return
	}

	for {
		if entry.Key == song.Title {
			entry.Value = song
			return
		}
		if entry.Next == nil {
			break
		}
		entry = entry.Next
	}

	entry.Next = &HashMapEntry{
		Key:   song.Title,
		Value: song,
		Next:  nil,
	}
	shm.Size++

	if shm.Size > shm.Capacity*2 {
		shm.resize()
	}
}

// Get retrieves a song by ID
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) Get(songID string) (*models.Song, error) {
	if songID == "" {
		return nil, fmt.Errorf("song ID cannot be empty")
	}

	index := shm.hash(songID)
	entry := shm.Buckets[index]

	// Traverse the chain to find the key
	for entry != nil {
		if entry.Key == songID {
			return entry.Value, nil
		}
		entry = entry.Next
	}

	return nil, fmt.Errorf("song with ID '%s' not found", songID)
}

// GetByTitle retrieves a song by title
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) GetByTitle(title string) (*models.Song, error) {
	if title == "" {
		return nil, fmt.Errorf("song title cannot be empty")
	}

	index := shm.hash(title)
	entry := shm.Buckets[index]

	for entry != nil {
		if entry.Key == title {
			return entry.Value, nil
		}
		entry = entry.Next
	}

	return nil, fmt.Errorf("song with title '%s' not found", title)
}

// Delete removes a song from the hash map by ID
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) Delete(songID string) (*models.Song, error) {
	if songID == "" {
		return nil, fmt.Errorf("song ID cannot be empty")
	}

	index := shm.hash(songID)
	entry := shm.Buckets[index]

	if entry == nil {
		return nil, fmt.Errorf("song with ID '%s' not found", songID)
	}

	// If first entry matches
	if entry.Key == songID {
		song := entry.Value
		shm.Buckets[index] = entry.Next
		shm.Size--
		return song, nil
	}

	// Search in the chain
	for entry.Next != nil {
		if entry.Next.Key == songID {
			song := entry.Next.Value
			entry.Next = entry.Next.Next
			shm.Size--
			return song, nil
		}
		entry = entry.Next
	}

	return nil, fmt.Errorf("song with ID '%s' not found", songID)
}

// Contains checks if a song exists in the hash map
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) Contains(songID string) bool {
	_, err := shm.Get(songID)
	return err == nil
}

// ContainsByTitle checks if a song exists by title
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) ContainsByTitle(title string) bool {
	_, err := shm.GetByTitle(title)
	return err == nil
}

// GetAllSongs returns all songs in the hash map
// Time Complexity: O(n + capacity)
// Space Complexity: O(n)
func (shm *SongHashMap) GetAllSongs() []*models.Song {
	songs := make([]*models.Song, 0, shm.Size)

	for i := 0; i < shm.Capacity; i++ {
		entry := shm.Buckets[i]
		for entry != nil {
			songs = append(songs, entry.Value)
			entry = entry.Next
		}
	}

	return songs
}

// GetAllKeys returns all keys (song IDs) in the hash map
// Time Complexity: O(n + capacity)
// Space Complexity: O(n)
func (shm *SongHashMap) GetAllKeys() []string {
	keys := make([]string, 0, shm.Size)

	for i := 0; i < shm.Capacity; i++ {
		entry := shm.Buckets[i]
		for entry != nil {
			keys = append(keys, entry.Key)
			entry = entry.Next
		}
	}

	return keys
}

// GetMetadata returns metadata for a song by ID
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) GetMetadata(songID string) (map[string]interface{}, error) {
	song, err := shm.Get(songID)
	if err != nil {
		return nil, err
	}
	return song.GetMetadata(), nil
}

// UpdateSong updates an existing song's information
// Time Complexity: O(1) average, O(n) worst case
// Space Complexity: O(1)
func (shm *SongHashMap) UpdateSong(song *models.Song) error {
	if song == nil || song.ID == "" {
		return fmt.Errorf("invalid song or empty ID")
	}

	index := shm.hash(song.ID)
	entry := shm.Buckets[index]

	for entry != nil {
		if entry.Key == song.ID {
			entry.Value = song
			return nil
		}
		entry = entry.Next
	}

	return fmt.Errorf("song with ID '%s' not found", song.ID)
}

// GetSize returns the number of songs in the hash map
// Time Complexity: O(1)
// Space Complexity: O(1)
func (shm *SongHashMap) GetSize() int {
	return shm.Size
}

// GetCapacity returns the current capacity of the hash map
// Time Complexity: O(1)
// Space Complexity: O(1)
func (shm *SongHashMap) GetCapacity() int {
	return shm.Capacity
}

// GetLoadFactor returns the load factor of the hash map
// Time Complexity: O(1)
// Space Complexity: O(1)
func (shm *SongHashMap) GetLoadFactor() float64 {
	return float64(shm.Size) / float64(shm.Capacity)
}

// IsEmpty checks if the hash map is empty
// Time Complexity: O(1)
// Space Complexity: O(1)
func (shm *SongHashMap) IsEmpty() bool {
	return shm.Size == 0
}

// Clear removes all songs from the hash map
// Time Complexity: O(capacity)
// Space Complexity: O(1)
func (shm *SongHashMap) Clear() {
	for i := 0; i < shm.Capacity; i++ {
		shm.Buckets[i] = nil
	}
	shm.Size = 0
}

// resize increases the capacity of the hash map and rehashes all entries
// Time Complexity: O(n)
// Space Complexity: O(new_capacity)
func (shm *SongHashMap) resize() {
	oldBuckets := shm.Buckets
	oldCapacity := shm.Capacity

	// Double the capacity
	shm.Capacity *= 2
	shm.Buckets = make([]*HashMapEntry, shm.Capacity)
	shm.Size = 0

	// Rehash all entries
	for i := 0; i < oldCapacity; i++ {
		entry := oldBuckets[i]
		for entry != nil {
			shm.Put(entry.Value)
			entry = entry.Next
		}
	}
}

// GetBucketDistribution returns the distribution of entries across buckets
// Useful for analyzing hash function performance
// Time Complexity: O(capacity)
// Space Complexity: O(1)
func (shm *SongHashMap) GetBucketDistribution() map[string]interface{} {
	distribution := make([]int, shm.Capacity)
	maxChainLength := 0
	emptyBuckets := 0

	for i := 0; i < shm.Capacity; i++ {
		chainLength := 0
		entry := shm.Buckets[i]

		if entry == nil {
			emptyBuckets++
		} else {
			for entry != nil {
				chainLength++
				entry = entry.Next
			}
			if chainLength > maxChainLength {
				maxChainLength = chainLength
			}
		}
		distribution[i] = chainLength
	}

	return map[string]interface{}{
		"capacity":         shm.Capacity,
		"size":             shm.Size,
		"load_factor":      shm.GetLoadFactor(),
		"empty_buckets":    emptyBuckets,
		"max_chain_length": maxChainLength,
		"distribution":     distribution,
	}
}

// String returns a string representation of the hash map
// Time Complexity: O(n + capacity)
// Space Complexity: O(n)
func (shm *SongHashMap) String() string {
	if shm.IsEmpty() {
		return "Empty Song HashMap"
	}

	result := fmt.Sprintf("Song HashMap (Size: %d, Capacity: %d, Load Factor: %.2f):\n",
		shm.Size, shm.Capacity, shm.GetLoadFactor())

	songs := shm.GetAllSongs()
	for i, song := range songs {
		result += fmt.Sprintf("%d. %s - %s by %s\n", i+1, song.ID, song.Title, song.Artist)
	}

	return result
}
