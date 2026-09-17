---
name: parallel-pipeline-download-and-two-pass-log-processor
description: Autonomously implement fallback to previous pipeline runs, concurrent section/job log downloads for single commits, and two-pass non-mutating parallel line execution and filtering algorithms across GitMap.
---

# Parallel Pipeline Download and Two-Pass Log Processor

Autonomously implement and verify:
1. Fallback to previous pipeline run / pipeline DB records when no new pipeline or failing run exists on the current commit.
2. Concurrent parallel downloading of pipeline sections, workflows, and job logs for single commits using bounded worker pools.
3. Two-pass non-mutating parallel log line processing algorithm:
   - Pass 1: Concurrency-safe parallel line marking (keep/remove mask) without slice mutation.
   - Pass 2: Exact pre-allocation and fast parallel/block materialization into clean arrays.
4. Concurrency util functions for safe parallel slice and log operations with thread-safe tracking.
