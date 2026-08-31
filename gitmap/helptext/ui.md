# Browser UI

Open the local browser interface for settings management, pipeline telemetry monitoring, command documentation browsing, and interactive web terminal execution.

## Usage

    gitmap ui [flags]

## Options

| Flag     | Type    | Default | Description                                                     |
|----------|---------|---------|-----------------------------------------------------------------|
| --port   | integer | 8080    | Port number for the local HTTP server                           |
| --no-open| boolean | false   | Start the server without automatically opening the default browser |

## Features

- **Settings Management (`/settings`)**: Configure workspace temp directory, cache promotion thresholds, and polling frequencies.
- **Pipeline Monitor (`/pipeline`)**: Live GitHub Actions workflow state, active job status, ETA countdown, and error log viewer.
- **Interactive Web Terminal**: Docked slide-out terminal connected to live SSE streams for interactive command execution.
- **Dark Mode Support**: High-contrast rendering (`text-white`, `text-slate-100`) across all help cards and terminal widgets.
- **Interactive Documentation (`/docs`)**: Full searchable command reference with runnable examples.

## Examples

### Open Default Browser UI

```bash
$ gitmap ui
  Starting local UI server at http://localhost:8080...
  Opened browser at http://localhost:8080/settings
```

### Start Server on Custom Port Without Opening Browser

```bash
$ gitmap ui --port 3000 --no-open
  Starting local UI server at http://localhost:3000...
```
