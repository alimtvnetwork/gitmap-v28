# Fix fmt.Errorf (Part 7)

Total items: 45

## Files to Modify

- `.\gitmap\txn\revert.go:151`: `return fmt.Errorf("transaction revert create: %w", err)`
- `.\gitmap\txn\revert.go:155`: `return fmt.Errorf("transaction revert copy: %w", err)`
- `.\gitmap\txn\snapshot.go:52`: `return fmt.Errorf("transaction snapshot stat %q: %w", absPath, err)`
- `.\gitmap\txn\snapshot.go:111`: `return "", fmt.Errorf("transaction backup mkdir: %w", err)`
- `.\gitmap\txn\snapshot.go:115`: `return "", fmt.Errorf("transaction backup open src: %w", err)`
- `.\gitmap\txn\snapshot.go:126`: `return "", fmt.Errorf("transaction backup create: %w", err)`
- `.\gitmap\txn\snapshot.go:131`: `return "", fmt.Errorf("transaction backup copy: %w", err)`
- `.\gitmap\visibility\exclude.go:63`: `return fmt.Errorf("Error: empty exclusion token at position %d (operation: parse-exclusion, reason: blank between commas)", tokIdx)`
- `.\gitmap\visibility\exclude.go:72`: `return fmt.Errorf("Error: non-numeric exclusion token %q at position %d (operation: parse-exclusion, reason: %s)", tok, tokIdx, err.Error())`
- `.\gitmap\visibility\exclude.go:88`: `return fmt.Errorf("Error: malformed range %q at position %d (operation: parse-exclusion, reason: range bounds must be integers)", tok, tokIdx)`
- `.\gitmap\visibility\exclude.go:91`: `return fmt.Errorf("Error: descending range %q at position %d (operation: parse-exclusion, reason: hi < lo)", tok, tokIdx)`
- `.\gitmap\visibility\exclude.go:113`: `return fmt.Errorf("Error: exclusion index %d out of range in token %q (operation: parse-exclusion, reason: valid range is 1..%d)", n, tok, totalCount)`
- `.\gitmap\visibility\pattern.go:39`: `return Pattern{}, fmt.Errorf("Error: empty pattern (operation: parse-pattern, reason: blank token)")`
- `.\gitmap\visibility\pattern.go:42`: `return Pattern{}, fmt.Errorf("Error: bare '*' pattern is refused at %q (operation: parse-pattern, reason: would match every repo under the owner)", raw)`
- `.\gitmap\visibility\pattern.go:47`: `return Pattern{}, fmt.Errorf("Error: pattern %q has no literal segments (operation: parse-pattern, reason: only wildcards)", raw)`
- `.\gitmap\visibility\pattern.go:106`: `return nil, fmt.Errorf("Error: empty pattern list (operation: parse-pattern-list, reason: arg is blank)")`
- `.\gitmap\visibility\pattern.go:115`: `return nil, fmt.Errorf("Error: empty pattern at token %d (operation: parse-pattern-list, reason: blank between commas)", i+1)`
- `.\gitmap\visibility\pattern.go:124`: `return nil, fmt.Errorf("Error: token %d %q: %w", i+1, trimmed, err)`
- `.\gitmap\vscodepm\io.go:29`: `return fmt.Errorf(constants.ErrVSCodePMRenameFailed,`
- `.\gitmap\vscodepm\io.go:39`: `return fmt.Errorf(constants.ErrVSCodePMWriteTempFailed, path, err)`
- `.\gitmap\vscodepm\io.go:44`: `return fmt.Errorf(constants.ErrVSCodePMWriteTempFailed, path, err)`
- `.\gitmap\vscodepm\io.go:51`: `return fmt.Errorf(constants.ErrVSCodePMWriteTempFailed, path, err)`
- `.\gitmap\vscodepm\mergemode.go:75`: `return MergeModeUnion, fmt.Errorf(`
- `.\gitmap\vscodepm\overwrite.go:47`: `return fmt.Errorf("%s: %w", constants.VSCodePMProjectsFile, err)`
- `.\gitmap\vscodepm\sync.go:102`: `return nil, fmt.Errorf(constants.ErrVSCodePMReadFailed, path, err)`
- `.\gitmap\vscodepm\sync.go:111`: `return nil, fmt.Errorf(constants.ErrVSCodePMParseFailed, path, err)`
- `.\gitmap\vscodepm\update_path.go:34`: `return fmt.Errorf("project not found with rootPath %s", oldPath)`
- `.\gitmap\vscodeworkspace\build.go:64`: `return nil, fmt.Errorf("vscode-workspace: encode: %w", err)`
- `.\gitmap\vscodeworkspace\build.go:81`: `firstErr = fmt.Errorf("relativize %q against %q: %w", f.Path, baseDir, err)`
- `.\gitmap\vscodeworkspace\write.go:23`: `return fmt.Errorf(constants.ErrVSCodeWorkspaceWriteTemp, tmpPath, err)`
- `.\gitmap\vscodeworkspace\write.go:26`: `return fmt.Errorf(constants.ErrVSCodeWorkspaceWriteTemp, tmpPath, err)`
- `.\gitmap\vscodeworkspace\write.go:31`: `return fmt.Errorf(constants.ErrVSCodeWorkspaceRename, outPath, err)`
- `.\gitmap-updater\cmd\github.go:21`: `return "", fmt.Errorf("failed to create request: %w", err)`
- `.\gitmap-updater\cmd\github.go:29`: `return "", fmt.Errorf("network error: %w", err)`
- `.\gitmap-updater\cmd\github.go:36`: `return "", fmt.Errorf(`
- `.\gitmap-updater\cmd\github.go:51`: `return "", fmt.Errorf("failed to parse release JSON: %w", err)`
- `.\gitmap-updater\cmd\worker.go:53`: `return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)`
- `.\scripts\changelog\internal\gitlog\gitlog.go:92`: `return "", fmt.Errorf("git tag failed: %w", err)`
- `.\scripts\changelog\internal\gitlog\gitlog.go:115`: `return nil, fmt.Errorf("git log failed: %w", err)`
- `.\scripts\changelog\internal\runner\args.go:62`: `return Args{}, fmt.Errorf("invalid -mode %q (want write or check)", *mode)`
- `.\scripts\changelog\internal\runner\execute.go:20`: `return 0, fmt.Errorf("collecting commits: %w", err)`
- `.\scripts\changelog\internal\runner\execute.go:86`: `return 0, fmt.Errorf("unhandled mode %q", mode)`
- `.\scripts\changelog\internal\writer\writer.go:24`: `return fmt.Errorf("CHANGELOG.md: %w", err)`
- `.\scripts\changelog\internal\writer\writer.go:29`: `return fmt.Errorf("src/data/changelog.ts: %w", err)`
- `.\scripts\changelog\internal\writer\writer.go:55`: `return fmt.Errorf("marker %q not found", marker)`
