# Plugin Model — Adding a New Language

> **Version:** 1.0.0
> **Updated:** 2026-04-19

Each check script is a **language plugin**. Adding support for a new
language means writing one Python file under
`linters-cicd/checks/<rule>/<language>.py` and registering it in
`linters-cicd/checks/registry.json` — no orchestrator changes.

---

## File layout

```
linters-cicd/
├── checks/
│   ├── _lib/
│   │   ├── sarif.py          # SARIF 2.1.0 emitter (shared)
│   │   ├── walker.py         # File discovery with .gitignore respect
│   │   └── language.py       # Language detection by extension
│   ├── nested-if/
│   │   ├── go.py             # Go AST walker
│   │   ├── typescript.py     # TS AST walker (uses tree-sitter-typescript via py)
│   │   └── readme.md         # Rule rationale + spec link
│   ├── magic-strings/
│   │   ├── go.py
│   │   └── typescript.py
│   ├── function-length/
│   │   ├── go.py
│   │   └── typescript.py
│   └── registry.json         # rule-id → languages → script path
├── run-all.sh                # POSIX orchestrator
├── action.yml                # GitHub composite Action
├── ci/                       # platform-specific templates
└── install.sh                # one-liner installer
```

---

## Plugin contract

Every check script:

1. Is invoked as `python3 <script> --path <dir> [--format sarif|text]`.
2. Walks files matching its language extensions only.
3. Emits SARIF 2.1.0 per [`01-sarif-contract.md`](./01-sarif-contract.md).
4. Exits `0` (clean) / `1` (findings) / `2` (tool error).
5. Has a sibling `<plugin>_test.py` with at least one bad fixture and
   one good fixture under `linters-cicd/checks/<rule>/fixtures/`.

---

## Adding a new language (example: Rust for `nested-if`)

1. `linters-cicd/checks/nested-if/rust.py` — implement using
   `tree-sitter-rust` or regex-based detection.
2. Add fixtures in `linters-cicd/checks/nested-if/fixtures/rust/`.
3. Append to `registry.json`:

   ```json
   {
     "CODE-RED-001": {
       "go": "checks/nested-if/go.py",
       "typescript": "checks/nested-if/typescript.py",
       "rust": "checks/nested-if/rust.py"
     }
   }
   ```

4. CI picks it up automatically — no orchestrator change.

---

## Why Python (not Go) for plugins

- Tree-sitter has first-class Python bindings for every target language.
- Python 3 ships on every CI runner without extra install.
- AST walks are I/O-bound; Python is fast enough.
- Each plugin is < 150 lines — language-design concerns dominate, not
  performance.

---

*Part of [CI/CD Integration](./00-overview.md)*
