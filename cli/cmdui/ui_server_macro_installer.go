package cmdui

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func mountMacroAndInstallerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/macros", handleAPIMacros)
	mux.HandleFunc("/api/macros/save", handleAPIMacroSave)
	mux.HandleFunc("/api/macros/delete", handleAPIMacroDelete)
	mux.HandleFunc("/api/installers/profiles", handleAPIInstallerProfiles)
}

func handleAPIMacros(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	res := macro.ListMacros()
	if !res.IsSuccess() {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": res.AppError().Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "macros": res.Data})
}

func handleAPIMacroSave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var m macro.Macro
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	if strings.TrimSpace(m.Name) == "" {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "macro name required"})
		return
	}

	if err := macro.SaveMacro(&m); err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "name": m.Name})
}

func handleAPIMacroDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": "name parameter required"})
		return
	}

	if err := macro.DeleteMacro(name); err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "name": name})
}

// ProfileTreeDTO formats an InstallProfile for UI display.
type ProfileTreeDTO struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`
	IsInstalled bool     `json:"isInstalled"`
}

func handleAPIInstallerProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)

	profiles := cmdinstall.AllInstallProfiles()
	var dtos []ProfileTreeDTO
	for _, p := range profiles {
		dtos = append(dtos, ProfileTreeDTO{
			Name:        p.Name,
			Title:       p.Title,
			Description: p.Description,
			Tools:       p.Tools,
		})
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"profiles": dtos,
	})
}
