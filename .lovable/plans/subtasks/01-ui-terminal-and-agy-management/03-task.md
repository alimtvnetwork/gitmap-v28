# Task 3: SQLite Suggestion Engine

Implement SQLite-backed suggestions for `rm` and `mv`.
- Add `GetRepoSuggestions` in `gitmap/store`.
- If the user types a partial name (e.g. `GH`), suggest the matching repos.
- Ensure semantic variable naming (no `data`, `tmp`).
