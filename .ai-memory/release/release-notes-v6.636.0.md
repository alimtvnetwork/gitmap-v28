# Release Notes: GitMap v6.636.0

## Release Details
- **Version:** `v6.636.0`
- **Date:** 2026-10-31
- **Architect & Author:** MD ALIM UL KARIM (`alimtvnetwork`)
- **Sponsor:** RISEUP ASIA LLC (`https://riseup-asia.com`)

---

## 🌟 Highlights

1. **Advanced `gitmap cat` Range, Query, and Line Inspection Engine:**
   - `-L, --lines <start:end>`: Slices line ranges (e.g. `gitmap cat file.go -L 140:220`).
   - `-q, --query <pattern>`: Case-insensitive query filter matching lines with original line numbers (e.g. `gitmap cat skill.md -q "ban" -l 15`).
   - `-n, --head <count>` & `--tail <count>`: Direct head and tail line slicing.
   - `-num, --number`: Displays 1-based gutter line numbers (`  140 | ...`).
   - `--raw`: Emits clean, borderless output.

2. **Total Ban on Ad-Hoc Python File Reading:**
   - Enshrined non-negotiable Rule 19 in `.agents/skills/gitmap/skill.md`.
   - Banned `gitmap py -c "open('...').read()"` and `lines=open('...').readlines()`.
   - Provided universal replacement matrix across all workflows.

3. **Sync Script Dual Versioning (V6 vs V7 Engine):**
   - V6 engine: `38-sync-prompts-skills-scripts.py` (unrestricted fleet sync).
   - V7 engine: `38-sync-prompts-skills-scripts-v7.py` (enforces calendar-aware morning commit cutoff < 08:00 AM on weekdays, auto-window scheduling, and bot identity governance).

4. **Fleet Synchronization:**
   - Hardened skills and prompts synchronized across all 43 connected repositories.

---

## Verification
- Unit tests: PASS (`cmdmacro` test suite 100% green).
- Live CLI verification: `gitmap cat <file> -L <start>:<end>`, `-q <query>`, and `--num` verified in terminal.
