# Subtask Plan 01: Spec and plan for Spec 277 (Muse Memory IO and GitMap Integration)

- **Version:** 1.0.0
- **Updated:** 2026-10-17
- **Owner:** Worker 01
- **Status:** DONE
- **Subtask Slug:** `01-spec-and-plan`
- **Parent Task:** `muse-memory-io-and-gitmap-integration` (Spec 277)

---

## Goal

Write the Spec 277 spec tree (readme + scope/plan + acceptance criteria) that
fully pins the `gitmap llm train muse [--days N]` contract: digest contract
consumed, clone-or-pull semantics, N-day diff, terminal summary + AI traversal
checklist, missing-digest tolerance, and exact per-criterion evidence —
before any implementation starts.

## Inputs

- Research brief: digest root at repo-secrets `Muse/README.md` + sequenced
  `Muse/01-*.md` + `Muse/chat-seeds/`; digest authored by another worker; the
  command must tolerate its absence.
- Highest existing spec number: 277 = next free (276 is highest existing;
  auto-discovered from `02-spec/21-app/`).
- Existing `llm train` code: `cli/cmd/llm/llm_train.go` (`RunTrain`,
  `parseTrainFlags`), `cli/cmd/llm/llm_types.go` (`TrainOptions`).
- Store API: `cli/store` `OpenSpecialReposSplitDB().GetSpecialRepo("rs")`.

## Steps

1. Create `02-spec/21-app/277-muse-memory-io-and-gitmap-integration/`
   directory.
2. Write `readme.md`: H1 `# Spec 277 — Muse Memory IO and GitMap
   Integration`, Status line, file map.
3. Write `01-scope-and-plan.md`: metadata header (Version, Updated, Spec 277,
   Parent program, Owner) + numbered IN SCOPE / OUT OF SCOPE + command design
   table + file plan table.
4. Write `02-acceptance-criteria-and-evidence.md`: numbered AC-01 … AC-14,
   each naming exact evidence (file paths, function names).
5. Keep style gates: lowercase file names, two-digit prefixes, kebab-case.

## Acceptance criteria

- [x] `readme.md` has the H1, Status line, and file map.
- [x] `01-scope-and-plan.md` has the metadata header and numbered IN SCOPE /
      OUT OF SCOPE lists.
- [x] `02-acceptance-criteria-and-evidence.md` has AC-01 … AC-14 with exact
      evidence per criterion.
- [x] Spec number 277 verified as next-free (276 is highest existing).

## Out of scope

- Digest authoring (`Muse/README.md`, sequenced files, `chat-seeds/`).
- Implementation, docs, changelog (Subtask 02).
- Version bumps, releases.
