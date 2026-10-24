package cmdworkgroup

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func isHelpTrigger(s string) bool {
	return s == "help" || s == "--help" || s == "-h"
}

func tryExtractFlagValue(arg, flag string) (string, bool) {
	prefix := flag + "="
	if strings.HasPrefix(arg, prefix) {
		return strings.TrimPrefix(arg, prefix), true
	}

	return "", false
}

func parseKeyValueFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, "--ssh-key"); ok {
		opts.SSHKeyPath = val

		return idx
	}
	if arg == "--ssh-key" && idx+1 < len(args) {
		opts.SSHKeyPath = args[idx+1]

		return idx + 1
	}
	if val, ok := tryExtractFlagValue(arg, "--desc"); ok {
		opts.Description = val

		return idx
	}
	if arg == "--desc" && idx+1 < len(args) {
		opts.Description = args[idx+1]

		return idx + 1
	}

	return idx
}

func parseHostFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, "--host"); ok {
		opts.Host = val

		return idx
	}
	if arg == "--host" && idx+1 < len(args) {
		opts.Host = args[idx+1]

		return idx + 1
	}

	return idx
}

func parseUserHostFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, "--user-name"); ok {
		opts.UserName = val

		return idx
	}
	if arg == "--user-name" && idx+1 < len(args) {
		opts.UserName = args[idx+1]

		return idx + 1
	}
	if val, ok := tryExtractFlagValue(arg, "--user-email"); ok {
		opts.UserEmail = val

		return idx
	}
	if arg == "--user-email" && idx+1 < len(args) {
		opts.UserEmail = args[idx+1]

		return idx + 1
	}

	return parseHostFlag(args, idx, opts)
}

func parseBoolFlag(arg string, opts *WorkgroupOptions) bool {
	if arg == "-j" || arg == "--json" {
		opts.IsJSON = true

		return true
	}
	if arg == "--default" {
		opts.IsDefault = true

		return true
	}

	return false
}

func parseFlagsLoop(args []string, opts *WorkgroupOptions) []string {
	pos := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if parseBoolFlag(args[i], opts) {
			continue
		}
		newI := parseKeyValueFlag(args, i, opts)
		if newI != i {
			i = newI

			continue
		}
		newI2 := parseUserHostFlag(args, i, opts)
		if newI2 != i {
			i = newI2

			continue
		}
		pos = append(pos, args[i])
	}

	return pos
}

func assignPositionalWorkgroup(pos []string, opts *WorkgroupOptions) {
	if len(pos) == 0 {
		return
	}
	opts.Subcommand = strings.ToLower(pos[0])
	if len(pos) > 1 {
		opts.Name = pos[1]
		opts.Path = pos[1]
	}
	if len(pos) > 2 {
		opts.Path = pos[2]
	}
}

// ParseWorkgroupOptions extracts flags and positional arguments.
func ParseWorkgroupOptions(args []string) WorkgroupOptions {
	var opts WorkgroupOptions
	pos := parseFlagsLoop(args, &opts)
	assignPositionalWorkgroup(pos, &opts)

	return opts
}

func isHelpSubcommand(sub string) bool {
	return sub == "" || sub == "help" || sub == "--help" || sub == "-h"
}

func dispatchWorkgroupAction(opts WorkgroupOptions) error {
	if isHelpSubcommand(opts.Subcommand) {
		RenderWorkgroupHelp()

		return nil
	}
	switch opts.Subcommand {
	case "ls", "list":
		return runWorkgroupList(opts)
	case "add":
		return runWorkgroupAdd(opts)
	case "bind":
		return runWorkgroupBind(opts)
	case "unbind":
		return runWorkgroupUnbind(opts)
	default:
		return dispatchWorkgroupSecondaryAction(opts)
	}
}

func dispatchWorkgroupSecondaryAction(opts WorkgroupOptions) error {
	switch opts.Subcommand {
	case "show":
		return runWorkgroupShow(opts)
	case "test-ssh", "test":
		return runWorkgroupTestSSH(opts)
	case "rm", "delete", "remove":
		return runWorkgroupDelete(opts)
	default:
		return runWorkgroupDefaultShow(opts)
	}
}

func outputJSON(v any) error {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal json")
	}
	fmt.Println(string(bytes))

	return nil
}

func runWorkgroupList(opts WorkgroupOptions) error {
	records, err := ListWorkgroups()
	if err != nil {
		return err
	}
	if opts.IsJSON {
		return outputJSON(records)
	}
	RenderWorkgroupTable(records)

	return nil
}

func runWorkgroupAdd(opts WorkgroupOptions) error {
	if opts.Name == "" || opts.Path == "" {
		return apperror.NewSimple("Usage: gitmap wg add <name> <dir> [flags]", "E9003")
	}
	absPath, _ := filepath.Abs(opts.Path)
	rec := WorkgroupRecord{
		Name:          opts.Name,
		RootDirectory: absPath,
		SSHKeyPath:    opts.SSHKeyPath,
		Description:   opts.Description,
		GitUserName:   opts.UserName,
		GitUserEmail:  opts.UserEmail,
		IsDefault:     opts.IsDefault,
	}
	created, err := AddWorkgroup(rec)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Registered workgroup '%s' -> %s\n", created.Name, created.RootDirectory)

	return nil
}

func runWorkgroupBind(opts WorkgroupOptions) error {
	if opts.Name == "" || opts.SSHKeyPath == "" {
		return apperror.NewSimple("Usage: gitmap wg bind <dir> <ssh-key-path> [flags]", "E9004")
	}
	absPath, _ := filepath.Abs(opts.Name)
	rec := FolderBindingRecord{
		FolderPath:   absPath,
		SSHKeyPath:   opts.SSHKeyPath,
		GitUserName:  opts.UserName,
		GitUserEmail: opts.UserEmail,
	}
	created, err := BindFolder(rec)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Bound folder '%s' to SSH key '%s'\n", created.FolderPath, created.SSHKeyPath)

	return nil
}

func runWorkgroupUnbind(opts WorkgroupOptions) error {
	target := opts.Name
	if target == "" {
		target = "."
	}
	absPath, _ := filepath.Abs(target)
	if err := UnbindFolder(absPath); err != nil {
		return err
	}
	fmt.Printf("✓ Unbound folder: %s\n", absPath)

	return nil
}

func runWorkgroupShow(opts WorkgroupOptions) error {
	target := opts.Name
	if target == "" {
		target = "."
	}
	rec, err := GetWorkgroup(target)
	if err != nil {
		return apperror.WrapNotFound(err, "workgroup not found: "+target)
	}
	eff, _ := ResolveWorkgroupForPath(rec.RootDirectory)
	if opts.IsJSON {
		return outputJSON(rec)
	}
	RenderWorkgroupDetailCard(rec, eff)

	return nil
}

func runWorkgroupDefaultShow(opts WorkgroupOptions) error {
	rec, err := GetWorkgroup(opts.Subcommand)
	if err != nil {
		RenderWorkgroupHelp()

		return nil
	}
	eff, _ := ResolveWorkgroupForPath(rec.RootDirectory)
	RenderWorkgroupDetailCard(rec, eff)

	return nil
}

func runWorkgroupTestSSH(opts WorkgroupOptions) error {
	target := opts.Name
	if target == "" {
		target = "."
	}
	eff, err := ResolveWorkgroupForPath(target)
	if err != nil || !eff.HasSSHKey {
		return apperror.NewSimple("No bound SSH key found for target: "+target, "E9005")
	}
	fmt.Printf("Testing SSH handshake for key: %s (host: %s)...\n", eff.SSHKeyPath, resolveSSHHost(opts.Host))
	res, errTest := TestSSHConnection(eff.SSHKeyPath, opts.Host)
	if errTest != nil {
		return errTest
	}
	fmt.Printf("✓ Handshake response:\n  %s\n", res)

	return nil
}

func runWorkgroupDelete(opts WorkgroupOptions) error {
	if opts.Name == "" {
		return apperror.NewSimple("Usage: gitmap wg rm <name|id>", "E9006")
	}
	if err := DeleteWorkgroup(opts.Name); err != nil {
		return err
	}
	fmt.Printf("✓ Removed workgroup: %s\n", opts.Name)

	return nil
}

// RunWorkgroup dispatches gitmap workgroup / wg / group CLI commands.
func RunWorkgroup(args []string) error {
	opts := ParseWorkgroupOptions(args)

	return dispatchWorkgroupAction(opts)
}
