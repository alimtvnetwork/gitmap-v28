# Subtask 01: Mandatory AI Name Validation & On-The-Fly Teaching

- **Parent Plan:** `360-mandatory-ai-name-and-settings-ui-configuration`
- **Slug:** `mandatory-ai-name-and-settings-ui-configuration - Task 01`
- **Status:** `completed`
- **Assigned Subagent:** Worker 01

---

## Deliverables
1. Update `cli/cmdlogin/login_cmd.go` and `cli/cmdlogin/login_ai_auto.go`:
   - If `gitmap login set ai-auto` is called without an agent name (and without `--global`/`global`):
     - Check configured `ai_auto.default_agent` from settings.
     - If empty, return a validation error showing on-the-fly teaching examples (`muse`, `gemini`, `cursor-model`, `antigravity`, `claude`).
   - If `global` is requested without an agent name:
     - Check configured `ai_auto.default_name` and `ai_auto.default_email` from settings.
2. Update `cli/cmdcommit/commit_flags.go`:
   - If committing under an AI bot identity without `--ai-agent`:
     - Resolve configured fallback from `ai_auto.default_agent`.
     - Display on-the-fly teaching guidance demonstrating `--ai-agent <name>`.
3. Add unit tests in `cli/cmdlogin/` and `cli/cmdcommit/` verifying mandatory validation and teaching feedback.
