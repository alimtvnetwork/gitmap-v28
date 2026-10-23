# Subtask 02: Component Contracts, Schemas & Examples (Spec 313)

- **Parent Plan:** [.ai-memory/plans/pending/313-gitmap-execution-and-optimization-instructions.md](../../pending/313-gitmap-execution-and-optimization-instructions.md)
- **Specification:** [02-spec/21-app/313-gitmap-execution-and-optimization-instructions/02-component-spec.md](../../../02-spec/21-app/313-gitmap-execution-and-optimization-instructions/02-component-spec.md)
- **Status:** `ACTIVE`

## Scope & Deliverables
1. Go data contracts and schemas:
   - `PullTaskRecord` and `PullItemHistoryRecord` in `cli/store/split_db_pull_types.go`.
   - `HeatmapActionPayload` in `cli/heatmap/heatmap_json_payload.go`.
   - `SshKey` encrypted private key fields in `cli/store/sshkey.go`.
   - `ParseIndexSelection` in `cli/cmdexclude/exclude_selector.go`.
2. CLI Command Handlers:
   - `gitmap scan disable-repos` and `gitmap scan enable-repos` with sequence selection (`1-3, 5`).
   - `gitmap heatmap log --json '<payload>'` supporting direct JSON unmarshaling and model/tool metadata (`Gemini`, `Antigravity IDE`, `Muse`).
3. REST API Routes:
   - `GET /api/repos/status`
   - `POST /api/repos/disable`
   - `POST /api/repos/enable`
4. React UI Controls:
   - `src/pages/RepoSettings.tsx` repository scan inclusion toggles and table view.
