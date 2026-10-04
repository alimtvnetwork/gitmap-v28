package cmdagy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type agyUIOptions struct {
	port      int
	noBrowser bool
}

func parseAgyUIOptions(args []string) agyUIOptions {
	opts := agyUIOptions{port: 7430, noBrowser: false}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--no-browser" || arg == "-n" {
			opts.noBrowser = true
			continue
		}
		if (arg == "--port" || arg == "-p") && i+1 < len(args) {
			if val, err := strconv.Atoi(args[i+1]); err == nil && val > 0 {
				opts.port = val
				i++
			}
		}
	}
	return opts
}

func newAgyUIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleAgyUIIndex)
	mux.HandleFunc("/api/status", handleAgyUIStatus)
	mux.HandleFunc("/api/prompts/send", handleAgyUIPromptSend)
	mux.HandleFunc("/api/prompts/enqueue", handleAgyUIPromptEnqueue)
	mux.HandleFunc("/api/prompts/history", handleAgyUIPromptHistory)
	mux.HandleFunc("/api/prompts/save", handleAgyUIPromptSave)
	return mux
}

func handleAgyUIIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(agyUIDashboardHTML))
}

func handleAgyUIStatus(w http.ResponseWriter, _ *http.Request) {
	status, err := GetAGYUIStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func handleAgyUIPromptSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var payload PromptPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	res, _ := SendPromptImmediate(payload)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func handleAgyUIPromptEnqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var payload PromptPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	rec, err := EnqueuePromptPayload(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rec)
}

func handleAgyUIPromptHistory(w http.ResponseWriter, _ *http.Request) {
	history, _ := GetPromptHistory()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(history)
}

func handleAgyUIPromptSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var rec PromptRecord
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	_ = SavePromptTemplate(rec)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func bindAgyUIListener(port int) (net.Listener, int, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err == nil {
		return ln, port, nil
	}
	dynLn, dynErr := net.Listen("tcp", "127.0.0.1:0")
	if dynErr != nil {
		return nil, 0, apperror.WrapSimple(dynErr, "bind agy ui listener")
	}
	dynPort := dynLn.Addr().(*net.TCPAddr).Port
	return dynLn, dynPort, nil
}

func printServerBanner(port int, url string) {
	fmt.Printf("\n%s● Antigravity Studio Web UI running on:%s %s\n", constants.ColorCyan, constants.ColorReset, url)
	fmt.Println("  Press Ctrl+C to terminate.")
}

// RunAgyUI starts the embedded Antigravity studio HTTP server.
func RunAgyUI(args []string) error {
	opts := parseAgyUIOptions(args)
	ln, activePort, err := bindAgyUIListener(opts.port)
	if err != nil {
		return err
	}
	defer ln.Close()

	url := fmt.Sprintf("http://127.0.0.1:%d", activePort)
	printServerBanner(activePort, url)
	if !opts.noBrowser {
		go func() {
			time.Sleep(300 * time.Millisecond)
			openSUGUIBrowser(url)
		}()
	}

	return http.Serve(ln, newAgyUIMux())
}

// IsAgyUICommand checks if the command string triggers AGY UI.
func IsAgyUICommand(cmd string) bool {
	low := strings.ToLower(cmd)
	return low == "ui" || low == "dashboard" || low == "studio"
}
