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
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/enums"
	appconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

type sshLoginOptions struct {
	alias       string
	customDir   string
	isOverwrite bool
}

// RunLoginSSH generates or reuses an Ed25519 SSH key, copies the public key to clipboard,
// displays provider links, registers the profile, and configures repo core.sshCommand.
func RunLoginSSH(args []string) error {
	opts := parseSSHArgs(args)

	sshDir, keyPath, pubPath, err := resolveSSHKeyPaths(opts.customDir, opts.alias)
	if err != nil {
		return err
	}

	if err := ensureSSHKey(sshDir, keyPath, pubPath, opts.alias, opts.isOverwrite); err != nil {
		return err
	}

	return completeSSHLoginFlow(opts.alias, keyPath, pubPath)
}

func completeSSHLoginFlow(alias, keyPath, pubPath string) error {
	if err := copyAndDisplaySSHKey(pubPath); err != nil {
		return err
	}

	if err := recordSSHProfile(alias, keyPath); err != nil {
		fmt.Printf("  %s! Warning: could not update git profiles: %v%s\n", appconstants.ColorYellow, err, appconstants.ColorReset)
	}

	bindRepoSSHCommandIfApplicable(alias, keyPath)

	return nil
}

func copyAndDisplaySSHKey(pubPath string) error {
	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return apperror.WrapSimple(err, "read public key")
	}

	pubKey := strings.TrimSpace(string(pubBytes))
	_, clipErr := copyToClipboard(pubKey)
	isCopied := clipErr == nil

	printSSHPublicKeyBanner(pubKey, isCopied)

	return nil
}

func parseSSHArgs(args []string) sshLoginOptions {
	opts := sshLoginOptions{}
	for i := 0; i < len(args); i++ {
		step, isParsed := applySSHOption(args, i, &opts)
		if isParsed {
			i += step
			continue
		}

		if !strings.HasPrefix(args[i], "-") && len(opts.alias) == 0 {
			opts.alias = strings.TrimSpace(args[i])
		}
	}

	return opts
}

func applySSHOption(args []string, i int, opts *sshLoginOptions) (int, bool) {
	arg := args[i]
	if isSSHCommandOrForceFlag(arg, opts) {
		return 0, true
	}

	return parseSSHDirFlag(args, i, opts)
}

func isSSHCommandOrForceFlag(arg string, opts *sshLoginOptions) bool {
	switch arg {
	case constants.SubCmdLoginSSH:
		return true
	case constants.FlagYesShortArg, constants.FlagYesArg, constants.FlagForceShortArg, constants.FlagForceArg:
		opts.isOverwrite = true
		return true
	default:
		return false
	}
}

func parseSSHDirFlag(args []string, i int, opts *sshLoginOptions) (int, bool) {
	arg := args[i]
	if (arg == constants.FlagDirShortArg || arg == constants.FlagDirArg) && i+1 < len(args) {
		opts.customDir = args[i+1]
		return 1, true
	}

	if strings.HasPrefix(arg, constants.FlagDirPrefix) {
		opts.customDir = strings.TrimPrefix(arg, constants.FlagDirPrefix)
		return 0, true
	}

	return 0, false
}

func resolveSSHKeyPaths(customDir, alias string) (string, string, string, error) {
	sshDir, err := resolveSSHDir(customDir)
	if err != nil {
		return "", "", "", err
	}

	keyFileName := resolveSSHKeyFileName(alias)
	keyPath := filepath.Join(sshDir, keyFileName)
	pubPath := keyPath + constants.PubExt

	return sshDir, keyPath, pubPath, nil
}

func resolveSSHDir(customDir string) (string, error) {
	if len(customDir) > 0 {
		return customDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve user home directory")
	}

	return filepath.Join(homeDir, constants.DefaultSSHDir), nil
}

func resolveSSHKeyFileName(alias string) string {
	if len(alias) > 0 {
		return constants.DefaultSSHKeyPrefix + alias
	}

	return constants.DefaultSSHKeyBaseName
}

func ensureSSHKey(sshDir, keyPath, pubPath, alias string, isOverwrite bool) error {
	_, err := os.Stat(keyPath)
	isKeyExists := err == nil
	if isKeyExists && !isOverwrite {
		fmt.Printf("%s✔ Reusing existing SSH key:%s %s\n", appconstants.ColorCyan, appconstants.ColorReset, keyPath)
		return nil
	}

	if isKeyExists && isOverwrite {
		_ = os.Remove(keyPath)
		_ = os.Remove(pubPath)
	}

	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return apperror.WrapSimple(err, "create ssh directory")
	}

	return generateEd25519Key(keyPath, alias)
}

func generateEd25519Key(keyPath, alias string) error {
	comment := constants.DefaultSSHComment
	if len(alias) > 0 {
		comment = alias
	}

	cmd := exec.Command(constants.BinSshKeygen, constants.SshKeygenTypeFlag, constants.DefaultSSHKeyType, constants.SshKeygenCommentFlag, comment, constants.SshKeygenFileFlag, keyPath, constants.SshKeygenPassFlag, "")
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.WrapSimple(err, "ssh-keygen failed").WithContext("output", strings.TrimSpace(string(out)))
	}

	fmt.Printf("%s✔ Generated modern Ed25519 SSH key:%s %s\n", appconstants.ColorGreen, appconstants.ColorReset, keyPath)

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
	case constants.PlatformWindows:
		return constants.ClipWindows, nil
	case constants.PlatformDarwin:
		return constants.ClipDarwin, nil
	}

	return resolveLinuxClipboardCommand()
}

func resolveLinuxClipboardCommand() (string, []string) {
	if path, err := exec.LookPath(constants.ClipWlCopy); err == nil {
		return path, nil
	}

	if path, err := exec.LookPath(constants.ClipXclip); err == nil {
		return path, []string{constants.ClipboardSelectionFlag, constants.ClipboardClipboardArg}
	}

	if path, err := exec.LookPath(constants.ClipXsel); err == nil {
		return path, []string{constants.ClipboardFlag, constants.ClipboardInputFlag}
	}

	return "", nil
}

func printSSHPublicKeyBanner(pubKey string, isCopied bool) {
	printBannerKeyBox(pubKey)
	printCopyStatus(isCopied)
	printProviderKeyLinks()
}

func printBannerKeyBox(pubKey string) {
	fmt.Println()
	fmt.Println(appconstants.ColorCyan + constants.BannerDivider + appconstants.ColorReset)
	fmt.Println(appconstants.ColorCyan + "  SSH Public Key:" + appconstants.ColorReset)
	fmt.Println("  " + pubKey)
	fmt.Println(appconstants.ColorCyan + constants.BannerDivider + appconstants.ColorReset)
}

func printCopyStatus(isCopied bool) {
	if isCopied {
		fmt.Printf("  %s✔ Copied public key to clipboard!%s\n\n", appconstants.ColorGreen, appconstants.ColorReset)
		return
	}

	fmt.Printf("  %s! Could not automatically copy to clipboard; copy manually above.%s\n\n", appconstants.ColorYellow, appconstants.ColorReset)
}

func printProviderKeyLinks() {
	fmt.Println("  Add this key to your Git provider:")
	fmt.Printf("    • GitHub: %s%s%s\n", appconstants.ColorCyan, constants.GithubNewSSHKeyURL, appconstants.ColorReset)
	fmt.Printf("    • GitLab: %s%s%s\n", appconstants.ColorCyan, constants.GitlabNewSSHKeyURL, appconstants.ColorReset)
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

	return constants.DefaultProfileName
}

func updateOrAppendSSHProfile(cfg *model.GitProfileConfig, name, keyPath, keyName string) {
	if tryUpdateExistingProfile(cfg, name, keyPath, keyName) {
		return
	}

	addNewSSHProfile(cfg, name, keyPath, keyName)
}

func tryUpdateExistingProfile(cfg *model.GitProfileConfig, name, keyPath, keyName string) bool {
	for idx := range cfg.Profiles {
		if strings.EqualFold(cfg.Profiles[idx].Name, name) || cfg.Profiles[idx].ID == name {
			cfg.Profiles[idx].AuthMethod = enums.LoginAuthMethodSSH.String()
			cfg.Profiles[idx].SSHKeyPath = keyPath
			cfg.Profiles[idx].SSHKeyName = keyName
			cfg.Profiles[idx].LastUsedAt = time.Now()
			return true
		}
	}

	return false
}

func addNewSSHProfile(cfg *model.GitProfileConfig, name, keyPath, keyName string) {
	newProf := model.GitProfile{
		ID:         name,
		Name:       name,
		Provider:   enums.GitAuthProviderGitHub.String(),
		Type:       constants.DefaultProfileType,
		AuthMethod: enums.LoginAuthMethodSSH.String(),
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
	if exec.Command(constants.GitBin, constants.GitRevParse, constants.GitInsideWorkTree).Run() != nil {
		return
	}

	cleanPath := filepath.ToSlash(keyPath)
	cmdStr := fmt.Sprintf(constants.SSHCommandFormat, cleanPath)
	if err := exec.Command(constants.GitBin, constants.GitConfig, constants.GitConfigCoreSSHCommand, cmdStr).Run(); err == nil {
		fmt.Printf("  %s✔ Configured core.sshCommand for current repository.%s\n", appconstants.ColorGreen, appconstants.ColorReset)
	}

	if len(alias) > 0 {
		_ = usercontext.BindProject(".", alias)
	}
}
