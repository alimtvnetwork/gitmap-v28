# Specification: Uppercase Semantic Commit Formatter (No Third Brackets) & `cfr` Alias

## 1. Context & User Directives
Recent commit history displayed bracket pollution ("third brackets" / square brackets `[...]` inside parentheses), such as:
```text
fix([workgroup]): support positional ssh key in bind and ...
docs([spec]): add bracketed module commit standards ...
feat([workgroup]): folder-specific workdirs ...
```
The user strictly instructed:
> "And do not write, let's say, commits like this third bracket, update the coding guideline for the execute prompt in V6. It should be done using the `gitmap` CFR CBF, and it should be in uppercase, and in bracket, there should be no third brackets for the commits. So that needs to be corrected everywhere with an example."

---

## 2. Authorized Commit Format Grammar

```text
CommitMessage ::= Type "(" Module ")" ":" <space> Summary
Type          ::= "FEAT" | "FIX" | "REFACTOR" | "CHORE" | "RELEASE" | "DOCS" | "TEST" | "PERF"
Module        ::= [A-Z0-9_\-\./]+  (Uppercase alphanumeric, hyphen, dot, or forward slash)
Summary       ::= [^\n\r]+
```

### Examples
- Feature: `FEAT(WORKGROUP): folder-specific workdirs and ssh key automation`
- Bugfix: `FIX(DB): resolve transaction lock in split-db`
- Refactor: `REFACTOR(ENGINE): optimize credential lookup cascade`
- Chore: `CHORE(SKILLS): sync gitmap skill across fleet repositories`
- Release: `RELEASE(V6.605.0): production version cut`

---

## 3. Strict Non-Negotiable Invariants
1. **Zero Third Brackets (`[` `]`):** No square brackets allowed in git commit messages.
2. **Automatic Sanitization:** If an agent or developer passes square brackets (e.g. `cpf "[auth]: add jwt"` or `cpf "feat([auth]): add jwt"`), the formatter automatically cleans and translates it to `FEAT(AUTH): add jwt`.
3. **Uppercase Types & Modules:** The verb and module are converted to uppercase.
4. **Command Verbs & Native `cfr` Registration:**
   - `gitmap cpf` $\rightarrow$ `FEAT(MODULE): summary`
   - `gitmap cpb` / `gitmap cbf` $\rightarrow$ `FIX(MODULE): summary`
   - `gitmap cfr` $\rightarrow$ `REFACTOR(MODULE): summary` (Commit Push Refactor)
   - `gitmap cpc` $\rightarrow$ `CHORE(MODULE): summary`
   - `gitmap cpr` $\rightarrow$ `RELEASE(MODULE): summary`
