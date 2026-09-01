package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/completion"
)

type termSession struct {
	cmdExec *exec.Cmd
	stdIn   io.WriteCloser
	stdOut  io.ReadCloser
	termMu  sync.Mutex
}

var termSessionsMap sync.Map

func mountTerminalHandlers(muxServer *http.ServeMux) {
	muxServer.HandleFunc("/api/terminal/stream", termStreamHandler)
	muxServer.HandleFunc("/api/terminal/input", termInputHandler)
	muxServer.HandleFunc("/api/terminal/autocomplete", termAutocompleteHandler)
	muxServer.HandleFunc("/api/command/exec", termCommandExecHandler)
}

func startTerminal(sessionID string) (*termSession, error) {
	safeID := getSafeSessionID(sessionID)
	existingSess, hasExisting := termSessionsMap.Load(safeID)
	if hasExisting {
		return castTerminal(existingSess)
	}
	return createTerminal(safeID)
}

func getSafeSessionID(sessionID string) string {
	hasSession := sessionID != ""
	if hasSession {
		return sessionID
	}
	return "default"
}

func castTerminal(existingSess any) (*termSession, error) {
	termSess, isValid := existingSess.(*termSession)
	if isValid {
		return termSess, nil
	}
	return nil, apperror.NewSimple("term_cast", "invalid_type")
}

func createTerminal(sessionID string) (*termSession, error) {
	newTerm, startErr := newTerminal()
	if startErr != nil {
		return nil, startErr
	}
	termSessionsMap.Store(sessionID, newTerm)
	return newTerm, nil
}

func newTerminal() (*termSession, error) {
	shellPath := getShell()
	cmdExec := exec.Command(shellPath)
	stdIn, stdOut, pipeErr := attachPipes(cmdExec)
	if pipeErr != nil {
		return nil, pipeErr
	}
	if startErr := cmdExec.Start(); startErr != nil {
		return nil, apperror.WrapSimple(startErr, "term_start")
	}
	return &termSession{cmdExec: cmdExec, stdIn: stdIn, stdOut: stdOut}, nil
}

func attachPipes(cmdExec *exec.Cmd) (io.WriteCloser, io.ReadCloser, error) {
	stdIn, inErr := cmdExec.StdinPipe()
	if inErr != nil {
		return nil, nil, apperror.WrapSimple(inErr, "term_stdin")
	}
	stdOut, outErr := cmdExec.StdoutPipe()
	if outErr != nil {
		return nil, nil, apperror.WrapSimple(outErr, "term_stdout")
	}
	cmdExec.Stderr = cmdExec.Stdout
	return stdIn, stdOut, nil
}

func getShell() string {
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		return "cmd.exe"
	}
	shellPath := os.Getenv("SHELL")
	hasShell := shellPath != ""
	if hasShell {
		return shellPath
	}
	return "/bin/sh"
}

func termStreamHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	sessionID := httpRequest.URL.Query().Get("session_id")
	termSess, startErr := startTerminal(sessionID)
	if startErr != nil {
		http.Error(httpWriter, startErr.Error(), http.StatusInternalServerError)
		return
	}
	setStreamHeaders(httpWriter)
	httpFlusher, isFlusher := httpWriter.(http.Flusher)
	if !isFlusher {
		return
	}
	streamOutput(termSess, httpWriter, httpFlusher)
}

func setStreamHeaders(httpWriter http.ResponseWriter) {
	httpWriter.Header().Set("Content-Type", "text/event-stream")
	httpWriter.Header().Set("Cache-Control", "no-cache")
	httpWriter.Header().Set("Connection", "keep-alive")
}

func streamOutput(termSess *termSession, httpWriter io.Writer, httpFlusher http.Flusher) {
	termBuf := make([]byte, 1024)
	for {
		readBytes, readErr := termSess.stdOut.Read(termBuf)
		hasBytes := readBytes > 0
		if hasBytes {
			writeStreamChunk(termBuf[:readBytes], httpWriter, httpFlusher)
		}
		if readErr != nil {
			break
		}
	}
}

func writeStreamChunk(chunkBytes []byte, httpWriter io.Writer, httpFlusher http.Flusher) {
	outputMsg := string(chunkBytes)
	cleanMsg := strings.ReplaceAll(outputMsg, "\n", "\\n")
	fmt.Fprintf(httpWriter, "data: %s\n\n", cleanMsg)
	httpFlusher.Flush()
}

func termInputHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	sessionID := httpRequest.URL.Query().Get("session_id")
	termSess, startErr := startTerminal(sessionID)
	if startErr != nil {
		http.Error(httpWriter, startErr.Error(), http.StatusInternalServerError)
		return
	}
	reqBody, readErr := readInputBody(httpRequest)
	if readErr != nil {
		http.Error(httpWriter, readErr.Error(), http.StatusBadRequest)
		return
	}
	writeErr := writeTermInput(termSess, reqBody)
	if writeErr != nil {
		http.Error(httpWriter, writeErr.Error(), http.StatusInternalServerError)
		return
	}
	httpWriter.WriteHeader(http.StatusOK)
}

func readInputBody(httpRequest *http.Request) ([]byte, error) {
	reqBody, readErr := io.ReadAll(httpRequest.Body)
	if readErr != nil {
		return nil, apperror.WrapSimple(readErr, "term_input_read")
	}
	return reqBody, nil
}

func writeTermInput(termSess *termSession, reqBody []byte) error {
	termSess.termMu.Lock()
	defer termSess.termMu.Unlock()
	_, writeErr := termSess.stdIn.Write(reqBody)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "term_input_write")
	}
	return nil
}

func termAutocompleteHandler(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	queryText := httpRequest.URL.Query().Get("q")
	completions := []string{}
	isGitmap := strings.HasPrefix(queryText, "gitmap ")
	if isGitmap {
		completions = completion.AllCommands()
	}
	jsonBytes, _ := json.Marshal(completions)
	httpWriter.Header().Set("Content-Type", "application/json")
	httpWriter.Write(jsonBytes)
}

type commandExecReq struct {
	Command string `json:"command"`
}

type commandExecResp struct {
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	ExitCode int    `json:"exitCode"`
	Error    string `json:"error,omitempty"`
}

func termCommandExecHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	var req commandExecReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp := executeCLIForAPI(req.Command)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func executeCLIForAPI(cmdStr string) commandExecResp {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	out, err := cmd.CombinedOutput()
	if err == nil {
		return commandExecResp{Success: true, Output: string(out), ExitCode: 0}
	}
	exitCode := 1
	exitErr, ok := err.(*exec.ExitError)
	if ok {
		exitCode = exitErr.ExitCode()
	}
	return commandExecResp{
		Success:  false,
		Output:   string(out),
		ExitCode: exitCode,
		Error:    err.Error(),
	}
}
