package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runSSHConfig regenerates and displays the managed SSH config block.
func runSSHConfig() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHConfig, sshConfigPath(), err)

		return
	}
	defer db.Close()

	updateSSHConfig(db)

	block := buildManagedBlock(db)
	if len(block) > 0 {
		fmt.Fprint(os.Stdout, constants.MsgSSHConfigShow)
		fmt.Println(block)
	}
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

	if startIdx >= 0 && endIdx >= 0 {
		endIdx += len(constants.SSHConfigMarkerEnd)
		before := content[:startIdx]
		after := content[endIdx:]

		if len(block) == 0 {
			return strings.TrimRight(before, "\n") + after
		}

		return before + block + after
	}

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
