# CI Issue 83: Sponsor/SEO Templates Minimum Count Under State DB Seeding RCA

## 1. Reproduction & Symptoms
- **CI Run ID:** [36279852280](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36279852280)
- **Failing Workflows:**
  - `Cross-Platform Build / ubuntu-latest / go build + test (Step: go test ./... (no cache))`
  - `Cross-Platform Build / macos-latest / go build + test (Step: go test ./... (no cache))`
- **Symptoms:**
  `TestCatalogTemplatesCount` in `cli/cmdprompttemplate` failed:
  ```text
  catalog_test.go:20: expected at least 20 sponsor/seo templates, got 1
  FAIL: github.com/alimtvnetwork/gitmap-v28/cli/cmdprompttemplate 0.028s
  ```

## 2. Root Cause Analysis
1. **Partial State DB Seeding**: In `cli/store/templates_split_db.go`, default schema initialization seeded 1 default SEO template (`tpl-seo-sponsor-default`).
2. **Short-Circuiting in `GetSponsorCategoryTemplates()`**: In `cli/cmdprompttemplate/catalog_sponsor.go`, `GetSponsorCategoryTemplates()` checked `if stateItems := loadStateDBSEOTemplates(); len(stateItems) > 0 { return stateItems }`. Because 1 template existed in the state database, `len(stateItems) > 0` was satisfied, returning only 1 item and preventing fallback to `buildGenericArchitectureTemplates()` (which contains the 20 required catalog entries).
3. **Unit Test Expectation**: `TestCatalogTemplatesCount` asserts that `len(sp) >= 20` for catalog completeness.

## 3. Code Fix
- In `cli/cmdprompttemplate/catalog_sponsor.go` (`GetSponsorCategoryTemplates`):
  - If `len(stateItems) >= 20`, return `stateItems`.
  - If `len(stateItems) == 0`, return `generic`.
  - If `len(stateItems) < 20`, append `stateItems` to `generic` so the catalog always satisfies the >= 20 contract while retaining user/state templates.

## 4. Prevention
- When integrating dynamic SQLite state stores with existing static catalog registries, ensure default fallbacks combine or supplement partial datasets rather than prematurely short-circuiting when item counts fall below contract minimums.
