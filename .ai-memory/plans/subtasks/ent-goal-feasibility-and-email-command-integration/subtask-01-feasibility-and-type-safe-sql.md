# Subtask 01: Ent Go Feasibility Analysis & Type-Safe SQL Audit

**Parent Plan:** `.ai-memory/plans/ent-goal-feasibility-and-email-command-integration.md`  
**State:** `COMPLETED`  
**Owned Files:** `02-spec/21-app/ent-goal-feasibility-and-email-command-integration/ent-feasibility.md`  

## Objectives
1. Provide definitive, benchmark-backed answers to the user's two primary questions regarding Ent Go:
   - Question 1: Will using Ent Go slow down the system? (Answered with latency, allocation, startup, and build metrics).
   - Question 2: Can a migration be done in one shot without any hassle? (Answered with structural blockers: PK naming conflict, Split-DB impedance, lack of cross-file graph edges).
2. Audit the codebase for dynamic SQL concatenations and mandate 100% compile-time parameterized query usage.

## Evidence & Verification
- Feasibility document completed in `02-spec/21-app/ent-goal-feasibility-and-email-command-integration/ent-feasibility.md`.
- Concludes that preserving `cli/dbengine/` with parameterized queries preserves performance, eliminates dynamic SQL risks, and avoids massive binary bloat and transaction deadlocks.
