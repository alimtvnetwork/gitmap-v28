---
name: auto-aliasing-and-error-storage-reset
description: Autonomously implement and verify repository auto-aliasing generation, separate aliases database table, bracket and tree-view alias rendering in scan tables, and error storage reset commands across GitMap.
---

# Auto Aliasing & Error Storage Reset Suite

## Overview
Implements automatic repository alias generation, persistence in a dedicated SQLite database table, bracketed and tree-view rendering in repository scan tables, and dedicated error storage reset commands across GitMap.

## Core Rules & Architecture
1. **Auto-Aliasing Rules**:
   - Word Boundaries: Split repository names by hyphens (`-`), underscores (`_`), and spaces (` `).
   - Acronym Formation: Combine first character of each word (e.g. `anti-gravity-manager` -> `agm`).
   - Compound Words: If single word without delimiters, find compound roots/syllables (e.g. `gitmap` -> `gm`).
   - WP Prefix: Any package with `wp-` / `wp` retains prefix (e.g. `wp-git-log` -> `wp-gl`, `wp-html-automated` -> `wp-ha`).
   - Presentations Prefix: Presentation repositories receive `prep-` prefix + short name (e.g. `presentations-repos/bsrm-presentation-hiltrax` -> `prep-bsrm`).
   - Single Word/Core Packages: Fall back gracefully if abbreviations cannot be cleanly derived.
2. **Database Persistence**:
   - Store aliases in a dedicated SQLite table (`repo_aliases`) linked by repo slug or path.
   - Automatically check, generate, and persist on scan or installation update.
3. **Table & Tree Rendering**:
   - Display primary alias in brackets (e.g. `[gm]`, `[wp-gl]`) in scan outputs.
   - Render multi-alias tree view when a repository has multiple registered aliases.
4. **Error Storage Reset**:
   - Implement `gitmap storage reset-errors` / `gitmap storage error-reset`.
   - Display clear command suggestion in storage and database inspect views.
5. **Coding Guidelines**:
   - Functions <= 8-15 lines, single return types, universal `*apperror.AppError`, affirmative booleans (`is*`, `has*`), Unix LF line endings.
