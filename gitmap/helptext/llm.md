# llm

Outputs the LLM specification, instructions, and capability guide for autonomous AI agents and coding assistants.

## Aliases

None (`gitmap llm` or `gitmap llm-docs`)

## Usage

    gitmap llm [--url|--instruction]

## Flags

| Flag            | Type    | Default | Description                                              |
|-----------------|---------|---------|----------------------------------------------------------|
| --url           | boolean | false   | Output the public URL to the LLM markdown specification  |
| --instruction   | boolean | false   | Output the full markdown instructions directly to stdout |

## Autonomous Agent Integration

AI coding agents should invoke `gitmap llm` to inspect core capabilities:
- **File Finding**: `gitmap find-files`, `gitmap find-files-any`, `gitmap find`
- **CI/CD Diagnostics**: `gitmap pipeline status --json`, `gitmap eta`, `gitmap error-logs --json --tempfile "err.json"`
- **Semantic Commits**: `gitmap commit-push-feature`, `gitmap commit-push-bug`, `gitmap commit-push-release`

## Examples

```bash
gitmap llm
gitmap llm --url
gitmap llm --instruction
```

## See Also

- [pipeline](pipeline.md) — Live CI/CD telemetry, ETA wait times, and failure logs
- [find-files](find-files.md) — Exact and substring filename matching with extension filters

