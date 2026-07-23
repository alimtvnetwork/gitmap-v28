# Project Type Detection — JSON Output

## Output Files

Each project type produces a dedicated JSON file in `.gitmap/output/`:

| File                    | Contents                     |
|-------------------------|------------------------------|
| `go-projects.json`      | All detected Go projects     |
| `node-projects.json`    | All detected Node.js projects|
| `react-projects.json`   | All detected React projects  |
| `cpp-projects.json`     | All detected C++ projects    |
| `csharp-projects.json`  | All detected C# projects     |

---

## Base JSON Record Schema

All project types share this base structure:

```json
{
  "id": "uuid",
  "repoId": "uuid",
  "repoName": "my-api",
  "projectType": "go",
  "projectName": "github.com/user/my-api",
  "absolutePath": "/home/user/repos/my-api",
  "repoPath": "/home/user/repos/my-api",
  "relativePath": ".",
  "primaryIndicator": "go.mod",
  "detectedAt": "2026-03-11T09:54:00Z"
}
```

Go projects include a `goMetadata` field (see `04-go-metadata.md`).
C# projects include a `csharpMetadata` field (see `05-csharp-metadata.md`).

---

## Write Behavior

- Files are **overwritten** on each scan (not merged).
- Records are **sorted** by `repoName` then `relativePath`.
- Files are **only written** if at least one project of that type
  was detected. Empty files are not created.
- Write failures are logged to stderr but do not abort the scan.
