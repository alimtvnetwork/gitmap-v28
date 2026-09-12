# Rule: Micro-Tasking & Batching

1. **5-8 File Bounded Batches:** All refactors and multi-file modifications must be bounded into micro-batches of 5 to 8 files.
2. **Verification per Batch:** Run local CI/CD verification or linters after each batch before proceeding to the next.
3. **Isolated Loops:** Complete each subtask boundary, verify, and state status before continuing.
