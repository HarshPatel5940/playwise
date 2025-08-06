package services

import (
	"fmt"
	"src/internal/datastructures"
	"src/internal/models"
	"strings"
	"testing"
	"time"
)

func TestNewPlaylistEngine(t *testing.T) {
	engine := NewPlaylistEngine("Test Playlist")

	if engine == nil {
		t.Fatal("Expected non-nil playlist engine")
	}
	if engine.playlistName != "Test Playlist" {
		t.Errorf("Expected playlist name 'Test Playlist', got %s", engine.playlistName)
	}
	if engine.currentPlaylist == nil {
		t.Error("Current playlist should be initialized")
	}
	if engine.playbackHistory == nil {
		t.Error("Playback history should be initialized")
	}
	if engine.ratingTree == nil {
		t.Error("Rating tree should be initialized")
	}
	if engine.songLookup == nil {
		t.Error("Song lookup should be initialized")
	}
	if engine.titleLookup == nil {
		t.Error("Title lookup should be initialized")
	}
	if engine.playlistTree == nil {
		t.Error("Playlist tree should be initialized")
	}
	if engine.sorter == nil {
		t.Error("Sorter should be initialized")
	}
	if engine.totalPlayTime != 0 {
		t.Error("Total play time should be initialized to 0")
	}
	if engine.createdAt.IsZero() {
		t.Error("Created at should be set")
	}
}

func TestAddSong(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Test valid song addition
	err := engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if engine.GetPlaylistSize() != 1 {
		t.Errorf("Expected playlist size 1, got %d", engine.GetPlaylistSize())
	}

	if engine.totalPlayTime != 240 {
		t.Errorf("Expected total play time 240, got %d", engine.totalPlayTime)
	}

	// Test duplicate song addition
	err = engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	if err == nil {
		t.Error("Expected error for duplicate song")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected duplicate error message, got %s", err.Error())
	}

	// Test empty title
	err = engine.AddSong("", "Artist", "Album", "Genre", "Subgenre", "Mood", 180, 100)
	if err == nil {
		t.Error("Expected error for empty title")
	}

	// Test empty artist
	err = engine.AddSong("Title", "", "Album", "Genre", "Subgenre", "Mood", 180, 100)
	if err == nil {
		t.Error("Expected error for empty artist")
	}

	// Test whitespace-only title and artist
	err = engine.AddSong("   ", "   ", "Album", "Genre", "Subgenre", "Mood", 180, 100)
	if err == nil {
		t.Error("Expected error for whitespace-only title and artist")
	}
}

func TestDeleteSong(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add some songs first
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	engine.AddSong("Song 3", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 300, 90)

	initialSize := engine.GetPlaylistSize()
	initialPlayTime := engine.totalPlayTime

	// Test valid deletion
	deletedSong, err := engine.DeleteSong(1) // Delete "Song 2"
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if deletedSong == nil {
		t.Fatal("Expected deleted song to be returned")
	}
	if deletedSong.Title != "Song 2" {
		t.Errorf("Expected deleted song title 'Song 2', got %s", deletedSong.Title)
	}

	if engine.GetPlaylistSize() != initialSize-1 {
		t.Errorf("Expected playlist size %d, got %d", initialSize-1, engine.GetPlaylistSize())
	}
	if engine.totalPlayTime != initialPlayTime-200 {
		t.Errorf("Expected total play time %d, got %d", initialPlayTime-200, engine.totalPlayTime)
	}

	// Test invalid index
	_, err = engine.DeleteSong(100)
	if err == nil {
		t.Error("Expected error for invalid index")
	}

	_, err = engine.DeleteSong(-1)
	if err == nil {
		t.Error("Expected error for negative index")
	}
}

func TestMoveSong(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	engine.AddSong("Song 3", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 300, 90)

	// Test valid move
	err := engine.MoveSong(0, 2) // Move first song to last position
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the move
	songs := engine.GetCurrentPlaylist()
	if songs[2].Title != "Song 1" {
		t.Errorf("Expected 'Song 1' at position 2, got %s", songs[2].Title)
	}

	// Test invalid indices
	err = engine.MoveSong(-1, 0)
	if err == nil {
		t.Error("Expected error for negative from index")
	}

	err = engine.MoveSong(0, 100)
	if err == nil {
		t.Error("Expected error for invalid to index")
	}
}

func TestReversePlaylist(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	engine.AddSong("Song 3", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 300, 90)

	originalSongs := engine.GetCurrentPlaylist()

	// Reverse playlist
	engine.ReversePlaylist()

	reversedSongs := engine.GetCurrentPlaylist()

	// Verify reversal
	if len(reversedSongs) != len(originalSongs) {
		t.Error("Playlist size should not change after reversal")
	}

	for i := 0; i < len(originalSongs); i++ {
		if originalSongs[i].ID != reversedSongs[len(reversedSongs)-1-i].ID {
			t.Error("Playlist not properly reversed")
			break
		}
	}
}

func TestPlaySong(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add a song
	engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	// Play the song
	playedSong, err := engine.PlaySong(0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if playedSong == nil {
		t.Fatal("Expected played song to be returned")
	}

	// Check that play count increased
	if playedSong.PlayCount != 1 {
		t.Errorf("Expected play count 1, got %d", playedSong.PlayCount)
	}

	// Check that song is in playback history
	recentSongs := engine.GetRecentlyPlayedSongs(1)
	if len(recentSongs) != 1 {
		t.Errorf("Expected 1 recent song, got %d", len(recentSongs))
	}
	if recentSongs[0].ID != playedSong.ID {
		t.Error("Recently played song should match played song")
	}

	// Test invalid index
	_, err = engine.PlaySong(100)
	if err == nil {
		t.Error("Expected error for invalid index")
	}
}

func TestUndoLastPlay(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Test undo with empty history
	_, err := engine.UndoLastPlay()
	if err == nil {
		t.Error("Expected error when undoing with empty history")
	}

	// Add and play a song
	engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	playedSong, _ := engine.PlaySong(0)

	// Undo last play
	undoSong, err := engine.UndoLastPlay()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if undoSong == nil {
		t.Fatal("Expected undo song to be returned")
	}
	if undoSong.ID != playedSong.ID {
		t.Error("Undo song should match last played song")
	}

	// Check that history is now empty
	recentSongs := engine.GetRecentlyPlayedSongs(1)
	if len(recentSongs) != 0 {
		t.Errorf("Expected 0 recent songs after undo, got %d", len(recentSongs))
	}
}

func TestRateSong(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add a song
	engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	songs := engine.GetCurrentPlaylist()
	songID := songs[0].ID

	// Test valid rating
	err := engine.RateSong(songID, 4)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify rating was set
	ratedSongs := engine.GetSongsByRating(4)
	if len(ratedSongs) != 1 {
		t.Errorf("Expected 1 song with rating 4, got %d", len(ratedSongs))
	}
	if ratedSongs[0].ID != songID {
		t.Error("Rated song ID should match")
	}

	// Test changing rating
	err = engine.RateSong(songID, 5)
	if err != nil {
		t.Errorf("Expected no error changing rating, got %v", err)
	}

	// Verify old rating bucket is empty and new rating bucket has the song
	oldRatingSongs := engine.GetSongsByRating(4)
	if len(oldRatingSongs) != 0 {
		t.Errorf("Expected 0 songs with old rating, got %d", len(oldRatingSongs))
	}

	newRatingSongs := engine.GetSongsByRating(5)
	if len(newRatingSongs) != 1 {
		t.Errorf("Expected 1 song with new rating, got %d", len(newRatingSongs))
	}

	// Test invalid rating
	err = engine.RateSong(songID, 0)
	if err == nil {
		t.Error("Expected error for rating 0")
	}

	err = engine.RateSong(songID, 6)
	if err == nil {
		t.Error("Expected error for rating 6")
	}

	// Test non-existent song
	err = engine.RateSong("nonexistent", 5)
	if err == nil {
		t.Error("Expected error for non-existent song")
	}
}

func TestSearchSongByID(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add a song
	engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	songs := engine.GetCurrentPlaylist()
	songID := songs[0].ID

	// Test valid search
	foundSong, err := engine.SearchSongByID(songID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if foundSong == nil {
		t.Fatal("Expected found song")
	}
	if foundSong.ID != songID {
		t.Error("Found song ID should match search ID")
	}

	// Test invalid search
	_, err = engine.SearchSongByID("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent song ID")
	}
}

func TestSearchSongByTitle(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add a song
	engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	// Test valid search
	foundSong, err := engine.SearchSongByTitle("Test Song")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if foundSong == nil {
		t.Fatal("Expected found song")
	}
	if foundSong.Title != "Test Song" {
		t.Error("Found song title should match search title")
	}

	// Test invalid search
	_, err = engine.SearchSongByTitle("Nonexistent Song")
	if err == nil {
		t.Error("Expected error for non-existent song title")
	}
}

func TestGetSongsByRating(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs with different ratings
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	engine.AddSong("Song 3", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 300, 90)

	songs := engine.GetCurrentPlaylist()
	engine.RateSong(songs[0].ID, 4)
	engine.RateSong(songs[1].ID, 4)
	engine.RateSong(songs[2].ID, 5)

	// Test getting songs by rating
	rating4Songs := engine.GetSongsByRating(4)
	if len(rating4Songs) != 2 {
		t.Errorf("Expected 2 songs with rating 4, got %d", len(rating4Songs))
	}

	rating5Songs := engine.GetSongsByRating(5)
	if len(rating5Songs) != 1 {
		t.Errorf("Expected 1 song with rating 5, got %d", len(rating5Songs))
	}

	// Test non-existent rating
	rating1Songs := engine.GetSongsByRating(1)
	if len(rating1Songs) != 0 {
		t.Errorf("Expected 0 songs with rating 1, got %d", len(rating1Songs))
	}
}

func TestGetSongsByRatingRange(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs with different ratings
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	engine.AddSong("Song 3", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 300, 90)
	engine.AddSong("Song 4", "Artist 4", "Album 4", "Blues", "Electric", "Intense", 280, 95)

	songs := engine.GetCurrentPlaylist()
	engine.RateSong(songs[0].ID, 2)
	engine.RateSong(songs[1].ID, 3)
	engine.RateSong(songs[2].ID, 4)
	engine.RateSong(songs[3].ID, 5)

	// Test range search
	rangeSongs := engine.GetSongsByRatingRange(3, 4)
	if len(rangeSongs) != 2 {
		t.Errorf("Expected 2 songs in range 3-4, got %d", len(rangeSongs))
	}

	// Test full range
	allRatedSongs := engine.GetSongsByRatingRange(1, 5)
	if len(allRatedSongs) != 4 {
		t.Errorf("Expected 4 songs in range 1-5, got %d", len(allRatedSongs))
	}

	// Test empty range
	emptySongs := engine.GetSongsByRatingRange(6, 7)
	if len(emptySongs) != 0 {
		t.Errorf("Expected 0 songs in range 6-7, got %d", len(emptySongs))
	}
}

func TestSortPlaylist(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs in unsorted order
	engine.AddSong("Zebra", "Artist C", "Album 3", "Rock", "Alternative", "Energetic", 300, 120)
	engine.AddSong("Alpha", "Artist A", "Album 1", "Pop", "Mainstream", "Happy", 180, 110)
	engine.AddSong("Beta", "Artist B", "Album 2", "Jazz", "Smooth", "Relaxed", 240, 90)

	// Test sorting by title
	engine.SortPlaylist(datastructures.SortByTitle, "merge")

	songs := engine.GetCurrentPlaylist()
	expectedTitles := []string{"Alpha", "Beta", "Zebra"}

	for i, expectedTitle := range expectedTitles {
		if songs[i].Title != expectedTitle {
			t.Errorf("Position %d: expected %s, got %s", i, expectedTitle, songs[i].Title)
		}
	}

	// Test sorting by duration (ascending)
	engine.SortPlaylist(datastructures.SortByDurationAsc, "quick")

	songs = engine.GetCurrentPlaylist()
	expectedDurations := []int{180, 240, 300}

	for i, expectedDuration := range expectedDurations {
		if songs[i].Duration != expectedDuration {
			t.Errorf("Position %d: expected duration %d, got %d", i, expectedDuration, songs[i].Duration)
		}
	}
}

func TestGetRecentlyPlayedSongs(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add and play multiple songs
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	engine.AddSong("Song 3", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 300, 90)

	engine.PlaySong(0)
	engine.PlaySong(1)
	engine.PlaySong(2)

	// Test getting recent songs
	recent := engine.GetRecentlyPlayedSongs(2)
	if len(recent) != 2 {
		t.Errorf("Expected 2 recent songs, got %d", len(recent))
	}

	// Should be in reverse order (most recent first)
	if recent[0].Title != "Song 3" {
		t.Errorf("Expected most recent song to be 'Song 3', got %s", recent[0].Title)
	}

	// Test getting more than available
	allRecent := engine.GetRecentlyPlayedSongs(10)
	if len(allRecent) != 3 {
		t.Errorf("Expected 3 recent songs, got %d", len(allRecent))
	}

	// Test getting zero
	zeroRecent := engine.GetRecentlyPlayedSongs(0)
	if len(zeroRecent) != 0 {
		t.Errorf("Expected 0 recent songs, got %d", len(zeroRecent))
	}
}

func TestPlaylistExplorerMethods(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs with different categories
	engine.AddSong("Rock Song 1", "Rock Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Rock Song 2", "Rock Artist 2", "Album 2", "Rock", "Classic Rock", "Epic", 280, 115)
	engine.AddSong("Pop Song", "Pop Artist", "Album 3", "Pop", "Mainstream", "Happy", 200, 125)

	// Test getting genres
	genres := engine.GetGenres()
	if len(genres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(genres))
	}

	// Test getting subgenres
	rockSubgenres := engine.GetSubgenres("Rock")
	if len(rockSubgenres) != 2 {
		t.Errorf("Expected 2 Rock subgenres, got %d", len(rockSubgenres))
	}

	// Test getting moods
	alternativeMoods := engine.GetMoods("Rock", "Alternative")
	if len(alternativeMoods) != 1 {
		t.Errorf("Expected 1 mood for Rock->Alternative, got %d", len(alternativeMoods))
	}

	// Test getting artists
	artists := engine.GetArtists("Rock", "Alternative", "Energetic")
	if len(artists) != 1 {
		t.Errorf("Expected 1 artist for Rock->Alternative->Energetic, got %d", len(artists))
	}

	// Test getting songs by explorer
	songs := engine.GetPlaylistByExplorer("Rock", "Alternative", "Energetic", "Rock Artist 1")
	if len(songs) != 1 {
		t.Errorf("Expected 1 song for specific path, got %d", len(songs))
	}
	if songs[0].Title != "Rock Song 1" {
		t.Errorf("Expected 'Rock Song 1', got %s", songs[0].Title)
	}
}

func TestGetSmartRecommendations(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs with similar characteristics
	engine.AddSong("Rock Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Rock Song 2", "Artist 2", "Album 2", "Rock", "Alternative", "Energetic", 250, 125)
	engine.AddSong("Pop Song", "Artist 3", "Album 3", "Pop", "Mainstream", "Happy", 200, 110)
	engine.AddSong("Jazz Song", "Artist 4", "Album 4", "Jazz", "Smooth", "Relaxed", 300, 90)

	// Test with no history
	recommendations := engine.GetSmartRecommendations(2)
	if len(recommendations) == 0 {
		t.Error("Should return some recommendations even with no history")
	}

	// Play some songs to create history
	engine.PlaySong(0) // Play "Rock Song 1"

	recommendations = engine.GetSmartRecommendations(3)
	if len(recommendations) == 0 {
		t.Error("Should return recommendations based on history")
	}

	// Verify that recently played song is not in recommendations
	for _, rec := range recommendations {
		if rec.Title == "Rock Song 1" {
			t.Error("Recently played song should not be in recommendations")
		}
	}

	// Test with count of 0 (should default to 10)
	recommendations = engine.GetSmartRecommendations(0)
	if len(recommendations) == 0 {
		t.Error("Should return recommendations with default count")
	}
}

func TestExportSnapshot(t *testing.T) {
	engine := NewPlaylistEngine("Test Playlist")

	// Add some songs and play them
	engine.AddSong("Long Song", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 400, 120)
	engine.AddSong("Short Song", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 150, 130)
	engine.AddSong("Medium Song", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 250, 100)

	songs := engine.GetCurrentPlaylist()
	engine.RateSong(songs[0].ID, 5)
	engine.RateSong(songs[1].ID, 3)

	engine.PlaySong(0)
	engine.PlaySong(1)

	snapshot := engine.ExportSnapshot()

	// Test playlist info
	playlistInfo, exists := snapshot["playlist_info"]
	if !exists {
		t.Error("Snapshot should contain playlist_info")
	}

	playlistInfoMap, ok := playlistInfo.(map[string]interface{})
	if !ok {
		t.Error("playlist_info should be a map")
	}

	if playlistInfoMap["name"] != "Test Playlist" {
		t.Error("Snapshot should contain correct playlist name")
	}

	if playlistInfoMap["total_songs"].(int) != 3 {
		t.Error("Snapshot should contain correct song count")
	}

	// Test top longest songs
	topLongest, exists := snapshot["top_longest_songs"]
	if !exists {
		t.Error("Snapshot should contain top_longest_songs")
	}

	topLongestSlice, ok := topLongest.([]*models.Song)
	if !ok {
		t.Error("top_longest_songs should be a song slice")
	}

	if len(topLongestSlice) == 0 {
		t.Error("Should have at least one longest song")
	}

	// First song should be the longest
	if topLongestSlice[0].Title != "Long Song" {
		t.Error("Longest song should be first in top_longest_songs")
	}

	// Test recently played
	_, exists = snapshot["recently_played"]
	if !exists {
		t.Error("Snapshot should contain recently_played")
	}

	// Test other sections exist
	_, exists = snapshot["rating_distribution"]
	if !exists {
		t.Error("Snapshot should contain rating_distribution")
	}

	_, exists = snapshot["genre_stats"]
	if !exists {
		t.Error("Snapshot should contain genre_stats")
	}

	_, exists = snapshot["hash_map_stats"]
	if !exists {
		t.Error("Snapshot should contain hash_map_stats")
	}
}

func TestGetPlaylistStats(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs with different artists and play counts
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 200, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 300, 130)
	engine.AddSong("Song 3", "Artist 1", "Album 3", "Rock", "Classic Rock", "Epic", 250, 110)

	// Play some songs
	engine.PlaySong(0)
	engine.PlaySong(0) // Play same song twice

	// Rate some songs
	songs := engine.GetCurrentPlaylist()
	engine.RateSong(songs[0].ID, 4)
	engine.RateSong(songs[1].ID, 5)

	stats := engine.GetPlaylistStats()

	// Test total songs
	if stats["total_songs"].(int) != 3 {
		t.Errorf("Expected 3 total songs, got %v", stats["total_songs"])
	}

	// Test total duration
	if stats["total_duration"].(int) != 750 { // 200 + 300 + 250
		t.Errorf("Expected total duration 750, got %v", stats["total_duration"])
	}

	// Test average song length
	avgLength := stats["average_song_length"].(float64)
	if avgLength != 250.0 { // 750 / 3
		t.Errorf("Expected average song length 250, got %f", avgLength)
	}

	// Test unique artists (should be 2: Artist 1, Artist 2)
	if stats["unique_artists"].(int) != 2 {
		t.Errorf("Expected 2 unique artists, got %v", stats["unique_artists"])
	}

	// Test that other stats exist
	_, exists := stats["total_play_count"]
	if !exists {
		t.Error("Stats should contain total_play_count")
	}

	_, exists = stats["unique_genres"]
	if !exists {
		t.Error("Stats should contain unique_genres")
	}

	_, exists = stats["rating_distribution"]
	if !exists {
		t.Error("Stats should contain rating_distribution")
	}
}

func TestPlaylistNameOperations(t *testing.T) {
	engine := NewPlaylistEngine("Original Name")

	// Test getting name
	if engine.GetPlaylistName() != "Original Name" {
		t.Error("Should return correct playlist name")
	}

	// Test setting name
	engine.SetPlaylistName("New Name")
	if engine.GetPlaylistName() != "New Name" {
		t.Error("Should update playlist name")
	}
}

func TestClearPlaylist(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add some songs and data
	engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)

	songs := engine.GetCurrentPlaylist()
	engine.RateSong(songs[0].ID, 4)
	engine.PlaySong(0)

	// Verify data exists
	if engine.GetPlaylistSize() == 0 {
		t.Fatal("Playlist should have songs before clear")
	}
	if engine.totalPlayTime == 0 {
		t.Fatal("Should have play time before clear")
	}

	// Clear playlist
	engine.ClearPlaylist()

	// Verify everything is cleared
	if engine.GetPlaylistSize() != 0 {
		t.Error("Playlist should be empty after clear")
	}
	if engine.totalPlayTime != 0 {
		t.Error("Total play time should be 0 after clear")
	}

	// Verify lookups are cleared
	_, err := engine.SearchSongByID(songs[0].ID)
	if err == nil {
		t.Error("Song lookup should be cleared")
	}

	// Verify tree is reset
	genres := engine.GetGenres()
	if len(genres) != 0 {
		t.Error("Genre tree should be cleared")
	}
}

func TestBenchmarkSort(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Add songs for benchmarking
	for i := 0; i < 10; i++ {
		engine.AddSong(
			fmt.Sprintf("Song %d", i),
			fmt.Sprintf("Artist %d", i),
			"Album",
			"Genre",
			"Subgenre",
			"Mood",
			200+i*10,
			120,
		)
	}

	benchmarks := engine.BenchmarkSort()

	// Verify benchmark results
	expectedAlgorithms := []string{"merge_sort", "quick_sort", "heap_sort"}
	for _, algorithm := range expectedAlgorithms {
		if _, exists := benchmarks[algorithm]; !exists {
			t.Errorf("Benchmark missing for %s", algorithm)
		}
		if benchmarks[algorithm] < 0 {
			t.Errorf("Benchmark time cannot be negative for %s", algorithm)
		}
	}
}

func TestGenerateSongID(t *testing.T) {
	engine := NewPlaylistEngine("Test")

	// Generate IDs for same song at different times
	id1 := engine.generateSongID("Test Song", "Test Artist")
	time.Sleep(1 * time.Millisecond) // Ensure different timestamp
	id2 := engine.generateSongID("Test Song", "Test Artist")

	// IDs should be different due to timestamp
	if id1 == id2 {
		t.Error("Generated IDs should be unique")
	}

	// IDs should contain normalized title and artist
	if !strings.Contains(id1, "test-song") {
		t.Error("ID should contain normalized title")
	}
	if !strings.Contains(id1, "test-artist") {
		t.Error("ID should contain normalized artist")
	}
}

func TestIntegrationScenario(t *testing.T) {
	// Full integration test simulating real usage
	engine := NewPlaylistEngine("My Awesome Playlist")

	// Add diverse songs
	songs := []struct {
		title, artist, album, genre, subgenre, mood string
		duration, bpm                               int
	}{
		{"Bohemian Rhapsody", "Queen", "A Night at the Opera", "Rock", "Classic Rock", "Epic", 355, 72},
		{"Smells Like Teen Spirit", "Nirvana", "Nevermind", "Rock", "Alternative", "Energetic", 301, 117},
		{"Billie Jean", "Michael Jackson", "Thriller", "Pop", "Dance Pop", "Upbeat", 294, 117},
		{"Take Five", "Dave Brubeck", "Time Out", "Jazz", "Cool Jazz", "Relaxed", 324, 86},
		{"Hotel California", "Eagles", "Hotel California", "Rock", "Classic Rock", "Mysterious", 391, 74},
	}

	for _, song := range songs {
		err := engine.AddSong(song.title, song.artist, song.album, song.genre, song.subgenre, song.mood, song.duration, song.bpm)
		if err != nil {
			t.Errorf("Failed to add song %s: %v", song.title, err)
		}
	}

	// Test playlist size
	if engine.GetPlaylistSize() != 5 {
		t.Errorf("Expected 5 songs, got %d", engine.GetPlaylistSize())
	}

	// Play some songs
	engine.PlaySong(1) // Nirvana
	engine.PlaySong(2) // Michael Jackson
	engine.PlaySong(0) // Queen

	// Rate songs
	allSongs := engine.GetCurrentPlaylist()
	engine.RateSong(allSongs[0].ID, 5) // Queen
	engine.RateSong(allSongs[1].ID, 4) // Nirvana
	engine.RateSong(allSongs[2].ID, 4) // Michael Jackson

	// Test search functionality
	queen, err := engine.SearchSongByTitle("Bohemian Rhapsody")
	if err != nil || queen.Artist != "Queen" {
		t.Error("Failed to search Queen song")
	}

	// Test rating queries
	fiveStarSongs := engine.GetSongsByRating(5)
	if len(fiveStarSongs) != 1 || fiveStarSongs[0].Artist != "Queen" {
		t.Error("Failed to get 5-star songs")
	}

	// Test sorting
	engine.SortPlaylist(datastructures.SortByDurationDesc, "merge")
	sortedSongs := engine.GetCurrentPlaylist()
	if sortedSongs[0].Title != "Hotel California" { // Longest song
		t.Error("Failed to sort by duration descending")
	}

	// Test explorer functionality
	rockSongs := engine.GetPlaylistByExplorer("Rock", "Classic Rock", "Epic", "Queen")
	if len(rockSongs) != 1 || rockSongs[0].Title != "Bohemian Rhapsody" {
		t.Error("Failed to navigate playlist tree")
	}

	// Test recommendations
	recommendations := engine.GetSmartRecommendations(3)
	if len(recommendations) == 0 {
		t.Error("Failed to get recommendations")
	}

	// Test snapshot export
	snapshot := engine.ExportSnapshot()
	if snapshot == nil {
		t.Error("Failed to export snapshot")
	}

	playlistInfo := snapshot["playlist_info"].(map[string]interface{})
	if playlistInfo["total_songs"].(int) != 5 {
		t.Error("Snapshot has incorrect song count")
	}

	// Test recent playback
	recent := engine.GetRecentlyPlayedSongs(3)
	if len(recent) != 3 {
		t.Error("Failed to get recent songs")
	}
	if recent[0].Artist != "Queen" { // Most recent
		t.Error("Recent songs not in correct order")
	}

	// Test comprehensive stats
	stats := engine.GetPlaylistStats()
	if stats["unique_artists"].(int) != 5 {
		t.Error("Incorrect unique artist count")
	}
	if stats["unique_genres"].(int) != 3 { // Rock, Pop, Jazz
		t.Error("Incorrect unique genre count")
	}
}

// Helper function for fmt.Sprintf in tests
func init() {
	// This ensures fmt is available for sprintf operations in tests
	_ = fmt.Sprintf
}
