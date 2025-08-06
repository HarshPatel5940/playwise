package server

import (
	"net/http"

	"src/cmd/web"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	e.GET("/playlist", echo.WrapHandler(http.HandlerFunc(web.PlaylistDashboardHandler)))

	e.GET("/", s.HelloWorldHandler)
	e.GET("/health", s.healthHandler)

	playlistHandlers := NewPlaylistHandlers()

	api := e.Group("/api")

	playlist := api.Group("/playlist")
	{
		playlist.GET("", playlistHandlers.GetPlaylist)                             // Get current playlist
		playlist.GET("/html", playlistHandlers.GetPlaylistHTML)                    // Get current playlist as HTML for HTMX
		playlist.POST("/songs", playlistHandlers.AddSong)                          // Add song to playlist
		playlist.DELETE("/songs/:index", playlistHandlers.DeleteSong)              // Delete song by index
		playlist.PUT("/songs/:fromIndex/move/:toIndex", playlistHandlers.MoveSong) // Move song
		playlist.POST("/reverse", playlistHandlers.ReversePlaylist)                // Reverse playlist order
		playlist.DELETE("", playlistHandlers.ClearPlaylist)                        // Clear entire playlist
		playlist.PUT("/name", playlistHandlers.SetPlaylistName)                    // Update playlist name

		playlist.POST("/songs/:index/play", playlistHandlers.PlaySong) // Play song by index
		playlist.POST("/undo", playlistHandlers.UndoLastPlay)          // Undo last play

		playlist.POST("/songs/:songId/rate", playlistHandlers.RateSong)    // Rate a song
		playlist.GET("/rating/:rating", playlistHandlers.GetSongsByRating) // Get songs by rating

		playlist.GET("/search", playlistHandlers.SearchSong) // Search by ID or title

		playlist.POST("/sort", playlistHandlers.SortPlaylist) // Sort playlist

		playlist.GET("/history", playlistHandlers.GetPlaybackHistory)         // Get playback history
		playlist.GET("/recommendations", playlistHandlers.GetRecommendations) // Get smart recommendations

		playlist.GET("/stats", playlistHandlers.GetStats)          // Get playlist statistics
		playlist.GET("/benchmark", playlistHandlers.BenchmarkSort) // Benchmark sorting algorithms

		playlist.POST("/sample-data", playlistHandlers.LoadSampleData) // Load sample data for demo
	}

	explorer := api.Group("/explorer")
	{
		explorer.GET("/genres", playlistHandlers.GetGenres)                                                 // Get all genres
		explorer.GET("/genres/html", playlistHandlers.GetGenresHTML)                                        // Get all genres as HTML for HTMX
		explorer.GET("/genres/:genre/subgenres", playlistHandlers.GetSubgenres)                             // Get subgenres for genre
		explorer.GET("/genres/:genre/subgenres/:subgenre/moods", playlistHandlers.GetMoods)                 // Get moods for genre+subgenre
		explorer.GET("/genres/:genre/subgenres/:subgenre/moods/:mood/artists", playlistHandlers.GetArtists) // Get artists for genre+subgenre+mood
		explorer.GET("/songs", playlistHandlers.GetSongsByExplorer)                                         // Get songs by hierarchical path
	}

	api.GET("/dashboard", playlistHandlers.GetDashboard)          // Get comprehensive dashboard snapshot
	api.GET("/dashboard/html", playlistHandlers.GetDashboardHTML) // Get dashboard as HTML for HTMX

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
