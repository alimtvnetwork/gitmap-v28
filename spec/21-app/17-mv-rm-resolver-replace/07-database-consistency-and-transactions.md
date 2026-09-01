# Specification 17 — Chapter 7: Database Consistency & Transactions

## 1. Atomicity Guarantee

Every multi-table modification in `gitmap mv` and `gitmap rm` executes inside a SQLite transaction:
- `BEGIN IMMEDIATE`
- Update / Delete `Repo`
- Update / Delete `Alias`
- Update / Delete `ScanFolder` association
- Record command history in `CommandHistory`
- `COMMIT`

## 2. Rollback Policy

If physical filesystem operations fail during `gitmap mv`, database changes are rolled back or restored to ensure no orphan database records point to non-existent paths.
