package datastructures

import (
	"fmt"
	"src/internal/models"
)

// PlaylistNode represents a node in the doubly linked list
// Each node contains a song and pointers to next and previous nodes
// Time Complexity: O(1) for all field operations
// Space Complexity: O(1) per node
type PlaylistNode struct {
	Song *models.Song
	Next *PlaylistNode
	Prev *PlaylistNode
}

// DoublyLinkedList represents a playlist using doubly linked list
// Supports efficient insertion, deletion, and traversal operations
// Time Complexity: O(1) for head/tail operations, O(n) for index-based operations
// Space Complexity: O(n) where n is the number of songs
type DoublyLinkedList struct {
	Head   *PlaylistNode
	Tail   *PlaylistNode
	Length int
}

// NewDoublyLinkedList creates a new empty playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{
		Head:   nil,
		Tail:   nil,
		Length: 0,
	}
}

// AddSong adds a song to the end of the playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) AddSong(song *models.Song) {
	newNode := &PlaylistNode{
		Song: song,
		Next: nil,
		Prev: dll.Tail,
	}

	if dll.Head == nil {
		// Empty list
		dll.Head = newNode
		dll.Tail = newNode
	} else {
		// Add to tail
		dll.Tail.Next = newNode
		dll.Tail = newNode
	}

	dll.Length++
}

// AddSongAtIndex adds a song at a specific index
// Time Complexity: O(n) where n is the index
// Space Complexity: O(1)
func (dll *DoublyLinkedList) AddSongAtIndex(song *models.Song, index int) error {
	if index < 0 || index > dll.Length {
		return fmt.Errorf("index out of bounds: %d", index)
	}

	if index == dll.Length {
		dll.AddSong(song)
		return nil
	}

	if index == 0 {
		dll.AddSongToBeginning(song)
		return nil
	}

	newNode := &PlaylistNode{
		Song: song,
	}

	current := dll.getNodeAtIndex(index)
	newNode.Next = current
	newNode.Prev = current.Prev
	current.Prev.Next = newNode
	current.Prev = newNode

	dll.Length++
	return nil
}

// AddSongToBeginning adds a song to the beginning of the playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) AddSongToBeginning(song *models.Song) {
	newNode := &PlaylistNode{
		Song: song,
		Next: dll.Head,
		Prev: nil,
	}

	if dll.Head == nil {
		// Empty list
		dll.Head = newNode
		dll.Tail = newNode
	} else {
		dll.Head.Prev = newNode
		dll.Head = newNode
	}

	dll.Length++
}

// DeleteSong removes a song at the specified index
// Time Complexity: O(n) where n is the index
// Space Complexity: O(1)
func (dll *DoublyLinkedList) DeleteSong(index int) (*models.Song, error) {
	if index < 0 || index >= dll.Length {
		return nil, fmt.Errorf("index out of bounds: %d", index)
	}

	if dll.Length == 1 {
		// Only one node
		song := dll.Head.Song
		dll.Head = nil
		dll.Tail = nil
		dll.Length = 0
		return song, nil
	}

	nodeToDelete := dll.getNodeAtIndex(index)
	song := nodeToDelete.Song

	if nodeToDelete == dll.Head {
		// Delete head
		dll.Head = dll.Head.Next
		dll.Head.Prev = nil
	} else if nodeToDelete == dll.Tail {
		// Delete tail
		dll.Tail = dll.Tail.Prev
		dll.Tail.Next = nil
	} else {
		// Delete middle node
		nodeToDelete.Prev.Next = nodeToDelete.Next
		nodeToDelete.Next.Prev = nodeToDelete.Prev
	}

	dll.Length--
	return song, nil
}

// MoveSong moves a song from one index to another
// Time Complexity: O(n) where n is max(fromIndex, toIndex)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) MoveSong(fromIndex, toIndex int) error {
	if fromIndex < 0 || fromIndex >= dll.Length || toIndex < 0 || toIndex >= dll.Length {
		return fmt.Errorf("index out of bounds")
	}

	if fromIndex == toIndex {
		return nil
	}

	// Remove the song from the original position
	song, err := dll.DeleteSong(fromIndex)
	if err != nil {
		return err
	}

	// Adjust toIndex if necessary (when moving forward, index shifts after deletion)
	if toIndex > fromIndex {
		toIndex--
	}

	// Insert at new position
	return dll.AddSongAtIndex(song, toIndex)
}

// ReversePlaylist reverses the entire playlist
// Time Complexity: O(n)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) ReversePlaylist() {
	if dll.Head == nil || dll.Head == dll.Tail {
		return
	}

	current := dll.Head
	var temp *PlaylistNode

	// Swap next and prev for all nodes
	for current != nil {
		temp = current.Prev
		current.Prev = current.Next
		current.Next = temp
		current = current.Prev // Move to next node (which is now prev)
	}

	// Swap head and tail
	dll.Head, dll.Tail = dll.Tail, dll.Head
}

// GetSong returns the song at the specified index
// Time Complexity: O(n) where n is the index
// Space Complexity: O(1)
func (dll *DoublyLinkedList) GetSong(index int) (*models.Song, error) {
	if index < 0 || index >= dll.Length {
		return nil, fmt.Errorf("index out of bounds: %d", index)
	}

	node := dll.getNodeAtIndex(index)
	return node.Song, nil
}

// getNodeAtIndex is a helper method to get node at specific index
// Time Complexity: O(n) where n is the index
// Space Complexity: O(1)
func (dll *DoublyLinkedList) getNodeAtIndex(index int) *PlaylistNode {
	var current *PlaylistNode

	// Optimize by starting from head or tail based on index
	if index < dll.Length/2 {
		// Start from head
		current = dll.Head
		for i := 0; i < index; i++ {
			current = current.Next
		}
	} else {
		// Start from tail
		current = dll.Tail
		for i := dll.Length - 1; i > index; i-- {
			current = current.Prev
		}
	}

	return current
}

// ToSlice returns all songs as a slice for easy iteration
// Time Complexity: O(n)
// Space Complexity: O(n)
func (dll *DoublyLinkedList) ToSlice() []*models.Song {
	songs := make([]*models.Song, 0, dll.Length)
	current := dll.Head

	for current != nil {
		songs = append(songs, current.Song)
		current = current.Next
	}

	return songs
}

// Size returns the number of songs in the playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) Size() int {
	return dll.Length
}

// IsEmpty checks if the playlist is empty
// Time Complexity: O(1)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) IsEmpty() bool {
	return dll.Length == 0
}

// Clear removes all songs from the playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) Clear() {
	dll.Head = nil
	dll.Tail = nil
	dll.Length = 0
}

// GetTotalDuration calculates total duration of all songs in playlist
// Time Complexity: O(n)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) GetTotalDuration() int {
	totalDuration := 0
	current := dll.Head

	for current != nil {
		totalDuration += current.Song.Duration
		current = current.Next
	}

	return totalDuration
}

// FindSongByID searches for a song by ID and returns its index
// Time Complexity: O(n)
// Space Complexity: O(1)
func (dll *DoublyLinkedList) FindSongByID(songID string) (int, error) {
	current := dll.Head
	index := 0

	for current != nil {
		if current.Song.ID == songID {
			return index, nil
		}
		current = current.Next
		index++
	}

	return -1, fmt.Errorf("song with ID %s not found", songID)
}

// String returns a string representation of the playlist
// Time Complexity: O(n)
// Space Complexity: O(n)
func (dll *DoublyLinkedList) String() string {
	if dll.IsEmpty() {
		return "Empty Playlist"
	}

	result := "Playlist:\n"
	current := dll.Head
	index := 0

	for current != nil {
		result += fmt.Sprintf("%d. %s - %s (%s)\n",
			index+1, current.Song.Title, current.Song.Artist, current.Song.DurationString())
		current = current.Next
		index++
	}

	return result
}
