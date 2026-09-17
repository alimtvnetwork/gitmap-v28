package cmdssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var daemonOnlyDirectives = map[string]bool{
	"authorizedkeysfile":          true,
	"authorizedkeyscommand":       true,
	"authorizedkeyscommanduser":   true,
	"authorizedprincipalsfile":    true,
	"authorizedprincipalscommand": true,
	"subsystem":                   true,
	"permitrootlogin":             true,
	"allowusers":                  true,
	"deniusers":                   true,
	"allowgroups":                 true,
	"denygroups":                  true,
	"clientaliveinterval":         true,
	"clientalivecountmax":         true,
	"strictmodes":                 true,
	"usepam":                      true,
	"maxauthtries":                true,
	"maxsessions":                 true,
	"maxstartups":                 true,
	"pidfile":                     true,
	"printmotd":                   true,
	"printlastlog":                true,
}

func isDaemonOnlyDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return false
	}
	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return false
	}
	directive := strings.ToLower(parts[0])
	directive = strings.TrimSuffix(directive, "=")
	return daemonOnlyDirectives[directive]
}

func sanitizeConfigLine(line string) (string, bool) {
	if !isDaemonOnlyDirective(line) {
		return line, false
	}
	return "# [gitmap removed server-only directive]: " + strings.TrimSpace(line), true
}

// SanitizeSSHConfigContent removes or comments out server-only directives.
func SanitizeSSHConfigContent(content string) (string, int) {
	lines := strings.Split(content, "\n")
	out := make([]string, 0, len(lines))
	removed := 0
	for _, ln := range lines {
		sanitized, isRemoved := sanitizeConfigLine(ln)
		if isRemoved {
			removed++
		}
		out = append(out, sanitized)
	}
	return strings.Join(out, "\n"), removed
}

// SanitizeSSHConfigFile cleans invalid server directives in an SSH config file.
func SanitizeSSHConfigFile(path string) (bool, *apperror.AppError) {
	if path == "" {
		path = sshConfigPath()
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, apperror.WrapSimple(err, "read ssh config for sanitize")
	}
	sanitized, count := SanitizeSSHConfigContent(string(data))
	if count == 0 {
		return false, nil
	}
	if writeErr := os.WriteFile(path, []byte(sanitized), 0o600); writeErr != nil {
		return false, apperror.WrapSimple(writeErr, "write sanitized ssh config")
	}
	fmt.Fprintf(os.Stdout, "  ✓ Sanitized SSH config at %s (removed %d server-only directive(s))\n", path, count)
	return true, nil
}

func matchesAnyFlag(arg string, flags []string) bool {
	for _, f := range flags {
		if arg == f {
			return true
		}
	}
	return false
}

func hasFlag(args []string, flags ...string) bool {
	for _, a := range args {
		if matchesAnyFlag(a, flags) {
			return true
		}
	}
	return false
}

// runSSHConfig regenerates and displays the managed SSH config block.
func runSSHConfig(args []string) error {
	hasSanitizeOnly := hasFlag(args, "--sanitize", "-s")
	if hasSanitizeOnly {
		_, err := SanitizeSSHConfigFile(sshConfigPath())
		return err
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHConfig, sshConfigPath(), err)

		return nil
	}

	defer db.Close()

	updateSSHConfig(db)

	block := buildManagedBlock(db)
	if len(block) > 0 {
		fmt.Fprint(os.Stdout, constants.MsgSSHConfigShow)
		fmt.Println(block)
	}

	return nil
}

// updateSSHConfig writes the managed block to ~/.ssh/config.
func updateSSHConfig(db *store.DB) {
	configPath := sshConfigPath()

	existing := ""
	if data, err := os.ReadFile(configPath); err == nil {
		existing = string(data)
	}

	block := buildManagedBlock(db)
	updated := replaceManagedBlock(existing, block)
	updated, _ = SanitizeSSHConfigContent(updated)

	if err := ensureSSHDir(sshDir()); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHConfig, configPath, err)

		return
	}

	if err := os.WriteFile(configPath, []byte(updated), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHConfig, configPath, err)

		return
	}

	fmt.Fprint(os.Stdout, constants.MsgSSHConfigDone)
}

// buildManagedBlock generates the managed SSH config block from DB keys.
func buildManagedBlock(db *store.DB) string {
	keys, err := db.ListSSHKeys()
	if err != nil || len(keys) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(constants.SSHConfigMarkerStart + "\n")

	for _, k := range keys {
		host := "github.com"
		if len(keys) > 1 || k.Name != constants.DefaultSSHKeyName {
			host = "github.com-" + k.Name
		}

		b.WriteString(fmt.Sprintf(constants.SSHConfigHostEntry, host, "github.com", k.PrivatePath))
		b.WriteString("\n")
	}

	b.WriteString(constants.SSHConfigMarkerEnd)

	return b.String()
}

// replaceManagedBlock replaces or appends the managed block in the config.
func replaceManagedBlock(content, block string) string {
	startIdx := strings.Index(content, constants.SSHConfigMarkerStart)
	endIdx := strings.Index(content, constants.SSHConfigMarkerEnd)

	hasMarkers := startIdx >= 0 && endIdx >= 0
	if hasMarkers {
		endIdx += len(constants.SSHConfigMarkerEnd)
		before := content[:startIdx]
		after := content[endIdx:]

		return insertBlock(before, after, block)
	}

	return appendBlock(content, block)
}

func insertBlock(before, after, block string) string {
	if len(block) == 0 {
		return strings.TrimRight(before, "\n") + after
	}

	return before + block + after
}

func appendBlock(content, block string) string {
	if len(block) == 0 {
		return content
	}

	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return content + "\n" + block + "\n"
}

// sshConfigPath returns the path to ~/.ssh/config.
func sshConfigPath() string {
	return filepath.Join(sshDir(), "config")
}

// sshDir returns the path to ~/.ssh.
func sshDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine home directory: %v\n", err)

		return ""
	}

	return filepath.Join(home, ".ssh")
}
