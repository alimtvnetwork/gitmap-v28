# 02 Data Contracts: GitMap Ignore & Cache Engine

This document outlines the API and structural contracts for the new ignore engine, PAS formula task queue, and commit-push review structures.

## 1. PAS Formula Task Queue Contract
The GitMap PAS Formula (Pull-All-SSH Standard) enforces a resilient async orchestration pattern for fleet-wide delegation (e.g., `gitmap pull all ssh` and `gitmap fix-ignores-all-ssh`). 

**Table: `TaskQueue`**
- `QueueId` (String, Primary Key)
- `Section` (String, e.g., 'ssh', 'ignore', 'pull')
- `Action` (String)
- `Target` (String)
- `ForwardPayload` (JSON)
- `InversePayload` (JSON)
- `Status` (Enum: 'pending', 'running', 'completed', 'failed')
- `CreatedAt` (Timestamp)
- `UpdatedAt` (Timestamp)

**Concurrency Rules:**
- Max 2 concurrent workers per SSH node.
- Max 2 async operations running simultaneously.

## 2. Ignore Engine API Contracts
The ignore suite performs deduplication and resolution of `.gitignore` configurations.

**Fix Ignore All / Scan Responses:**
```json
{
  "repoName": "string",
  "hasResumeTask": "boolean",
  "hasDuplicate": "boolean",
  "missingGitmapDir": "boolean",
  "ignoreCount": "number",
  "proposedRemediations": ["string"]
}
```

## 3. Commit-Push Review Structures
When executing `gitmap commit push all`, the system must support a review step (`-r` or `--review`) to display what will be committed.

**Review Summary Payload:**
```json
{
  "totalRepos": "number",
  "reposWithChanges": [
    {
      "repoId": "string",
      "path": "string",
      "status": "string",
      "fileChanges": [
        {
          "file": "string",
          "changeType": "modified | untracked | deleted"
        }
      ]
    }
  ]
}
```
