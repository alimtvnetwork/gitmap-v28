package cmdworkgroup

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

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
