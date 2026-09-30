# Hallucination Checks & Prevention

## 1. Defining "Hallucination" in Code

An AI agent "hallucinates" when it:

- Invents API endpoints that do not exist in the backend.
- Proposes database columns or tables that were never defined in the schema.
- Uses third-party library functions that do not exist or are deprecated.

## 2. Prevention Strategies

### The "Read Before Write" Protocol

Agents MUST explicitly search and read the relevant definition files (e.g., router definitions, database schemas, `package.json`) before calling functions or endpoints.

### Strict Null and Type Enforcement

Ensure that the codebase enforces strict typing (TypeScript `strict: true`, Go `Result[T]`, Python `mypy`). A hallucinated property access will immediately fail CI/CD build steps, catching the hallucination before it reaches production.

### Verification Step

Agents are required to verify their own work by running static analysis or build commands (`npm run build`, `go build`, `cargo check`) immediately after generating a block of code.
