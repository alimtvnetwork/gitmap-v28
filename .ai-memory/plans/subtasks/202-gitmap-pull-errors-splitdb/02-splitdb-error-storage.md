---
title: "SplitDB Two-Tier Error Storage"
status: "pending"
---

# Objective
Refactor cli/store/errors_split_db.go and cli/store/pull_split_db_errors.go to implement the two-tier SplitDB error storage architecture.

# Implementation Steps
1. **Schema Refactoring:**
   - Modify the initialization in ErrorsSplitDB to create a lightweight RootErrorIndex table.
   - Define schema for a new RepoErrorDB (per-repo) to hold stack traces, context JSON, and remediation commands.
2. **Data Ingestion (Write Path):**
   - Refactor LogInternalErrorRecord and InsertPullError.
   - Split the incoming record: write the metadata to the Root DB, and connect to the target repo's local DB to write the heavy telemetry.
3. **Data Retrieval (Read Path):**
   - Update ListErrors, GetError, and QueryLatestPullErrors to execute a federated read.
   - First, query the Root DB for recent ErrorIDs.
   - Second, fetch the corresponding heavy payloads from the localized per-repo DBs and stitch them together for the CLI output.
4. **Data Purging:**
   - Update ClearErrors to iterate through known repositories and clear their local error tables before wiping the Root DB.

# Validation
- Ensure no stack traces are written to the main gitmap-errors.db.
- Verify per-repo databases are correctly initialized on error occurrence.
