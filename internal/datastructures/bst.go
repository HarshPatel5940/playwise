package datastructures

import (
	"fmt"
	"src/internal/models"
)

// RatingBucket holds multiple songs with the same rating
// Time Complexity: O(1) for append, O(n) for remove operations
// Space Complexity: O(k) where k is the number of songs with same rating
type RatingBucket struct {
	Rating int
	Songs  []*models.Song
}

// NewRatingBucket creates a new rating bucket
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewRatingBucket(rating int) *RatingBucket {
	return &RatingBucket{
		Rating: rating,
		Songs:  make([]*models.Song, 0),
	}
}

// AddSong adds a song to the rating bucket
// Time Complexity: O(1)
// Space Complexity: O(1)
func (rb *RatingBucket) AddSong(song *models.Song) {
	rb.Songs = append(rb.Songs, song)
}

// RemoveSong removes a song from the rating bucket by ID
// Time Complexity: O(n) where n is the number of songs in the bucket
// Space Complexity: O(1)
func (rb *RatingBucket) RemoveSong(songID string) bool {
	for i, song := range rb.Songs {
		if song.ID == songID {
			// Remove song by swapping with last and reducing slice
			rb.Songs[i] = rb.Songs[len(rb.Songs)-1]
			rb.Songs = rb.Songs[:len(rb.Songs)-1]
			return true
		}
	}
	return false
}

// IsEmpty checks if the rating bucket has no songs
// Time Complexity: O(1)
// Space Complexity: O(1)
func (rb *RatingBucket) IsEmpty() bool {
	return len(rb.Songs) == 0
}

// GetSongs returns all songs in the bucket
// Time Complexity: O(1)
// Space Complexity: O(1)
func (rb *RatingBucket) GetSongs() []*models.Song {
	return rb.Songs
}

// BSTNode represents a node in the Binary Search Tree
// Each node contains a rating bucket and left/right children
// Time Complexity: O(1) for field access
// Space Complexity: O(1) per node
type BSTNode struct {
	Bucket *RatingBucket
	Left   *BSTNode
	Right  *BSTNode
}

// SongRatingBST represents a Binary Search Tree for song ratings
// Organizes songs by rating (1-5 stars) for fast lookup and manipulation
// Time Complexity: O(log n) average, O(n) worst case for search/insert/delete
// Space Complexity: O(n) where n is the number of unique ratings
type SongRatingBST struct {
	Root      *BSTNode
	NodeCount int
}

// NewSongRatingBST creates a new song rating BST
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewSongRatingBST() *SongRatingBST {
	return &SongRatingBST{
		Root:      nil,
		NodeCount: 0,
	}
}

// InsertSong inserts a song with its rating into the BST
// Time Complexity: O(log n) average, O(n) worst case
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) InsertSong(song *models.Song, rating int) {
	if rating < 1 || rating > 5 {
		return // Invalid rating
	}

	song.SetRating(rating) // Update song's rating
	bst.Root = bst.insertNode(bst.Root, song, rating)
}

// insertNode is a recursive helper for inserting nodes
// Time Complexity: O(log n) average, O(n) worst case
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) insertNode(node *BSTNode, song *models.Song, rating int) *BSTNode {
	if node == nil {
		// Create new node with rating bucket
		bucket := NewRatingBucket(rating)
		bucket.AddSong(song)
		bst.NodeCount++
		return &BSTNode{
			Bucket: bucket,
			Left:   nil,
			Right:  nil,
		}
	}

	if rating == node.Bucket.Rating {
		// Same rating, add to existing bucket
		node.Bucket.AddSong(song)
	} else if rating < node.Bucket.Rating {
		// Insert in left subtree
		node.Left = bst.insertNode(node.Left, song, rating)
	} else {
		// Insert in right subtree
		node.Right = bst.insertNode(node.Right, song, rating)
	}

	return node
}

// SearchByRating returns all songs with the specified rating
// Time Complexity: O(log n) average, O(n) worst case
// Space Complexity: O(1)
func (bst *SongRatingBST) SearchByRating(rating int) []*models.Song {
	if rating < 1 || rating > 5 {
		return []*models.Song{}
	}

	node := bst.searchNode(bst.Root, rating)
	if node != nil {
		return node.Bucket.GetSongs()
	}
	return []*models.Song{}
}

// searchNode is a recursive helper for searching nodes by rating
// Time Complexity: O(log n) average, O(n) worst case
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) searchNode(node *BSTNode, rating int) *BSTNode {
	if node == nil || node.Bucket.Rating == rating {
		return node
	}

	if rating < node.Bucket.Rating {
		return bst.searchNode(node.Left, rating)
	}
	return bst.searchNode(node.Right, rating)
}

// DeleteSong removes a song from the BST by song ID
// Time Complexity: O(log n + k) where k is songs in the rating bucket
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) DeleteSong(songID string) bool {
	song := bst.findSongByID(songID)
	if song == nil {
		return false
	}

	rating := song.Rating
	node := bst.searchNode(bst.Root, rating)
	if node == nil {
		return false
	}

	// Remove song from bucket
	removed := node.Bucket.RemoveSong(songID)
	if removed && node.Bucket.IsEmpty() {
		// If bucket is empty, remove the entire node
		bst.Root = bst.deleteNode(bst.Root, rating)
		bst.NodeCount--
	}

	return removed
}

// deleteNode is a recursive helper for deleting nodes
// Time Complexity: O(log n) average, O(n) worst case
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) deleteNode(node *BSTNode, rating int) *BSTNode {
	if node == nil {
		return nil
	}

	if rating < node.Bucket.Rating {
		node.Left = bst.deleteNode(node.Left, rating)
	} else if rating > node.Bucket.Rating {
		node.Right = bst.deleteNode(node.Right, rating)
	} else {
		// Node to be deleted found
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		}

		// Node has two children - find inorder successor
		successor := bst.findMinNode(node.Right)
		node.Bucket = successor.Bucket
		node.Right = bst.deleteNode(node.Right, successor.Bucket.Rating)
	}

	return node
}

// findMinNode finds the node with minimum rating in a subtree
// Time Complexity: O(log n) average, O(n) worst case
// Space Complexity: O(1)
func (bst *SongRatingBST) findMinNode(node *BSTNode) *BSTNode {
	for node.Left != nil {
		node = node.Left
	}
	return node
}

// findSongByID searches for a song by ID across all rating buckets
// Time Complexity: O(n * k) where n is nodes and k is average songs per bucket
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) findSongByID(songID string) *models.Song {
	return bst.findSongInSubtree(bst.Root, songID)
}

// findSongInSubtree is a recursive helper to find song in subtree
// Time Complexity: O(n * k) where n is nodes and k is average songs per bucket
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) findSongInSubtree(node *BSTNode, songID string) *models.Song {
	if node == nil {
		return nil
	}

	// Search in current bucket
	for _, song := range node.Bucket.Songs {
		if song.ID == songID {
			return song
		}
	}

	// Search in left subtree
	if song := bst.findSongInSubtree(node.Left, songID); song != nil {
		return song
	}

	// Search in right subtree
	return bst.findSongInSubtree(node.Right, songID)
}

// GetAllSongs returns all songs in the BST sorted by rating (ascending)
// Time Complexity: O(n * k) where n is nodes and k is average songs per bucket
// Space Complexity: O(n * k) for result slice + O(log n) for recursion
func (bst *SongRatingBST) GetAllSongs() []*models.Song {
	songs := make([]*models.Song, 0)
	bst.inorderTraversal(bst.Root, &songs)
	return songs
}

// inorderTraversal performs inorder traversal to get songs sorted by rating
// Time Complexity: O(n * k) where n is nodes and k is average songs per bucket
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) inorderTraversal(node *BSTNode, songs *[]*models.Song) {
	if node != nil {
		bst.inorderTraversal(node.Left, songs)
		*songs = append(*songs, node.Bucket.Songs...)
		bst.inorderTraversal(node.Right, songs)
	}
}

// GetSongsByRatingRange returns all songs within a rating range (inclusive)
// Time Complexity: O(n * k) where n is nodes and k is average songs per bucket
// Space Complexity: O(result_size)
func (bst *SongRatingBST) GetSongsByRatingRange(minRating, maxRating int) []*models.Song {
	if minRating > maxRating || minRating < 1 || maxRating > 5 {
		return []*models.Song{}
	}

	songs := make([]*models.Song, 0)
	bst.rangeSearch(bst.Root, minRating, maxRating, &songs)
	return songs
}

// rangeSearch is a recursive helper for range searching
// Time Complexity: O(n * k) in worst case
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) rangeSearch(node *BSTNode, minRating, maxRating int, songs *[]*models.Song) {
	if node == nil {
		return
	}

	// If current rating is in range, add songs
	if node.Bucket.Rating >= minRating && node.Bucket.Rating <= maxRating {
		*songs = append(*songs, node.Bucket.Songs...)
	}

	// Search left if there might be valid ratings
	if minRating < node.Bucket.Rating {
		bst.rangeSearch(node.Left, minRating, maxRating, songs)
	}

	// Search right if there might be valid ratings
	if maxRating > node.Bucket.Rating {
		bst.rangeSearch(node.Right, minRating, maxRating, songs)
	}
}

// GetRatingStats returns statistics about song ratings
// Time Complexity: O(n * k) where n is nodes and k is average songs per bucket
// Space Complexity: O(1)
func (bst *SongRatingBST) GetRatingStats() map[int]int {
	stats := make(map[int]int)
	bst.collectStats(bst.Root, stats)
	return stats
}

// collectStats is a recursive helper to collect rating statistics
// Time Complexity: O(n)
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) collectStats(node *BSTNode, stats map[int]int) {
	if node != nil {
		stats[node.Bucket.Rating] = len(node.Bucket.Songs)
		bst.collectStats(node.Left, stats)
		bst.collectStats(node.Right, stats)
	}
}

// IsEmpty checks if the BST is empty
// Time Complexity: O(1)
// Space Complexity: O(1)
func (bst *SongRatingBST) IsEmpty() bool {
	return bst.Root == nil
}

// GetNodeCount returns the number of rating nodes in the BST
// Time Complexity: O(1)
// Space Complexity: O(1)
func (bst *SongRatingBST) GetNodeCount() int {
	return bst.NodeCount
}

// GetTotalSongs returns the total number of songs across all ratings
// Time Complexity: O(n)
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) GetTotalSongs() int {
	return bst.countSongs(bst.Root)
}

// countSongs is a recursive helper to count total songs
// Time Complexity: O(n)
// Space Complexity: O(log n) due to recursion stack
func (bst *SongRatingBST) countSongs(node *BSTNode) int {
	if node == nil {
		return 0
	}
	return len(node.Bucket.Songs) + bst.countSongs(node.Left) + bst.countSongs(node.Right)
}

// Clear removes all nodes from the BST
// Time Complexity: O(1)
// Space Complexity: O(1)
func (bst *SongRatingBST) Clear() {
	bst.Root = nil
	bst.NodeCount = 0
}

// String returns a string representation of the BST
// Time Complexity: O(n * k * log(n))
// Space Complexity: O(n * k)
func (bst *SongRatingBST) String() string {
	if bst.IsEmpty() {
		return "Empty Rating Tree"
	}

	result := "Song Rating Tree (by rating):\n"
	songs := bst.GetAllSongs()

	currentRating := -1
	for _, song := range songs {
		if song.Rating != currentRating {
			currentRating = song.Rating
			result += fmt.Sprintf("\n%d Star Songs:\n", currentRating)
		}
		result += fmt.Sprintf("  - %s by %s (%s)\n", song.Title, song.Artist, song.DurationString())
	}

	return result
}
