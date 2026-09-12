# Subtask 03: Nuclear Package Modularization of cmd Monolith

## Objective
Decompose monolithic command domains from `gitmap/cmd` into smaller, independent, acyclic Go packages under `gitmap/` to shrink package size and dramatically reduce compilation times.

## Target Domains to Extract
1. `gitmap/chromeprofile` (30 files from `chromeprofile*.go`)
2. `gitmap/cmdagy` (29 files from `agy*.go`)
3. `gitmap/cmdprompt` (28 files from `prompt*.go`)

## Acyclic Rules
1. Extracted packages MUST NOT import `gitmap/cmd`.
2. Extracted packages only import leaf packages: `gitmap/constants`, `gitmap/model`, `gitmap/apperror`, `gitmap/cliexit`, `gitmap/fsutil`.
3. `gitmap/cmd/root.go` delegates CLI command dispatching to the extracted package handlers.
4. Verify with `go vet ./...` that zero import cycles exist.
