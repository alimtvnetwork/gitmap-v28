// Package cmd — ai_memory_server.go provides an ultra-fast in-memory HTTP responder across 4 unique candidate ports (47831..47834) and AUM CLI dispatch.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/searcher"
)

// UniqueAICandidatePorts defines the 4 dedicated unique ports for the GitMap In-Memory AI Ping/Telemetry Server.
var UniqueAICandidatePorts = []int{47831, 47832, 47833, 47834}

// StartInMemoryAIServerOnCandidatePorts binds an HTTP server to the first free port in UniqueAICandidatePorts.
func StartInMemoryAIServerOnCandidatePorts(ports []int) (net.Listener, int, *http.Server, error) {
	if len(ports) == 0 {
		ports = UniqueAICandidatePorts
	}
	var ln net.Listener
	var boundPort int
	var err error
	for _, p := range ports {
		ln, err = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err == nil {
			boundPort = ln.Addr().(*net.TCPAddr).Port
			break
		}
	}
	if ln == nil {
		ln, err = net.Listen("tcp", "127.0.0.1:0")
	}
	if err != nil {
		return nil, 0, nil, err
	}
	boundPort = ln.Addr().(*net.TCPAddr).Port
	srv := &http.Server{
		Handler:           buildAIMemoryMux(boundPort),
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		_ = srv.Serve(ln)
	}()
	return ln, boundPort, srv, nil
}

func buildAIMemoryMux(boundPort int) *http.ServeMux {
	bootTime := time.Now()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/ai/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]any{
			"status":          "ok",
			"engine":          "gitmap-ai-memory-v1",
			"bound_port":      boundPort,
			"candidate_ports": UniqueAICandidatePorts,
			"version":         constants.Version,
			"uptime_ms":       time.Since(bootTime).Milliseconds(),
			"mode":            "zero-disk-in-memory",
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/api/v1/ai/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		topSearches, _ := searcher.ListTopAUMSearches(5)
		resp := map[string]any{
			"status":              "ready",
			"bound_port":          boundPort,
			"candidate_ports":     UniqueAICandidatePorts,
			"version":             constants.Version,
			"aum_hot_cache_count": len(topSearches),
			"top_dh2d_queries":    topSearches,
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	mux.HandleFunc("/api/v1/ai/search-cache", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if q != "" {
			res, dh2d, isHot := searcher.LookupHotCachedSearch(q, "keyword")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"query":       q,
				"dh2d_id":     dh2d,
				"is_hot":      isHot,
				"match_count": len(res),
				"results":     res,
			})
			return
		}
		topSearches, _ := searcher.ListTopAUMSearches(20)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  "ok",
			"entries": topSearches,
		})
	})
	return mux
}

func runAIMemoryServerCmd(args []string) error {
	if len(args) > 0 && (args[0] == "search-history" || args[0] == "history" || args[0] == "dh2d") {
		return searcher.RenderAUMSearchHistoryTable(25)
	}
	ln, port, srv, err := StartInMemoryAIServerOnCandidatePorts(UniqueAICandidatePorts)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
		_ = ln.Close()
	}()

	url := fmt.Sprintf("http://127.0.0.1:%d/api/v1/ai/ping", port)
	resp, getErr := http.Get(url)
	if getErr != nil {
		return getErr
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("⚡ GitMap In-Memory AI Server active on port %d (candidates: %v)\n", port, UniqueAICandidatePorts)
	fmt.Println(string(body))
	return nil
}
