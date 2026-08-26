# Split Database Architecture — Acceptance Criteria

**Version:** 3.2.0  
**Last Updated:** 2026-04-16

---

## AC-01: Database Partitioning

- [ ] SQLite databases split correctly by domain (main, config, analytics)
- [ ] Cross-database queries use proper attach/detach patterns
- [ ] Migration system handles schema changes per database partition

## AC-02: Data Integrity

- [ ] Foreign key relationships maintained within each database
- [ ] Backup and restore operations handle all database partitions
- [ ] Connection pooling manages multiple SQLite file handles efficiently

---

## Cross-References

- [Overview](./00-overview.md)
