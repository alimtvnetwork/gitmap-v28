# 🔴 CODE RED: No Generated Code, Artifacts, or Test Results

**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 🚫 The Rule: NEVER Commit Generated Code or Artifacts

Under no circumstances should any generated code (e.g., ORM models, gRPC stubs, OpenAPI clients), compiled binary, test result, test report, or temporary test data ever be committed to the Git repository.

### What is Forbidden?

- **Generated Code:** `*.generated.*` or `*_generated.*` files produced by code generators (unless explicitly permitted as test fixtures).
- **Test Results & Reports:** Any outputs from test runs, such as HTML reports, coverage data, JSON summaries, `.test-report.*`, or CSV results.
- **Compiled Binaries & Executables:** `.exe`, `.dll`, `.so`, `.dylib`, `.out`, `.class`, or any pre-compiled binaries resulting from a build step.
- **Build Directories:** `build/`, `bin/`, `obj/`, `dist/` (unless explicitly allowed for a specific deployment pipeline in a separate spec).
- **Temporary Data:** Any mock data, database dumps, or local logs generated during a local run.

### Why?

1. **Repository Bloat:** Binary files and continuous test reports exponentially increase the Git repository size, slowing down clones and pulls for all team members and CI/CD pipelines.
2. **Noise in History:** Committing transient data pollutes the Git history and diffs, making it impossible to perform meaningful code reviews.
3. **Security Risks:** Test results or compiled binaries can accidentally leak sensitive environment data, secrets, or internal architectural layouts.

---

## AI Agent Directives

> 🛑 **MANDATORY INSTRUCTION FOR AI:**
>
> If you are instructed to run a test suite, compile a binary, or execute a script that produces an output file, you MUST NOT commit those outputs to the Git repository.
>
> If you notice that an action you performed has generated a new test report or binary, you must ensure that it is either ignored via `.gitignore` or deleted before you stage files using `git add`.
>
> **Never blindly run `git add .` without verifying that generated artifacts are safely ignored.**

## Enforcement

This rule is enforced globally via `.gitignore` patterns. If a new type of artifact is introduced, you must update `.gitignore` before committing anything else.
