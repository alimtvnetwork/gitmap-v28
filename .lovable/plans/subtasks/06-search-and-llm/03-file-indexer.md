# Subtask 3: File Indexing & Skipping Logic

1. Implement the recursive directory walker (skipping .git,
ode_modules and default dot folders unless forced).
2. Implement Delta Sync: Check modified times to only sync new/changed files into the Repo DB.
3. Identify files > 300KB and mark is_big=true without storing their content.
