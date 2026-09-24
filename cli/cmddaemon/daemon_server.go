package cmddaemon

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var daemonStartTime = time.Now()

// RunDaemon starts the authenticated REST listener.
func RunDaemon(port int, expectedToken string) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return apperror.WrapSimple(err, "daemon_listen")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(DaemonStatusResp{
			Status:    "online",
			Version:   constants.Version,
			OS:        runtime.GOOS,
			UptimeSec: int64(time.Since(daemonStartTime).Seconds()),
		})
	})

	mux.HandleFunc("/api/v1/exec", func(w http.ResponseWriter, r *http.Request) {
		handleDaemonExec(w, r, expectedToken)
	})

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("⚡ GitMap REST Triad Daemon active on port %d\n", port)

	return server.Serve(listener)
}

func handleDaemonExec(w http.ResponseWriter, r *http.Request, expectedToken string) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)

		return
	}

	if !isTokenAuthorized(r, expectedToken) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	var req DaemonExecReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = json.NewEncoder(w).Encode(DaemonExecResp{Status: "error", Error: err.Error()})

		return
	}

	resp := executeDaemonCommand(req)
	_ = json.NewEncoder(w).Encode(resp)
}

func isTokenAuthorized(r *http.Request, expectedToken string) bool {
	isTokenOptional := expectedToken == ""
	if isTokenOptional {
		return true
	}

	return r.Header.Get("X-GitMap-Token") == expectedToken
}

func executeDaemonCommand(req DaemonExecReq) DaemonExecResp {
	start := time.Now()
	cmd := buildDaemonCmd(req.Command, req.Args)
	out, err := cmd.CombinedOutput()
	dur := time.Since(start).Milliseconds()

	if err != nil {
		return DaemonExecResp{
			Status:     "error",
			Success:    false,
			ExitCode:   1,
			Stdout:     string(out),
			DurationMs: dur,
			Error:      err.Error(),
		}
	}

	return DaemonExecResp{
		Status:     "success",
		Success:    true,
		ExitCode:   0,
		Stdout:     string(out),
		DurationMs: dur,
	}
}

func buildDaemonCmd(cmdStr string, args []string) *exec.Cmd {
	hasArgs := len(args) > 0
	if hasArgs {
		return exec.Command(cmdStr, args...)
	}

	if runtime.GOOS == "windows" {
		return exec.Command("cmd.exe", "/c", cmdStr)
	}

	return exec.Command("sh", "-c", cmdStr)
}
