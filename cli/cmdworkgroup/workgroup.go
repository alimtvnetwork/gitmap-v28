package cmdworkgroup

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func isHelpTrigger(s string) bool {
	return s == "help" || s == string(constants.WorkgroupFlagHelp) || s == string(constants.WorkgroupFlagHelpShort)
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
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagSSHKey)); ok {
		opts.SSHKeyPath = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagSSHKey) && idx+1 < len(args) {
		opts.SSHKeyPath = args[idx+1]

		return idx + 1
	}
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagDesc)); ok {
		opts.Description = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagDesc) && idx+1 < len(args) {
		opts.Description = args[idx+1]

		return idx + 1
	}

	return idx
}

func parseHostFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagHost)); ok {
		opts.Host = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagHost) && idx+1 < len(args) {
		opts.Host = args[idx+1]

		return idx + 1
	}

	return idx
}

func parseUserHostFlag(args []string, idx int, opts *WorkgroupOptions) int {
	arg := args[idx]
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagUserName)); ok {
		opts.UserName = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagUserName) && idx+1 < len(args) {
		opts.UserName = args[idx+1]

		return idx + 1
	}
	if val, ok := tryExtractFlagValue(arg, string(constants.WorkgroupFlagUserEmail)); ok {
		opts.UserEmail = val

		return idx
	}
	if arg == string(constants.WorkgroupFlagUserEmail) && idx+1 < len(args) {
		opts.UserEmail = args[idx+1]

		return idx + 1
	}

	return parseHostFlag(args, idx, opts)
}

func parseBoolFlag(arg string, opts *WorkgroupOptions) bool {
	if arg == string(constants.WorkgroupFlagJSONShort) || arg == string(constants.WorkgroupFlagJSON) {
		opts.IsJSON = true

		return true
	}
	if arg == string(constants.WorkgroupFlagDefault) {
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
		assignSecondPositional(pos[2], opts)
	}
}

func assignSecondPositional(val string, opts *WorkgroupOptions) {
	if opts.Subcommand == string(constants.WorkgroupCmdBind) && opts.SSHKeyPath == "" {
		opts.SSHKeyPath = val

		return
	}
	opts.Path = val
}

// ParseWorkgroupOptions extracts flags and positional arguments.
func ParseWorkgroupOptions(args []string) WorkgroupOptions {
	var opts WorkgroupOptions
	pos := parseFlagsLoop(args, &opts)
	assignPositionalWorkgroup(pos, &opts)

	return opts
}

func isHelpSubcommand(sub string) bool {
	return sub == "" || sub == "help" || sub == string(constants.WorkgroupFlagHelp) || sub == string(constants.WorkgroupFlagHelpShort)
}

func dispatchWorkgroupAction(opts WorkgroupOptions) error {
	if isHelpSubcommand(opts.Subcommand) {
		RenderWorkgroupHelp()

		return nil
	}
	switch opts.Subcommand {
	case string(constants.WorkgroupCmdLs), string(constants.WorkgroupCmdList):
		return runWorkgroupList(opts)
	case string(constants.WorkgroupCmdAdd):
		return runWorkgroupAdd(opts)
	case string(constants.WorkgroupCmdBind):
		return runWorkgroupBind(opts)
	case string(constants.WorkgroupCmdUnbind):
		return runWorkgroupUnbind(opts)
	default:
		return dispatchWorkgroupSecondaryAction(opts)
	}
}

func dispatchWorkgroupSecondaryAction(opts WorkgroupOptions) error {
	switch opts.Subcommand {
	case string(constants.WorkgroupCmdShow):
		return runWorkgroupShow(opts)
	case string(constants.WorkgroupCmdTestSSH), string(constants.WorkgroupCmdTest):
		return runWorkgroupTestSSH(opts)
	case string(constants.WorkgroupCmdRm), string(constants.WorkgroupCmdDelete), string(constants.WorkgroupCmdRemove):
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

func createWorkgroupRecord(opts WorkgroupOptions) WorkgroupRecord {
	absPath, _ := filepath.Abs(opts.Path)

	return WorkgroupRecord{
		Name:          opts.Name,
		RootDirectory: absPath,
		SSHKeyPath:    opts.SSHKeyPath,
		Description:   opts.Description,
		GitUserName:   opts.UserName,
		GitUserEmail:  opts.UserEmail,
		IsDefault:     opts.IsDefault,
	}
}

func runWorkgroupAdd(opts WorkgroupOptions) error {
	if opts.Name == "" || opts.Path == "" {
		return apperror.NewSimple("Usage: gitmap wg add <name> <dir> [flags]", string(constants.ErrCodeWorkgroupAddUsage))
	}
	rec := createWorkgroupRecord(opts)
	created, err := AddWorkgroup(rec)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Registered workgroup '%s' -> %s\n", created.Name, created.RootDirectory)

	return nil
}

func createFolderBindingRecord(opts WorkgroupOptions) FolderBindingRecord {
	absPath, _ := filepath.Abs(opts.Name)

	return FolderBindingRecord{
		FolderPath:   absPath,
		SSHKeyPath:   opts.SSHKeyPath,
		GitUserName:  opts.UserName,
		GitUserEmail: opts.UserEmail,
	}
}

func runWorkgroupBind(opts WorkgroupOptions) error {
	if opts.Name == "" || opts.SSHKeyPath == "" {
		return apperror.NewSimple("Usage: gitmap wg bind <dir> <ssh-key-path> [flags]", string(constants.ErrCodeWorkgroupBindUsage))
	}
	rec := createFolderBindingRecord(opts)
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

func executeSSHTestHandshake(keyPath, host string) error {
	fmt.Printf("Testing SSH handshake for key: %s (host: %s)...\n", keyPath, resolveSSHHost(host))
	res, errTest := TestSSHConnection(keyPath, host)
	if errTest != nil {
		return errTest
	}
	fmt.Printf("✓ Handshake response:\n  %s\n", res)

	return nil
}

func runWorkgroupTestSSH(opts WorkgroupOptions) error {
	target := opts.Name
	if target == "" {
		target = "."
	}
	eff, err := ResolveWorkgroupForPath(target)
	if err != nil || !eff.HasSSHKey {
		return apperror.NewSimple("No bound SSH key found for target: "+target, string(constants.ErrCodeWorkgroupNoBoundKey))
	}

	return executeSSHTestHandshake(eff.SSHKeyPath, opts.Host)
}

func runWorkgroupDelete(opts WorkgroupOptions) error {
	if opts.Name == "" {
		return apperror.NewSimple("Usage: gitmap wg rm <name|id>", string(constants.ErrCodeWorkgroupDeleteUsage))
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
