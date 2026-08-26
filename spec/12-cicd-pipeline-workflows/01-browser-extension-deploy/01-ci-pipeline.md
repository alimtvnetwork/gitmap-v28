# Browser Extension — CI Pipeline

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The CI pipeline validates every push and pull request to the `main` branch. It runs linting, tests, and a full dependency-graph build to verify that the extension compiles successfully.

---

## Trigger and Concurrency

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

concurrency:
  group: ci-main-${{ github.sha }}
  cancel-in-progress: true

permissions:
  contents: read
```

---

## Job Graph

```
setup (lint + test)
  └── build-sdk
        ├── build-module-a (parallel)
        ├── build-module-b (parallel)
        └── build-module-c (parallel — no SDK dependency if independent)
              └── build-extension (depends on all above)
```

### Job 1: Setup — Lint — Test

The entry point job performs:

1. **Checkout** with `fetch-depth: 0` for full history
2. **Filename enforcement** — verify all `.md` files are lowercase kebab-case
3. **Node.js and package manager setup** — pinned versions
4. **Install dependencies** — root and extension-specific
5. **Lint** — run linter on root and extension workspaces
6. **Test** — run the full test suite

### Dependency Installation Pattern

When a project uses workspaces, the extension subdirectory may need its own install step:

```bash

# Remove workspace file to prevent cross-linking issues

rm -f pnpm-workspace.yaml
if [ -f pnpm-lock.yaml ]; then
  pnpm install --frozen-lockfile
else
  pnpm install --no-frozen-lockfile --lockfile=false
fi
```

**Why**: The extension may have its own `package.json` with dependencies that conflict with the root workspace. Removing the workspace file isolates the install.

---

## Job 2: Build SDK

The SDK is a shared library used by multiple downstream modules. It must be built first and uploaded as an artifact.

```yaml
build-sdk:
  needs: setup
  steps:
    - checkout
    - setup-node
    - setup-pnpm
    - install dependencies
    - run: pnpm run build:sdk
    - upload-artifact:
        name: sdk-dist
        path: standalone-scripts/<sdk-name>/dist/
        retention-days: 1
```

---

## Jobs 3a–3c: Build Standalone Modules (Parallel)

Multiple modules build in parallel after the SDK is available. Each job:

1. Checks out the repository
2. Sets up Node.js and pnpm
3. Installs dependencies
4. **Downloads the SDK artifact** into the expected path
5. Runs its build command
6. Uploads its own dist artifact

```yaml
build-module-a:
  needs: build-sdk
  steps:
    - download-artifact: sdk-dist → standalone-scripts/<sdk-name>/dist/
    - run: pnpm run build:module-a
    - upload-artifact: module-a-dist
```

Modules that do **not** depend on the SDK can use `needs: setup` instead of `needs: build-sdk` to start earlier.

---

## Job 4: Build Extension

The final assembly job depends on **all** upstream builds:

```yaml
build-extension:
  needs: [build-sdk, build-module-a, build-module-b, build-module-c]
```

It downloads every artifact into its expected directory, then runs the extension build command. The output is the complete extension `dist/` folder.

---

## Artifact Flow

```
build-sdk → sdk-dist artifact
  ├── build-module-a → module-a-dist artifact
  ├── build-module-b → module-b-dist artifact
  └── build-module-c → module-c-dist artifact (independent)

build-extension ← downloads ALL of the above
```

All inter-job artifacts use `retention-days: 1` since they are only needed within the same workflow run.

---

## Constraints

- SDK must be built before any module that depends on it
- Independent modules can build in parallel
- The extension build must depend on ALL upstream modules
- Workspace isolation may be needed for subdirectory installs
- No `cd` in CI steps — use `working-directory`
- All tool versions are pinned

---

*Browser extension CI pipeline — updated: 2026-04-09*
