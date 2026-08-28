import os, re

p = 'gitmap/cmd/clone.go'
with open(p, 'r', encoding='utf8') as f: c = f.read()

c = c.replace('func applySSHKey(name string) {', 'func applySSHKey(name string) *apperror.AppError {')
c = c.replace('if len(name) == 0 {\n\t\treturn\n\t}', 'if len(name) == 0 {\n\t\treturn nil\n\t}')
c = c.replace('fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)\n}', 'fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)\n\treturn nil\n}')
c = c.replace('applySSHKey(cf.SSHKey)', 'if err := applySSHKey(cf.SSHKey); err != nil { return err }')

c = c.replace('func executeDirectClone(url, folderName string, ghDesktopFlag, noReplace bool, output string, noVSCodeSync bool) {', 'func executeDirectClone(url, folderName string, ghDesktopFlag, noReplace bool, output string, noVSCodeSync bool) *apperror.AppError {')
c = re.sub(r'(fmt\.Fprintf\(os\.Stderr, "  \\033\[1;92m[^"]+Task %s spawned and running\.\\n", taskID\)\n\t\}\n\})', r'\1\n\treturn nil\n}', c)
c = c.replace('executeDirectClone(url, folderName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)\n\t\treturn nil', 'return executeDirectClone(url, folderName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)')
c = c.replace('executeDirectClone(rootName, rootName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)\n\treturn nil', 'return executeDirectClone(rootName, rootName, cf.GithubDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)')

c = c.replace('func validateShorthandPath(resolved string) string {', 'func validateShorthandPath(resolved string) (string, *apperror.AppError) {')
c = c.replace('return resolved\n\t}', 'return resolved, nil\n\t}')
c = c.replace('return apperror.New(constants.ErrShorthandNotFound, "E9000", nil)\n\n\treturn ""\n}', 'return "", apperror.New(constants.ErrShorthandNotFound, "E9000", nil)\n}')
c = c.replace('return validateShorthandPath(resolved)', 'path, _ := validateShorthandPath(resolved)\n\t\treturn path')

c = c.replace('func executeClone(source, targetDir string, safePull, ghDesktop bool, maxConcurrency int, defaultBranch string, noVSCodeSync bool, clean bool, missingOnly bool) {', 'func executeClone(source, targetDir string, safePull, ghDesktop bool, maxConcurrency int, defaultBranch string, noVSCodeSync bool, clean bool, missingOnly bool) *apperror.AppError {')
c = re.sub(r'(fmt\.Fprintf\(os\.Stderr, "  \\033\[1;92m[^"]+Task %s spawned and running\.\\n", taskID\)\n\})', r'\1\n\treturn nil\n}', c)
c = c.replace('executeClone(source, target, cf.SafePull, cf.GithubDesktop, cf.MaxConcurrency, cf.DefaultBranch, cf.NoVSCodeSync, cf.Clean, cf.MissingOnly)\n\treturn nil', 'return executeClone(source, target, cf.SafePull, cf.GithubDesktop, cf.MaxConcurrency, cf.DefaultBranch, cf.NoVSCodeSync, cf.Clean, cf.MissingOnly)')

with open(p, 'w', encoding='utf8') as f: f.write(c)
