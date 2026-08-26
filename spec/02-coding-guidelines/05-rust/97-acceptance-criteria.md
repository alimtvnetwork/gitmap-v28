# Rust Coding Standards — Acceptance Criteria

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

---

## AC-01: Naming & Style

- [ ] PascalCase serialization uses `#[serde(rename_all = "PascalCase")]` consistently
- [ ] Module and file naming follows Rust snake_case conventions
- [ ] Exported types use descriptive, domain-specific names

## AC-02: Error Handling

- [ ] Domain errors use `thiserror` with project error codes
- [ ] Application logic uses `anyhow` with contextual messages
- [ ] All error variants map to documented error codes

## AC-03: Async Patterns

- [ ] Tokio async runtime is used consistently
- [ ] Channel usage follows bounded MPSC patterns
- [ ] Cancellation-safe patterns are applied in async code

## AC-04: Memory Safety

- [ ] `unsafe` blocks include `// SAFETY:` justification comments
- [ ] Borrowing and `Arc` are preferred over cloning
- [ ] No unnecessary `Clone` derives on large structs

## AC-05: Testing

- [ ] Unit tests follow AAA (Arrange-Act-Assert) pattern
- [ ] OS dependencies use trait-based mocking
- [ ] Integration tests are separated from unit tests

## AC-06: Platform Abstraction

- [ ] Cross-platform code uses `PlatformApi` traits
- [ ] Conditional compilation uses `#[cfg(target_os)]` correctly
- [ ] FFI boundaries include safety documentation

---

## Cross-References

- [Overview](./00-overview.md)
- [Naming Conventions](./01-naming-conventions.md)
- [Error Handling](./02-error-handling.md)
- [Async Patterns](./03-async-patterns.md)
- [Memory Safety](./04-memory-safety.md)
- [Testing Standards](./05-testing-standards.md)
- [FFI & Platform](./06-ffi-platform.md)
