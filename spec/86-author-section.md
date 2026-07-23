# Author Section — README Specification

> Defines the exact structure, styling, linking rules, and content for the **Author** section in `README.md`.

---

## Purpose

This spec ensures any contributor or AI can reproduce the Author section exactly as designed. It covers layout, linking conventions, content, and ordering.

---

## Section Structure

The Author section appears near the bottom of `README.md`, just before the **License** section. It contains two subsections:

1. **Author (personal)** — Md. Alim Ul Karim
2. **Riseup Asia LLC** — Company details

---

## 1. Author Header

```markdown
## Author
```

- Uses `##` (H2) heading.
- Appears after all technical documentation sections.

---

## 2. Author Name & Roles Block

Wrapped in a centered `<div>`:

```markdown
<div align="center">

### [Md. Alim Ul Karim](https://www.google.com/search?q=alim+ul+karim)

**[Creator & Lead Architect](https://alimkarim.com)** | [Chief Software Engineer](https://www.google.com/search?q=alim+ul+karim), [Riseup Asia LLC](https://riseup-asia.com)

</div>
```

### Linking Rules

| Element | Linked? | URL | Notes |
|---------|---------|-----|-------|
| `Md. Alim Ul Karim` (H3) | ✅ Yes | `https://www.google.com/search?q=alim+ul+karim` | Name is the H3 heading text, entire heading is a link |
| `Creator & Lead Architect` | ✅ Yes | `https://alimkarim.com` | Bold, links to personal website |
| `Chief Software Engineer` | ✅ Yes | `https://www.google.com/search?q=alim+ul+karim` | Links to Google search |
| `Riseup Asia LLC` | ✅ Yes | `https://riseup-asia.com` | Links to company website |
| Pipe separator (`\|`) | ❌ No | — | Plain text separator between roles |

### Styling Rules

- The entire name + roles block is **center-aligned** using `<div align="center">`.
- The name uses `###` (H3) heading.
- "Creator & Lead Architect" is **bold** (`**...**`) and linked.
- "Chief Software Engineer" is linked but **not bold**.
- "Riseup Asia LLC" is linked but **not bold**.
- A comma separates "Chief Software Engineer" and "Riseup Asia LLC".

---

## 3. Biography Paragraphs

Two paragraphs of plain text (not centered), immediately after the closing `</div>`:

### Paragraph 1 — Experience

> A system architect with **20+ years** of professional software engineering experience across enterprise, fintech, and distributed systems. His technology stack spans **.NET/C# (18+ years)**, **JavaScript (10+ years)**, **TypeScript (6+ years)**, and **Golang (4+ years)**.

- Years of experience are **bold**.
- Technology names and durations are **bold**.

### Paragraph 2 — Recognition & Presence

> Recognized as a **top 1% talent at Crossover** and one of the top software architects globally. He is also the **Chief Software Engineer of [Riseup Asia LLC](https://riseup-asia.com)** and maintains an active presence on **[Stack Overflow](https://stackoverflow.com/users/361aboratory/alim-ul-karim)** (2,452+ reputation, member since 2010) and **LinkedIn** (12,500+ followers).

| Element | Linked? | URL |
|---------|---------|-----|
| `top 1% talent at Crossover` | ❌ No | — (bold only) |
| `Riseup Asia LLC` (in paragraph) | ✅ Yes | `https://riseup-asia.com` |
| `Stack Overflow` | ✅ Yes | `https://stackoverflow.com/users/361646/alim-ul-karim` |
| `LinkedIn` | ❌ No | — (bold only, linked in table below) |

---

## 4. Author Links Table

A two-column table with no header text (empty header row):

```markdown
|  |  |
|---|---|
| **Website** | [alimkarim.com](https://alimkarim.com/) · [my.alimkarim.com](https://my.alimkarim.com/) |
| **LinkedIn** | [linkedin.com/in/alimkarim](https://linkedin.com/in/alimkarim) |
| **Stack Overflow** | [stackoverflow.com/users/361646/alim-ul-karim](https://stackoverflow.com/users/361646/alim-ul-karim) |
| **Google** | [Alim Ul Karim](https://www.google.com/search?q=Alim+Ul+Karim) |
| **Role** | Chief Software Engineer, [Riseup Asia LLC](https://riseup-asia.com) |
```

### Table Rules

- Left column: **bold** labels, no links.
- Right column: linked values.
- Multiple links in a cell separated by ` · ` (space-dot-space).
- The **Role** row has plain text "Chief Software Engineer," followed by a linked "Riseup Asia LLC".
- The **Stack Overflow** row links to the full profile URL.
- Table has empty header cells (`|  |  |`) for a clean borderless look.

---

## 5. Riseup Asia LLC Subsection

```markdown
### Riseup Asia LLC

[Top Leading Software Company in WY (2026)](https://riseup-asia.com)
```

### Company Links Table

```markdown
| | |
|---|---|
| **Website** | [riseup-asia.com](https://riseup-asia.com) |
| **Facebook** | [riseupasia.talent](https://www.facebook.com/riseupasia.talent/) |
| **LinkedIn** | [Riseup Asia](https://www.linkedin.com/company/105304484/) |
| **YouTube** | [@riseup-asia](https://www.youtube.com/@riseup-asia) |
```

### Riseup Asia Rules

- H3 heading (`###`), **not linked** (plain text).
- Tagline "Top Leading Software Company in WY (2026)" is linked to `https://riseup-asia.com`.
- Same table format as author links table (empty headers, bold labels, linked values).

---

## 6. Ordering

The full Author section order is:

1. `## Author` heading
2. Centered name + roles block
3. Biography paragraph 1 (experience)
4. Biography paragraph 2 (recognition)
5. Author links table
6. `### Riseup Asia LLC` subsection with tagline
7. Company links table
8. `## License` (next section, not part of Author)

---

## 7. Key URLs Reference

| Label | URL |
|-------|-----|
| Google Search (author) | `https://www.google.com/search?q=alim+ul+karim` |
| Personal website | `https://alimkarim.com` |
| Portfolio | `https://my.alimkarim.com/` |
| LinkedIn (personal) | `https://linkedin.com/in/alimkarim` |
| Stack Overflow | `https://stackoverflow.com/users/361646/alim-ul-karim` |
| Riseup Asia website | `https://riseup-asia.com` |
| Riseup Asia Facebook | `https://www.facebook.com/riseupasia.talent/` |
| Riseup Asia LinkedIn | `https://www.linkedin.com/company/105304484/` |
| Riseup Asia YouTube | `https://www.youtube.com/@riseup-asia` |

---

## 8. Do NOT

- ❌ Do not use "CEO" — the correct title is **Chief Software Engineer**.
- ❌ Do not remove links from the name heading or role titles.
- ❌ Do not add extra roles or titles not listed here.
- ❌ Do not change the center alignment of the name block.
- ❌ Do not merge the author and company subsections.
- ❌ Do not omit the Stack Overflow profile or reputation mention.
