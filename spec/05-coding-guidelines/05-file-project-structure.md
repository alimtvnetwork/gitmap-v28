# File & Project Structure — Organization Standards

## Overview

Rules for directory layout, module boundaries, file sizing, and
import ordering. Applies to TypeScript and Go projects.

---

## 1. Directory Layout — Group by Responsibility

Organize files by what they do, not by what they are.

### TypeScript (React)

```
src/
  components/       # Reusable UI components
    projects/       # Feature-scoped component group
      types.ts      # Named types for this feature
      TypeBadge.tsx
      ProjectDetailDialog.tsx
  constants/        # Shared enums and constants
    index.ts
  data/             # Static data and configuration
  hooks/            # Custom React hooks
  lib/              # Utility functions
  pages/            # Route-level page components
  types/            # Shared cross-feature types
```

### Go (CLI)

```
cmd/                # CLI routing and flag parsing
config/             # Config loading and merging
constants/          # All shared string literals
scanner/            # Directory walking
mapper/             # Data transformation
formatter/          # Output rendering
model/              # Shared data structures
store/              # Database access
```

---

## 2. One Responsibility Per File

Each file owns a single concern. Signals to split:

| Signal | Action |
|--------|--------|
| File exceeds 200 lines | Split by responsibility |
| File has 2+ unrelated function groups | Separate files |
| File mixes types and logic | Extract `types.ts` or `model.go` |
| File mixes constants and components | Extract `constants.ts` |

---

## 3. One Responsibility Per Directory

Each directory owns a single domain. Never mix unrelated features
in the same folder.

- `components/projects/` — only project-related components.
- `components/spec/` — only spec index components.
- A new feature gets its own directory.

---

## 4. Import Ordering

Group imports in a consistent order separated by blank lines.

### TypeScript

```ts
// 1. External libraries
import { useState, useCallback } from "react";
import { motion } from "framer-motion";

// 2. Internal absolute imports (aliases)
import { ProjectTypes } from "@/components/projects/TypeBadge";
import { RepoStatus } from "@/constants";
import type { ProjectTypeConfig } from "@/components/projects/types";

// 3. Relative imports
import { formatLabel } from "./utils";
import type { LocalType } from "./types";
```

### Go

```go
// 1. Standard library
import (
    "fmt"
    "os"
    "path/filepath"
)

// 2. External packages
import (
    "github.com/spf13/cobra"
)

// 3. Internal packages
import (
    "gitmap-v28/constants"
    "gitmap-v28/model"
    "gitmap-v28/store"
)
```

---

## 5. Index Files — Barrel Exports

Use `index.ts` barrel files sparingly and only for public APIs
of a module. Never re-export everything blindly.

```ts
// constants/index.ts — curated public API
export { RepoStatus } from "./repoStatus";
export { TerminalLineType } from "./terminalLineType";
export { FILTER_ALL } from "./filters";
```

Avoid deep barrel chains (barrel importing from another barrel).

---

## 6. Co-location Rules

Keep related files together:

| File | Location |
|------|----------|
| Component types | Same directory as component (`types.ts`) |
| Component constants | Same directory or `constants/` if shared |
| Component tests | Same directory (`Component.test.tsx`) |
| Shared types | `src/types/` or `model/` |
| Shared constants | `src/constants/` or `constants/` |

---

## 7. Naming Conventions for Files

| Language | Convention | Example |
|----------|-----------|---------|
| TypeScript components | PascalCase | `TypeBadge.tsx` |
| TypeScript utilities | camelCase | `specData.ts` |
| Go source files | lowercase, single word | `terminal.go` |
| Spec documents | kebab-case with numeric prefix | `01-overview.md` |

---

## 8. Leaf Packages

`model`, `constants`, and `types` are leaf packages:

- Imported by many other packages.
- Never import project-specific packages themselves.
- Breaking this rule creates circular dependencies.

---

## References

- Code Quality Improvement: `spec/05-coding-guidelines/01-code-quality-improvement.md`
- Go Code Style: `spec/05-coding-guidelines/02-go-code-style.md`
- Naming Conventions: `spec/05-coding-guidelines/03-naming-conventions.md`
