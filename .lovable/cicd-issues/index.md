- [14-top-level-cmd-registry.md](14-top-level-cmd-registry.md): TestTopLevelCmdRegistryMatchesAST fails when new top-level Cmd constants are not added to topLevelCmds() registry.
- [15-clone-sync-examples.md](15-clone-sync-examples.md): TestEveryHelpFileHasExamples fails when using 4-space indents instead of fenced code blocks.
- [16-github-actions-yaml-parse-fatal.md](16-github-actions-yaml-parse-fatal.md): Workflows stopped triggering due to fatal YAML syntax errors (missing runs, email obfuscation artifacts, missing required inputs).
- [17-cicd-fixes-go-version-drift.md](./17-cicd-fixes-go-version-drift.md) - Go Version Drift, Linter Failure, and Generate Desync

- [18-ci-multi-regression-drift.md](18-ci-multi-regression-drift.md): Fixed go generate drift, setup-go-cached checkout order, python syntax corruption, changelog sync, installer regex, and golangci-lint go version mismatch.

- [19-bump-script-unicode-and-regex-corruption.md](19-bump-script-unicode-and-regex-corruption.md): Fixed UnicodeDecodeError in python script and greedy regex corruption of readme.md.

- [19-ci-gitmap-open-error-refactor.md](19-ci-gitmap-open-error-refactor.md): Fixed Python cp1252 decode errors, completed mass AST refactoring of 80+ commands to return typed errors instead of os.Exit(1), and implemented gitmap open cross-platform.
- [20-cicd-four-headed-hydra.md](./20-cicd-four-headed-hydra.md): Fix for env platform return types, gofmt, lint script key errors, and legacy gitmap-v6 markdown references.
- [21-go-generate-drift.md](./21-go-generate-drift.md): Fix generated files drift caused by modified CLI constants.
