# Python Array Constants, Dynamic Enum Generation & Config Compilation — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-python-arrays`, `cg-python-enums`, `python-array-constants`, `python-dynamic-enum`, `compile-array-constants`, `audit python arrays`, `cors-origins-constants`, `python-config-constants`

> [!IMPORTANT]
> Prompt Version: 2.3.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this prompt, header alert block, or incoming user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and fix hardcoded array combinations, cartesian string permutations (e.g. host-port combinations), mutable module-level lists, and untyped configurations in Python code. Modifying source files directly, enforce immutable `Final` tuples, orthogonal configuration classes, deterministic generator functions, dynamic runtime Enums ending with the mandatory `*Type` suffix, and compact functions (target <= 8 lines, hard cap <= 15 lines) until 100% compliant.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory hardcoded permutation lists, mutable constants, and repeated origin/host arrays without truncation.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/readme.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring: replace hardcoded combinatorial arrays with orthogonal constant classes (`DevHosts`, `DevPorts`), generator functions (`build_http_origins`), and immutable `Final[Tuple[...]]` collections.
6. [ ] /goal Phase 2 (Step B): Enforce dynamic Enum naming with the mandatory `*Type` suffix (`LocalhostOriginType`), providing strict typing and IDE autocompletion for generated items.
7. [ ] /goal Phase 2 (Step C): Enforce <= 8–15 line function decomposition on all consuming functions and middleware config builders.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/14-constants-enums.md` for constants and enums standards.
15. [ ] /learn Ingest `02-spec/02-coding-guidelines/` for domain-specific architectural specifications.
16. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
17. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Codebase for Combinatorial Arrays, Write .ai-memory/plans/pending/ Spec, Create .ai-memory/plans/subtasks/, Verify/Create Linter Hook)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit Code, Extract Parameterized Constants, Build Dynamic Enums, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Python Array Constants & Dynamic Enum Architecture

Hardcoding cartesian products of configuration values (e.g. `http://localhost:5173`, `http://127.0.0.1:5173`, `http://localhost:5174`, `http://127.0.0.1:5174`) directly into functions leads to fragile maintenance. Adding or updating a port requires manually modifying repeated string literals across multiple call sites.

---

### 1. Mandatory Array Constant & Enum Standards

1. **Immutability Mandate (Tuple Over List for Constants):**
   - NEVER use mutable `list` definitions (`origins = ["http://..."]`) for module-level constants.
   - All constant sequences MUST be defined as immutable `tuple` wrapped in `typing.Final`:
     ```python
     from typing import Final, Tuple

     # ✅ REQUIRED: Immutable tuple for constant sequences
     ACTIVE_PORTS: Final[Tuple[int, ...]] = (5173, 5174)
     ```

2. **Orthogonal Separation of Concerns:**
   - Decompose multi-part configurations into their primitive dimensions:
     - Hostnames in `DevHosts`
     - Port numbers in `DevPorts`
     - Schemes (e.g. `http`, `https`)
     - External/Extension wildcards in `ExtensionOrigin`

3. **Deterministic Generator Functions:**
   - Combine orthogonal configurations using pure generator functions with typed signatures:
     ```python
     def build_http_origins(
         ports: Iterable[int],
         hosts: Iterable[str] = DevHosts.ALL,
         scheme: str = "http",
     ) -> Tuple[str, ...]:
         return tuple(f"{scheme}://{host}:{port}" for port in ports for host in hosts)
     ```

4. **Dynamic Enum Construction with Mandatory `*Type` Suffix:**
   - Every Enum type name MUST end with `Type` (e.g. `LocalhostOriginType`), adhering to repository-wide enum conventions.
   - When generating an Enum dynamically from compiled array constants, use the functional `Enum(...)` constructor:
     ```python
     LocalhostOriginType = Enum(
         "LocalhostOriginType",
         {
             f"PORT_{origin.split(':')[-1]}_{'LOCALHOST' if 'localhost' in origin else 'LOOPBACK'}": origin
             for origin in COMPILED_DEV_ORIGIN_ITEMS
         },
         type=str,
     )
     ```

5. **Function Length Cap (Target <= 8 lines, Hard Cap <= 15 lines):**
   - Consuming functions (like CORS middleware configurators) MUST remain concise (<= 8 lines of body logic).
   - Functions must focus purely on merging and business decisions, offloading array generation to pre-compiled constants.

---

## 2. Code Review Reference: Before vs. After Architecture

### 2.1 Before: Fragile Hardcoded Combinatorial Arrays

```python
# ❌ ANTI-PATTERN: Hardcoded string permutations repeating host + port combinations
def configure_cors_origins_before(settings) -> list[str]:
    origins = list(settings.cors_origins) if settings.cors_origins else []

    # Problem: Adding or changing a port requires manually updating every host permutation.
    for origin in (
        "http://localhost:5173",
        "http://127.0.0.1:5173",
        "http://localhost:5174",
        "http://127.0.0.1:5174",
    ):
        if origin not in origins:
            origins.append(origin)

    if "chrome-extension://*" not in origins:
        origins.append("chrome-extension://*")

    return origins
```

---

### 2.2 Constants & Dynamic Enum Builder (`constants/cors.py`)

```python
"""CORS origin constants and dynamic enum configuration."""

from enum import Enum
from typing import Final, Iterable, Tuple

class DevHosts:
    """Standard loopback hostnames used in local development."""
    LOCALHOST: Final[str] = "localhost"
    IPV4_LOOPBACK: Final[str] = "127.0.0.1"

    ALL: Final[Tuple[str, ...]] = (LOCALHOST, IPV4_LOOPBACK)

class DevPorts:
    """Configurable ports for local frontends (Vite, Next.js, Storybook, etc.)."""
    VITE_PRIMARY: Final[int] = 5173
    VITE_SECONDARY: Final[int] = 5174

    # Switchable/extensible array of active development ports
    ACTIVE_PORTS: Final[Tuple[int, ...]] = (
        VITE_PRIMARY,
        VITE_SECONDARY,
    )

class ExtensionOrigin:
    """Browser extension origins."""
    CHROME_WILDCARD: Final[str] = "chrome-extension://*"

def build_http_origins(
    ports: Iterable[int],
    hosts: Iterable[str] = DevHosts.ALL,
    scheme: str = "http",
) -> Tuple[str, ...]:
    """Generates a normalized, unique list of HTTP origin URLs from hosts and ports."""
    return tuple(
        f"{scheme}://{host}:{port}"
        for port in ports
        for host in hosts
    )

# Step A: Generate the array of origin items from the active port list
COMPILED_DEV_ORIGIN_ITEMS: Final[Tuple[str, ...]] = build_http_origins(DevPorts.ACTIVE_PORTS)

# Step B: Dynamically construct an Enum with mandatory *Type suffix for strict typing & autocomplete
LocalhostOriginType = Enum(
    "LocalhostOriginType",
    {
        f"PORT_{origin.split(':')[-1]}_{'LOCALHOST' if 'localhost' in origin else 'LOOPBACK'}": origin
        for origin in COMPILED_DEV_ORIGIN_ITEMS
    },
    type=str,
)

class CorsDefaults:
    """Centralized, immutable defaults ready for injection into middleware."""
    DEFAULT_DEV_ORIGINS: Final[Tuple[str, ...]] = COMPILED_DEV_ORIGIN_ITEMS
    DEFAULT_EXTENSION_ORIGINS: Final[Tuple[str, ...]] = (ExtensionOrigin.CHROME_WILDCARD,)

    # Master compiled collection of all default origins
    ALL_DEFAULT_ORIGINS: Final[Tuple[str, ...]] = (
        DEFAULT_DEV_ORIGINS + DEFAULT_EXTENSION_ORIGINS
    )
```

---

### 2.3 After: Decoupled Consuming Configuration (< 8 Lines)

```python
def configure_cors_origins(settings) -> list[str]:
    """Attach CORS config ensuring default localhost and extension origins are allowed.

    Adding or removing ports now requires updating only `DevPorts.ACTIVE_PORTS`.
    """
    origins: list[str] = list(settings.cors_origins) if settings.cors_origins else []

    for default_origin in CorsDefaults.ALL_DEFAULT_ORIGINS:
        if default_origin not in origins:
            origins.append(default_origin)

    return origins
```

---

### 2.4 Modular Decomposition for AI Clarity

To ensure absolute clarity for AI agents during refactoring, notice how responsibilities are cleanly isolated:

| Component | Responsibility | Sizing Tier |
|---|---|---|
| `DevHosts` / `DevPorts` | Raw, single-purpose primitives | $\le 10$ lines |
| `build_http_origins` | Pure combinatorial compilation function | $\le 6$ lines |
| `LocalhostOriginType` | Dynamic enum with `*Type` suffix for type safety | $\le 8$ lines |
| `CorsDefaults` | Composed immutable master origin collections | $\le 8$ lines |
| `configure_cors_origins` | Consumer function performing idempotent merge | $\le 8$ lines |

---

## 3. Canonical Golang Equivalent Architecture

For cross-language parity, here is how the same dynamic origin compilation is represented in Go following repository standards:

```go
package cors

import "fmt"

const (
	HostLocalhost    = "localhost"
	HostIPv4Loopback = "127.0.0.1"

	PortVitePrimary   = 5173
	PortViteSecondary = 5174

	ChromeExtensionWildcard = "chrome-extension://*"
)

// ActiveDevPorts lists all currently active frontend development ports.
var ActiveDevPorts = []int{
	PortVitePrimary,
	PortViteSecondary,
}

// DevHosts lists all standard local loopback hostnames.
var DevHosts = []string{
	HostLocalhost,
	HostIPv4Loopback,
}

// BuildHTTPOrigins generates normalized HTTP origin URLs from hosts and ports.
func BuildHTTPOrigins(hosts []string, ports []int) []string {
	origins := make([]string, 0, len(hosts)*len(ports))

	for _, port := range ports {
		for _, host := range hosts {
			origins = append(origins, fmt.Sprintf("http://%s:%d", host, port))
		}
	}

	return origins
}
```

---

## 4. Verification Checklist

- [x] **Zero Hardcoded Permutation Lists:** Cartesian combinations assembled via parameterized generators.
- [x] **Immutable Sequence Constants:** Replaced mutable `list` constants with typed `Final[Tuple[...]]`.
- [x] **Mandatory `*Type` Enum Suffix:** Dynamic Enum strictly named `LocalhostOriginType`.
- [x] **Decoupled Extension Origins:** Wildcards and third-party origins isolated from localhost primitives.
- [x] **Function Length Limit:** Consuming functions and helpers strictly <= 8–15 lines.
