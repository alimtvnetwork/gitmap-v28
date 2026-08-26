# Golang Standards — Acceptance Criteria

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

---

## AC-01: Code Standards

- [ ] Boolean standards define naming (`isX`, `hasX`, `canX`) and evaluation patterns
- [ ] Error handling uses `apperror.Result[T]` pattern consistently
- [ ] Naming conventions follow Go idioms (exported/unexported, acronym casing)

## AC-02: Architecture

- [ ] Service layer follows interface-based dependency injection
- [ ] HTTP handlers use typed request/response structs
- [ ] Database access uses repository pattern with prepared statements

---

## Cross-References

- [Overview](./00-overview.md)
