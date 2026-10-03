# SARIF 2.1.0 Output Contract (AI Execution Prompt)

> **/goal** Emit strictly compliant SARIF 2.1.0 JSON payloads from every check script to integrate seamlessly with GitHub Actions, GitLab, and Azure DevOps.
> **/learn** Master SARIF top-level schemas (`runs`, `tool.driver`, `rules`, `results`), relative URI resolution, and exact severity level mappings (`error`, `warning`, `note`).

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Ensure every emitted SARIF report adheres to the official SARIF 2.1.0 schema with required driver, rule, and result fields.
- [ ] `/learn` Never output absolute filesystem paths; always use relative workspace paths in `artifactLocation.uri`.
- [ ] `/goal` Map CODE RED findings to SARIF `error` level to guarantee pipeline merge blocking.
- [ ] `/learn` Validate SARIF outputs using `linters-cicd/scripts/validate-sarif.py` and directory compliance via autofixer.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Version:** 1.0.0
> **Updated:** 2026-04-19

Every check script in `linters-cicd/checks/` MUST emit SARIF 2.1.0
conforming to this contract when run with `--format sarif`.

---

## Top-level shape

```json
{
  "$schema": "https://json.schemastore.org/sarif-2.1.0.json",
  "version": "2.1.0",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "coding-guidelines-<check-id>",
          "version": "1.0.0",
          "informationUri": "https://github.com/alimtvnetwork/coding-guidelines-v24",
          "rules": [
            {
              "id": "CODE-RED-001",
              "name": "NoNestedIf",
              "shortDescription": { "text": "Nested if statements are forbidden" },
              "helpUri": "https://github.com/alimtvnetwork/coding-guidelines-v24/blob/main/02-spec/02-coding-guidelines/01-cross-language/04-code-style/readme.md"
            }
          ]
        }
      },
      "results": [
        {
          "ruleId": "CODE-RED-001",
          "level": "error",
          "message": { "text": "Nested if at depth 3 — extract guard clauses." },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": { "uri": "src/payments/charge.ts" },
                "region": { "startLine": 42, "startColumn": 5 }
              }
            }
          ]
        }
      ]
    }
  ]
}
```

---

## Required fields per result

| Field | Required | Notes |
|-------|----------|-------|
| `ruleId` | ✅ | `CODE-RED-NNN` from `06-rules-mapping.md` |
| `level` | ✅ | `error` for CODE RED, `warning` for STYLE |
| `message.text` | ✅ | Human sentence ending in a period |
| `locations[].physicalLocation.artifactLocation.uri` | ✅ | Path **relative** to the scanned root |
| `locations[].physicalLocation.region.startLine` | ✅ | 1-indexed |
| `locations[].physicalLocation.region.startColumn` | ⚠️ | Optional but recommended |

---

## Severity mapping

| Spec severity | SARIF `level` | CI behavior |
|---------------|---------------|-------------|
| 🔴 CODE RED | `error` | Block merge |
| 🟡 STYLE | `warning` | Annotate, do not block |
| ℹ️ INFO | `note` | Surface in report only |

---

## Merging multiple check outputs

`run-all.sh` calls each check, then merges outputs into a single SARIF
file with one `runs[]` entry per tool. Consumers pick this single file
up via `upload-sarif` (GitHub) or `Code Quality: sarif` (GitLab via
sarif-to-codequality converter).

---

## Validation

The `linters-cicd/scripts/validate-sarif.py` script validates every
emitted file against the official SARIF 2.1.0 schema. CI runs this on
every PR to the linter pack itself.

---

*Part of [CI/CD Integration](./readme.md)*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-cicd-integration/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-CI-002: SARIF 2.1.0 Contract Conformance and Validation

**Given** CI/CD pipeline infrastructure and linter configurations.
**When** Audited against this integration specification.
**Then** Zero configuration or SARIF contract defects exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration --check-only
```
**Expected:** exit 0. Zero violations.
