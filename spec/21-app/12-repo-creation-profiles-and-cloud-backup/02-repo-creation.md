# Repository Creation Engine

## Architecture & Lifecycle

The repository creation engine (`gitmap create`) handles zero-setup repository provisioning:

```
+---------------------------+
| 1. User Input             |  gitmap create my-api [--org auktvgo] [--public]
+---------------------------+
              |
              v
+---------------------------+
| 2. Profile Resolution     |  Loads default or specified profile/org from git_profiles.json
+---------------------------+
              |
              v
+---------------------------+
| 3. Local Scaffolding      |  os.MkdirAll, git init -b main, README.md, .gitignore
+---------------------------+
              |
              v
+---------------------------+
| 4. Initial Commit         |  git add . && git commit -m "feat: initial commit"
+---------------------------+
              |
              v
+---------------------------+
| 5. Remote Provisioning    |  gh repo create <target> --<visibility> --source=. --push
+---------------------------+
              |
              v
+---------------------------+
| 6. Catalog Tracking       |  Increments profile UsageCount, updates LastUsedAt
+---------------------------+
```

---

## Command Signatures

```bash
gitmap create [repo] <name> [flags]
```

### Safety & Invariants
- **Private By Default**: Unless `--public` is explicitly provided, all created repositories are created with `--private` visibility.
- **Local-Only Isolation**: Supplying `--no-remote` skips all network provisioning, ensuring offline safety.
- **AST Compliance**: Command constant `CmdCreate = "create"` registered under `constants_cli.go` and verified by `cmd_constants_test.go`.
