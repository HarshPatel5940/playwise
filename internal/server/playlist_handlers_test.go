package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNewPlaylistHandlers(t *testing.T) {
	handlers := NewPlaylistHandlers()

	if handlers == nil {
		t.Fatal("Expected non-nil handlers")
	}
	if handlers.engine == nil {
		t.Error("Engine should be initialized")
	}
}

func setupTestEcho() (*echo.Echo, *PlaylistHandlers) {
	e := echo.New()
	handlers := NewPlaylistHandlers()
	return e, handlers
}

func TestGetPlaylist(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a test song first
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodGet, "/playlist", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetPlaylist(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	songs, exists := response["songs"]
	if !exists {
		t.Error("Response should contain songs")
	}

	songsSlice, ok := songs.([]interface{})
	if !ok {
		t.Error("Songs should be an array")
	}

	if len(songsSlice) != 1 {
		t.Errorf("Expected 1 song, got %d", len(songsSlice))
	}
}

func TestAddSong(t *testing.T) {
	e, handlers := setupTestEcho()

	// Test valid song addition
	requestBody := map[string]interface{}{
		"title":    "Test Song",
		"artist":   "Test Artist",
		"album":    "Test Album",
		"genre":    "Rock",
		"subgenre": "Alternative",
		"mood":     "Energetic",
		"duration": 240,
		"bpm":      120,
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/playlist/songs", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.AddSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "Song added successfully" {
		t.Error("Response should contain success message")
	}

	// Verify song was actually added
	if handlers.engine.GetPlaylistSize() != 1 {
		t.Error("Song should have been added to engine")
	}
}

func TestAddSongInvalidJSON(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodPost, "/playlist/songs", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.AddSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestAddSongMissingFields(t *testing.T) {
	e, handlers := setupTestEcho()

	// Test missing title
	requestBody := map[string]interface{}{
		"artist":   "Test Artist",
		"duration": 240,
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/playlist/songs", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.AddSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestDeleteSong(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a song first
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodDelete, "/playlist/songs/0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("index")
	c.SetParamValues("0")

	err := handlers.DeleteSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify song was deleted
	if handlers.engine.GetPlaylistSize() != 0 {
		t.Error("Song should have been deleted from engine")
	}
}

func TestDeleteSongInvalidIndex(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodDelete, "/playlist/songs/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("index")
	c.SetParamValues("invalid")

	err := handlers.DeleteSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestDeleteSongOutOfRange(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodDelete, "/playlist/songs/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("index")
	c.SetParamValues("999")

	err := handlers.DeleteSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

func TestMoveSong(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add multiple songs
	handlers.engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	handlers.engine.AddSong("Song 3", "Artist 3", "Album 3", "Jazz", "Smooth", "Relaxed", 300, 90)

	requestBody := map[string]interface{}{
		"fromIndex": 0,
		"toIndex":   2,
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/playlist/move", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.MoveSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify the move
	songs := handlers.engine.GetCurrentPlaylist()
	if songs[2].Title != "Song 1" {
		t.Error("Song should have been moved to new position")
	}
}

func TestMoveSongInvalidJSON(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodPut, "/playlist/move", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.MoveSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestReversePlaylist(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs
	handlers.engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)

	originalSongs := handlers.engine.GetCurrentPlaylist()
	originalFirst := originalSongs[0].Title

	req := httptest.NewRequest(http.MethodPut, "/playlist/reverse", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.ReversePlaylist(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify reversal
	reversedSongs := handlers.engine.GetCurrentPlaylist()
	if reversedSongs[0].Title == originalFirst {
		t.Error("Playlist should have been reversed")
	}
}

func TestPlaySong(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodPost, "/playlist/play/0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("index")
	c.SetParamValues("0")

	err := handlers.PlaySong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify song was played (check play count)
	songs := handlers.engine.GetCurrentPlaylist()
	if songs[0].PlayCount != 1 {
		t.Error("Song play count should have increased")
	}
}

func TestPlaySongInvalidIndex(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodPost, "/playlist/play/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("index")
	c.SetParamValues("invalid")

	err := handlers.PlaySong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestUndoLastPlay(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add and play a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.PlaySong(0)

	req := httptest.NewRequest(http.MethodPost, "/playlist/undo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.UndoLastPlay(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestUndoLastPlayEmpty(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodPost, "/playlist/undo", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.UndoLastPlay(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestRateSong(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	songs := handlers.engine.GetCurrentPlaylist()
	songID := songs[0].ID

	requestBody := map[string]interface{}{
		"rating": 4,
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/playlist/rate/"+songID, bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("songId")
	c.SetParamValues(songID)

	err := handlers.RateSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify rating was set
	ratedSongs := handlers.engine.GetSongsByRating(4)
	if len(ratedSongs) != 1 {
		t.Error("Song should have been rated")
	}
}

func TestRateSongInvalidRating(t *testing.T) {
	e, handlers := setupTestEcho()

	requestBody := map[string]interface{}{
		"rating": 6, // Invalid rating
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/playlist/rate/songid", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("songId")
	c.SetParamValues("songid")

	err := handlers.RateSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestSearchSong(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	// Test search by title
	req := httptest.NewRequest(http.MethodGet, "/playlist/search?type=title&q=Test+Song", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SearchSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	data, exists := response["data"]
	if !exists {
		t.Error("Response should contain data")
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		t.Error("Data should be an object")
	}

	song, exists := dataMap["song"]
	if !exists {
		t.Error("Data should contain song")
	}

	songMap, ok := song.(map[string]interface{})
	if !ok {
		t.Error("Song should be an object")
	}

	if songMap["title"] != "Test Song" {
		t.Error("Found song should match search query")
	}
}

func TestSearchSongInvalidType(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodGet, "/playlist/search?type=invalid&q=test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SearchSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestSearchSongNotFound(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodGet, "/playlist/search?type=title&q=Nonexistent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SearchSong(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rec.Code)
	}
}

func TestGetSongsByRating(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add and rate a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)
	songs := handlers.engine.GetCurrentPlaylist()
	handlers.engine.RateSong(songs[0].ID, 4)

	req := httptest.NewRequest(http.MethodGet, "/playlist/rating/4", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("rating")
	c.SetParamValues("4")

	err := handlers.GetSongsByRating(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	songs_response, exists := response["songs"]
	if !exists {
		t.Error("Response should contain songs")
	}

	songsSlice, ok := songs_response.([]interface{})
	if !ok {
		t.Error("Songs should be an array")
	}

	if len(songsSlice) != 1 {
		t.Errorf("Expected 1 song with rating 4, got %d", len(songsSlice))
	}
}

func TestGetSongsByRatingInvalid(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodGet, "/playlist/rating/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("rating")
	c.SetParamValues("invalid")

	err := handlers.GetSongsByRating(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestSortPlaylist(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs in unsorted order
	handlers.engine.AddSong("Zebra", "Artist Z", "Album Z", "Rock", "Alternative", "Energetic", 300, 120)
	handlers.engine.AddSong("Alpha", "Artist A", "Album A", "Pop", "Mainstream", "Happy", 200, 110)

	requestBody := map[string]interface{}{
		"criteria":  "title",
		"algorithm": "merge",
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/playlist/sort", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SortPlaylist(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify sorting
	songs := handlers.engine.GetCurrentPlaylist()
	if songs[0].Title != "Alpha" {
		t.Error("Playlist should be sorted by title")
	}
}

func TestSortPlaylistInvalidCriteria(t *testing.T) {
	e, handlers := setupTestEcho()

	requestBody := map[string]interface{}{
		"criteria":  "invalid",
		"algorithm": "merge",
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/playlist/sort", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SortPlaylist(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestGetPlaybackHistory(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add and play songs
	handlers.engine.AddSong("Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Song 2", "Artist 2", "Album 2", "Pop", "Mainstream", "Happy", 200, 110)
	handlers.engine.PlaySong(0)
	handlers.engine.PlaySong(1)

	req := httptest.NewRequest(http.MethodGet, "/playlist/history?count=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetPlaybackHistory(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	history, exists := response["history"]
	if !exists {
		t.Error("Response should contain history")
	}

	historySlice, ok := history.([]interface{})
	if !ok {
		t.Error("History should be an array")
	}

	if len(historySlice) != 2 {
		t.Errorf("Expected 2 songs in history, got %d", len(historySlice))
	}
}

func TestGetGenres(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs with different genres
	handlers.engine.AddSong("Rock Song", "Rock Artist", "Album", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Pop Song", "Pop Artist", "Album", "Pop", "Mainstream", "Happy", 200, 110)

	req := httptest.NewRequest(http.MethodGet, "/playlist/genres", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetGenres(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	genres, exists := response["genres"]
	if !exists {
		t.Error("Response should contain genres")
	}

	genresSlice, ok := genres.([]interface{})
	if !ok {
		t.Error("Genres should be an array")
	}

	if len(genresSlice) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(genresSlice))
	}
}

func TestGetSubgenres(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs with subgenres
	handlers.engine.AddSong("Song 1", "Artist 1", "Album", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Song 2", "Artist 2", "Album", "Rock", "Classic Rock", "Epic", 280, 115)

	req := httptest.NewRequest(http.MethodGet, "/playlist/subgenres?genre=Rock", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetSubgenres(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	subgenres, exists := response["subgenres"]
	if !exists {
		t.Error("Response should contain subgenres")
	}

	subgenresSlice, ok := subgenres.([]interface{})
	if !ok {
		t.Error("Subgenres should be an array")
	}

	if len(subgenresSlice) != 2 {
		t.Errorf("Expected 2 subgenres, got %d", len(subgenresSlice))
	}
}

func TestGetMoods(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs with moods
	handlers.engine.AddSong("Song 1", "Artist 1", "Album", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Song 2", "Artist 2", "Album", "Rock", "Alternative", "Melancholic", 250, 100)

	req := httptest.NewRequest(http.MethodGet, "/playlist/moods?genre=Rock&subgenre=Alternative", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetMoods(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	moods, exists := response["moods"]
	if !exists {
		t.Error("Response should contain moods")
	}

	moodsSlice, ok := moods.([]interface{})
	if !ok {
		t.Error("Moods should be an array")
	}

	if len(moodsSlice) != 2 {
		t.Errorf("Expected 2 moods, got %d", len(moodsSlice))
	}
}

func TestGetArtists(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs with same path but different artists
	handlers.engine.AddSong("Song 1", "Artist 1", "Album", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Song 2", "Artist 2", "Album", "Rock", "Alternative", "Energetic", 250, 125)

	req := httptest.NewRequest(http.MethodGet, "/playlist/artists?genre=Rock&subgenre=Alternative&mood=Energetic", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetArtists(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	artists, exists := response["artists"]
	if !exists {
		t.Error("Response should contain artists")
	}

	artistsSlice, ok := artists.([]interface{})
	if !ok {
		t.Error("Artists should be an array")
	}

	if len(artistsSlice) != 2 {
		t.Errorf("Expected 2 artists, got %d", len(artistsSlice))
	}
}

func TestGetSongsByExplorer(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodGet, "/playlist/explorer?genre=Rock&subgenre=Alternative&mood=Energetic&artist=Test Artist", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetSongsByExplorer(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	songs, exists := response["songs"]
	if !exists {
		t.Error("Response should contain songs")
	}

	songsSlice, ok := songs.([]interface{})
	if !ok {
		t.Error("Songs should be an array")
	}

	if len(songsSlice) != 1 {
		t.Errorf("Expected 1 song, got %d", len(songsSlice))
	}
}

func TestGetRecommendations(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs and play some
	handlers.engine.AddSong("Rock Song 1", "Artist 1", "Album 1", "Rock", "Alternative", "Energetic", 240, 120)
	handlers.engine.AddSong("Rock Song 2", "Artist 2", "Album 2", "Rock", "Alternative", "Energetic", 250, 125)
	handlers.engine.PlaySong(0)

	req := httptest.NewRequest(http.MethodGet, "/playlist/recommendations?count=3", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetRecommendations(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	recommendations, exists := response["recommendations"]
	if !exists {
		t.Error("Response should contain recommendations")
	}

	recsSlice, ok := recommendations.([]interface{})
	if !ok {
		t.Error("Recommendations should be an array")
	}

	// Should return at least some recommendations
	if len(recsSlice) == 0 {
		t.Error("Should return some recommendations")
	}
}

func TestGetDashboard(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add some data for dashboard
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodGet, "/playlist/dashboard", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetDashboard(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should contain dashboard data
	_, exists := response["playlist_info"]
	if !exists {
		t.Error("Dashboard should contain playlist_info")
	}
}

func TestGetStats(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add some data for stats
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodGet, "/playlist/stats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetStats(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should contain stats
	if response["total_songs"].(float64) != 1 {
		t.Error("Stats should show correct song count")
	}
}

func TestBenchmarkSort(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs for benchmarking
	for i := 0; i < 5; i++ {
		handlers.engine.AddSong(fmt.Sprintf("Song %d", i), "Artist", "Album", "Genre", "Subgenre", "Mood", 200+i*10, 120)
	}

	req := httptest.NewRequest(http.MethodGet, "/playlist/benchmark", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.BenchmarkSort(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	benchmarks, exists := response["benchmarks"]
	if !exists {
		t.Error("Response should contain benchmarks")
	}

	benchmarksMap, ok := benchmarks.(map[string]interface{})
	if !ok {
		t.Error("Benchmarks should be a map")
	}

	// Should contain benchmark results for different algorithms
	if len(benchmarksMap) == 0 {
		t.Error("Should contain benchmark results")
	}
}

func TestClearPlaylist(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	// Verify song exists
	if handlers.engine.GetPlaylistSize() != 1 {
		t.Fatal("Song should exist before clear")
	}

	req := httptest.NewRequest(http.MethodDelete, "/playlist/clear", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.ClearPlaylist(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify playlist is cleared
	if handlers.engine.GetPlaylistSize() != 0 {
		t.Error("Playlist should be empty after clear")
	}
}

func TestSetPlaylistName(t *testing.T) {
	e, handlers := setupTestEcho()

	requestBody := map[string]interface{}{
		"name": "New Playlist Name",
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/playlist/name", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SetPlaylistName(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify name was changed
	if handlers.engine.GetPlaylistName() != "New Playlist Name" {
		t.Error("Playlist name should have been updated")
	}
}

func TestSetPlaylistNameInvalidJSON(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodPut, "/playlist/name", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SetPlaylistName(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestSetPlaylistNameEmpty(t *testing.T) {
	e, handlers := setupTestEcho()

	requestBody := map[string]interface{}{
		"name": "",
	}

	jsonData, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/playlist/name", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.SetPlaylistName(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

func TestLoadSampleData(t *testing.T) {
	e, handlers := setupTestEcho()

	req := httptest.NewRequest(http.MethodPost, "/playlist/sample", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.LoadSampleData(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify sample data was loaded
	if handlers.engine.GetPlaylistSize() == 0 {
		t.Error("Sample data should have been loaded")
	}
}

func TestGetPlaylistHTML(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add a song
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetPlaylistHTML(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Should return HTML content
	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Error("Should return HTML content type")
	}

	// Should contain some expected HTML structure
	body := rec.Body.String()
	if !strings.Contains(body, "html") {
		t.Error("Should contain HTML structure")
	}
}

func TestGetGenresHTML(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add songs with genres
	handlers.engine.AddSong("Rock Song", "Rock Artist", "Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetGenresHTML(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Should return HTML content
	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Error("Should return HTML content type")
	}
}

func TestGetDashboardHTML(t *testing.T) {
	e, handlers := setupTestEcho()

	// Add some data for dashboard
	handlers.engine.AddSong("Test Song", "Test Artist", "Test Album", "Rock", "Alternative", "Energetic", 240, 120)

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetDashboardHTML(c)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Should return HTML content
	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Error("Should return HTML content type")
	}
}

// Integration test for multiple operations
func TestHandlersIntegration(t *testing.T) {
	e, handlers := setupTestEcho()

	// Test full workflow
	// 1. Add songs
	songs := []map[string]interface{}{
		{"title": "Song 1", "artist": "Artist 1", "album": "Album 1", "genre": "Rock", "subgenre": "Alternative", "mood": "Energetic", "duration": 240, "bpm": 120},
		{"title": "Song 2", "artist": "Artist 2", "album": "Album 2", "genre": "Pop", "subgenre": "Mainstream", "mood": "Happy", "duration": 200, "bpm": 130},
		{"title": "Song 3", "artist": "Artist 3", "album": "Album 3", "genre": "Jazz", "subgenre": "Smooth", "mood": "Relaxed", "duration": 300, "bpm": 90},
	}

	for _, song := range songs {
		jsonData, _ := json.Marshal(song)
		req := httptest.NewRequest(http.MethodPost, "/playlist/songs", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handlers.AddSong(c)

		if rec.Code != http.StatusCreated {
			t.Errorf("Failed to add song: %v", song["title"])
		}
	}

	// 2. Verify playlist size
	if handlers.engine.GetPlaylistSize() != 3 {
		t.Error("Should have 3 songs in playlist")
	}

	// 3. Play a song
	req := httptest.NewRequest(http.MethodPost, "/playlist/play/0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("index")
	c.SetParamValues("0")
	handlers.PlaySong(c)

	if rec.Code != http.StatusOK {
		t.Error("Failed to play song")
	}

	// 4. Rate a song
	playlistSongs := handlers.engine.GetCurrentPlaylist()
	songID := playlistSongs[0].ID
	ratingData := map[string]interface{}{"rating": 5}
	jsonData, _ := json.Marshal(ratingData)
	req = httptest.NewRequest(http.MethodPut, "/playlist/rate/"+songID, bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("songId")
	c.SetParamValues(songID)
	handlers.RateSong(c)

	if rec.Code != http.StatusOK {
		t.Error("Failed to rate song")
	}

	// 5. Sort playlist
	sortData := map[string]interface{}{"criteria": "title", "algorithm": "merge"}
	jsonData, _ = json.Marshal(sortData)
	req = httptest.NewRequest(http.MethodPut, "/playlist/sort", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	handlers.SortPlaylist(c)

	if rec.Code != http.StatusOK {
		t.Error("Failed to sort playlist")
	}

	// 6. Get dashboard data
	req = httptest.NewRequest(http.MethodGet, "/playlist/dashboard", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	handlers.GetDashboard(c)

	if rec.Code != http.StatusOK {
		t.Error("Failed to get dashboard")
	}

	// Verify dashboard contains expected data
	var dashboardResponse map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &dashboardResponse)

	playlistInfo, exists := dashboardResponse["playlist_info"]
	if !exists {
		t.Error("Dashboard should contain playlist info")
	}

	playlistInfoMap := playlistInfo.(map[string]interface{})
	if playlistInfoMap["total_songs"].(float64) != 3 {
		t.Error("Dashboard should show correct song count")
	}
}

// Helper function to convert int to string for URL parameters
func intToString(i int) string {
	return strconv.Itoa(i)
}
