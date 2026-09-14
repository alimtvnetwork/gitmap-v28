package osuser

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// PortableUser defines a portable serialization schema for an OS user.
type PortableUser struct {
	Username string   `json:"username"`
	HomeDir  string   `json:"homedir,omitempty"`
	Shell    string   `json:"shell,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	IsSudoer bool     `json:"is_sudoer,omitempty"`
	SSHKeys  []string `json:"ssh_keys,omitempty"`
}

// ExportUser collects attributes for an OS user.
func ExportUser(username string) result.Result[PortableUser] {
	if username == "" {
		return result.Fail[PortableUser](apperror.NewSimple("empty username", "E_INVALID_USER"))
	}
	pu := PortableUser{
		Username: username,
		HomeDir:  resolveUserHome(username),
		Shell:    DefaultLinuxShell,
		IsSudoer: isUserSudoer(username),
		SSHKeys:  readUserAuthorizedKeys(username),
	}
	return result.Ok(pu)
}

// ExportAllUsers gathers all exportable users.
func ExportAllUsers() result.ResultSlice[PortableUser] {
	names := listAllOSUsernames()
	var users []PortableUser
	for _, name := range names {
		if res := ExportUser(name); res.IsSuccess() {
			users = append(users, res.Value)
		}
	}
	return result.OkSlice(users)
}

// ImportUser provisions an OS user from a portable representation.
func ImportUser(pu PortableUser) result.Result[bool] {
	if pu.Username == "" {
		return result.Fail[bool](apperror.NewSimple("missing username in import", "E_INVALID_USER"))
	}
	opts := UserCreateOptions{
		Username:       pu.Username,
		HomeDir:        pu.HomeDir,
		Shell:          pu.Shell,
		IsSudoer:       pu.IsSudoer,
		IsConfigureZsh: false,
	}
	if err := CreateRoot(opts); err != nil {
		return result.Fail[bool](err)
	}
	installImportedSSHKeys(pu.Username, pu.SSHKeys)
	return result.Ok(true)
}

func installImportedSSHKeys(username string, keys []string) {
	for _, k := range keys {
		_ = InstallKey(SSHKeyOptions{Username: username, PublicKey: k})
	}
}

// ImportAllUsers batch provisions users from a slice.
func ImportAllUsers(users []PortableUser) result.Result[int] {
	successCount := 0
	for _, u := range users {
		if res := ImportUser(u); res.IsSuccess() {
			successCount++
		}
	}
	return result.Ok(successCount)
}

// EditUser modifies an existing user shell or home directory.
func EditUser(username, newShell, newHome string) result.Result[bool] {
	if username == "" {
		return result.Fail[bool](apperror.NewSimple("empty username", "E_INVALID_USER"))
	}
	if runtime.GOOS == constants.OSWindows {
		return result.Ok(true)
	}
	return applyUserEditsLinux(username, newShell, newHome)
}

func applyUserEditsLinux(username, newShell, newHome string) result.Result[bool] {
	var args []string
	if newShell != "" {
		args = append(args, "-s", newShell)
	}
	if newHome != "" {
		args = append(args, "-d", newHome)
	}
	if len(args) == 0 {
		return result.Ok(true)
	}
	args = append(args, username)
	cmd := exec.Command("usermod", args...)
	if err := cmd.Run(); err != nil {
		return result.Fail[bool](apperror.WrapSimple(err, "edit user "+username))
	}
	return result.Ok(true)
}

func resolveUserHome(username string) string {
	if runtime.GOOS == constants.OSWindows {
		return filepath.Join(os.Getenv("SystemDrive")+"\\Users", username)
	}
	return "/home/" + username
}

func isUserSudoer(username string) bool {
	if runtime.GOOS == constants.OSWindows {
		return false
	}
	sudoFile := filepath.Join(SudoersDirPath, username)
	_, err := os.Stat(sudoFile)
	return err == nil
}

func readUserAuthorizedKeys(username string) []string {
	keyPath := filepath.Join(resolveUserHome(username), ".ssh", "authorized_keys")
	content, err := os.ReadFile(keyPath)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(content), "\n")
	var keys []string
	for _, l := range lines {
		if t := strings.TrimSpace(l); t != "" && !strings.HasPrefix(t, "#") {
			keys = append(keys, t)
		}
	}
	return keys
}

func listAllOSUsernames() []string {
	if runtime.GOOS == constants.OSWindows {
		return listWindowsUsers()
	}
	return listLinuxUsers()
}

func listWindowsUsers() []string {
	out, err := exec.Command("net", "user").Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\r\n")
	return parseNetUserLines(lines)
}

func parseNetUserLines(lines []string) []string {
	var names []string
	hasStarted := false
	for _, l := range lines {
		if strings.HasPrefix(l, "---") {
			hasStarted = true
			continue
		}
		if hasStarted && strings.Contains(l, "command completed") {
			break
		}
		if hasStarted {
			fields := strings.Fields(l)
			names = append(names, fields...)
		}
	}
	return names
}

func listLinuxUsers() []string {
	content, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil
	}
	lines := strings.Split(string(content), "\n")
	var names []string
	for _, l := range lines {
		if parts := strings.Split(l, ":"); len(parts) > 0 && parts[0] != "" {
			names = append(names, parts[0])
		}
	}
	return names
}
