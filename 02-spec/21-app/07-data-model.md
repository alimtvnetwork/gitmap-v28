# Data Model

## ScanRecord

| Field            | Type   | Required | Default | Notes                          |
|------------------|--------|----------|---------|--------------------------------|
| id               | string | yes      | UUID    | Unique record identifier       |
| slug             | string | yes      | —       | Derived from HTTPS URL         |
| repoName         | string | yes      | —       | Derived from remote URL        |
| httpsUrl         | string | yes      | —       | HTTPS clone URL                |
| sshUrl           | string | yes      | —       | SSH clone URL                  |
| branch           | string | yes      | "main"  | Current checked-out branch     |
| relativePath     | string | yes      | —       | Path relative to scan root     |
| absolutePath     | string | yes      | —       | Full filesystem path           |
| cloneInstruction | string | yes      | —       | Full `git clone` command       |
| notes            | string | no       | ""      | User or system notes           |

**Slug generation:** Populated by `mapper.BuildRecords` during scan.
Algorithm: last path segment of HTTPS URL, strip `.git`, lowercase.
Falls back to `repoName` when HTTPS URL is empty.

## Config

See [06-config.md](./06-config.md).

## CloneResult

| Field   | Type       | Description                        |
|---------|------------|------------------------------------|
| Record  | ScanRecord | The repo record                    |
| Success | bool       | Whether the clone succeeded        |
| Error   | string     | Error message (empty on success)   |

## CloneSummary

| Field     | Type          | Description                          |
|-----------|---------------|--------------------------------------|
| Succeeded | int           | Number of successful clones          |
| Failed    | int           | Number of failed clones              |
| Cloned    | []CloneResult | Successfully cloned repos            |
| Errors    | []CloneResult | Failed clone operations with reasons |

## Group

| Field       | Type   | Required | Default | Notes                          |
|-------------|--------|----------|---------|--------------------------------|
| id          | string | yes      | UUID    | Unique group identifier        |
| name        | string | yes      | —       | Display name (unique)          |
| description | string | no       | ""      | Optional description           |
| color       | string | no       | ""      | Terminal color (e.g. "green")  |
| createdAt   | string | yes      | now     | Creation timestamp             |

## GroupRepo

| Field   | Type   | Description                        |
|---------|--------|------------------------------------|
| groupId | string | FK → Group.id                      |
| repoId  | string | FK → ScanRecord.id                 |

See [16-database.md](./16-database.md) for schema details and
[17-repo-grouping.md](./17-repo-grouping.md) for CLI commands.

## AmendmentRecord

| Field          | Type          | Required | Default | Notes                              |
|----------------|---------------|----------|---------|------------------------------------|
| id             | string        | yes      | UUID    | Unique amendment identifier        |
| timestamp      | string        | yes      | now     | ISO 8601 UTC timestamp             |
| branch         | string        | yes      | —       | Target branch name                 |
| fromCommit     | string        | yes      | —       | First commit SHA in range          |
| toCommit       | string        | yes      | —       | Last commit SHA (HEAD at amend)    |
| totalCommits   | int           | yes      | —       | Number of commits rewritten        |
| previousAuthor | AuthorInfo    | yes      | —       | Original author name/email         |
| newAuthor      | AuthorInfo    | yes      | —       | Replacement author name/email      |
| mode           | string        | yes      | —       | `all`, `range`, or `head`          |
| forcePushed    | bool          | yes      | false   | Whether force-push was executed    |
| commits        | []CommitEntry | yes      | —       | List of amended commits            |

## AuthorInfo

| Field | Type   | Description        |
|-------|--------|--------------------|
| name  | string | Author name        |
| email | string | Author email       |

## CommitEntry

| Field   | Type   | Description              |
|---------|--------|--------------------------|
| sha     | string | Full commit SHA          |
| message | string | Commit message (subject) |

See [24-amend-author.md](./24-amend-author.md) for the amend command spec.
