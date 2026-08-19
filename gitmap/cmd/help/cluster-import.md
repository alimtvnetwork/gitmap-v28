# gitmap cluster import

Import a cluster node registry file (JSON or CSV).

## Usage

```
gitmap cluster import <file>
```

## Description

Merges node records from the file into the local SQLite `ClusterNode` table.

**Merge Semantics:**
- Records with an existing `NodeId` are updated.
- Records with a new `NodeId` are inserted as new nodes.
- Unchanged records are safely skipped.

## Examples

```
gitmap cluster import nodes.json
gitmap cluster import backup.csv
```
