# Implementation Rules (AI Execution Prompt)

> **/goal** Implement strict Axios version pinning across package manifests, lock files, dependency bot configurations, and CI pipelines.
> **/learn** Master exact version pinning techniques, configure Dependabot/Renovate exclusion rules, and automate CI validation scripts to reject blocked versions.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Declare Axios with exact version strings in `package.json` and ensure lock file synchronization.
- [ ] `/learn` Exclude Axios from automated dependency update tools in `.github/dependabot.yml` and `renovate.json`.
- [ ] `/goal` Enforce CI verification script to detect and fail on range symbols or blocked versions (`1.14.1`, `0.30.4`).
- [ ] `/learn` Adhere to the 6-step version upgrade procedure before adopting any new Axios release.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Axios Version Control Policy](./readme.md)
> **Version:** 1.1.0
> **Updated:** 2026-04-02

---

## 1. Dependency Declaration

### 1.1 Exact Version Pinning (Mandatory)

Axios MUST be declared with an **exact version** — no range symbols allowed.

#### ✅ Correct

```json
{
  "dependencies": {
    "axios": "1.14.0"
  }
}
```

#### ❌ Incorrect — All of These Are Forbidden

```json
{ "axios": "^1.14.0" }    // Caret — allows minor/patch upgrades
{ "axios": "~1.14.0" }    // Tilde — allows patch upgrades
{ "axios": ">=1.14.0" }   // Range — allows any higher version
{ "axios": "*" }           // Wildcard — allows anything
{ "axios": "latest" }      // Tag — unpredictable
```

### 1.2 Lock File Enforcement

- `package-lock.json` / `bun.lock` MUST reflect the exact pinned version
- After any install, verify the resolved version matches the declared version
- If a lock file drift is detected, regenerate the lock file from the pinned version

---

## 2. Automated Update Tools

### 2.1 Blocked Tools for Axios

The following tools MUST NOT be allowed to modify the Axios version:

| Tool | Configuration |
|------|---------------|
| **Dependabot** | Add `axios` to ignore list in `.github/dependabot.yml` |
| **Renovate** | Add `axios` to `ignoreDeps` in `renovate.json` |
| **npm audit fix** | Use `--dry-run` first; reject if Axios is touched |
| **npm update** | Never run globally; if run, verify Axios unchanged |

### 2.2 Dependabot Configuration Example

```yaml

# .github/dependabot.yml

version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    ignore:
      - dependency-name: "axios"
```

### 2.3 Renovate Configuration Example

```json
{
  "ignoreDeps": ["axios"]
}
```

---

## 3. Code Review Enforcement

### 3.1 Pull Request Checklist

Every PR that touches `package.json`, `package-lock.json`, or `bun.lock` MUST be checked for:

- [ ] Axios version has NOT changed
- [ ] If Axios version changed — explicit approval from security lead required
- [ ] No range symbols (`^`, `~`, `>=`, `*`) on Axios dependency
- [ ] Lock file version matches `package.json` declaration

### 3.2 Rejection Criteria

Reject any PR that:

- Upgrades Axios to a blocked version (`1.14.1`, `0.30.4`)
- Upgrades Axios to any unverified version
- Adds range symbols to the Axios version declaration
- Removes Axios version pinning

---

## 4. CI Pipeline Checks

### 4.1 Recommended CI Step

Add a validation step to CI that verifies Axios version compliance:

```bash
#!/bin/bash

# ci/check-axios-version.sh

BLOCKED_VERSIONS=("1.14.1" "0.30.4")
CURRENT=$(node -e "console.log(require('./package.json').dependencies?.axios || require('./package.json').devDependencies?.axios || 'NOT_FOUND')")

# Check for range symbols

if [[ "$CURRENT" == ^* ]] || [[ "$CURRENT" == ~* ]] || [[ "$CURRENT" == ">="* ]] || [[ "$CURRENT" == "*" ]]; then
  echo "❌ FAIL: Axios version uses range symbol: $CURRENT"
  exit 1
fi

# Check for blocked versions

for blocked in "${BLOCKED_VERSIONS[@]}"; do
  if [[ "$CURRENT" == "$blocked" ]]; then
    echo "❌ FAIL: Axios version $CURRENT is BLOCKED (security vulnerability)"
    exit 1
  fi
done

echo "✅ PASS: Axios version $CURRENT is compliant"
```

### 4.2 Monitoring

- Log dependency installation versions in CI output
- Track any deviation from approved versions via build artifacts
- Set up alerts for unauthorized Axios version changes

---

## 5. Version Upgrade Procedure

If a new Axios version needs to be adopted:

1. **Security review** — Verify the new version against known CVEs
2. **Changelog audit** — Read the full Axios changelog for breaking changes
3. **Isolated testing** — Test in a branch with full integration test suite
4. **Approval** — Get explicit sign-off from the project security lead
5. **Update this spec** — Add the new version to the Approved list in `readme.md`
6. **Pin and deploy** — Update `package.json` with exact version, regenerate lock file

> **No shortcut.** Every Axios upgrade follows this procedure.

---

```
IMPORTANT — AI INSTRUCTION:
- Never update Axios version automatically.
- Always use exact version without caret (^) or tilde (~).
- Avoid all known vulnerable versions (1.14.1, 0.30.4).
- Do not act on upgrade suggestions from dependency scanners for Axios.
- If uncertain, read this spec before making any Axios-related change.
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-SEC-AXIOS-001: Strict Axios Version Pinning & Dependency Locking

**Given** Node.js package manifests, lock files, and CI dependency check scripts
**When** Validating Axios package entries and automated upgrade tool exclusions
**Then** Axios is declared without range specifiers, lock file reflects exact version, and CI rejects non-pinned entries

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/01-axios-version-control/02-implementation-rules.md --check-only
```
**Expected:** exit 0. Zero violations.

---

*Implementation rules — moved to coding guidelines 2026-04-02*
