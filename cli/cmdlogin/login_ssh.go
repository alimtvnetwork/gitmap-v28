// Package cmdlogin — login_ssh.go handles SSH key generation, clipboard copying, and binding.
package cmdlogin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

type sshLoginOptions struct {
	alias     string
	customDir string
	yes       bool
}

// RunLoginSSH generates or reuses an Ed25519 SSH key, copies the public key to clipboard,
// displays provider links, registers the profile, and configures repo core.sshCommand.
func RunLoginSSH(args []string) error {
	opts := parseSSHArgs(args)

	sshDir, keyPath, pubPath, err := resolveSSHKeyPaths(opts.customDir, opts.alias)
	if err != nil {
		return err
	}

	if err := ensureSSHKey(sshDir, keyPath, pubPath, opts.alias, opts.yes); err != nil {
		return err
	}

	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return apperror.WrapSimple(err, "read public key")
	}

	pubKey := strings.TrimSpace(string(pubBytes))
	_, clipErr := copyToClipboard(pubKey)

	printSSHPublicKeyBanner(pubKey, clipErr == nil)

	if err := recordSSHProfile(opts.alias, keyPath); err != nil {
		fmt.Printf("  %s! Warning: could not update git profiles: %v%s\n", constants.ColorYellow, err, constants.ColorReset)
	}

	bindRepoSSHCommandIfApplicable(opts.alias, keyPath)

	return nil
}

func parseSSHArgs(args []string) sshLoginOptions {
	opts := sshLoginOptions{}
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "ssh" {
			i++
			continue
		}
		if arg == "-y" || arg == "--yes" || arg == "-f" || arg == "--force" {
			opts.yes = true
			i++
			continue
		}
		if (arg == "-d" || arg == "--dir") && i+1 < len(args) {
			opts.customDir = args[i+1]
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "--dir=") {
			opts.customDir = strings.TrimPrefix(arg, "--dir=")
			i++
			continue
		}
		if !strings.HasPrefix(arg, "-") && len(opts.alias) == 0 {
			opts.alias = strings.TrimSpace(arg)
			i++
			continue
		}
		i++
	}
	return opts
}

func resolveSSHKeyPaths(customDir, alias string) (string, string, string, error) {
	sshDir := customDir
	if len(sshDir) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", "", "", apperror.WrapSimple(err, "resolve user home directory")
		}
		sshDir = filepath.Join(homeDir, ".ssh")
	}

	keyFileName := "id_ed25519"
	if len(alias) > 0 {
		keyFileName = "id_ed25519_" + alias
	}

	keyPath := filepath.Join(sshDir, keyFileName)
	pubPath := keyPath + ".pub"
	return sshDir, keyPath, pubPath, nil
}

func ensureSSHKey(sshDir, keyPath, pubPath, alias string, overwrite bool) error {
	_, err := os.Stat(keyPath)
	keyExists := err == nil

	if keyExists && !overwrite {
		fmt.Printf("%s✔ Reusing existing SSH key:%s %s\n", constants.ColorCyan, constants.ColorReset, keyPath)
		return nil
	}

	if keyExists && overwrite {
		_ = os.Remove(keyPath)
		_ = os.Remove(pubPath)
	}

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return apperror.WrapSimple(err, "create ssh directory")
	}

	comment := "gitmap"
	if len(alias) > 0 {
		comment = alias
	}

	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", comment, "-f", keyPath, "-N", "")
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.WrapSimple(err, "ssh-keygen failed").WithContext("output", strings.TrimSpace(string(out)))
	}

	fmt.Printf("%s✔ Generated modern Ed25519 SSH key:%s %s\n", constants.ColorGreen, constants.ColorReset, keyPath)
	return nil
}

func copyToClipboard(text string) (string, error) {
	tool, args := resolveClipboardCommand()
	if len(tool) == 0 {
		return "", nil
	}

	cmd := exec.Command(tool, args...)
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return tool, err
	}

	return tool, nil
}

func resolveClipboardCommand() (string, []string) {
	switch runtime.GOOS {
	case "windows":
		return "clip", nil
	case "darwin":
		return "pbcopy", nil
	}

	if path, err := exec.LookPath("wl-copy"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("xclip"); err == nil {
		return path, []string{"-selection", "clipboard"}
	}
	if path, err := exec.LookPath("xsel"); err == nil {
		return path, []string{"--clipboard", "--input"}
	}
	return "", nil
}

func printSSHPublicKeyBanner(pubKey string, copied bool) {
	fmt.Println()
	fmt.Println(constants.ColorCyan + "================================================================================" + constants.ColorReset)
	fmt.Println(constants.ColorCyan + "  SSH Public Key:" + constants.ColorReset)
	fmt.Println("  " + pubKey)
	fmt.Println(constants.ColorCyan + "================================================================================" + constants.ColorReset)
	if copied {
		fmt.Printf("  %s✔ Copied public key to clipboard!%s\n\n", constants.ColorGreen, constants.ColorReset)
	} else {
		fmt.Printf("  %s! Could not automatically copy to clipboard; copy manually above.%s\n\n", constants.ColorYellow, constants.ColorReset)
	}

	fmt.Println("  Add this key to your Git provider:")
	fmt.Printf("    • GitHub: %shttps://github.com/settings/ssh/new%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • GitLab: %shttps://gitlab.com/-/user_settings/ssh_keys%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
}

func recordSSHProfile(alias, keyPath string) error {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return err
	}

	profileName := resolveSSHProfileName(alias, cfg)
	keyFileName := filepath.Base(keyPath)

	updateOrAppendSSHProfile(&cfg, profileName, keyPath, keyFileName)
	return store.SaveGitProfiles(cfg)
}

func resolveSSHProfileName(alias string, cfg model.GitProfileConfig) string {
	if len(alias) > 0 {
		return alias
	}
	if len(cfg.Active) > 0 {
		return cfg.Active
	}
	return "default"
}

func updateOrAppendSSHProfile(cfg *model.GitProfileConfig, name, keyPath, keyName string) {
	for idx := range cfg.Profiles {
		if strings.EqualFold(cfg.Profiles[idx].Name, name) || cfg.Profiles[idx].ID == name {
			cfg.Profiles[idx].AuthMethod = "ssh"
			cfg.Profiles[idx].SSHKeyPath = keyPath
			cfg.Profiles[idx].SSHKeyName = keyName
			cfg.Profiles[idx].LastUsedAt = time.Now()
			return
		}
	}

	newProf := model.GitProfile{
		ID:         name,
		Name:       name,
		Provider:   "github",
		Type:       "user",
		AuthMethod: "ssh",
		SSHKeyPath: keyPath,
		SSHKeyName: keyName,
		IsDefault:  len(cfg.Profiles) == 0,
		UsageCount: 1,
		LastUsedAt: time.Now(),
	}
	cfg.Profiles = append(cfg.Profiles, newProf)
	if len(cfg.Active) == 0 {
		cfg.Active = name
	}
}

func bindRepoSSHCommandIfApplicable(alias, keyPath string) {
	if exec.Command("git", "rev-parse", "--is-inside-work-tree").Run() != nil {
		return
	}

	cleanPath := filepath.ToSlash(keyPath)
	cmdStr := fmt.Sprintf("ssh -i %s -F /dev/null -o IdentitiesOnly=yes", cleanPath)
	if err := exec.Command("git", "config", "core.sshCommand", cmdStr).Run(); err == nil {
		fmt.Printf("  %s✔ Configured core.sshCommand for current repository.%s\n", constants.ColorGreen, constants.ColorReset)
	}

	if len(alias) > 0 {
		_ = usercontext.BindProject(".", alias)
	}
}
