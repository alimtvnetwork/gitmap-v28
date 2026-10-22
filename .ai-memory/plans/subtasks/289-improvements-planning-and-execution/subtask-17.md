# Subtask 17: Coding Guidelines Help Menu Tokyo Night Modernization

> **Subtask ID:** Subtask-17  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/cmdcg/cg_help.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Modernize `cli/cmdcg/cg_help.go` help screen (`gitmap cg --help`).
2. Replace hardcoded legacy ANSI color constants with Tokyo Night theme tokens (`termout.RoleSGR`, `termout.TokyoPrimary`, `termout.TokyoAccent`, etc.).
3. Enrich the menu with sectional badges (`[CORE]`, `[AI]`, `[SSOT]`) utilizing the new `BadgeText` and `BadgeColor` fields.
4. Align the banner styling with the modern GitMap visual identity.

---

## 2. Detailed Technical Plan

### 2.1 Update `buildCGHelpMenu`
- Set `BannerColor` to Tokyo Primary SGR or hex:
  ```go
  func buildCGHelpMenu() termout.HelpMenu {
  	return termout.HelpMenu{
  		Title:       "Coding Guidelines, Prompts & Quality Standards (gitmap cg)",
  		BannerColor: termout.RoleSGRBold(termout.RolePrimary),
  		UsageLines:  buildCGUsageLines(),
  		Sections: []termout.HelpSection{
  			buildCGOpsSection(),
  			buildCGPromptSection(),
  			buildCGVersioningSection(),
  		},
  		FooterFlags: buildCGFooterFlags(),
  		Tips:        buildCGTips(),
  	}
  }
  ```

### 2.2 Update Section Builders with Badges & Colors
- `buildCGOpsSection`:
  ```go
  func buildCGOpsSection() termout.HelpSection {
  	return termout.HelpSection{
  		Title:      "Coding Guidelines Operations",
  		Color:      termout.RoleSGRBold(termout.RoleOK),
  		BadgeText:  "CORE",
  		BadgeColor: termout.RoleSGR(termout.RoleFGSubtle),
  		Entries: []termout.CommandEntry{
  			{Command: "status, stat, ls", Description: "View tabulated Coding Guidelines install status & version across repos"},
  			{Command: "audit, check", Description: "Deep audit of coding standards compliance (specs, rules, memory)"},
  			{Command: "doctor, doc", Description: "Health check and configuration repair diagnostics for coding guidelines"},
  			{Command: "install", Description: "Install official Coding Guidelines standards into target repositories"},
  			{Command: "update", Description: "In-place update of Coding Guidelines for installed repositories"},
  			{Command: "diff, drift", Description: "Compare local guidelines version with latest upstream canonical standards"},
  		},
  	}
  }
  ```
- `buildCGPromptSection`:
  ```go
  func buildCGPromptSection() termout.HelpSection {
  	return termout.HelpSection{
  		Title:      "Prompt Architect & AI Standards",
  		Color:      termout.RoleSGRBold(termout.RoleAccent),
  		BadgeText:  "AI",
  		BadgeColor: termout.RoleSGR(termout.RoleAccent),
  		Entries: []termout.CommandEntry{
  			{Command: "install-prompts", Description: "Install Prompt Architect v2 suite and canonical prompt library"},
  			{Command: "update-prompts", Description: "In-place update of Prompt Architect templates across target repos"},
  			{Command: "prompts-status", Description: "Check installed status and version of Prompt Architect"},
  			{Command: "prompts-version", Description: "Show Prompt Architect version installed in repositories"},
  			{Command: "sync, sy", Description: "Unified sync of both Coding Guidelines and Prompt Architect"},
  		},
  	}
  }
  ```
- `buildCGVersioningSection`:
  ```go
  func buildCGVersioningSection() termout.HelpSection {
  	return termout.HelpSection{
  		Title:      "Canonical Versioning & SSOT",
  		Color:      termout.RoleSGRBold(termout.RoleWarn),
  		BadgeText:  "SSOT",
  		BadgeColor: termout.RoleSGR(termout.RoleWarn),
  		Entries: []termout.CommandEntry{
  			{Command: "version, ver, -v", Description: "Display installed Coding Guidelines version"},
  			{Command: "install-version-json, init-version", Description: "Initialize canonical version.json SSOT with version tag"},
  			{Command: "general, general-install", Description: "Multi-repo installation suite across discovered workspaces"},
  		},
  	}
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Function Limit:** All section builder functions remain <= 15 lines.
- **Constant Decoupling:** Replaces raw color constants with theme-aware `termout.RoleSGR` calls.
- **Return Formatting:** Single blank line after each closing brace and before return.

---

## 4. Verification Criteria
- [x] Executing `gitmap cg -h` displays Tokyo Night styled banner.
- [x] Section headers display bracketed tags: `Coding Guidelines Operations [OPS]:`, `Prompt Architect & AI Standards [AI]:`, `Canonical Versioning & SSOT [SSOT]:`.
- [x] Colors conform to active `--theme` selection.
- [x] Clean suppression of colors when piped to files or grep.
