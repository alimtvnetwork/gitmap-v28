# Subtask 03: Suggestion Type Object Domain Model & Multi-Tier Resolution Engine

> **Subtask ID:** Subtask-03  
> **Parent Plan:** `.ai-memory/plans/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design.md`  
> **Target Path:** `.ai-memory/plans/subtasks/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design/03-suggestion-type-object-and-engine-design.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Define strongly typed suggestion objects in `cli/suggestion/types.go` (`Suggestion`, `SuggestionGroup`, `SuggestionCategory`, `ActionType`).
2. Implement 4-tier resolution engine (`cli/suggestion/engine.go`):
   - Tier 1: Canonical Alias Registry.
   - Tier 2: Levenshtein distance & QGram fuzzy matching (distance <= 3).
   - Tier 3: Subsequence / Prefix matcher.
   - Tier 4: Semantic Context & Synonym Dictionary.
3. Implement pluggable renderers (`cli/suggestion/render.go`):
   - Catppuccin terminal box format (`RenderFormatBox`).
   - Compact inline format (`RenderFormatCompact`).
   - Machine-readable JSON format (`RenderFormatJSON`).

---

## 2. Detailed Technical Plan

### 2.1 Suggestion Domain Types (`cli/suggestion/types.go`)
- Create package `suggestion`.
- Enums for `SuggestionCategory` (`typo`, `subcommand`, `flag`, `remediation`, `intent`, `resource`).
- Enums for `ActionType` (`auto-run`, `manual`, `copy`).
- Struct `Suggestion`:
  ```go
  type Suggestion struct {
      Command     string             `json:"command"`
      Description string             `json:"description"`
      Category    SuggestionCategory `json:"category"`
      Confidence  float64            `json:"confidence"`
      ActionType  ActionType         `json:"actionType"`
      Context     map[string]any     `json:"context,omitempty"`
  }
  ```
- Struct `SuggestionGroup`:
  ```go
  type SuggestionGroup struct {
      Title       string       `json:"title"`
      Reason      string       `json:"reason"`
      InputToken  string       `json:"inputToken,omitempty"`
      Suggestions []Suggestion `json:"suggestions"`
  }
  ```

### 2.2 Multi-Tier Matcher Implementation (`cli/suggestion/engine.go`)
- Maintain command catalog and alias map.
- Fast Levenshtein distance calculation with early exit when threshold > 3.
- Confidence scoring:
  $$\text{Confidence} = 1.0 - \left(\frac{\text{Distance}}{\max(\text{len}(A), \text{len}(B))}\right)$$
- Semantic synonym mapping:
  * `"docker"` -> `gitmap cluster` / `gitmap nodes`
  * `"pull error"` -> `gitmap pe`
  * `"grep"` -> `gitmap aum search`
  * `"remove"` -> `gitmap rm` / `gitmap clean`

### 2.3 Terminal Box Rendering (`cli/suggestion/render.go`)
- Uses ANSI color codes and Unicode box-drawing characters (`╭─`, `│`, `╰─`).
- Highlights closest match in bright cyan or Catppuccin blue.
- Displays confidence percentages and short descriptions.

---

## 3. Verification Criteria
- [ ] Unit tests in `cli/suggestion/engine_test.go` cover exact matches, typos with distance 1-3, and unknown tokens with distance > 3.
- [ ] Confidence scores correctly computed between 0.0 and 1.0.
- [ ] Renderers output clean, non-overlapping box borders and valid JSON envelopes.
