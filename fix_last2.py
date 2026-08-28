import os, re

def r(p, o, n):
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace(o, n)
    with open(p, 'w', encoding='utf8') as f: f.write(c)

# changelog.go
r('gitmap/cmd/changelog.go', 'func printSingleVersion(entries []release.ChangelogEntry, version string, pretty bool) {', 'func printSingleVersion(entries []release.ChangelogEntry, version string, pretty bool) *apperror.AppError {')
r('gitmap/cmd/changelog.go', 'printChangelogEntry(entry, pretty)\n}', 'printChangelogEntry(entry, pretty)\n\treturn nil\n}')
r('gitmap/cmd/changelog.go', 'printSingleVersion(entries, version, pretty)\n\n\t\treturn nil', 'return printSingleVersion(entries, version, pretty)')

# changeloggen.go
r('gitmap/cmd/changeloggen.go', 'fmt.Printf(constants.MsgChangelogGenWritten, constants.ChangelogFile)\n}', 'fmt.Printf(constants.MsgChangelogGenWritten, constants.ChangelogFile)\n\treturn nil\n}')

# chrome_backup.go
r('gitmap/cmd/chrome_backup.go', 'func doDryRun(src, dst string) {', 'func doDryRun(src, dst string) *apperror.AppError {')
r('gitmap/cmd/chrome_backup.go', 'fmt.Printf("\\033[1;93m chrome restore (dry-run)\\033[0m  %d file(s) would land in \\033[1;96m%s\\033[0m\\n", n, dst)\n}', 'fmt.Printf("\\033[1;93m chrome restore (dry-run)\\033[0m  %d file(s) would land in \\033[1;96m%s\\033[0m\\n", n, dst)\n\treturn nil\n}')
r('gitmap/cmd/chrome_backup.go', 'doDryRun(src, dst)\n\t\treturn nil', 'return doDryRun(src, dst)')

# chrome_bookmarks.go
r('gitmap/cmd/chrome_bookmarks.go', 'return apperror.Wrap(match, "chrome export-bookmarks: ERROR no bookmarks matched --match=%q --title=%q within the selected subtree  hint: --match is a case-insensitive substring; --title is an exact title match", nil)', 'return apperror.New("chrome export-bookmarks: ERROR no bookmarks matched", "E9000", nil)')

# clone.go
r('gitmap/cmd/clone.go', 'func setupCloneSSH(name string) {', 'func setupCloneSSH(name string) *apperror.AppError {')
r('gitmap/cmd/clone.go', 'if len(name) == 0 {\n\t\treturn\n\t}', 'if len(name) == 0 {\n\t\treturn nil\n\t}')
r('gitmap/cmd/clone.go', 'fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)\n}', 'fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)\n\treturn nil\n}')
r('gitmap/cmd/clone.go', 'setupCloneSSH(cf.SSHKey)', 'if err := setupCloneSSH(cf.SSHKey); err != nil { return err }')

r('gitmap/cmd/clone.go', 'func dispatchCloneToTemp(folderName string, noReplace bool) {', 'func dispatchCloneToTemp(folderName string, noReplace bool) *apperror.AppError {')
r('gitmap/cmd/clone.go', 'fmt.Fprintf(os.Stderr, "  \\033[1;92m✔\\033[0m  Task %s spawned and running.\\n", taskID)\n\t}\n}', 'fmt.Fprintf(os.Stderr, "  \\033[1;92m✔\\033[0m  Task %s spawned and running.\\n", taskID)\n\t}\n\treturn nil\n}')
r('gitmap/cmd/clone.go', 'dispatchCloneToTemp(cf.Folder, cf.NoReplace)\n\t\treturn nil', 'return dispatchCloneToTemp(cf.Folder, cf.NoReplace)')
r('gitmap/cmd/clone.go', 'dispatchCloneToTemp(rootName, cf.NoReplace)\n\treturn nil', 'return dispatchCloneToTemp(rootName, cf.NoReplace)')

r('gitmap/cmd/clone.go', 'func validateShorthandPath(resolved string) string {', 'func validateShorthandPath(resolved string) (string, *apperror.AppError) {')
r('gitmap/cmd/clone.go', 'return resolved\n\t}', 'return resolved, nil\n\t}')
r('gitmap/cmd/clone.go', 'return apperror.New(constants.ErrShorthandNotFound, "E9000", nil)\n\n\treturn ""\n}', 'return "", apperror.New(constants.ErrShorthandNotFound, "E9000", nil)\n}')
r('gitmap/cmd/clone.go', 'target = validateShorthandPath(target)', 'target, err = validateShorthandPath(target)\n\tif err != nil { return err }')
