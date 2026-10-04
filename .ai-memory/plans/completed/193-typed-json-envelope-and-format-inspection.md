# Plan 193: Typed JSON Envelope Architecture, Format Inspection (`which-format`), and Repo-Secrets Normalization

## User Request (Verbatim)

```text
/goal I think we have different types of JSON, right? Import, export, deploy, many things. So what we should do is JSON will have a type. Okay? Output as a type. And when we output JSON or any data type, we will have two distinct sections on top. First is the attributes, which actually contains the data type, where it is coming from, how it is coming from. The other is the data. The data contains most of the information in the node. So we should have a command in the CLI that would actually tell us, if we do the CLI space, which format, and then pass the JSONs. Or if we don't pass any JSON and folder has multiple JSONs, it will automatically take the JSONs as a format. It will understand and share a command that we can use to run these commands to import. And what will be the end result, what it will change, it will give a summary as well. So let's say a folder contains five JSON. Okay? Three of the JSON really matches the format that we have in the system. So it will tell us, these are the three. If you want to import each one of them, this is the command you should type, and you will get it, and this will change, these things. Now, the last two does not match, so it will say, these are the two things we could not match, so we will not import it. And we will also provide a single line of command that user can use to import all of these. Okay? And if they have to prompt something, they could also say, or the suggestion also could say, you could do "-y" to bypass all the prompting to yes. Okay? So this is one of the techniques that I want in all over the code base. I want to make a huge change. So please make detailed planning so that we don't miss anything. We have a consistency. And according to this, I also want you to go inside the repo secrets and change all the JSONs that we have to this order and this format. Do you understand? Do you have any question and confusion? If you have, let me know. Other than that, please start working on it. And at the end, you verify all these types that we have. Do the end-to-end testing here as much as possible. Okay? Okay. And add some samples. Do not pass any security-based data like password or anything to CI/CD or commit this data. Remember that. So any secret things, it should be inside the repo secrets folder, nothing outside. Remember that. Very important. Can you please apply these changes as I've mentioned? Is it clear?
```

## Spec Reference
- `02-spec/21-app/183-typed-json-envelope-and-format-inspection.md`

## Actionable Deliverables & Task Breakdown

- **Task-01: Typed JSON Envelope Architecture & Core Package (`cli/jsonenvelope`)**
  - Implement universal envelope model: `attributes` (`type`, `source`, `how`, `timestamp`, `version`) and `data` (payload).
  - Provide encoding, decoding, type detection, and backward-compatible unmarshaling.
  - Implement registry of known GitMap data types (`ssh-nodes`, `macro`, `commit-pull-config`, `ui-settings`, `test-inventory`, `pipeline-config`).

- **Task-02: Format Inspection CLI Command (`gitmap which-format`)**
  - Implement `gitmap which-format [files...]` and directory auto-scanning when no files are provided.
  - Register aliases: `gitmap which format`, `gitmap format which`, `gitmap format inspect`.
  - For matched files: report type, attributes, exact suggested import command, and system impact.
  - For unmatched files: clearly report reason for non-match and non-import status.
  - Render actionable single-line batch command to import all matched JSONs with `-y` bypass guidance.

- **Task-03: Subsystem Importer Dual-Format Ingestion**
  - Update SSH cluster node importer (`cli/cmdssh`) to accept both typed envelope and legacy array/object formats.
  - Update Macro importer (`cli/cmdmacro`) to accept both typed envelope and legacy array formats.
  - Update Commit-In config reader (`cli/cmd/commitin`) to accept both typed envelope and legacy config formats.
  - Update UI settings reader (`cli/cmdui`) to accept both typed envelope and legacy formats.

- **Task-04: Repo-Secrets JSON Manifests Normalization**
  - Pull `./repo-secrets`.
  - Convert all JSON manifests in `01-gitmap/` and machine folders to the standard typed envelope.
  - Guard secrets: ensure zero plaintext credentials leak into `gitmap` codebase or test logs.
  - Commit and push `./repo-secrets`.

- **Task-05: Automated Unit Tests, Safe Fixtures & E2E Validation**
  - Create safe, mocked test fixtures in `cli/jsonenvelope/fixtures/`.
  - Add comprehensive unit tests for `which-format` and typed envelope serialization/deserialization.
  - Enforce zero-nesting, positive booleans, and zero lint warnings.

- **Task-06: Cache Cleanup, Minor Release Bump (`v6.404.0`) & Remote CI/CD Verification**
  - Run `python 03-ai-scripts/42-clean-test-and-build-caches.py`.
  - Author verification prompt `01-prompts/24-verify-typed-json-envelope-and-which-format.md`.
  - Consolidate plan and subtasks.
  - Bump minor version to `v6.404.0` and verify green CI/CD on GitHub Actions.
