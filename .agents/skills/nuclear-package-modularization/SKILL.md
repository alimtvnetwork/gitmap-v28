---
name: nuclear-package-modularization
description: Autonomously modularize monolithic Go packages into smaller acyclic subpackages, isolate heavy subprocess tests into gitmap/tests/heavy_test, and estimate test inventory durations.
---

# Nuclear Package Modularization & Heavy Test Isolation

## Purpose & Overview
Autonomously decompose large Go packages into small, acyclic packages (DAG architecture) to eliminate compiler bloat and test latency, while isolating heavy subprocess tests into blackbox test packages and maintaining an accurate test inventory with duration estimation.

## Core Rules & Execution Flow
1. **Strict DAG Architecture**:
   - Zero Go circular dependencies (import cycle not allowed).
   - Root packages (cmd, main) dispatch downwards to domain subpackages.
   - Domain subpackages only import leaf packages (constants, model, store, apperror, fsutil, cliexit).
2. **Heavy Test Isolation**:
   - Tests invoking exec.Command, external git processes, sockets, or time.Sleep belong in gitmap/tests/heavy_test/ (package heavy_test).
   - Routine package tests must be 100% in-memory fast unit tests (<0.01s).
3. **Test Inventory Synchronization**:
   - Synchronize .lovable/test-inventory.json using 03-ai-scripts/33-test-inventory-generator.py.
   - Record modified files under lock via --record.
4. **Coding Guidelines & Quality Gates**:
   - Functions <= 15 lines, blank lines before return statements.
   - Affirmative boolean prefixes (is*, has*), no negative names.
   - Universal AppError wrapping.
   - Verify with go vet ./... and go build ./... (exit 0).
