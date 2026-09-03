# CI/CD Implementation Guide for AI Agents

When instructing an AI to implement a robust, exactly similar CI/CD architecture for another repository (whether it is Golang, React, TypeScript, Rust, etc.), you must feed the AI the following core reference files from this repository.

These files demonstrate the "decoupled versioning" and automated release pipeline patterns used here.

## Files to Feed to the AI

### 1. `version.json` (Root directory)

**Why feed it:** This file demonstrates the single-source-of-truth pattern for versioning. By keeping the version in a JSON file rather than hardcoded in source files, you prevent brittle regex string replacements during CI/CD. The AI needs to see this to understand how versioning is centralized.

### 2. `.github/workflows/release.yml`

**Why feed it:** This is the backbone of the automated deployment. It shows the AI how to:
- Extract the version safely from `version.json` (using tools like `jq`).
- Inject the version dynamically during the compilation phase (e.g., using `-ldflags` in Golang, build flags in Rust, or environment variables in React/Node).
- Automatically generate a GitHub release and upload the compiled binaries.

### 3. `.github/workflows/ci.yml`

**Why feed it:** This demonstrates the Continuous Integration gates. It shows the AI how to set up strict linting rules, run unit tests, and configure caching correctly for the language being used, ensuring no broken code reaches the release phase.

### 4. `.lovable/memory/tech/01-version-json-architecture.md`

**Why feed it:** This memory file acts as the "Architecture Decision Record". It explicitly instructs the AI on *why* we use `version.json` and explicitly forbids the AI from trying to sweep through source code to bump version numbers. This is crucial for maintaining a reliable pipeline.

### 5. `changelog.md`

**Why feed it:** This file demonstrates the standardized markdown structure required so that the CI/CD pipeline can automatically parse out the release notes for the latest version and attach them to the GitHub Release.

## Instructions to give the AI (Prompt Template)

When providing the above files to the AI in the new repository, accompany them with this instruction:

> *"Attached are the core CI/CD pipeline files and architecture documents from our reference repository. I want you to implement the exact same CI/CD architecture for this [Insert Language] project.*
>
> *Key requirements:*
> *1. Establish a `version.json` at the root as the sole source of truth.*
> *2. Modify the build scripts/pipeline to inject this version dynamically at build time (do not hardcode it in the source files).*
> *3. Create a GitHub Actions workflow that reads `version.json`, compiles the code with the injected version, and automates the GitHub release.*
> *4. Ensure the CI workflow runs strict linting and testing on PRs."*
