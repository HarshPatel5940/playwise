package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"src/internal/datastructures"
	"src/internal/models"
	"src/internal/services"

	"github.com/labstack/echo/v4"
)

// PlaylistHandlers contains all playlist-related HTTP handlers
type PlaylistHandlers struct {
	engine *services.PlaylistEngine
}

// NewPlaylistHandlers creates a new playlist handlers instance
func NewPlaylistHandlers() *PlaylistHandlers {
	return &PlaylistHandlers{
		engine: services.NewPlaylistEngine("My Playlist"),
	}
}

// GetPlaylist returns the current playlist
// GET /api/playlist
func (ph *PlaylistHandlers) GetPlaylist(c echo.Context) error {
	songs := ph.engine.GetCurrentPlaylist()

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"name":  ph.engine.GetPlaylistName(),
			"size":  ph.engine.GetPlaylistSize(),
			"songs": songs,
		},
	}

	return c.JSON(http.StatusOK, response)
}

// AddSong adds a new song to the playlist
// POST /api/playlist/songs
func (ph *PlaylistHandlers) AddSong(c echo.Context) error {
	// Check if it's an HTMX request
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	// Parse request body
	var req struct {
		Title    string `json:"title" validate:"required"`
		Artist   string `json:"artist" validate:"required"`
		Album    string `json:"album"`
		Genre    string `json:"genre"`
		SubGenre string `json:"subgenre"`
		Mood     string `json:"mood"`
		Duration int    `json:"duration" validate:"min=1"`
		BPM      int    `json:"bpm"`
	}

	// Handle form data for HTMX requests
	if isHTMX {
		req.Title = c.FormValue("title")
		req.Artist = c.FormValue("artist")
		req.Album = c.FormValue("album")
		req.Genre = c.FormValue("genre")
		req.SubGenre = c.FormValue("subgenre")
		req.Mood = c.FormValue("mood")
		if duration := c.FormValue("duration"); duration != "" {
			if d, err := strconv.Atoi(duration); err == nil {
				req.Duration = d
			}
		}
		if bpm := c.FormValue("bpm"); bpm != "" {
			if b, err := strconv.Atoi(bpm); err == nil {
				req.BPM = b
			}
		}
	} else {
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "Invalid request format",
			})
		}
	}

	// Validate required fields
	if req.Title == "" || req.Artist == "" {
		if isHTMX {
			return c.HTML(http.StatusBadRequest, `<div class="text-red-500">Title and Artist are required</div>`)
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Title and Artist are required",
		})
	}

	// Set default duration if not provided
	if req.Duration == 0 {
		req.Duration = 180 // 3 minutes default
	}

	// Add song to playlist
	err := ph.engine.AddSong(
		req.Title, req.Artist, req.Album,
		req.Genre, req.SubGenre, req.Mood,
		req.Duration, req.BPM,
	)

	if err != nil {
		if isHTMX {
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf(`<div class="text-red-500">Error: %s</div>`, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	if isHTMX {
		// Return updated playlist HTML
		return ph.GetPlaylistHTML(c)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Song added successfully",
	})
}

// DeleteSong removes a song from the playlist by index
// DELETE /api/playlist/songs/:index
func (ph *PlaylistHandlers) DeleteSong(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid index format",
		})
	}

	deletedSong, err := ph.engine.DeleteSong(index)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Song deleted successfully",
		"data": map[string]interface{}{
			"deleted_song":  deletedSong,
			"playlist_size": ph.engine.GetPlaylistSize(),
		},
	})
}

// MoveSong moves a song from one position to another
// PUT /api/playlist/songs/:fromIndex/move/:toIndex
func (ph *PlaylistHandlers) MoveSong(c echo.Context) error {
	fromIndexStr := c.Param("fromIndex")
	toIndexStr := c.Param("toIndex")

	fromIndex, err := strconv.Atoi(fromIndexStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid fromIndex format",
		})
	}

	toIndex, err := strconv.Atoi(toIndexStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid toIndex format",
		})
	}

	err = ph.engine.MoveSong(fromIndex, toIndex)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Song moved successfully",
	})
}

// ReversePlaylist reverses the order of songs in the playlist
// POST /api/playlist/reverse
func (ph *PlaylistHandlers) ReversePlaylist(c echo.Context) error {
	ph.engine.ReversePlaylist()

	// Check if it's an HTMX request
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	if isHTMX {
		// Return updated playlist HTML
		return ph.GetPlaylistHTML(c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Playlist reversed successfully",
	})
}

// PlaySong simulates playing a song
// POST /api/playlist/songs/:index/play
func (ph *PlaylistHandlers) PlaySong(c echo.Context) error {
	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid index format",
		})
	}

	song, err := ph.engine.PlaySong(index)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Song played successfully",
		"data": map[string]interface{}{
			"song": song,
		},
	})
}

// UndoLastPlay undoes the last played song
// POST /api/playlist/undo
func (ph *PlaylistHandlers) UndoLastPlay(c echo.Context) error {
	song, err := ph.engine.UndoLastPlay()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Last play undone successfully",
		"data": map[string]interface{}{
			"song": song,
		},
	})
}

// RateSong assigns a rating to a song
// POST /api/playlist/songs/:songId/rate
func (ph *PlaylistHandlers) RateSong(c echo.Context) error {
	songID := c.Param("songId")

	var req struct {
		Rating int `json:"rating" validate:"required,min=1,max=5"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	err := ph.engine.RateSong(songID, req.Rating)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Song rated successfully",
	})
}

// SearchSong searches for a song by ID or title
// GET /api/playlist/search
func (ph *PlaylistHandlers) SearchSong(c echo.Context) error {
	searchType := c.QueryParam("type") // "id" or "title"
	query := c.QueryParam("q")

	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Search query is required",
		})
	}

	var song *models.Song
	var err error

	switch searchType {
	case "id":
		song, err = ph.engine.SearchSongByID(query)
	case "title":
		song, err = ph.engine.SearchSongByTitle(query)
	default:
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Search type must be 'id' or 'title'",
		})
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"song": song,
		},
	})
}

// GetSongsByRating returns songs with a specific rating
// GET /api/playlist/rating/:rating
func (ph *PlaylistHandlers) GetSongsByRating(c echo.Context) error {
	ratingStr := c.Param("rating")
	rating, err := strconv.Atoi(ratingStr)
	if err != nil || rating < 1 || rating > 5 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Rating must be between 1 and 5",
		})
	}

	songs := ph.engine.GetSongsByRating(rating)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"rating": rating,
			"songs":  songs,
			"count":  len(songs),
		},
	})
}

// SortPlaylist sorts the playlist using specified criteria and algorithm
// SortPlaylist sorts the playlist by specified criteria
// POST /api/playlist/sort
func (ph *PlaylistHandlers) SortPlaylist(c echo.Context) error {
	// Check if it's an HTMX request
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	var req struct {
		Criteria  string `json:"criteria" validate:"required"`
		Algorithm string `json:"algorithm"`
	}

	if isHTMX {
		// Handle form data or URL params for HTMX requests
		req.Criteria = c.FormValue("criteria")
		if req.Criteria == "" {
			req.Criteria = c.QueryParam("criteria")
		}
		req.Algorithm = c.FormValue("algorithm")
		if req.Algorithm == "" {
			req.Algorithm = c.QueryParam("algorithm")
		}
	} else {
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "Invalid request format",
			})
		}
	}

	if req.Algorithm == "" {
		req.Algorithm = "merge" // Default to merge sort
	}

	// Map string criteria to enum
	var criteria datastructures.SortCriteria
	switch req.Criteria {
	case "title":
		criteria = datastructures.SortByTitle
	case "artist":
		criteria = datastructures.SortByArtist
	case "duration_asc":
		criteria = datastructures.SortByDurationAsc
	case "duration_desc":
		criteria = datastructures.SortByDurationDesc
	case "recently_added":
		criteria = datastructures.SortByRecentlyAdded
	case "oldest_added":
		criteria = datastructures.SortByOldestAdded
	case "rating":
		criteria = datastructures.SortByRating
	case "play_count":
		criteria = datastructures.SortByPlayCount
	default:
		if isHTMX {
			return c.HTML(http.StatusBadRequest, `<div class="text-red-500">Invalid sort criteria</div>`)
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid sort criteria",
		})
	}

	ph.engine.SortPlaylist(criteria, req.Algorithm)

	if isHTMX {
		// Return updated playlist HTML
		return ph.GetPlaylistHTML(c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Playlist sorted by %s using %s sort", req.Criteria, req.Algorithm),
	})
}

// GetPlaybackHistory returns the playback history
// GET /api/playlist/history
func (ph *PlaylistHandlers) GetPlaybackHistory(c echo.Context) error {
	countStr := c.QueryParam("count")
	count := 10 // Default count

	if countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
			count = parsedCount
		}
	}

	songs := ph.engine.GetRecentlyPlayedSongs(count)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"history": songs,
			"count":   len(songs),
		},
	})
}

// GetGenres returns all available genres
// GET /api/explorer/genres
func (ph *PlaylistHandlers) GetGenres(c echo.Context) error {
	genres := ph.engine.GetGenres()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"genres": genres,
			"count":  len(genres),
		},
	})
}

// GetSubgenres returns subgenres for a specific genre
// GET /api/explorer/genres/:genre/subgenres
func (ph *PlaylistHandlers) GetSubgenres(c echo.Context) error {
	genre := c.Param("genre")
	subgenres := ph.engine.GetSubgenres(genre)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"genre":     genre,
			"subgenres": subgenres,
			"count":     len(subgenres),
		},
	})
}

// GetMoods returns moods for a specific genre and subgenre
// GET /api/explorer/genres/:genre/subgenres/:subgenre/moods
func (ph *PlaylistHandlers) GetMoods(c echo.Context) error {
	genre := c.Param("genre")
	subgenre := c.Param("subgenre")
	moods := ph.engine.GetMoods(genre, subgenre)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"genre":    genre,
			"subgenre": subgenre,
			"moods":    moods,
			"count":    len(moods),
		},
	})
}

// GetArtists returns artists for a specific genre, subgenre, and mood
// GET /api/explorer/genres/:genre/subgenres/:subgenre/moods/:mood/artists
func (ph *PlaylistHandlers) GetArtists(c echo.Context) error {
	genre := c.Param("genre")
	subgenre := c.Param("subgenre")
	mood := c.Param("mood")
	artists := ph.engine.GetArtists(genre, subgenre, mood)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"genre":    genre,
			"subgenre": subgenre,
			"mood":     mood,
			"artists":  artists,
			"count":    len(artists),
		},
	})
}

// GetSongsByExplorer returns songs for a specific path in the explorer
// GET /api/explorer/songs
func (ph *PlaylistHandlers) GetSongsByExplorer(c echo.Context) error {
	genre := c.QueryParam("genre")
	subgenre := c.QueryParam("subgenre")
	mood := c.QueryParam("mood")
	artist := c.QueryParam("artist")

	songs := ph.engine.GetPlaylistByExplorer(genre, subgenre, mood, artist)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"path": map[string]string{
				"genre":    genre,
				"subgenre": subgenre,
				"mood":     mood,
				"artist":   artist,
			},
			"songs": songs,
			"count": len(songs),
		},
	})
}

// GetRecommendations returns smart recommendations
// GET /api/playlist/recommendations
func (ph *PlaylistHandlers) GetRecommendations(c echo.Context) error {
	countStr := c.QueryParam("count")
	count := 10 // Default count

	if countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
			count = parsedCount
		}
	}

	recommendations := ph.engine.GetSmartRecommendations(count)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"recommendations": recommendations,
			"count":           len(recommendations),
		},
	})
}

// GetDashboard returns a comprehensive dashboard snapshot
// GET /api/dashboard
func (ph *PlaylistHandlers) GetDashboard(c echo.Context) error {
	snapshot := ph.engine.ExportSnapshot()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    snapshot,
	})
}

// GetStats returns playlist statistics
// GET /api/playlist/stats
func (ph *PlaylistHandlers) GetStats(c echo.Context) error {
	stats := ph.engine.GetPlaylistStats()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// BenchmarkSort compares sorting algorithm performance
// GET /api/playlist/benchmark
func (ph *PlaylistHandlers) BenchmarkSort(c echo.Context) error {
	benchmarks := ph.engine.BenchmarkSort()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"benchmarks": benchmarks,
		},
	})
}

// ClearPlaylist removes all songs from the playlist
// DELETE /api/playlist
func (ph *PlaylistHandlers) ClearPlaylist(c echo.Context) error {
	ph.engine.ClearPlaylist()

	// Check if it's an HTMX request
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	if isHTMX {
		// Return updated playlist HTML
		return ph.GetPlaylistHTML(c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Playlist cleared successfully",
	})
}

// SetPlaylistName updates the playlist name
// PUT /api/playlist/name
func (ph *PlaylistHandlers) SetPlaylistName(c echo.Context) error {
	var req struct {
		Name string `json:"name" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	if strings.TrimSpace(req.Name) == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Playlist name cannot be empty",
		})
	}

	ph.engine.SetPlaylistName(req.Name)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Playlist name updated successfully",
		"data": map[string]interface{}{
			"name": req.Name,
		},
	})
}

// LoadSampleData loads sample songs into the playlist for demonstration
// POST /api/playlist/sample-data
func (ph *PlaylistHandlers) LoadSampleData(c echo.Context) error {
	// Check if it's an HTMX request
	isHTMX := c.Request().Header.Get("HX-Request") == "true"

	// Clear existing playlist first
	ph.engine.ClearPlaylist()

	// Load sample data
	sampleLoader := services.NewSampleDataLoader()
	err := sampleLoader.LoadSampleData(ph.engine)
	if err != nil {
		if isHTMX {
			return c.HTML(http.StatusInternalServerError, fmt.Sprintf(`<div class="text-red-500">Failed to load sample data: %s</div>`, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Failed to load sample data: " + err.Error(),
		})
	}

	if isHTMX {
		// Return updated playlist HTML
		return ph.GetPlaylistHTML(c)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sample data loaded successfully",
		"data": map[string]interface{}{
			"songsLoaded": ph.engine.GetPlaylistSize(),
		},
	})
}

// HTMX Handlers - Return HTML fragments instead of JSON

// GetPlaylistHTML returns the playlist as HTML for HTMX
func (ph *PlaylistHandlers) GetPlaylistHTML(c echo.Context) error {
	songs := ph.engine.GetCurrentPlaylist()

	if len(songs) == 0 {
		html := `
		<div class="text-center py-8 text-gray-500">
			<p class="mb-4">Your playlist is empty</p>
			<button onclick="loadSampleData()" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
				📦 Load Sample Data
			</button>
		</div>`
		return c.HTML(http.StatusOK, html)
	}

	var html strings.Builder
	for i, song := range songs {
		html.WriteString(fmt.Sprintf(`
		<div class="playlist-item bg-gray-50 p-3 rounded-lg border mb-2" data-index="%d">
			<div class="flex justify-between items-start">
				<div class="flex-1 min-w-0">
					<div class="flex items-center gap-2 mb-1">
						<h4 class="font-semibold text-gray-800 truncate">%s</h4>
						<span class="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">%s</span>
					</div>
					<p class="text-gray-600 text-sm mb-1">%s%s</p>
					<div class="flex flex-wrap gap-2 text-xs text-gray-500">
						<span>%s</span>
						%s
						%s
						<span>• %d:%02d</span>
						%s
					</div>
					%s
				</div>
				<div class="flex flex-col gap-1 ml-4">
					<button
						hx-post="/api/playlist/songs/%d/play"
						hx-target="#history-container"
						class="bg-green-500 hover:bg-green-600 text-white px-2 py-1 rounded text-xs">
						▶️ Play
					</button>
					<button
						hx-delete="/api/playlist/songs/%d"
						hx-target="#playlist-container"
						hx-confirm="Delete this song?"
						class="bg-red-500 hover:bg-red-600 text-white px-2 py-1 rounded text-xs">
						🗑️
					</button>
				</div>
			</div>
		</div>`,
			i,
			song.Title,
			song.ID,
			song.Artist,
			func() string {
				if song.Album != "" {
					return " • " + song.Album
				}
				return ""
			}(),
			song.Genre,
			func() string {
				if song.SubGenre != "" {
					return "<span>• " + song.SubGenre + "</span>"
				}
				return ""
			}(),
			func() string {
				if song.Mood != "" {
					return "<span>• " + song.Mood + "</span>"
				}
				return ""
			}(),
			song.Duration/60, song.Duration%60,
			func() string {
				if song.BPM > 0 {
					return fmt.Sprintf("<span>• %d BPM</span>", song.BPM)
				}
				return ""
			}(),
			func() string {
				if song.Rating > 0 {
					return fmt.Sprintf(`<div class="mt-1">%s</div>`, strings.Repeat("⭐", song.Rating))
				}
				return ""
			}(),
			i,
			i,
		))
	}

	return c.HTML(http.StatusOK, html.String())
}

// GetGenresHTML returns genres as HTML for HTMX
func (ph *PlaylistHandlers) GetGenresHTML(c echo.Context) error {
	genres := ph.engine.GetGenres()

	if len(genres) == 0 {
		return c.HTML(http.StatusOK, `<div class="text-gray-500 text-sm">No genres available</div>`)
	}

	var html strings.Builder
	for _, genre := range genres {
		html.WriteString(fmt.Sprintf(`
		<button
			hx-get="/api/explorer/genres/%s/subgenres-html"
			hx-target="#subgenres-list"
			class="block w-full text-left px-2 py-1 rounded hover:bg-gray-100 text-sm">
			%s
		</button>`, genre, genre))
	}

	return c.HTML(http.StatusOK, html.String())
}

// GetDashboardHTML returns dashboard stats as HTML for HTMX
func (ph *PlaylistHandlers) GetDashboardHTML(c echo.Context) error {
	snapshot := ph.engine.ExportSnapshot()

	// Extract data from snapshot structure
	playlistInfo := snapshot["playlist_info"].(map[string]interface{})
	totalSongs := playlistInfo["total_songs"].(int)
	totalDuration := playlistInfo["total_duration"].(int)

	// Get additional stats
	stats := ph.engine.GetPlaylistStats()
	uniqueArtists := stats["unique_artists"].(int)

	// Get genre count from genre stats
	genreStats := snapshot["genre_stats"].(map[string]interface{})
	totalGenres := len(genreStats)

	html := fmt.Sprintf(`
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
		<div class="bg-gradient-to-r from-blue-500 to-blue-600 text-white p-4 sm:p-6 rounded-lg">
			<h3 class="text-lg font-semibold mb-2">Total Songs</h3>
			<div class="text-3xl font-bold">%d</div>
		</div>
		<div class="bg-gradient-to-r from-green-500 to-green-600 text-white p-4 sm:p-6 rounded-lg">
			<h3 class="text-lg font-semibold mb-2">Total Duration</h3>
			<div class="text-3xl font-bold">%d:%02d</div>
		</div>
		<div class="bg-gradient-to-r from-purple-500 to-purple-600 text-white p-4 sm:p-6 rounded-lg">
			<h3 class="text-lg font-semibold mb-2">Unique Artists</h3>
			<div class="text-3xl font-bold">%d</div>
		</div>
		<div class="bg-gradient-to-r from-orange-500 to-orange-600 text-white p-4 sm:p-6 rounded-lg">
			<h3 class="text-lg font-semibold mb-2">Genres</h3>
			<div class="text-3xl font-bold">%d</div>
		</div>
	</div>`,
		totalSongs,
		totalDuration/60, totalDuration%60,
		uniqueArtists,
		totalGenres,
	)

	return c.HTML(http.StatusOK, html)
}
