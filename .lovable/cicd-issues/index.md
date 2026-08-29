- [14-top-level-cmd-registry.md](14-top-level-cmd-registry.md): TestTopLevelCmdRegistryMatchesAST fails when new top-level Cmd constants are not added to topLevelCmds() registry.
- [15-clone-sync-examples.md](15-clone-sync-examples.md): TestEveryHelpFileHasExamples fails when using 4-space indents instead of fenced code blocks.
- [16-github-actions-yaml-parse-fatal.md](16-github-actions-yaml-parse-fatal.md): Workflows stopped triggering due to fatal YAML syntax errors (missing runs, email obfuscation artifacts, missing required inputs).
- [17-cicd-fixes-go-version-drift.md](./17-cicd-fixes-go-version-drift.md) - Go Version Drift, Linter Failure, and Generate Desync

- [18-ci-multi-regression-drift.md](18-ci-multi-regression-drift.md): Fixed go generate drift, setup-go-cached checkout order, python syntax corruption, changelog sync, installer regex, and golangci-lint go version mismatch.

- [19-bump-script-unicode-and-regex-corruption.md](19-bump-script-unicode-and-regex-corruption.md): Fixed UnicodeDecodeError in python script and greedy regex corruption of readme.md.

- [19-ci-gitmap-open-error-refactor.md](19-ci-gitmap-open-error-refactor.md): Fixed Python cp1252 decode errors, completed mass AST refactoring of 80+ commands to return typed errors instead of os.Exit(1), and implemented gitmap open cross-platform.
- [20-cicd-four-headed-hydra.md](./20-cicd-four-headed-hydra.md): Fix for env platform return types, gofmt, lint script key errors, and legacy gitmap-v6 markdown references. <!-- gitmap-legacy-ref-allow -->
- [21-go-generate-drift.md](./21-go-generate-drift.md): Fix generated files drift caused by modified CLI constants.
- [22-go-mod-bump-linter-break.md](./22-go-mod-bump-linter-break.md): Fix go.mod version bump causing typecheck failures in pinned linter.
- [23-linter-guideline-conflict.md](./23-linter-guideline-conflict.md): Linter guideline conflict resolution and explicit boolean rule protection.
- [24-equalfold-and-us-english-fix.md](./24-equalfold-and-us-english-fix.md): EqualFold reversion fix and US English misspell linter enforcement.
- [25-go-125-dependency-linter-break.md](./25-go-125-dependency-linter-break.md): Revert go.mod and golang.org/x/* dependencies to Go 1.24 compatibility to fix CI typecheck crash.
- [26-go-125-gofmt-and-goimports-drift.md](./26-go-125-gofmt-and-goimports-drift.md): Fix Go 1.25 gofmt and goimports drift across all repository packages.
- [27-installctx-e2e-with-explain-type-mismatch.md](./27-installctx-e2e-with-explain-type-mismatch.md): Fix runInstallCtxMac and runInstallCtxLinux type mismatch in withExplain cross-platform test suites.
- [28-legacy-refs-bare-err-and-lint-keyerror.md](./28-legacy-refs-bare-err-and-lint-keyerror.md): Fix legacy refs whitelist in issue notes, bare stderr prints in cmd, and is_fail KeyError in lint scripts.
- [29-lfs-zip-drift-changelog-sync-and-jq-diff-argjson.md](./29-lfs-zip-drift-changelog-sync-and-jq-diff-argjson.md): Fix Git LFS binary zip false-positive in generate drift check, sync constants.Version with changelog.md, and sanitize jq --argjson in linter diff scripts.
