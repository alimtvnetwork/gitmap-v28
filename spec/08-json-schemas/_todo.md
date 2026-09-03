# JSON Schema migration TODO

The following CLI commands currently emit JSON via `json.MarshalIndent(struct, ...)`,
which means **field order is reflection-defined and NOT contractual**. Each one
needs:

1. Migration of the encoder to `gitmap-v28/stablejson` (so ordering becomes a
   compile-time decision instead of a reflection accident).
2. A hand-written `<command>.schema.json` next to this file.
3. A `<command>_jsonschema_contract_test.go` in `gitmap-v28/cmd/` pinning the schema
   against the actual encoder output.

Until a command appears in the table in `README.md`, downstream consumers should
treat its JSON output as **shape-stable but key-order-unstable**.

## Pending commands

(Discovered via `rg -n "json.NewEncoder|json.Marshal" gitmap-v28/cmd/`. Order is
roughly by perceived consumer impact — high-traffic / scripting-friendly first.)

| Priority | Command (file) | Notes |
|---|---|---|
| ✅ done | ~~`gitmap-v28 list-releases --json` (`listreleases.go`, `listreleasesallrepos.go`)~~ | Migrated to `gitmap-v28/stablejson` via `listreleasesrender.go`. Schemas: [`list-releases.schema.json`](list-releases.schema.json) (per-repo, lowerCamel) + [`list-releases-all-repos.schema.json`](list-releases-all-repos.schema.json) (joined --all-repos, PascalCase preserved from legacy `MarshalIndent`). Pinned by `gitmap-v28/cmd/listreleases_jsonschema_contract_test.go` (9 tests incl. byte-compat with legacy output). |
| ✅ done | ~~`gitmap-v28 history --json` (`history.go`)~~ | Migrated to stablejson via `historyrender.go` (v5.64.0). Schema: [`history.schema.json`](history.schema.json). Pinned by `gitmap/cmd/history_jsonschema_contract_test.go` + `historyjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 watch --json` (`watch.go`)~~ | Migrated to stablejson via `watchrender.go` (v5.65.0). Nested repos + summary pre-rendered in compact mode. Schema: [`watch.schema.json`](watch.schema.json). Pinned by `gitmap/cmd/watch_jsonschema_contract_test.go` + `watchjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 probe-report` (`probereport.go`)~~ | Migrated to stablejson via `proberender.go` (v5.66.0). Schema: [`probe-report.schema.json`](probe-report.schema.json). Pinned by `gitmap/cmd/proberepor_jsonschema_contract_test.go` + `probereporjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 amend list --json` (`amendlist.go`)~~ | Migrated to stablejson via `amendlistrender.go` (v5.67.0). Schema: [`amend-list.schema.json`](amend-list.schema.json). Pinned by `gitmap/cmd/amendlist_jsonschema_contract_test.go` + `amendlistjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 amend audit` (`amendaudit.go`)~~ | Migrated to stablejson via `amendauditrender.go` (v5.68.0). Schema: [`amend-audit.schema.json`](amend-audit.schema.json). Pinned by `gitmap/cmd/amendaudit_jsonschema_contract_test.go` + `amendauditjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 diff-profiles --json` (`diffprofiles.go`)~~ | Migrated to stablejson via `diffprofilesrender.go` (v5.69.0). Nested arrays pre-rendered in compact mode. Schema: [`diff-profiles.schema.json`](diff-profiles.schema.json). Pinned by `gitmap/cmd/diffprofiles_jsonschema_contract_test.go` + `diffprofilesjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 bookmark list --json` (`bookmarklist.go`)~~ | Migrated to stablejson via `bookmarklistrender.go` (v5.70.0). Schema: [`bookmark-list.schema.json`](bookmark-list.schema.json). Pinned by `gitmap/cmd/bookmarklist_jsonschema_contract_test.go` + `bookmarklistjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 project repos --json` (`projectreposoutput.go`)~~ | Migrated to stablejson via `projectreposrender.go` (v5.71.0). Schema: [`project-repos.schema.json`](project-repos.schema.json). Pinned by `gitmap/cmd/projectrepos_jsonschema_contract_test.go` + `projectreposjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 ssh list --json` (`sshlist.go`)~~ | Migrated to stablejson via `sshlistrender.go` (v5.74.0). Schema: [`ssh-list.schema.json`](ssh-list.schema.json). Pinned by `gitmap/cmd/sshlist_jsonschema_contract_test.go` + `sshlistjson_contract_test.go`. |
| med | `gitmap-v28 env-registry` (`envregistry.go`) | No `--json` stdout flag; only reads/writes `env-registry.json` file. Schema migration does not apply. |
| ✅ done | ~~`gitmap-v28 export` (`export.go`)~~ | Migrated to stablejson via `exportrender.go` (v5.81.0). Top-level key order pinned: `version`, `exportedAt`, `repos`, `groups`, `releases`, `history`, `bookmarks`. Per-record property sets pinned in `export.schema.json` v2 (v5.82.0+) for all 5 nested arrays: `ScanRecord`, `GroupExport`, `ReleaseRecord`, `CommandHistoryRecord`, `BookmarkRecord`. Empty arrays emit `[]` (never `null`). Pinned by `gitmap/cmd/export_jsonschema_contract_test.go` + `export_nested_jsonschema_contract_test.go` + `exportjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 find-next --json` (`findnext.go`)~~ | Migrated to `gitmap-v28/stablejson` via `findnextrender.go` (v5.63.0). Schema: [`find-next.schema.json`](find-next.schema.json) — nested `repo` allows passthrough. Pinned by `gitmap-v28/cmd/findnext_jsonschema_contract_test.go` (top-level shape + encoder-keys ⊂ schema.properties) on top of the pre-existing tag-order golden in `findnextjson_contract_test.go`. |
| med | `gitmap-v28 rescan` (`rescan.go`) | No `--json` stdout flag; replays scan from `last-scan.json` cache. Schema migration does not apply. |
| ✅ done | ~~`gitmap-v28 latest-branch --json` (`latestbranchoutput.go`)~~ | Migrated to stablejson via `latestbranchrender.go` (v5.72.0). Nested top-N array pre-rendered in compact mode. Schema: [`latest-branch.schema.json`](latest-branch.schema.json). Pinned by `gitmap/cmd/latestbranch_jsonschema_contract_test.go` + `latestbranchjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 llm-docs` (`llmdocs.go`)~~ | Migrated to stablejson via `llmdocsrender.go` (v5.80.0). Nested command groups + per-command rows pre-rendered in compact mode. Optional top-level sections and per-command `example` conditionally omitted. Schema: [`llm-docs.schema.json`](llm-docs.schema.json). Pinned by `gitmap/cmd/llmdocs_jsonschema_contract_test.go` + `llmdocsjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 list-versions --json` (`listversionsutil.go`)~~ | Migrated to stablejson via `listversionsrender.go` (v5.73.0). Optional `source`/`changelog` are conditionally appended to preserve legacy omitempty shape. Schema: [`list-versions.schema.json`](list-versions.schema.json). Pinned by `gitmap/cmd/listversions_jsonschema_contract_test.go` + `listversionsjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 version-history --json` (`versionhistory.go`)~~ | Migrated to stablejson via `versionhistoryrender.go` (v5.76.0). Schema: [`version-history.schema.json`](version-history.schema.json). Pinned by `gitmap/cmd/versionhistory_jsonschema_contract_test.go` + `versionhistoryjson_contract_test.go`. |
| med | `gitmap-v28 seo write` (`seowritecreate.go`) | Sample/template output; no `--json` stdout flag. Schema migration does not apply. |
| ✅ done | ~~`gitmap-v28 scan-project` (`scanprojectoutput.go`)~~ | Per-type file output shape pinned v5.84.0. Schema: [`scan-project.schema.json`](scan-project.schema.json) covers all 5 sibling files; top-level record keys are PascalCase (`Project`, `GoMeta`, `Csharp`) — contractual since v1. Pinned by `gitmap/cmd/scanproject_jsonschema_contract_test.go` (file-map registry + record-keys ⊂ schema.properties). Not piped to stdout, but consumers (downstream tooling, CI) get a stable on-the-wire shape. |
| ✅ done | ~~`gitmap-v28 stats --json` (`stats.go`)~~ | Migrated to stablejson via `statsrender.go` (v5.75.0). Top-level object + nested compact `commands` array (`json.RawMessage`). Schema: [`stats.schema.json`](stats.schema.json). Pinned by `gitmap/cmd/stats_jsonschema_contract_test.go` + `statsjson_contract_test.go`. |
| ✅ done | ~~`gitmap-v28 temp-releaselist --json` (`tempreleaselist.go`)~~ | Migrated to stablejson via `tempreleaselistrender.go` (v5.77.0). Schema: [`temp-release-list.schema.json`](temp-release-list.schema.json). Pinned by `gitmap/cmd/tempreleaselist_jsonschema_contract_test.go` + `tempreleaselistjson_contract_test.go`. |

## Estimated effort

~30-60 min per command (encoder migration + schema + test). Total ~10-20 hours
of focused work. Recommend tackling in a single sprint to keep the
consumer-contract surface consistent rather than dribbling out one at a time.
