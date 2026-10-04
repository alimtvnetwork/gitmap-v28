package cmdinstall

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// BatchInstallResponse represents the JSON response envelope for batch tool installations.
type BatchInstallResponse struct {
	IsSuccess      bool                `json:"success"`
	Tools          []ToolInstallResult `json:"tools"`
	TotalCount     int                 `json:"totalCount"`
	InstalledCount int                 `json:"installedCount"`
	FailedCount    int                 `json:"failedCount"`
	DurationMs     int64               `json:"durationMs"`
}

// ToolInstallResult records the outcome of a single tool installation in a batch.
type ToolInstallResult struct {
	Tool       string `json:"tool"`
	Status     string `json:"status"` // "installed", "already_installed", "failed"
	Version    string `json:"version,omitempty"`
	DurationMs int64  `json:"durationMs"`
	Error      string `json:"error,omitempty"`
	IsSuccess  bool   `json:"isSuccess"`
}

// ParseToolsTokens parses a comma-separated string of tool names into normalized tokens.
func ParseToolsTokens(raw string) []string {
	var tokens []string
	for _, part := range strings.Split(raw, ",") {
		clean := strings.TrimSpace(part)
		if clean != "" {
			tokens = append(tokens, clean)
		}
	}
	return tokens
}

// ResolveBatchTools resolves tools either from the --tools flag or from variadic positional arguments.
func ResolveBatchTools(toolsFlag string, posArgs []string) []string {
	var rawCandidates []string

	if strings.TrimSpace(toolsFlag) != "" {
		rawCandidates = append(rawCandidates, ParseToolsTokens(toolsFlag)...)
	}

	if len(rawCandidates) == 0 && len(posArgs) > 1 {
		rawCandidates = append(rawCandidates, posArgs...)
	}

	if len(rawCandidates) == 0 && len(posArgs) == 1 && strings.Contains(posArgs[0], ",") {
		rawCandidates = append(rawCandidates, ParseToolsTokens(posArgs[0])...)
	}

	var canonicalList []string
	seen := make(map[string]bool)
	for _, cand := range rawCandidates {
		canonical := resolveToolAlias(cand)
		if canonical != "" && !seen[canonical] {
			seen[canonical] = true
			canonicalList = append(canonicalList, canonical)
		}
	}

	return canonicalList
}

// ExecuteBatchInstall runs sequential installation of multiple tools with optional JSON output.
func ExecuteBatchInstall(opts installOptions) error {
	startTime := time.Now()
	toolList := opts.ToolList
	if len(toolList) == 0 && opts.Tool != "" {
		toolList = []string{opts.Tool}
	}

	response := BatchInstallResponse{
		IsSuccess:  true,
		Tools:      make([]ToolInstallResult, 0, len(toolList)),
		TotalCount: len(toolList),
	}

	for _, tool := range toolList {
		toolStart := time.Now()
		result := ToolInstallResult{
			Tool:      tool,
			IsSuccess: true,
		}

		if !opts.IsJson {
			fmt.Fprintf(os.Stderr, "\n%s=== Installing %s ===%s\n", constants.ColorCyan, tool, constants.ColorReset)
		}

		singleOpts := opts
		singleOpts.Tool = tool
		singleOpts.ToolList = nil

		var execErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					execErr = fmt.Errorf("panic during install of %s: %v", tool, r)
				}
			}()
			executeInstall(singleOpts)
		}()

		result.DurationMs = time.Since(toolStart).Milliseconds()

		if execErr != nil {
			result.IsSuccess = false
			result.Status = "failed"
			result.Error = execErr.Error()
			response.FailedCount++

			if !opts.HasIgnoreErrors {
				response.IsSuccess = false
				response.Tools = append(response.Tools, result)
				response.DurationMs = time.Since(startTime).Milliseconds()
				if opts.IsJson {
					return outputBatchJSON(response)
				}
				return fmt.Errorf("failed installing %s: %w", tool, execErr)
			}
		} else {
			existingVer := detectInstalledVersion(tool)
			if existingVer != "" {
				result.Status = "installed"
				result.Version = existingVer
			} else {
				result.Status = "installed"
			}
			response.InstalledCount++
		}

		response.Tools = append(response.Tools, result)
	}

	response.DurationMs = time.Since(startTime).Milliseconds()
	response.IsSuccess = response.FailedCount == 0

	if opts.IsJson {
		return outputBatchJSON(response)
	}

	return nil
}

func outputBatchJSON(resp BatchInstallResponse) error {
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
