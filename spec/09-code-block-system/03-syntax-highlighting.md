# Syntax Highlighting — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Library

- **highlight.js** v11+ (core only — tree-shakeable)
- Theme CSS: `highlight.js/styles/github-dark.css`
- Custom token colors override the theme (see Styling spec)

---

## Registered Languages

Languages are registered individually for bundle efficiency:

| Registration Name(s) | highlight.js Module | Notes |
|-----------------------|---------------------|-------|
| `typescript`, `ts`, `tsx`, `javascript`, `js` | `typescript` | JS uses TS grammar (superset) |
| `go`, `golang` | `go` | |
| `php` | `php` | |
| `css` | `css` | |
| `json` | `json` | |
| `bash`, `sh`, `shell` | `bash` | |
| `sql` | `sql` | |
| `rust` | `rust` | |
| `html`, `xml` | `xml` | HTML uses XML grammar |
| `yaml`, `yml` | `yaml` | |
| `markdown`, `md` | `markdown` | |

---

## Language Resolution Flow

```
Input: (code, lang) from markdown fence

1. normalizeLang(lang)
   - Trim, lowercase
   - Check against known groups: TS, JS, Go, all supported, plaintext
   - Return normalized or empty string

2. resolveDisplayLang(code, lang)
   - If no lang AND code looks like tree → return "tree"
   - Otherwise return normalized lang

3. highlightCode(code, lang)
   - If no lang AND tree-like → highlightAsTree()
   - If lang is registered → hljs.highlight(code, { language })
   - Otherwise → hljs.highlightAuto(code)
     - If auto returns plaintext AND tree-like → highlightAsTree()
   - All paths have try/catch → fallback to escapeHtml()
```

---

## Tree Structure Detection

A code block is detected as a "tree" (folder/file structure) when ANY of these match:

| Pattern | Regex | Example |
|---------|-------|---------|
| Box-drawing characters | `/[├└│─]/` | `├── src/` |
| Directory line | `/^\s*[A-Za-z0-9{}._-]+\/$/m` | `components/` |
| File line | `/^\s*[A-Za-z0-9{}._-]+\.[A-Za-z0-9_-]+\s*$/m` | `index.ts` |

---

## Tree Rendering

Each line is processed individually by `highlightTreeLine()`:

1. **Split comments**: anything after `#` is extracted as a comment
2. **Escape HTML** on the content portion
3. **Apply spans** in order:
   - Box-drawing chars → `<span class="tree-guide">` (50% opacity muted)
   - Ellipsis `...` → `<span class="tree-ellipsis">` (accent/pink color)
   - Directory names (ending `/`) → `<span class="tree-dir">📁 {name}</span>` (bold white)
   - File names (with extension) → `<span class="tree-file">📄 {name}</span>` (85% opacity)
4. **Comment** (if present) → `<span class="tree-comment">` (italic muted)

---

## HTML Escaping

`escapeHtml()` converts: `&` → `&amp;`, `<` → `&lt;`, `>` → `&gt;`

Used for:
- Raw code before tree highlighting
- Data attributes (`data-code`) for copy/download
- Checklist markdown encoding

---

## Syntax Token Colors

All colors use CSS custom properties (HSL format):

| Token Type | CSS Class(es) | Color Variable |
|------------|---------------|----------------|
| Keywords, types, built-ins | `.hljs-keyword`, `.hljs-type`, `.hljs-built_in` | `--primary` (purple) |
| Function/class names | `.hljs-title`, `.hljs-section` | `--foreground / 0.85` |
| Strings, attributes | `.hljs-string`, `.hljs-attr`, `.hljs-property` | `--accent` (pink) |
| Numbers, symbols | `.hljs-number`, `.hljs-symbol`, `.hljs-regexp` | `--warning` (amber) |
| Comments | `.hljs-comment`, `.hljs-quote` | `--muted-foreground` (italic) |
| Tree guides | `.tree-guide` | `--muted-foreground / 0.5` |
| Tree directories | `.tree-dir` | `--foreground` (bold) |
| Tree files | `.tree-file` | `--foreground / 0.85` |
| Tree ellipsis | `.tree-ellipsis` | `--accent` (pink) |
| Tree comments | `.tree-comment` | `--muted-foreground` (italic) |

---

*Syntax Highlighting — updated: 2026-04-08*
