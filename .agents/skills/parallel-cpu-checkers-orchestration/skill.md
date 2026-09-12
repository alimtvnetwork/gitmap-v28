---
name: parallel-cpu-checkers-orchestration
description: Autonomous orchestration for parallel multi-core quality checkers, upfront file discovery, and real-time progress percentages in CI/CD local runner.
---

# Parallel CPU Checkers & Real-Time Progress Engine Orchestration

## Overview
Autonomously transforms sequential single-threaded quality checkers into high-throughput multi-core parallel engines utilizing 100% available CPU capacity with real-time percentage progress.

## Core Rules
1. **Upfront File Listing**: Pre-gather all target files using `03-ai-scripts/02-shared-engine.py` (`stream_directory_files` / `process_repository_files`) or `03-file-manipulator.py`.
2. **Full Multi-Core CPU Saturation**: Partition files across `os.cpu_count()` workers using `concurrent.futures.ThreadPoolExecutor` / `ProcessPoolExecutor`.
3. **Live Progress Percentages**: Emit dynamic progress percentage callbacks (e.g., `[ 450/1500 ] 30% | 16 workers | 420 files/sec`) so execution never stalls at 0%.
4. **Zero-Stop Autonomous Loop**: Decompose tasks into `.lovable/plans/` and transition from planning to execution without waiting for user intervention.
