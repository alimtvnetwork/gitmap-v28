# Project Type Detection — Acceptance Criteria

## Detection

1. Go projects detected by `go.mod` presence.
2. Node.js projects detected by `package.json` presence.
3. React projects detected by `package.json` with `react` dependency.
4. C++ projects detected by `CMakeLists.txt`, `*.vcxproj`, or
   `meson.build` presence.
5. C# projects detected by `*.csproj` or `*.sln` presence.
6. Multiple projects in one repo are all detected.
7. Exclusion directories are never scanned.

## Go Metadata

1. `go.mod` path and `go.sum` path recorded.
2. Module name parsed from `go.mod`.
3. Go version parsed from `go.mod`.
4. `cmd/` subdirectories scanned for `main.go` runnables.
5. Root-level `main.go` recorded as runnable.
6. Runnables stored in `GoRunnableFiles` table.
7. Metadata stored in `GoProjectMetadata` table.

## C# Metadata

1. `.sln` path and name recorded.
2. `global.json` SDK version parsed.
3. All `.csproj` files discovered and parsed.
4. Target framework, output type, and SDK extracted from `.csproj`.
5. Key files (`nuget.config`, `Directory.Build.props`, etc.) recorded.
6. Metadata stored in `CSharpProjectMetadata` table.
7. Project files stored in `CSharpProjectFiles` table.
8. Key files stored in `CSharpKeyFiles` table.

## Data Model

1. `ProjectTypes` table seeded with all supported types.
2. `DetectedProjects` uses `ProjectTypeId` FK (not text type).
3. Go metadata in separate `GoProjectMetadata` table.
4. Go runnables in separate `GoRunnableFiles` table.
5. C# metadata in separate `CSharpProjectMetadata` table.
6. C# project files in separate `CSharpProjectFiles` table.
7. C# key files in separate `CSharpKeyFiles` table.
8. All foreign keys properly named with `Id` suffix.
9. Cascade deletes from parent to child tables.

## JSON Export

1. Each type produces a dedicated JSON file.
2. Go JSON includes `goMetadata` with runnables.
3. C# JSON includes `csharpMetadata` with project files and key files.
4. No duplicates on repeated scans.
5. Empty types do not produce files.

## Commands

1. `gitmap-v27 go-repos` returns Go projects from DB.
2. `gitmap-v27 node-repos` returns Node.js projects from DB.
3. `gitmap-v27 react-repos` returns React projects from DB.
4. `gitmap-v27 cpp-repos` returns C++ projects from DB.
5. `gitmap-v27 csharp-repos` returns C# projects from DB.
6. `--json` flag outputs JSON format.
7. `--count` flag outputs count only.

## Reliability

1. Scan completes even if one repo's detection fails.
2. Errors logged with repo path and indicator file.
3. Excluded directories skipped.
4. Extensible for future project types.

---

## Optional Enhancements (Future)

1. `gitmap-v27 projects` — unified command listing all types grouped.
2. `--type` flag for filtering: `gitmap-v27 projects --type go,react`.
3. Summary line after scan: `"Detected: 5 Go, 3 Node, 2 React"`.
4. Confidence score per detection.
5. Configurable detection rules via `config.json`.
6. Monorepo workspace detection (npm/yarn/pnpm workspaces).
7. Dry-run mode: `gitmap-v27 scan --detect-only`.
