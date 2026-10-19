package cmddb

import (
	"net/http"
)

// HandleVisualizerUI returns an HTTP handler for serving the Database Visualizer UI.
func HandleVisualizerUI(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(VisualizerUIHtml))
	}
}

// RegisterVisualizerRoutes registers the Visualizer UI routes onto the provided ServeMux.
func RegisterVisualizerRoutes(mux *http.ServeMux, mgr *MultiDbManager) {
	handler := HandleVisualizerUI(mgr)
	mux.HandleFunc("GET /ui", handler)
	mux.HandleFunc("GET /ui/", handler)
}
