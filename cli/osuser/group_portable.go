package osuser

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// PortableGroup defines a portable serialization schema for an OS group.
type PortableGroup struct {
	Name    string   `json:"name"`
	GID     string   `json:"gid,omitempty"`
	Members []string `json:"members,omitempty"`
}

// ExportGroup gathers group details for an OS group.
func ExportGroup(name string) result.Result[PortableGroup] {
	if name == "" {
		return result.Fail[PortableGroup](apperror.NewSimple("empty group name", "E_INVALID_GROUP"))
	}
	pg := PortableGroup{
		Name:    name,
		Members: listGroupMembers(name),
	}
	return result.Ok(pg)
}

// ExportAllGroups gathers all exportable OS groups.
func ExportAllGroups() result.ResultSlice[PortableGroup] {
	names := listAllGroupNames()
	var groups []PortableGroup
	for _, name := range names {
		if res := ExportGroup(name); res.IsSuccess() {
			groups = append(groups, res.Value)
		}
	}
	return result.OkSlice(groups)
}

// ImportGroup creates an OS group from a portable representation.
func ImportGroup(pg PortableGroup) result.Result[bool] {
	if pg.Name == "" {
		return result.Fail[bool](apperror.NewSimple("missing group name", "E_INVALID_GROUP"))
	}
	cmd := buildCreateGroupCmd(pg.Name)
	if err := cmd.Run(); err != nil {
		return result.Fail[bool](apperror.WrapSimple(err, "import group "+pg.Name))
	}
	return result.Ok(true)
}

// ImportAllGroups batch provisions groups from a slice.
func ImportAllGroups(groups []PortableGroup) result.Result[int] {
	successCount := 0
	for _, g := range groups {
		if res := ImportGroup(g); res.IsSuccess() {
			successCount++
		}
	}
	return result.Ok(successCount)
}

// EditGroup renames an OS group.
func EditGroup(oldName, newName string) result.Result[bool] {
	if oldName == "" || newName == "" {
		return result.Fail[bool](apperror.NewSimple("empty group name for edit", "E_INVALID_GROUP"))
	}
	if runtime.GOOS == constants.OSWindows {
		return result.Ok(true)
	}
	cmd := exec.Command("groupmod", "-n", newName, oldName)
	if err := cmd.Run(); err != nil {
		return result.Fail[bool](apperror.WrapSimple(err, "rename group"))
	}
	return result.Ok(true)
}

func buildCreateGroupCmd(name string) *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return exec.Command("net", "localgroup", name, "/add")
	}
	if runtime.GOOS == constants.OSDarwin {
		return exec.Command("dscl", ".", "-create", "/Groups/"+name)
	}
	return exec.Command("groupadd", name)
}

func listGroupMembers(name string) []string {
	if runtime.GOOS == constants.OSWindows {
		return listWindowsGroupMembers(name)
	}
	return listLinuxGroupMembers(name)
}

func listWindowsGroupMembers(name string) []string {
	out, err := exec.Command("net", "localgroup", name).Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\r\n")
	return parseNetUserLines(lines)
}

func listLinuxGroupMembers(name string) []string {
	content, err := os.ReadFile("/etc/group")
	if err != nil {
		return nil
	}
	for _, l := range strings.Split(string(content), "\n") {
		parts := strings.Split(l, ":")
		if len(parts) >= 4 && parts[0] == name && parts[3] != "" {
			return strings.Split(parts[3], ",")
		}
	}
	return nil
}

func listAllGroupNames() []string {
	if runtime.GOOS == constants.OSWindows {
		return listWindowsGroups()
	}
	return listLinuxGroups()
}

func listWindowsGroups() []string {
	out, err := exec.Command("net", "localgroup").Output()
	if err != nil {
		return nil
	}
	return parseNetGroupLines(strings.Split(string(out), "\r\n"))
}

func parseNetGroupLines(lines []string) []string {
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
		if hasStarted && strings.HasPrefix(l, "*") {
			names = append(names, strings.TrimSpace(strings.TrimPrefix(l, "*")))
		}
	}
	return names
}

func listLinuxGroups() []string {
	content, err := os.ReadFile("/etc/group")
	if err != nil {
		return nil
	}
	var names []string
	for _, l := range strings.Split(string(content), "\n") {
		if parts := strings.Split(l, ":"); len(parts) > 0 && parts[0] != "" {
			names = append(names, parts[0])
		}
	}
	return names
}
