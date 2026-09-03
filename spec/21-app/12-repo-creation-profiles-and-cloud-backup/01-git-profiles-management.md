# Multi-Account Git Profiles Management

## Data Schema & Storage

Git profiles are persisted in `data/git_profiles.json` relative to the binary's co-located data directory (`store.BinaryDataDir()`):

```json
{
  "profiles": [
    {
      "id": "prof_1",
      "name": "alimtvnetwork",
      "provider": "github",
      "type": "user",
      "authMethod": "gh-cli",
      "isDefault": true,
      "usageCount": 14,
      "lastUsedAt": "2026-09-03T18:45:12Z"
    },
    {
      "id": "auktvgo",
      "name": "auktvgo",
      "provider": "github",
      "type": "organization",
      "authMethod": "gh-cli",
      "isDefault": false,
      "usageCount": 3,
      "lastUsedAt": "2026-09-02T12:00:00Z"
    }
  ],
  "active": "alimtvnetwork",
  "default": "alimtvnetwork",
  "updatedAt": "2026-09-03T18:45:12Z"
}
```

---

## Auto-Discovery Protocol

When `data/git_profiles.json` is missing or empty, GitMap triggers auto-discovery:
1. Queries authenticated GitHub user via `gh api user --jq .login`.
2. Queries associated GitHub organizations via `gh api user/orgs --jq .[].login`.
3. Queries local Git committer identity via `git config user.name` and `user.email`.
4. Writes initial configured profiles and designates the primary authenticated user as default.

---

## Interactive Sequence Picker (`1`, `2`, `3`...)

When invoked in an interactive terminal (`isInteractiveStdin() == true`) without explicit arguments:
1. Renders indexed menu `[1]`, `[2]`, `[3]` with name, provider, and account type.
2. Prompts user for numeric choice: `Select default profile (Enter number 1-N): `.
3. Validates numeric input within bounds `[1, len(profiles)]`.
4. Updates `Default` and `Active` properties and persists to disk.
5. In non-interactive runners (`CI=1` or `GITMAP_NON_INTERACTIVE=1`), immediate error or bypass is returned without blocking.
