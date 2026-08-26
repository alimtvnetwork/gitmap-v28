# File & Folder Naming — Go

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Go has strict conventions enforced by `gofmt`, `golint`, and community standards. Package names ARE directory names.

---

## File Naming Rules

### 1. Files — `snake_case.go`

```
✅ http_handler.go
✅ user_service.go
✅ response_writer.go
❌ httpHandler.go
❌ http-handler.go
❌ HttpHandler.go
```

### 2. Test Files — `*_test.go`

```
✅ http_handler_test.go
✅ user_service_test.go
❌ http_handler.test.go
❌ httpHandlerTest.go
```

### 3. Platform-Specific Files

Go build tags use OS/architecture suffixes:

```
✅ file_linux.go
✅ file_windows.go
✅ file_darwin_amd64.go
```

### 4. Main Entry Point

```
✅ main.go       (in cmd/{app}/)
✅ doc.go        (package documentation)
```

---

## Folder / Package Naming Rules

### 1. Packages — `lowercase` (no hyphens, no underscores)

```
✅ internal/handlers/
✅ internal/enums/providertype/
✅ pkg/httputil/
❌ internal/Handlers/
❌ internal/http-util/
❌ internal/http_util/
```

### 2. Standard Go Project Layout

```
my-cli/                          ← kebab-case repo name (OK)
├── cmd/
│   └── mycli/                   ← lowercase, no hyphens
│       └── main.go
├── internal/                    ← private packages
│   ├── handlers/
│   │   ├── search_handler.go
│   │   └── search_handler_test.go
│   ├── services/
│   │   └── search_service.go
│   └── enums/
│       ├── providertype/        ← enum package (type suffix)
│       │   └── variant.go
│       └── enginetype/
│           └── variant.go
├── pkg/                         ← public packages
│   └── httputil/
│       └── client.go
├── go.mod
└── go.sum
```

### 3. Enum Package Convention

Enum packages MUST end with `type` suffix:

```
✅ providertype/
✅ enginetype/
✅ searchmodetype/
❌ provider/
❌ engine_type/
❌ search-mode-type/
```

---

## Forbidden Patterns

| Pattern | Why |
|---------|-----|
| `camelCase` folders | Go packages must be lowercase single words |
| Hyphens in packages | `go build` treats `-` as invalid in import paths |
| Underscores in packages | Discouraged by Go conventions (except `_test`) |
| `PascalCase` files | Go files are always `snake_case` |
| Mixed case folders | Causes import path issues across OS |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Golang Standards | [../03-golang/00-overview.md](../03-golang/00-overview.md) |
| Enum Specification | [../03-golang/01-enum-specification/00-overview.md](../03-golang/01-enum-specification/00-overview.md) |
| Cross-Language Rules | [./01-cross-language.md](./01-cross-language.md) |
