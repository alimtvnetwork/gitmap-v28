# SOLID Principles Reference

> **Version:** 1.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** All languages

---

## Overview

The SOLID principles guide maintainable, extensible software design. This spec defines how each principle applies across the project's languages with concrete examples and anti-patterns.

---

## Keywords

`solid` · `single-responsibility` · `open-closed` · `liskov` · `interface-segregation` · `dependency-inversion` · `architecture` · `design-principles` · `decoupling`

---

## S — Single Responsibility Principle (SRP)

A class, struct, or module should have **one reason to change**. If a type handles multiple concerns, split it.

```go
// ❌ FORBIDDEN: PluginService handles CRUD + validation + formatting
type PluginService struct { ... }
func (s *PluginService) Create(...) { ... }
func (s *PluginService) Validate(...) { ... }
func (s *PluginService) FormatResponse(...) { ... }

// ✅ REQUIRED: Separate concerns
type PluginService struct { ... }    // CRUD only
type PluginValidator struct { ... }  // validation only
type PluginFormatter struct { ... }  // formatting only
```

```typescript
// ❌ FORBIDDEN: Component does fetching + rendering + formatting
const UserProfile = () => {
    const [user, setUser] = useState(null);
    useEffect(() => { fetch('/api/user').then(...) }, []);
    const formatDate = (d: Date) => { ... };
    return <div>...</div>;
};

// ✅ REQUIRED: Separated
const useUser = () => { /* fetch logic */ };
const formatDate = (d: Date) => { /* formatting */ };
const UserProfile = () => {
    const user = useUser();
    return <div>...</div>;
};
```

---

## O — Open/Closed Principle (OCP)

Types should be **open for extension, closed for modification**. Add behavior via new types or composition, not by editing existing code.

```go
// ❌ FORBIDDEN: Adding new export types requires modifying existing function
func Export(data []Item, format string) error {
    if format == "csv" { ... }
    if format == "json" { ... }
    // Adding "xml" means editing this function
}

// ✅ REQUIRED: Interface-based extension
type Exporter interface {
    Export(data []Item) error
}

type CsvExporter struct{}
type JsonExporter struct{}
type XmlExporter struct{} // New format = new type, no edits to existing code
```

---

## L — Liskov Substitution Principle (LSP)

Subtypes must be substitutable for their base types without breaking behavior. In Go, this means interface implementations must honor the contract.

```go
// ❌ FORBIDDEN: Implementation that violates the interface contract
type Store interface {
    GetById(ctx context.Context, id int64) apperror.Result[Plugin]
}

type CachedStore struct{}
func (s *CachedStore) GetById(ctx context.Context, id int64) apperror.Result[Plugin] {
    // ❌ Returns stale data without error — violates expected freshness contract
}

// ✅ REQUIRED: Honor the contract or document cache behavior
func (s *CachedStore) GetById(ctx context.Context, id int64) apperror.Result[Plugin] {
    // Returns cached data if fresh, falls back to DB if stale
}
```

---

## I — Interface Segregation Principle (ISP)

Clients should not depend on methods they don't use. Prefer **small, focused interfaces** over large ones.

```go
// ❌ FORBIDDEN: Fat interface forces implementers to stub unused methods
type PluginManager interface {
    Create(ctx context.Context, p Plugin) apperror.Result[Plugin]
    Update(ctx context.Context, p Plugin) apperror.Result[Plugin]
    Delete(ctx context.Context, id int64) *apperror.AppError
    Export(ctx context.Context, id int64) apperror.Result[[]byte]
    SendNotification(ctx context.Context, id int64) *apperror.AppError
}

// ✅ REQUIRED: Segregated interfaces — one method per concern when possible
type PluginCreator interface {
    Create(ctx context.Context, p Plugin) apperror.Result[Plugin]
}

type PluginExporter interface {
    Export(ctx context.Context, id int64) apperror.Result[[]byte]
}

type PluginNotifier interface {
    SendNotification(ctx context.Context, id int64) *apperror.AppError
}
```

```typescript
// ❌ FORBIDDEN: Component accepts unused props
interface ButtonProps {
    label: string;
    onClick: () => void;
    onHover: () => void;     // Most consumers don't need this
    analytics: AnalyticsConfig; // Only used in 1 place
}

// ✅ REQUIRED: Core props only, extend when needed
interface ButtonProps {
    label: string;
    onClick: () => void;
}

interface AnalyticsButtonProps extends ButtonProps {
    analytics: AnalyticsConfig;
}
```

---

## D — Dependency Inversion Principle (DIP)

High-level modules must not depend on low-level modules. Both should depend on **abstractions** (interfaces).

```go
// ❌ FORBIDDEN: Service directly depends on concrete database implementation
type PluginService struct {
    db *sql.DB  // concrete dependency
}

// ✅ REQUIRED: Depend on interface
type PluginStore interface {
    GetById(ctx context.Context, id int64) apperror.Result[Plugin]
}

type PluginService struct {
    store PluginStore  // abstraction — can be DB, cache, mock, etc.
}
```

```typescript
// ❌ FORBIDDEN: Hook hardcoded to fetch
const usePlugins = () => {
    return fetch('/api/plugins').then(r => r.json());
};

// ✅ REQUIRED: Injectable data source
const usePlugins = (fetcher: () => Promise<Plugin[]>) => {
    return useQuery({ queryKey: ['plugins'], queryFn: fetcher });
};
```

---

## Checklist Summary

```
[ ] Each type/module has a single responsibility — one reason to change
[ ] New behavior added via extension, not modification of existing code
[ ] Interface implementations honor the full contract
[ ] Interfaces are small and focused — no unused methods forced on implementers
[ ] High-level modules depend on abstractions, not concrete implementations
```

---

## Cross-References

- [Function Naming](./10-function-naming.md) — Verb-led, single-purpose naming
- [Code Style §R6](./04-code-style/00-overview.md) — Max 15 lines per function (SRP at function level)
- [Code Style §R17](./04-code-style/00-overview.md) — Max 120 lines per struct/class (SRP at type level)
- [DRY Principles](./08-dry-principles.md) — Avoid duplication (supports OCP)
- [Strict Typing](./13-strict-typing.md) — Type safety supports LSP

---

*SOLID principles reference v1.0.0 — 2026-03-31*
