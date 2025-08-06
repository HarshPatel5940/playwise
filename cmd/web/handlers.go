package web

import (
	"net/http"
)

// PlaylistDashboardHandler serves the playlist dashboard page
func PlaylistDashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type to HTML
	w.Header().Set("Content-Type", "text/html")

	// Render the playlist dashboard template
	component := PlaylistDashboard()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
