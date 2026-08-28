import os

def rewrite(path, changes):
    with open(path, 'r', encoding='utf-8') as f:
        c = f.read()
    for o, n in changes:
        c = c.replace(o, n)
    with open(path, 'w', encoding='utf-8') as f:
        f.write(c)

rewrite('gitmap/cmd/addignoreattrs.go', [
    ('func runAddIgnore(args []string) error', 'func runAddIgnore(args []string) *apperror.AppError'),
    ('func runAddAttributes(args []string) error', 'func runAddAttributes(args []string) *apperror.AppError'),
    ('func executeAddTemplate(spec addTemplateSpec, args []string) {', 'func executeAddTemplate(spec addTemplateSpec, args []string) *apperror.AppError {'),
    ('flags, langs := parseAddTemplateArgs(spec, args)', 'flags, langs, err := parseAddTemplateArgs(spec, args)\n\tif err != nil { return err }'),
    ('fmt.Fprintln(os.Stderr, "  ? Not inside a Git repository.")\n\t\tfmt.Fprintln(os.Stderr, "    Run this from the root of a repo (where .git/ lives).")\n\t\tos.Exit(1)', 'return apperror.New("Not inside a Git repository", "E9000", nil)'),
    ('resolved, err := resolveAddTemplates(spec.kind, langs)', 'resolved, err2 := resolveAddTemplates(spec.kind, langs)'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err2 != nil {\n\t\treturn apperror.Wrap(err2, "?", nil)\n\t}'),
    ('target, err := repoFilePath(spec.targetName)', 'target, err3 := repoFilePath(spec.targetName)'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not locate repo root: %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err3 != nil {\n\t\treturn apperror.Wrap(err3, "? Could not locate repo root", nil)\n\t}'),
    ('if mergeErr != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not merge into %s: %v\\n", target, mergeErr)\n\t\tos.Exit(1)\n\t}', 'if mergeErr != nil {\n\t\treturn apperror.Wrap(mergeErr, "? Could not merge into", nil)\n\t}'),
    ('printAddTemplateDryRun(target, tag, body)\n\n\t\treturn\n\t}', 'printAddTemplateDryRun(target, tag, body)\n\n\t\treturn nil\n\t}'),
    ('printAddTemplateSummary(spec, res)\n}', 'printAddTemplateSummary(spec, res)\n\treturn nil\n}'),
    ('func parseAddTemplateArgs(spec addTemplateSpec, args []string) (addTemplateFlags, []string) {', 'func parseAddTemplateArgs(spec addTemplateSpec, args []string) (addTemplateFlags, []string, *apperror.AppError) {'),
    ('if err := fs.Parse(args); err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not parse flags: %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err := fs.Parse(args); err != nil {\n\t\treturn addTemplateFlags{}, nil, apperror.Wrap(err, "? Could not parse flags", nil)\n\t}'),
    ('return addTemplateFlags{dryRun: *dryRun}, normalizeLangs(fs.Args())\n}', 'return addTemplateFlags{dryRun: *dryRun}, normalizeLangs(fs.Args()), nil\n}'),
    ('executeAddTemplate(addTemplateSpec{\n\t\tkind:        "ignore",\n\t\tsubcommand:  "ignore",\n\t\ttargetName:  ".gitignore",\n\t\tbannerLabel: "merge curated .gitignore template block",\n\t}, args)\n\treturn nil', 'return executeAddTemplate(addTemplateSpec{\n\t\tkind:        "ignore",\n\t\tsubcommand:  "ignore",\n\t\ttargetName:  ".gitignore",\n\t\tbannerLabel: "merge curated .gitignore template block",\n\t}, args)'),
    ('executeAddTemplate(addTemplateSpec{\n\t\tkind:        "attributes",\n\t\tsubcommand:  "attributes",\n\t\ttargetName:  ".gitattributes",\n\t\tbannerLabel: "merge curated .gitattributes template block",\n\t}, args)\n\treturn nil', 'return executeAddTemplate(addTemplateSpec{\n\t\tkind:        "attributes",\n\t\tsubcommand:  "attributes",\n\t\ttargetName:  ".gitattributes",\n\t\tbannerLabel: "merge curated .gitattributes template block",\n\t}, args)')
])

rewrite('gitmap/cmd/addlfsinstall.go', [
    ('func runAddLFSInstall(args []string) error {', 'func runAddLFSInstall(args []string) *apperror.AppError {'),
    ('flags := parseAddLFSInstallFlags(args)', 'flags, err := parseAddLFSInstallFlags(args)\n\tif err != nil { return err }'),
    ('resolved, err := templates.Resolve("lfs", "common")', 'resolved, err2 := templates.Resolve("lfs", "common")'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not resolve lfs/common template: %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err2 != nil {\n\t\treturn apperror.Wrap(err2, "? Could not resolve lfs/common template", nil)\n\t}'),
    ('target, err := gitattributesPath()', 'target, err3 := gitattributesPath()'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not locate repo root: %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err3 != nil {\n\t\treturn apperror.Wrap(err3, "? Could not locate repo root", nil)\n\t}'),
    ('res, err := templates.Merge(target, addLFSInstallTag, resolved.Content)', 'res, err4 := templates.Merge(target, addLFSInstallTag, resolved.Content)'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not merge template into %s: %v\\n", target, err)\n\t\tos.Exit(1)\n\t}', 'if err4 != nil {\n\t\treturn apperror.Wrap(err4, "? Could not merge template", nil)\n\t}'),
    ('func parseAddLFSInstallFlags(args []string) addLFSInstallFlags {', 'func parseAddLFSInstallFlags(args []string) (addLFSInstallFlags, *apperror.AppError) {'),
    ('if err := fs.Parse(args); err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ✗ Could not parse flags: %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err := fs.Parse(args); err != nil {\n\t\treturn addLFSInstallFlags{}, apperror.Wrap(err, "✗ Could not parse flags", nil)\n\t}'),
    ('return addLFSInstallFlags{dryRun: *dryRun}\n}', 'return addLFSInstallFlags{dryRun: *dryRun}, nil\n}')
])

rewrite('gitmap/cmd/agy_cmd.go', [
    ('func runAgyAdd(args []string) error {', 'func runAgyAdd(args []string) *apperror.AppError {'),
    ('return fmt.Errorf("requires id and name")', 'return apperror.New("requires id and name", "E9000", nil)'),
    ('return createProjectFile(args[0], args[1])', 'if err := createProjectFile(args[0], args[1]); err != nil { return apperror.Wrap(err, "create", nil) }; return nil'),
    ('func runAgyRm(args []string) error {', 'func runAgyRm(args []string) *apperror.AppError {'),
    ('return fmt.Errorf("requires id")', 'return apperror.New("requires id", "E9000", nil)'),
    ('return deleteProjectFile(args[0])', 'if err := deleteProjectFile(args[0]); err != nil { return apperror.Wrap(err, "delete", nil) }; return nil'),
    ('func runAgyLs() error {', 'func runAgyLs() *apperror.AppError {'),
    ('return pathErr', 'return apperror.Wrap(pathErr, "path error", nil)'),
    ('return readErr', 'return apperror.Wrap(readErr, "read error", nil)'),
    ('return printProjectEntries(entries)', 'if err := printProjectEntries(entries); err != nil { return apperror.Wrap(err, "print entries", nil) }; return nil'),
    ('func runAgyStats() error {', 'func runAgyStats() *apperror.AppError {'),
    ('return printProjectStats(entries)', 'if err := printProjectStats(entries); err != nil { return apperror.Wrap(err, "print stats", nil) }; return nil'),
    ('func runAgyUpdate(args []string) error {', 'func runAgyUpdate(args []string) *apperror.AppError {'),
    ('return updateProjectFile(args[0])', 'if err := updateProjectFile(args[0]); err != nil { return apperror.Wrap(err, "update error", nil) }; return nil')
])

rewrite('gitmap/cmd/amend.go', [
    ('func runAmend(args []string) error {', 'func runAmend(args []string) *apperror.AppError {'),
    ('opts, args := parseAmendOpts(args)', 'opts, args, err := parseAmendOpts(args)\n\tif err != nil { return err }'),
    ('func parseAmendOpts(args []string) (amendOptions, []string) {', 'func parseAmendOpts(args []string) (amendOptions, []string, *apperror.AppError) {'),
    ('if err := fs.Parse(args); err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err := fs.Parse(args); err != nil {\n\t\treturn amendOptions{}, nil, apperror.Wrap(err, "parse flags", nil)\n\t}'),
    ('return amendOptions{\n\t\tlimit:   *limit,\n\t\tbranch:  *branch,\n\t\tdryRun:  *dryRun,\n\t\tedit:    *edit,\n\t\tverbose: *verbose,\n\t}, fs.Args()\n}', 'return amendOptions{\n\t\tlimit:   *limit,\n\t\tbranch:  *branch,\n\t\tdryRun:  *dryRun,\n\t\tedit:    *edit,\n\t\tverbose: *verbose,\n\t}, fs.Args(), nil\n}'),
    ('commits := loadRecentCommits(opts.limit)', 'commits, err := loadRecentCommits(opts.limit)\n\tif err != nil { return err }'),
    ('func loadRecentCommits(limit int) []model.CommitEntry {', 'func loadRecentCommits(limit int) ([]model.CommitEntry, *apperror.AppError) {'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not load branch history: %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err != nil {\n\t\treturn nil, apperror.Wrap(err, "load branch history", nil)\n\t}'),
    ('return commits\n}', 'return commits, nil\n}')
])

rewrite('gitmap/cmd/amendexec.go', [
    ('func amendCommit(branch string, entry model.CommitEntry, body string, dryRun bool) {', 'func amendCommit(branch string, entry model.CommitEntry, body string, dryRun bool) *apperror.AppError {'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? %v\\n", err)\n\t\tos.Exit(1)\n\t}', 'if err != nil {\n\t\treturn apperror.Wrap(err, "git checkout", nil)\n\t}'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? %s: %v\\n", constants.ErrGitCommitFailed, err)\n\t\tos.Exit(1)\n\t}', 'if err != nil {\n\t\treturn apperror.Wrap(err, constants.ErrGitCommitFailed, nil)\n\t}'),
    ('if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not return to %s: %v\\n", branch, err)\n\t\t// don\'t exit, we at least committed it\n\t}', 'if err != nil {\n\t\tfmt.Fprintf(os.Stderr, "  ? Could not return to %s: %v\\n", branch, err)\n\t\t// don\'t exit, we at least committed it\n\t}'),
    ('amendCommit(opts.branch, chosen, message, opts.dryRun)\n\treturn nil', 'return amendCommit(opts.branch, chosen, message, opts.dryRun)'),
    ('printAmendSuccess(opts, chosen, message)\n}', 'printAmendSuccess(opts, chosen, message)\n\treturn nil\n}')
])

