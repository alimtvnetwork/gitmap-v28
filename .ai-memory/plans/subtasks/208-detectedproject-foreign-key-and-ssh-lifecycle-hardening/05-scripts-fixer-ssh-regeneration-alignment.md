# Subtask 78.5: Scripts-Fixer Project SSH Key Regeneration Alignment

## 1. Context & Objective
The user noted that in the `scripts-fixer` project (or related scripts), SSH key regeneration is not regenerating existing keys properly. We must audit `scripts-fixer` and all SSH regeneration scripts in `scripts/` and `03-ai-scripts/` to ensure they pass `-y` or `--force` or use GitMap's updated non-interactive key generation options.

---

## 2. Technical Scope
- `scripts-fixer/`
- `03-ai-scripts/`
- `scripts/`

---

## 3. Remediation Checklist
- [ ] Locate references to `gitmap ssh create` or `ssh-keygen` in `scripts-fixer` and automation scripts.
- [ ] Ensure any automated SSH regeneration flags pass `-y` or `--confirm` to bypass interactive prompts.
- [ ] Document proper non-interactive regeneration syntax in CLI help and docs.
