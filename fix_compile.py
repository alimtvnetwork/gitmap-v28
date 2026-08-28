import os

def fix_addignoreattrs():
    p = 'gitmap/cmd/addignoreattrs.go'
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace('func executeAddTemplate(spec addTemplateSpec, args []string) {', 'func executeAddTemplate(spec addTemplateSpec, args []string) *apperror.AppError {')
    c = c.replace('flags, langs := parseAddTemplateArgs(spec, args)', 'flags, langs, err := parseAddTemplateArgs(spec, args)\n\tif err != nil { return err }')
    c = c.replace('return apperror.New("fatal error", "E9000", nil)', 'return apperror.New("Not inside a Git repository", "E9000", nil)')
    c = c.replace('target, err := repoFilePath', 'target, err2 := repoFilePath')
    c = c.replace('if err != nil {\n\t\treturn apperror.Wrap(err, "? Could not locate repo root:", nil)\n\t}', 'if err2 != nil {\n\t\treturn apperror.Wrap(err2, "? Could not locate repo root:", nil)\n\t}')
    c = c.replace('resolved, err := resolveAddTemplates', 'resolved, err3 := resolveAddTemplates')
    c = c.replace('if err != nil {\n\t\treturn apperror.Wrap(err, "?", nil)\n\t}', 'if err3 != nil {\n\t\treturn apperror.Wrap(err3, "?", nil)\n\t}')
    c = c.replace('return apperror.Wrap(target, "? Could not merge into :", nil)', 'return apperror.Wrap(mergeErr, "merge error", nil)')
    c = c.replace('printAddTemplateDryRun(target, tag, body)\n\n\t\treturn\n\t}', 'printAddTemplateDryRun(target, tag, body)\n\n\t\treturn nil\n\t}')
    c = c.replace('printAddTemplateSummary(spec, res)\n}', 'printAddTemplateSummary(spec, res)\n\treturn nil\n}')
    c = c.replace('func parseAddTemplateArgs(spec addTemplateSpec, args []string) (addTemplateFlags, []string) {', 'func parseAddTemplateArgs(spec addTemplateSpec, args []string) (addTemplateFlags, []string, *apperror.AppError) {')
    c = c.replace('return apperror.Wrap(err, "? Could not parse flags:", nil)', 'return addTemplateFlags{}, nil, apperror.Wrap(err, "parse flags", nil)')
    c = c.replace('return addTemplateFlags{dryRun: *dryRun}, normalizeLangs(fs.Args())\n}', 'return addTemplateFlags{dryRun: *dryRun}, normalizeLangs(fs.Args()), nil\n}')
    with open(p, 'w', encoding='utf8') as f: f.write(c)

def fix_addlfsinstall():
    p = 'gitmap/cmd/addlfsinstall.go'
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace('return apperror.Wrap(target, "? Could not merge template into :", nil)', 'return apperror.Wrap(err, "merge template", nil)')
    c = c.replace('func parseAddLFSInstallFlags(args []string) addLFSInstallFlags {', 'func parseAddLFSInstallFlags(args []string) (addLFSInstallFlags, *apperror.AppError) {')
    c = c.replace('return apperror.Wrap(err, "✗ Could not parse flags:", nil)', 'return addLFSInstallFlags{}, apperror.Wrap(err, "parse flags", nil)')
    c = c.replace('return addLFSInstallFlags{dryRun: *dryRun}\n}', 'return addLFSInstallFlags{dryRun: *dryRun}, nil\n}')
    c = c.replace('flags := parseAddLFSInstallFlags(args)', 'flags, err := parseAddLFSInstallFlags(args)\n\tif err != nil { return err }')
    with open(p, 'w', encoding='utf8') as f: f.write(c)

def fix_alias():
    p = 'gitmap/cmd/alias.go'
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace('return apperror.Wrap(sub, "constants.ErrUnknownCommand", nil)', 'return apperror.New(fmt.Sprintf(constants.ErrUnknownCommand, sub), "E9000", nil)')
    with open(p, 'w', encoding='utf8') as f: f.write(c)

def fix_agy_cmd():
    p = 'gitmap/cmd/agy_cmd.go'
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace('return fmt.Errorf("requires id and name")', 'return apperror.New("requires id and name", "E9000", nil)')
    c = c.replace('return createProjectFile(args[0], args[1])', 'if err := createProjectFile(args[0], args[1]); err != nil { return apperror.Wrap(err, "create", nil) }; return nil')
    c = c.replace('return fmt.Errorf("requires id")', 'return apperror.New("requires id", "E9000", nil)')
    c = c.replace('return deleteProjectFile(args[0])', 'if err := deleteProjectFile(args[0]); err != nil { return apperror.Wrap(err, "delete", nil) }; return nil')
    c = c.replace('return pathErr', 'return apperror.Wrap(pathErr, "path error", nil)')
    c = c.replace('return readErr', 'return apperror.Wrap(readErr, "read error", nil)')
    c = c.replace('return printProjectEntries(entries)', 'if err := printProjectEntries(entries); err != nil { return apperror.Wrap(err, "print error", nil) }; return nil')
    c = c.replace('return printProjectStats(entries)', 'if err := printProjectStats(entries); err != nil { return apperror.Wrap(err, "print stats error", nil) }; return nil')
    c = c.replace('return updateProjectFile(args[0])', 'if err := updateProjectFile(args[0]); err != nil { return apperror.Wrap(err, "update error", nil) }; return nil')
    with open(p, 'w', encoding='utf8') as f: f.write(c)

def fix_aliasresolve():
    p = 'gitmap/cmd/aliasresolve.go'
    with open(p, 'r', encoding='utf8') as f: c = f.read()
    c = c.replace('return apperror.New("constants.ErrAliasEmpty", "E9000", nil)\n\treturn "", nil', 'return "", args')
    c = c.replace('func resolveAliasContext(aliasName string) {', 'func resolveAliasContext(aliasName string) *apperror.AppError {')
    c = c.replace('resolveAliasContext(args[0])\n\treturn nil', 'return resolveAliasContext(args[0])')
    c = c.replace('return apperror.Wrap(aliasName, "?", nil)', 'return apperror.Wrap(err, "resolve alias", nil)')
    with open(p, 'w', encoding='utf8') as f: f.write(c)

fix_addignoreattrs()
fix_addlfsinstall()
fix_alias()
fix_agy_cmd()
fix_aliasresolve()
