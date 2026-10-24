package cmddb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"net/http"
)

// HandleVisualizerUI returns an HTTP handler for serving the Database Visualizer UI.
func HandleVisualizerUI(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(constants.HeaderContentType, "text/html; charset=utf-8")
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
