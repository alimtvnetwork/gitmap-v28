# Project Type Detection — Detection Rules

## Go

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `go.mod` file exists               | High       |
| Secondary | `go.sum` file exists               | Medium     |
| Secondary | `*.go` files present               | Low        |

**Project name:** Parsed from `module` directive in `go.mod` (first
line starting with `module `). Fall back to directory name.

**False positive prevention:** Ignore `go.mod` inside `vendor/` or
`testdata/` directories.

---

## Node.js

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `package.json` file exists         | High       |
| Secondary | `package-lock.json` exists         | Medium     |
| Secondary | `yarn.lock` exists                 | Medium     |
| Secondary | `pnpm-lock.yaml` exists            | Medium     |
| Secondary | `bun.lockb` or `bun.lock` exists   | Medium     |

**Project name:** Parsed from `name` field in `package.json`.
Fall back to directory name.

**False positive prevention:** Ignore `package.json` inside
`node_modules/`, `vendor/`, or `testdata/` directories.

**Classification upgrade:** If `package.json` contains React or
ReactJS dependencies, the project is classified as `react` instead
of `node`. See React rules below.

---

## React

A Node.js project is reclassified as React when **any** of these
conditions are true:

| Condition                                           | Check Location        |
|-----------------------------------------------------|-----------------------|
| `react` in `dependencies`                           | `package.json`        |
| `react` in `devDependencies`                        | `package.json`        |
| `@types/react` in `devDependencies`                 | `package.json`        |
| `react-scripts` in `dependencies`                   | `package.json` (CRA)  |
| `next` in `dependencies`                            | `package.json` (Next) |
| `gatsby` in `dependencies`                          | `package.json`        |
| `remix` or `@remix-run/react` in `dependencies`     | `package.json`        |

**Project name:** Same as Node.js (from `package.json` `name` field).

**Note:** React projects do **not** also appear as Node.js. They are
classified exclusively as `react`.

---

## C++

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `CMakeLists.txt` file exists       | High       |
| Primary   | `*.vcxproj` file exists            | High       |
| Primary   | `meson.build` file exists          | High       |
| Secondary | `Makefile` with C++ content        | Medium     |
| Secondary | `*.cpp`, `*.cc`, `*.cxx` files     | Medium     |
| Secondary | `*.hpp`, `*.hh`, `*.hxx` files     | Medium     |
| Tertiary  | `conanfile.txt` or `conanfile.py`  | Medium     |
| Tertiary  | `vcpkg.json` exists                | Medium     |

**Project name:** For CMake projects, parsed from `project()`
directive in `CMakeLists.txt`. For others, fall back to directory name.

**Makefile disambiguation:** A `Makefile` alone does not trigger C++
detection. It must be accompanied by at least one C++ source file
or a `CMakeLists.txt` / `*.vcxproj` in the same directory tree.

**False positive prevention:** Ignore `build/`, `cmake-build-*/`,
`out/`, and `target/` directories.

---

## C#

| Priority  | Indicator                          | Confidence |
|-----------|------------------------------------|------------|
| Primary   | `*.csproj` file exists             | High       |
| Primary   | `*.sln` file exists                | High       |
| Secondary | `*.fsproj` file exists             | Medium     |
| Secondary | `global.json` file exists          | Medium     |
| Secondary | `*.cs` source files present        | Low        |

**Project name:** For `.csproj` projects, parsed from the filename
(e.g., `MyApp.csproj` → `MyApp`). For `.sln` solutions, parsed from
the filename. Fall back to directory name.

**False positive prevention:** Ignore `bin/`, `obj/`, and `packages/`
directories.

### Solution vs Project Scope

When a `.sln` file is found at the repo root (or any directory), it
defines a **single** `DetectedProject` entry at that path. The `.sln`
is the `PrimaryIndicator`. Individual `.csproj` files beneath it are
**not** recorded as separate `DetectedProject` rows — they are stored
as `CSharpProjectFiles` child records under the solution's metadata.

A standalone `.csproj` (no parent `.sln` in any ancestor directory) is
recorded as its own `DetectedProject` with the `.csproj` as the
`PrimaryIndicator`.

**Precedence rule:** `.sln` takes priority. When walking the tree,
if a `.sln` is found, mark that directory as a C# project and skip
creating separate project entries for `.csproj` files anywhere below
that `.sln` directory. If no `.sln` exists, each `.csproj` becomes
its own detected project.

---

## Detection Scope

### Where Detection Runs

Detection scans **inside** each discovered Git repository. The scanner
walks the repo directory tree (excluding standard exclusion dirs) and
looks for indicator files.

### Monorepo Handling

A single Git repository may contain **multiple** detected projects.
For example, a monorepo with `backend/` (Go) and `frontend/` (React)
produces two separate `DetectedProject` records, both linked to the
same repo.

### Nested Projects

If a Node.js project at `./` contains a React project at `./web/`,
both are recorded. The more specific classification wins at each
path level — the root is `node`, the `web/` subdirectory is `react`.

### Exclusion Directories

The project detector skips the following directories:

| Directory            | Reason                          |
|----------------------|---------------------------------|
| `node_modules`       | Dependencies, not source        |
| `vendor`             | Vendored dependencies           |
| `.git`               | Git internals                   |
| `dist`               | Build output                    |
| `build`              | Build output                    |
| `target`             | Build output (Rust/Java/C++)    |
| `bin`                | Binary output                   |
| `obj`                | .NET build output               |
| `out`                | Generic build output            |
| `cmake-build-*`      | CMake build directories         |
| `testdata`           | Test fixtures                   |
| `packages`           | NuGet packages                  |
| `.venv`              | Python virtual environments     |
| `.cache`             | Cache directories               |
