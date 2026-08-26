# Constants & Maps — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Language Labels (`LANG_LABELS`)

Maps fence identifiers to display names shown in the badge:

| Key(s) | Display Label |
|--------|---------------|
| `ts`, `tsx`, `typescript` | TypeScript |
| `js`, `javascript` | JavaScript |
| `go`, `golang` | Go |
| `php` | PHP |
| `css` | CSS |
| `json` | JSON |
| `bash` | Bash |
| `sh`, `shell` | Shell |
| `sql` | SQL |
| `rust` | Rust |
| `html` | HTML |
| `xml` | XML |
| `yaml`, `yml` | YAML |
| `md`, `markdown` | Markdown |
| `tree` | Structure |
| `text`, `""` | Plain Text |

---

## Language Accent Colors (`LANG_COLORS`)

HSL values (hue saturation% lightness%) used for `--lang-accent` and `--badge-color`:

| Language | HSL |
|----------|-----|
| TypeScript / TS / TSX | `99 83% 62%` |
| JavaScript / JS | `53 93% 54%` |
| Go / Golang | `194 66% 55%` |
| PHP | `234 45% 60%` |
| CSS | `264 55% 58%` |
| JSON | `38 92% 50%` |
| Bash / SH / Shell | `120 40% 55%` |
| SQL | `200 70% 55%` |
| Rust | `25 85% 55%` |
| HTML / XML | `12 80% 55%` |
| YAML / YML | `0 75% 55%` |
| Markdown / MD / Tree | `252 85% 60%` |
| **Default** (unlisted) | `220 10% 50%` |

---

## Language Extensions (`LANG_EXTENSIONS`)

File extensions used for the download filename (`code.{ext}`):

| Language | Extension |
|----------|-----------|
| TypeScript | `ts` |
| TSX | `tsx` |
| JavaScript / JS | `js` |
| Go | `go` |
| PHP | `php` |
| CSS | `css` |
| JSON | `json` |
| Bash / SH / Shell | `sh` |
| SQL | `sql` |
| Rust | `rs` |
| HTML | `html` |
| XML | `xml` |
| YAML / YML | `yaml` / `yml` |
| Markdown / MD | `md` |
| Tree / Text / "" | `txt` |

---

## Language Group Constants

Used for normalization:

| Constant | Values |
|----------|--------|
| `TYPESCRIPT_LANGS` | `["typescript", "ts", "tsx"]` |
| `JAVASCRIPT_LANGS` | `["javascript", "js"]` |
| `GO_LANGS` | `["go", "golang"]` |
| `PLAINTEXT_LANGS` | `["text", "plaintext", "plain", "tree"]` |
| `ALL_SUPPORTED_LANGS` | `["php", "css", "json", "bash", "sh", "shell", "sql", "rust", "html", "xml", "yaml", "yml", "markdown", "md"]` |

---

## Font Size Constants

| Constant | Value |
|----------|-------|
| `DEFAULT_FONT_SIZE` | `18` (px) |
| `MIN_FONT_SIZE` | `12` (px) |
| `MAX_FONT_SIZE` | `32` (px) |
| `FONT_SIZE_STEP` | `2` (px) |

---

## Timing Constants

| Constant | Value | Usage |
|----------|-------|-------|
| `COPY_FEEDBACK_DELAY` | `2000` (ms) | Duration of "Copied!" feedback |

---

## DOM Selectors (`codeBlockDomHelpers.ts`)

| Constant | Selector |
|----------|----------|
| `SELECTOR_COPY_BTN` | `.copy-code-btn` |
| `SELECTOR_DOWNLOAD_BTN` | `.download-code-btn` |
| `SELECTOR_FULLSCREEN_BTN` | `.fullscreen-code-btn` |
| `SELECTOR_CHECKLIST_COPY_BTN` | `.checklist-copy-btn` |
| `SELECTOR_CHECKLIST_EXPORT_BTN` | `.checklist-export-btn` |
| `SELECTOR_COPY_SELECTED_BTN` | `.copy-selected-btn` |
| `SELECTOR_CLEAR_SELECTED_BTN` | `.clear-selected-btn` |
| `SELECTOR_CODE_WRAPPER` | `.code-block-wrapper` |
| `SELECTOR_TOOL_EXCLUSIONS` | `.code-tool-btn, .code-font-controls, .copy-selected-bar` |
| `SELECTOR_LINE_NUMBER` | `.code-line-number` |
| `SELECTOR_CODE_LINE` | `.code-line` |

---

## CSS Class Constants

| Constant | Class |
|----------|-------|
| `LINE_PINNED_CLASS` | `line-pinned` |
| `LINE_HIGHLIGHT_CLASS` | `line-highlight` |
| `COPIED_FEEDBACK_COLOR` | `hsl(152 70% 50%)` |

---

*Constants & Maps — updated: 2026-04-08*
