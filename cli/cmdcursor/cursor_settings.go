package cmdcursor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func getSecondaryCursorSettingsPath() string {
	if runtime.GOOS != "windows" {
		return ""
	}
	username := os.Getenv("USERNAME")
	return filepath.Join(getWindowsSystemDriveRoot(), "Users", username, "AppData", "Roaming", "Cursor", "User", "settings.json")
}

func getCursorSettingsPath() (string, error) {
	root, err := getCursorUserDataRoot()
	if err != nil {
		return "", err
	}
	primary := filepath.Join(root, "User", "settings.json")
	if isPathPresent(primary) {
		return primary, nil
	}
	sec := getSecondaryCursorSettingsPath()
	if sec != "" && isPathPresent(sec) {
		return sec, nil
	}
	return primary, nil
}

func readSettingsMap(path string) (map[string]interface{}, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return make(map[string]interface{}), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read cursor settings.json")
	}
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return make(map[string]interface{}), nil
	}
	return settings, nil
}

func viewCursorSettings() error {
	path, err := getCursorSettingsPath()
	if err != nil {
		return err
	}
	settings, err := readSettingsMap(path)
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	fmt.Printf("\n%s● Cursor Settings:%s %s\n", constants.ColorCyan, constants.ColorReset, path)
	fmt.Println(string(data))
	return nil
}

func applyCursorSettings() error {
	path, err := getCursorSettingsPath()
	if err != nil {
		return err
	}
	settings, err := readSettingsMap(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		bakPath := fmt.Sprintf("%s.bak.%d", path, time.Now().Unix())
		_ = os.Rename(path, bakPath)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir cursor user dir")
	}
	settings["workbench.colorTheme"] = "Dracula Theme"
	settings["editor.fontFamily"] = "'JetBrains Mono', 'Fira Code', Consolas, monospace"
	settings["editor.fontSize"] = 14
	settings["files.autoSave"] = "afterDelay"
	settings["files.autoSaveDelay"] = 1000
	settings["files.eol"] = "\n"
	settings["files.insertFinalNewline"] = true
	settings["files.trimTrailingWhitespace"] = true
	settings["editor.renderWhitespace"] = "selection"

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal cursor settings")
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return apperror.WrapSimple(err, "write cursor settings")
	}
	if sec := getSecondaryCursorSettingsPath(); sec != "" && sec != path {
		_ = os.MkdirAll(filepath.Dir(sec), 0755)
		_ = os.WriteFile(sec, data, 0644)
	}
	fmt.Printf("%s✔ Injected Dracula Theme and coding invariants:%s %s\n", constants.ColorGreen, constants.ColorReset, path)
	return nil
}

func syncCursorSettings(targetNode string) error {
	if targetNode == "" {
		targetNode = "u1"
	}
	path, err := getCursorSettingsPath()
	if err != nil {
		return err
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return apperror.WrapSimple(readErr, "read cursor settings.json")
	}
	fmt.Printf("%s● Synchronizing Cursor settings to remote node:%s %s -> %s\n", constants.ColorCyan, constants.ColorReset, path, targetNode)
	normalized := strings.ReplaceAll(string(data), "\r\n", "\n")
	remoteCmd := fmt.Sprintf("mkdir -p ~/.config/Cursor/User && cat << 'EOF' > ~/.config/Cursor/User/settings.json\n%s\nEOF", strings.TrimSpace(normalized))
	if errExec := cmdssh.RunSSHExec([]string{targetNode, remoteCmd}); errExec != nil {
		fmt.Printf("%s⚠ Note: Remote execution to %s: %v%s\n", constants.ColorYellow, targetNode, errExec, constants.ColorReset)
		return nil
	}
	fmt.Printf("%s✔ Synchronized Cursor settings to remote node:%s %s\n", constants.ColorGreen, constants.ColorReset, targetNode)
	return nil
}

// RunCursorSettings handles `gitmap cursor settings [view|apply|sync]`.
func RunCursorSettings(args []string) error {
	sub := "view"
	if len(args) > 0 {
		sub = strings.ToLower(args[0])
	}
	switch sub {
	case "view", "v", "show":
		return viewCursorSettings()
	case "apply", "a", "set":
		return applyCursorSettings()
	case "sync", "s":
		target := ""
		if len(args) > 1 {
			target = args[1]
		}
		return syncCursorSettings(target)
	default:
		return apperror.NewWithDetails("cursor.settings", "E1031", fmt.Sprintf("unknown settings action '%s'", sub), "cmdcursor", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}
}
