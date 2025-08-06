package services

import (
	"fmt"
	"src/internal/datastructures"
	"src/internal/models"
	"strings"
	"time"
)

// PlaylistEngine represents the core music playlist management system
// Integrates all data structures: DoublyLinkedList, Stack, BST, HashMap, Sorting
// Time Complexity: Varies by operation, documented per method
// Space Complexity: O(n) where n is the total number of songs
type PlaylistEngine struct {
	// Main playlist storage
	currentPlaylist *datastructures.DoublyLinkedList

	// Playback history management
	playbackHistory *datastructures.PlaybackHistoryStack

	// Song rating system
	ratingTree *datastructures.SongRatingBST

	// Fast song lookup
	songLookup  *datastructures.SongHashMap
	titleLookup *datastructures.SongHashMap

	// Playlist organization
	playlistTree *datastructures.PlaylistExplorerTree

	// Sorting functionality
	sorter *datastructures.PlaylistSorter

	// Engine metadata
	playlistName  string
	totalPlayTime int
	createdAt     time.Time
}

// NewPlaylistEngine creates a new playlist engine instance
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewPlaylistEngine(playlistName string) *PlaylistEngine {
	return &PlaylistEngine{
		currentPlaylist: datastructures.NewDoublyLinkedList(),
		playbackHistory: datastructures.NewPlaybackHistoryStack(100), // Keep last 100 played songs
		ratingTree:      datastructures.NewSongRatingBST(),
		songLookup:      datastructures.NewSongHashMap(64),
		titleLookup:     datastructures.NewSongHashMap(64),
		playlistTree:    datastructures.NewPlaylistExplorerTree(),
		sorter:          datastructures.NewPlaylistSorter(datastructures.SortByTitle),
		playlistName:    playlistName,
		totalPlayTime:   0,
		createdAt:       time.Now(),
	}
}

// AddSong adds a song to the playlist with full synchronization across all data structures
// Time Complexity: O(1) average for most operations, O(log n) for BST insertion
// Space Complexity: O(1)
func (pe *PlaylistEngine) AddSong(title, artist, album, genre, subgenre, mood string, duration, bpm int) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(artist) == "" {
		return fmt.Errorf("title and artist are required")
	}

	// Generate unique ID for the song
	songID := pe.generateSongID(title, artist)

	// Create new song
	song := models.NewSong(songID, title, artist, album, genre, subgenre, mood, duration, bpm)

	// Check if song already exists
	if pe.songLookup.Contains(songID) {
		return fmt.Errorf("song already exists in playlist")
	}

	// Add to playlist (doubly linked list)
	pe.currentPlaylist.AddSong(song)

	// Add to hash maps for fast lookup
	pe.songLookup.Put(song)
	pe.titleLookup.PutByTitle(song)

	// Add to playlist explorer tree
	pe.playlistTree.AddSong(song)

	// Add to rating tree with default rating of 0 (will be updated when user rates)
	if song.Rating > 0 {
		pe.ratingTree.InsertSong(song, song.Rating)
	}

	// Update total play time
	pe.totalPlayTime += duration

	return nil
}

// DeleteSong removes a song from the playlist by index
// Time Complexity: O(n) for playlist deletion, O(1) average for hash map operations
// Space Complexity: O(1)
func (pe *PlaylistEngine) DeleteSong(index int) (*models.Song, error) {
	// Remove from playlist
	song, err := pe.currentPlaylist.DeleteSong(index)
	if err != nil {
		return nil, err
	}

	// Remove from hash maps
	pe.songLookup.Delete(song.ID)
	// Note: We don't remove from titleLookup as there might be multiple songs with same title

	// Remove from rating tree if it was rated
	if song.Rating > 0 {
		pe.ratingTree.DeleteSong(song.ID)
	}

	// Remove from playlist tree
	pe.playlistTree.RemoveSong(song.ID)

	// Update total play time
	pe.totalPlayTime -= song.Duration

	return song, nil
}

// MoveSong moves a song from one position to another in the playlist
// Time Complexity: O(n) where n is max(fromIndex, toIndex)
// Space Complexity: O(1)
func (pe *PlaylistEngine) MoveSong(fromIndex, toIndex int) error {
	return pe.currentPlaylist.MoveSong(fromIndex, toIndex)
}

// ReversePlaylist reverses the entire playlist order
// Time Complexity: O(n)
// Space Complexity: O(1)
func (pe *PlaylistEngine) ReversePlaylist() {
	pe.currentPlaylist.ReversePlaylist()
}

// PlaySong simulates playing a song and adds it to playback history
// Time Complexity: O(n) for finding song by index, O(1) for history operations
// Space Complexity: O(1)
func (pe *PlaylistEngine) PlaySong(index int) (*models.Song, error) {
	song, err := pe.currentPlaylist.GetSong(index)
	if err != nil {
		return nil, err
	}

	// Update song's play statistics
	song.Play()

	// Add to playback history
	pe.playbackHistory.Push(song)

	// Update in hash maps to reflect new play statistics
	pe.songLookup.UpdateSong(song)
	pe.titleLookup.UpdateSong(song)

	return song, nil
}

// UndoLastPlay removes the last played song from history and returns it
// Time Complexity: O(1)
// Space Complexity: O(1)
func (pe *PlaylistEngine) UndoLastPlay() (*models.Song, error) {
	return pe.playbackHistory.UndoLastPlay()
}

// RateSong assigns a rating to a song and updates the rating tree
// Time Complexity: O(log n) for BST operations, O(1) average for hash map updates
// Space Complexity: O(1)
func (pe *PlaylistEngine) RateSong(songID string, rating int) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	song, err := pe.songLookup.Get(songID)
	if err != nil {
		return fmt.Errorf("song not found: %v", err)
	}

	oldRating := song.Rating

	// Remove from old rating bucket if previously rated
	if oldRating > 0 {
		pe.ratingTree.DeleteSong(songID)
	}

	// Update song rating
	song.SetRating(rating)

	// Add to new rating bucket
	pe.ratingTree.InsertSong(song, rating)

	// Update in hash maps
	pe.songLookup.UpdateSong(song)
	pe.titleLookup.UpdateSong(song)

	return nil
}

// SearchSongByID provides O(1) song lookup by ID
// Time Complexity: O(1) average
// Space Complexity: O(1)
func (pe *PlaylistEngine) SearchSongByID(songID string) (*models.Song, error) {
	return pe.songLookup.Get(songID)
}

// SearchSongByTitle provides O(1) song lookup by title
// Time Complexity: O(1) average
// Space Complexity: O(1)
func (pe *PlaylistEngine) SearchSongByTitle(title string) (*models.Song, error) {
	return pe.titleLookup.GetByTitle(title)
}

// GetSongsByRating returns all songs with a specific rating
// Time Complexity: O(log n) average for BST search
// Space Complexity: O(k) where k is the number of songs with that rating
func (pe *PlaylistEngine) GetSongsByRating(rating int) []*models.Song {
	return pe.ratingTree.SearchByRating(rating)
}

// GetSongsByRatingRange returns songs within a rating range
// Time Complexity: O(n) worst case for range search
// Space Complexity: O(k) where k is the number of matching songs
func (pe *PlaylistEngine) GetSongsByRatingRange(minRating, maxRating int) []*models.Song {
	return pe.ratingTree.GetSongsByRatingRange(minRating, maxRating)
}

// SortPlaylist sorts the current playlist using specified criteria and algorithm
// Time Complexity: O(n log n)
// Space Complexity: O(n)
func (pe *PlaylistEngine) SortPlaylist(criteria datastructures.SortCriteria, algorithm string) {
	pe.sorter.SetCriteria(criteria)
	pe.sorter.SortPlaylist(pe.currentPlaylist, algorithm)
}

// GetRecentlyPlayedSongs returns recently played songs from history
// Time Complexity: O(min(n, count))
// Space Complexity: O(min(n, count))
func (pe *PlaylistEngine) GetRecentlyPlayedSongs(count int) []*models.Song {
	return pe.playbackHistory.GetRecentSongs(count)
}

// GetPlaylistByExplorer returns songs from the hierarchical explorer
// Time Complexity: O(1) for navigation
// Space Complexity: O(1)
func (pe *PlaylistEngine) GetPlaylistByExplorer(genre, subgenre, mood, artist string) []*models.Song {
	return pe.playlistTree.GetSongs(genre, subgenre, mood, artist)
}

// GetGenres returns all available genres from the explorer tree
// Time Complexity: O(g) where g is the number of genres
// Space Complexity: O(g)
func (pe *PlaylistEngine) GetGenres() []string {
	return pe.playlistTree.GetGenres()
}

// GetSubgenres returns subgenres for a specific genre
// Time Complexity: O(s) where s is the number of subgenres
// Space Complexity: O(s)
func (pe *PlaylistEngine) GetSubgenres(genre string) []string {
	return pe.playlistTree.GetSubgenres(genre)
}

// GetMoods returns moods for a specific genre and subgenre
// Time Complexity: O(m) where m is the number of moods
// Space Complexity: O(m)
func (pe *PlaylistEngine) GetMoods(genre, subgenre string) []string {
	return pe.playlistTree.GetMoods(genre, subgenre)
}

// GetArtists returns artists for a specific genre, subgenre, and mood
// Time Complexity: O(a) where a is the number of artists
// Space Complexity: O(a)
func (pe *PlaylistEngine) GetArtists(genre, subgenre, mood string) []string {
	return pe.playlistTree.GetArtists(genre, subgenre, mood)
}

// GetSmartRecommendations returns songs similar to recently played but not played recently
// Time Complexity: O(n * h) where n is total songs and h is history size
// Space Complexity: O(k) where k is the number of recommendations
func (pe *PlaylistEngine) GetSmartRecommendations(count int) []*models.Song {
	if count <= 0 {
		count = 10
	}

	recommendations := make([]*models.Song, 0, count)
	recentSongs := pe.playbackHistory.GetRecentSongs(20) // Look at last 20 played songs

	if len(recentSongs) == 0 {
		// No history, return random songs from playlist
		allSongs := pe.currentPlaylist.ToSlice()
		maxReturn := min(count, len(allSongs))
		return allSongs[:maxReturn]
	}

	allSongs := pe.currentPlaylist.ToSlice()
	recentSongIDs := make(map[string]bool)

	// Create set of recently played song IDs
	for _, song := range recentSongs {
		recentSongIDs[song.ID] = true
	}

	// Find similar songs that haven't been played recently
	for _, song := range allSongs {
		if len(recommendations) >= count {
			break
		}

		// Skip if recently played
		if recentSongIDs[song.ID] {
			continue
		}

		// Check similarity with recent songs
		for _, recentSong := range recentSongs {
			if song.IsSimilar(recentSong) {
				recommendations = append(recommendations, song)
				break
			}
		}
	}

	// If not enough similar songs, fill with unplayed songs
	if len(recommendations) < count {
		for _, song := range allSongs {
			if len(recommendations) >= count {
				break
			}

			if !recentSongIDs[song.ID] && !pe.containsSong(recommendations, song.ID) {
				recommendations = append(recommendations, song)
			}
		}
	}

	return recommendations
}

// ExportSnapshot generates a live dashboard snapshot of the playlist state
// Time Complexity: O(n) for statistics collection
// Space Complexity: O(n) for the snapshot data
func (pe *PlaylistEngine) ExportSnapshot() map[string]interface{} {
	// Get top 5 longest songs
	allSongs := pe.currentPlaylist.ToSlice()
	pe.sorter.SetCriteria(datastructures.SortByDurationDesc)
	sortedByDuration := pe.sorter.MergeSort(allSongs)

	top5Longest := make([]*models.Song, 0, 5)
	for i := 0; i < min(5, len(sortedByDuration)); i++ {
		top5Longest = append(top5Longest, sortedByDuration[i])
	}

	// Get most recently played songs
	recentlyPlayed := pe.playbackHistory.GetRecentSongs(10)

	// Get song count by rating
	ratingStats := pe.ratingTree.GetRatingStats()

	// Get playlist tree statistics
	treeStats := pe.playlistTree.GetStats()

	// Get playback statistics
	playbackStats := pe.playbackHistory.GetPlaybackStats()

	return map[string]interface{}{
		"playlist_info": map[string]interface{}{
			"name":           pe.playlistName,
			"total_songs":    pe.currentPlaylist.Size(),
			"total_duration": pe.totalPlayTime,
			"created_at":     pe.createdAt,
			"last_updated":   time.Now(),
		},
		"top_longest_songs":   top5Longest,
		"recently_played":     recentlyPlayed,
		"rating_distribution": ratingStats,
		"genre_stats":         treeStats,
		"playback_stats":      playbackStats,
		"hash_map_stats": map[string]interface{}{
			"song_lookup_size":  pe.songLookup.GetSize(),
			"song_lookup_load":  pe.songLookup.GetLoadFactor(),
			"title_lookup_size": pe.titleLookup.GetSize(),
			"title_lookup_load": pe.titleLookup.GetLoadFactor(),
		},
	}
}

// GetPlaylistStats returns comprehensive statistics about the playlist
// Time Complexity: O(n)
// Space Complexity: O(1)
func (pe *PlaylistEngine) GetPlaylistStats() map[string]interface{} {
	return map[string]interface{}{
		"total_songs":         pe.currentPlaylist.Size(),
		"total_duration":      pe.totalPlayTime,
		"average_song_length": pe.getAverageSongLength(),
		"total_play_count":    pe.getTotalPlayCount(),
		"unique_artists":      pe.getUniqueArtistCount(),
		"unique_genres":       pe.playlistTree.GetStats()["genres"],
		"rating_distribution": pe.ratingTree.GetRatingStats(),
		"history_size":        pe.playbackHistory.GetSize(),
	}
}

// Helper methods

// generateSongID creates a unique ID for a song
func (pe *PlaylistEngine) generateSongID(title, artist string) string {
	return fmt.Sprintf("%s-%s-%d",
		strings.ReplaceAll(strings.ToLower(title), " ", "-"),
		strings.ReplaceAll(strings.ToLower(artist), " ", "-"),
		time.Now().UnixNano())
}

// getAverageSongLength calculates the average song duration
func (pe *PlaylistEngine) getAverageSongLength() float64 {
	if pe.currentPlaylist.Size() == 0 {
		return 0
	}
	return float64(pe.totalPlayTime) / float64(pe.currentPlaylist.Size())
}

// getTotalPlayCount sums up play counts for all songs
func (pe *PlaylistEngine) getTotalPlayCount() int {
	total := 0
	songs := pe.currentPlaylist.ToSlice()
	for _, song := range songs {
		total += song.PlayCount
	}
	return total
}

// getUniqueArtistCount counts unique artists in the playlist
func (pe *PlaylistEngine) getUniqueArtistCount() int {
	artistSet := make(map[string]bool)
	songs := pe.currentPlaylist.ToSlice()
	for _, song := range songs {
		artistSet[song.Artist] = true
	}
	return len(artistSet)
}

// containsSong checks if a song ID exists in a slice of songs
func (pe *PlaylistEngine) containsSong(songs []*models.Song, songID string) bool {
	for _, song := range songs {
		if song.ID == songID {
			return true
		}
	}
	return false
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetCurrentPlaylist returns all songs in the current playlist
// Time Complexity: O(n)
// Space Complexity: O(n)
func (pe *PlaylistEngine) GetCurrentPlaylist() []*models.Song {
	return pe.currentPlaylist.ToSlice()
}

// GetPlaylistSize returns the size of the current playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (pe *PlaylistEngine) GetPlaylistSize() int {
	return pe.currentPlaylist.Size()
}

// GetPlaylistName returns the name of the playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (pe *PlaylistEngine) GetPlaylistName() string {
	return pe.playlistName
}

// SetPlaylistName updates the playlist name
// Time Complexity: O(1)
// Space Complexity: O(1)
func (pe *PlaylistEngine) SetPlaylistName(name string) {
	pe.playlistName = name
}

// ClearPlaylist removes all songs from the playlist
// Time Complexity: O(1)
// Space Complexity: O(1)
func (pe *PlaylistEngine) ClearPlaylist() {
	pe.currentPlaylist.Clear()
	pe.ratingTree.Clear()
	pe.songLookup.Clear()
	pe.titleLookup.Clear()
	pe.playlistTree = datastructures.NewPlaylistExplorerTree()
	pe.totalPlayTime = 0
}

// BenchmarkSort compares the performance of different sorting algorithms
// Time Complexity: O(n log n) for each algorithm tested
// Space Complexity: O(n) for creating copies
func (pe *PlaylistEngine) BenchmarkSort() map[string]time.Duration {
	songs := pe.currentPlaylist.ToSlice()
	return pe.sorter.BenchmarkSort(songs)
}
