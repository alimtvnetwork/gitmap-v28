package cmdui

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcontextmenu"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ContextMenuStatusResponse represents context menu state for the Web UI.
type ContextMenuStatusResponse struct {
	IsActive           bool   `json:"isActive"`
	Status             string `json:"status"`
	IsCleaned          bool   `json:"isCleaned"`
	IsGitActionsHidden bool   `json:"isGitActionsHidden"`
	OS                 string `json:"os"`
	Success            bool   `json:"success"`
}

func mountContextMenuRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/context-menu/status", handleAPIContextMenuStatus)
	mux.HandleFunc("/api/context-menu/clean-powershell", handleAPIContextMenuCleanPowerShell)
	mux.HandleFunc("/api/context-menu/install", handleAPIContextMenuInstall)
	mux.HandleFunc("/api/context-menu/uninstall", handleAPIContextMenuUninstall)
	mux.HandleFunc("/api/context-menu/ide-setup", handleAPIContextMenuIDESetup)
	mux.HandleFunc("/api/context-menu/toggle-git-actions", handleAPIContextMenuToggleGitActions)
}

func handleAPIContextMenuIDESetup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method == http.MethodGet {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"target":  "all",
			"actions": []string{"status", "presets", "apply", "test-connection"},
		})

		return
	}

	handleIDESetupApply(w, r)
}

func handleAPIContextMenuToggleGitActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)

		return
	}

	isHidden := cmdinstall.IsGitActionsHidden()
	err := cmdinstall.SetGitActionsHidden(!isHidden)
	if err == nil {
		err = cmdinstall.RunInstallCtx(false)
	}
	encodeOperationResponse(w, err, "Git actions context menu visibility updated")
}

func handleAPIContextMenuStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)

	isActive := cmdinstall.IsContextMenuActive()
	statusText := "inactive"
	if isActive {
		statusText = "active"
	}

	st := cmdcontextmenu.DetectPowerShellStatus()
	resp := ContextMenuStatusResponse{
		IsActive:           isActive,
		Status:             statusText,
		IsCleaned:          st.IsCleaned,
		IsGitActionsHidden: cmdinstall.IsGitActionsHidden(),
		OS:                 runtime.GOOS,
		Success:            true,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func handleAPIContextMenuCleanPowerShell(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)

		return
	}

	err := cmdcontextmenu.CleanDuplicatePowerShell(false, false)
	encodeCleanResponse(w, err)
}

func encodeCleanResponse(w http.ResponseWriter, err error) {
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})

		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Duplicate legacy PowerShell context menu hidden successfully",
	})
}

func handleAPIContextMenuInstall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	err := cmdinstall.RunInstallCtx(false)
	encodeOperationResponse(w, err, "Context menu installed successfully")
}

func handleAPIContextMenuUninstall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	cmdinstall.RunUninstallCtx()
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Context menu uninstalled successfully",
	})
}

func encodeOperationResponse(w http.ResponseWriter, err error, msg string) {
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})

		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "message": msg})
}
