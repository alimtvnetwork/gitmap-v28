# Acceptance Criteria: Static Analysis & Linter Enforcement

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Must Have

- [ ] Every supported language (8) has a dedicated linter spec
- [ ] All specs enforce identical thresholds: 15-line functions, 3 params, complexity ≤ 10
- [ ] SonarQube rule IDs mapped for every enforceable rule per language
- [ ] Integration checklist uses standardized table format with 🔲 status
- [ ] CI pipeline spec defines a unified quality gate referencing all 8 languages
- [ ] Cross-language rule matrix covers all 7 SonarQube rules across all 8 languages
- [ ] TypeScript ESLint spec cross-referenced from overview (lives in `02-typescript/`)

## Should Have

- [ ] Each spec includes Keywords and Scoring sections
- [ ] Exemption/suppression syntax documented per language
- [ ] GitHub Actions and GitLab CI templates provided in CI spec

## Won't Have (This Version)

- Pre-commit hook configurations (deferred to future iteration)
- IDE-specific settings files (`.vscode/`, `.idea/`)
