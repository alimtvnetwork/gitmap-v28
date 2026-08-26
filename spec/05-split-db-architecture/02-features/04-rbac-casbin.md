# Split DB Architecture: Role-Based Access Control (RBAC) with Casbin

**Version:** 3.2.0  
**Created:** 2026-03-09  
**Status:** Active  
**Parent:** [00-overview.md](../00-overview.md)

---

## Overview

This document defines the Role-Based Access Control (RBAC) pattern using [Casbin](https://casbin.org/) for Go-based CLI tools. Casbin provides a flexible, policy-based authorization framework that integrates seamlessly with SQLite via the GORM adapter.

---

## RBAC Scope Levels

The RBAC system can be implemented at three levels depending on your application architecture:

| Level | Scope | Database Location | Use Case |
|-------|-------|-------------------|----------|
| **Root Level** | Global across all apps | `data/rbac.db` | Multi-tenant platform |
| **App Level** | Scoped to single app | `data/{appName}/rbac.db` | Single-application isolation |
| **Company Level** | Scoped to company within app | `data/{appName}/companies/{companySlug}/rbac.db` | Enterprise multi-tenant |

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           CASBIN RBAC ARCHITECTURE                                   │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│   ┌────────────────────┐    ┌────────────────────┐    ┌────────────────────┐        │
│   │    ROOT LEVEL      │    │     APP LEVEL      │    │   COMPANY LEVEL    │        │
│   │   data/rbac.db     │    │ data/{app}/rbac.db │    │ data/{app}/co/     │        │
│   │                    │    │                    │    │ {company}/rbac.db  │        │
│   └─────────┬──────────┘    └─────────┬──────────┘    └─────────┬──────────┘        │
│             │                         │                         │                    │
│             └─────────────────────────┴─────────────────────────┘                    │
│                                       │                                              │
│                                       ▼                                              │
│                         ┌─────────────────────────┐                                  │
│                         │     CasbinRule Table    │                                  │
│                         ├─────────────────────────┤                                  │
│                         │ Id, Ptype, V0-V5        │                                  │
│                         │ ─────────────────────── │                                  │
│                         │ p, alice, data1, read   │  ← Policy                       │
│                         │ g, alice, admin         │  ← Role Assignment              │
│                         │ g2, admin, superadmin   │  ← Role Hierarchy               │
│                         └─────────────────────────┘                                  │
│                                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────────────┐   │
│   │                          REQUEST FLOW                                        │   │
│   │                                                                              │   │
│   │   User Request ──▶ Middleware ──▶ Casbin Enforce() ──▶ Allow/Deny           │   │
│   │                         │                                                    │   │
│   │                         ▼                                                    │   │
│   │              Load policies from SQLite                                       │   │
│   │              Match against model rules                                       │   │
│   │              Return authorization result                                     │   │
│   │                                                                              │   │
│   └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Casbin Model Configuration

### RBAC Model (`rbac_model.conf`)

```ini
[request_definition]
r = Sub, Obj, Act

[policy_definition]
p = Sub, Obj, Act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.Sub, p.Sub) && keyMatch2(r.Obj, p.Obj) && r.Act == p.Act
```

### RBAC with Resource Hierarchy (`rbac_model_hierarchy.conf`)

```ini
[request_definition]
r = Sub, Dom, Obj, Act

[policy_definition]
p = Sub, Dom, Obj, Act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.Sub, p.Sub, r.Dom) && r.Dom == p.Dom && keyMatch2(r.Obj, p.Obj) && r.Act == p.Act
```

---

## SQLite Database Schema

### CasbinRule Table (Auto-created by Adapter)

```sql
-- Casbin automatically creates this table via GORM adapter
CREATE TABLE CasbinRule (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Ptype TEXT NOT NULL,
    V0 TEXT DEFAULT '',
    V1 TEXT DEFAULT '',
    V2 TEXT DEFAULT '',
    V3 TEXT DEFAULT '',
    V4 TEXT DEFAULT '',
    V5 TEXT DEFAULT ''
);

CREATE UNIQUE INDEX IdxCasbinRule ON CasbinRule(Ptype, V0, V1, V2, V3, V4, V5);
```

### Seed Policies (Initial Data)

```sql
-- Default roles (p = policy, g = role assignment)

-- Policies: role, resource, action
INSERT INTO CasbinRule (Ptype, V0, V1, V2) VALUES
-- Admin can do everything
('p', 'admin', '*', '*'),
-- Manager can read and write
('p', 'manager', '/api/*', 'read'),
('p', 'manager', '/api/*', 'write'),
-- Editor can read and write documents
('p', 'editor', '/api/documents/*', 'read'),
('p', 'editor', '/api/documents/*', 'write'),
-- Viewer can only read
('p', 'viewer', '/api/*', 'read');

-- Role hierarchy: g, child_role, parent_role
INSERT INTO CasbinRule (Ptype, V0, V1) VALUES
('g', 'superadmin', 'admin'),
('g', 'admin', 'manager'),
('g', 'manager', 'editor'),
('g', 'editor', 'viewer');
```

---

## Go Implementation

### Installation

```bash
go get github.com/casbin/casbin/v2
go get github.com/casbin/gorm-adapter/v3
go get gorm.io/driver/sqlite
go get gorm.io/gorm
```

### RBAC Manager

```go
package rbac

import (
    "fmt"
    "path/filepath"
    "sync"

    "github.com/casbin/casbin/v2"
    gormadapter "github.com/casbin/gorm-adapter/v3"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// RbacLevel defines the scope of RBAC enforcement
type RbacLevel string

const (
    RbacLevelRoot    RbacLevel = "root"
    RbacLevelApp     RbacLevel = "app"
    RbacLevelCompany RbacLevel = "company"
)

// RbacManager handles role-based access control
type RbacManager struct {
    enforcer *casbin.Enforcer
    adapter  *gormadapter.Adapter
    level    RbacLevel
    dataDir  string
    mu       sync.RWMutex
}

// RbacConfig defines configuration for RBAC initialization
type RbacConfig struct {
    DataDir     string
    Level       RbacLevel
    AppName     string    // Required for App/Company level
    CompanySlug string    // Required for Company level
    ModelPath   string    // Path to RBAC model file (optional, uses default)
    AutoMigrate bool      // Auto-create tables (default: true)
}

// NewRbacManager creates a new RBAC manager
func NewRbacManager(cfg RbacConfig) apperror.Result[*RbacManager] {
    // Determine database path based on level
    dbPath := buildRbacDbPath(cfg)
    
    // Open SQLite with GORM
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        return nil, apperror.Wrap(
            err,
            ErrDbOpen,
            "open rbac database",
        ).WithPath(dbPath)
    }
    
    // Create GORM adapter
    adapter, err := gormadapter.NewAdapterByDb(db)
    if err != nil {
        return nil, apperror.Wrap(
            err,
            ErrRbacAdapterCreate,
            "create GORM adapter for RBAC",
        )
    }
    
    // Determine model path
    modelPath := cfg.ModelPath
    if modelPath == "" {
        modelPath = filepath.Join(cfg.DataDir, "rbac_model.conf")
    }
    
    // Create enforcer
    enforcer, err := casbin.NewEnforcer(modelPath, adapter)
    if err != nil {
        return nil, apperror.Wrap(
            err,
            ErrRbacEnforcerCreate,
            "create casbin enforcer",
        ).WithPath(modelPath)
    }
    
    // Enable auto-save for policy changes
    enforcer.EnableAutoSave(true)
    
    // Load policies from database
    if err := enforcer.LoadPolicy(); err != nil {
        return nil, apperror.Wrap(
            err,
            ErrRbacPolicyLoad,
            "load RBAC policies",
        )
    }
    
    return &RbacManager{
        enforcer: enforcer,
        adapter:  adapter,
        level:    cfg.Level,
        dataDir:  cfg.DataDir,
    }, nil
}

func buildRbacDbPath(cfg RbacConfig) string {
    switch cfg.Level {
    case RbacLevelRoot:
        return filepath.Join(cfg.DataDir, "rbac.db")
    case RbacLevelApp:
        return filepath.Join(cfg.DataDir, cfg.AppName, "rbac.db")
    case RbacLevelCompany:
        return filepath.Join(cfg.DataDir, cfg.AppName, "companies", cfg.CompanySlug, "rbac.db")
    default:
        return filepath.Join(cfg.DataDir, "rbac.db")
    }
}

// Enforce checks if a user has permission to perform an action on a resource
func (m *RbacManager) Enforce(sub, obj, act string) apperror.Result[bool] {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.enforcer.Enforce(sub, obj, act)
}

// EnforceWithDomain checks permission with domain (for company-scoped RBAC)
func (m *RbacManager) EnforceWithDomain(sub, dom, obj, act string) apperror.Result[bool] {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.enforcer.Enforce(sub, dom, obj, act)
}

// AddRoleForUser assigns a role to a user
func (m *RbacManager) AddRoleForUser(user, role string) apperror.Result[bool] {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.enforcer.AddRoleForUser(user, role)
}

// AddRoleForUserInDomain assigns a role to a user within a domain
func (m *RbacManager) AddRoleForUserInDomain(user, role, domain string) apperror.Result[bool] {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.enforcer.AddRoleForUserInDomain(user, role, domain)
}

// RemoveRoleForUser removes a role from a user
func (m *RbacManager) RemoveRoleForUser(user, role string) apperror.Result[bool] {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.enforcer.DeleteRoleForUser(user, role)
}

// GetRolesForUser returns all roles for a user
func (m *RbacManager) GetRolesForUser(user string) apperror.Result[[]string] {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.enforcer.GetRolesForUser(user)
}

// GetUsersForRole returns all users with a specific role
func (m *RbacManager) GetUsersForRole(role string) apperror.Result[[]string] {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.enforcer.GetUsersForRole(role)
}

// AddPolicy adds a permission policy
func (m *RbacManager) AddPolicy(sub, obj, act string) apperror.Result[bool] {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.enforcer.AddPolicy(sub, obj, act)
}

// RemovePolicy removes a permission policy
func (m *RbacManager) RemovePolicy(sub, obj, act string) apperror.Result[bool] {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.enforcer.RemovePolicy(sub, obj, act)
}

// HasRole checks if a user has a specific role
func (m *RbacManager) HasRole(user, role string) apperror.Result[bool] {
    m.mu.RLock()
    defer m.mu.RUnlock()
    roles, err := m.enforcer.GetRolesForUser(user)
    if err != nil {
        return false, err
    }
    for _, r := range roles {
        if r == role {
            return true, nil
        }
    }
    return false, nil
}

// GetAllRoles returns all defined roles
func (m *RbacManager) GetAllRoles() apperror.Result[[]string] {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.enforcer.GetAllRoles()
}

// GetAllPolicies returns all defined policies
func (m *RbacManager) GetAllPolicies() apperror.Result[[][]string] {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.enforcer.GetPolicy()
}

// ReloadPolicy reloads policies from database
func (m *RbacManager) ReloadPolicy() error {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.enforcer.LoadPolicy()
}

// Close closes the RBAC manager
func (m *RbacManager) Close() error {
    // GORM adapter doesn't have explicit close
    return nil
}
```

### HTTP Middleware

```go
package middleware

import (
    "net/http"
    "strings"

    "yourapp/rbac"
)

// RbacMiddleware creates HTTP middleware for authorization
func RbacMiddleware(manager *rbac.RbacManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract user from context (set by auth middleware)
            user := r.Context().Value("userId").(string)
            if user == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            // Get resource and action
            resource := r.URL.Path
            action := methodToAction(r.Method)
            
            // Check permission
            allowed, err := manager.Enforce(user, resource, action)
            if err != nil {
                http.Error(w, "Authorization error", http.StatusInternalServerError)
                return
            }
            
            if !allowed {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func methodToAction(method string) string {
    switch strings.ToUpper(method) {
    case "GET", "HEAD", "OPTIONS":
        return "read"
    case "POST", "PUT", "PATCH":
        return "write"
    case "DELETE":
        return "delete"
    default:
        return "read"
    }
}
```

---

## Usage Examples

### Root Level RBAC (Platform-wide)

```go
// Initialize root-level RBAC
manager, err := rbac.NewRbacManager(rbac.RbacConfig{
    DataDir: "./data",
    Level:   rbac.RbacLevelRoot,
})
if err != nil {
    log.Fatal(err)
}
defer manager.Close()

// Assign admin role to user
manager.AddRoleForUser("user_123", "admin")

// Check permission
allowed, _ := manager.Enforce("user_123", "/api/users", "delete")
fmt.Println("Can delete users:", allowed) // true (admin has * permission)
```

### App Level RBAC (Single Application)

```go
// Initialize app-level RBAC for "gsearch"
manager, err := rbac.NewRbacManager(rbac.RbacConfig{
    DataDir: "./data",
    Level:   rbac.RbacLevelApp,
    AppName: "gsearch",
})
if err != nil {
    log.Fatal(err)
}

// Add custom policy for this app
manager.AddPolicy("researcher", "/api/search/*", "read")
manager.AddRoleForUser("user_456", "researcher")
```

### Company Level RBAC (Multi-tenant)

```go
// Initialize company-level RBAC
manager, err := rbac.NewRbacManager(rbac.RbacConfig{
    DataDir:     "./data",
    Level:       rbac.RbacLevelCompany,
    AppName:     "aibridge",
    CompanySlug: "acme-corp",
})
if err != nil {
    log.Fatal(err)
}

// Assign role within company domain
manager.AddRoleForUserInDomain("user_789", "manager", "acme-corp")

// Check permission with domain
allowed, _ := manager.EnforceWithDomain("user_789", "acme-corp", "/api/projects", "write")
```

---

## Default Role Definitions

| Role | Priority | Permissions | Description |
|------|----------|-------------|-------------|
| `superadmin` | 100 | `*` on `*` | Full system access |
| `admin` | 90 | `*` on `/api/*` | Full API access |
| `manager` | 70 | `read`, `write` on `/api/*` | Read/write API access |
| `editor` | 50 | `read`, `write` on `/api/documents/*` | Document management |
| `viewer` | 10 | `read` on `/api/*` | Read-only access |

---

## Best Practices

### 1. Use Wildcards Sparingly

```go
// ❌ Too permissive
manager.AddPolicy("user", "*", "*")

// ✅ Specific permissions
manager.AddPolicy("user", "/api/documents/*", "read")
```

### 2. Leverage Role Hierarchy

```go
// Roles inherit permissions from parent roles
// superadmin → admin → manager → editor → viewer

// Adding viewer role gives read access
// Adding manager role gives read + write access
```

### 3. Reload Policies After Changes

```go
// After bulk policy updates
manager.ReloadPolicy()
```

### 4. Use Domain-based RBAC for Multi-tenancy

```go
// Separate permissions per company
manager.EnforceWithDomain(user, company, resource, action)
```

---

## Integration with Split DB

The RBAC database follows the Split DB naming conventions:

| Level | Path Pattern |
|-------|--------------|
| Root | `data/rbac.db` |
| App | `data/{appName}/rbac.db` |
| Company | `data/{appName}/companies/{companySlug}/rbac.db` |

### Database Registry Entry

```sql
INSERT INTO DbRegistry (Id, Category, EntityId, SequenceNum, Path) VALUES
('rbac_root', 'rbac', 'root', 1, 'rbac.db'),
('rbac_gsearch', 'rbac', 'gsearch', 1, 'gsearch/rbac.db');
```

---

## References

- [Casbin Documentation](https://casbin.org/docs/overview)
- [GORM Adapter](https://github.com/casbin/gorm-adapter)
- [Casbin Go SDK](https://github.com/casbin/casbin)
- [Model Examples](https://casbin.org/docs/model-storage)
