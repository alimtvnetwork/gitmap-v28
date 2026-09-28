# Subtask 02: Settings, Web UI Help, CLI Helptext (`repo-secrets.md`, `repo-cache.md`) & V2 Prompts

## Allowed Target Files (Strict Disjoint Bounding Box)
- `cli/cmdagy/agy_telegram_email_settings.go`
- `cli/cmdui/ui_assets.go`
- `cli/helptext/repo-secrets.md`
- `cli/helptext/repo-cache.md`
- `cli/helptext/catalog.go`
- `01-prompts/special-repos-secrets-and-cache.md`

## Deliverables
1. Expose `special_repos.secrets_name` (`default: "repo-secrets"`) and `special_repos.cache_name` (`default: "repo-cache"`) in `gitmap settings` (`cli/cmdagy/agy_telegram_email_settings.go`) and Web UI `#tab-settings` / `#tab-help` (`cli/cmdui/ui_assets.go`).
2. Author `cli/helptext/repo-secrets.md` and `cli/helptext/repo-cache.md`, and register `rs`, `repo-secrets`, `rc`, `repo-cache`, `repo-storage` in `cli/helptext/catalog.go`.
3. Author `01-prompts/special-repos-secrets-and-cache.md` instructing AI models to store repository secrets in `repo-secrets` (`gitmap rs`) and reusable temporary scripts in `repo-cache` (`gitmap rc`).
