# gitmap settings

Manage unified speed settings, AI Auto bot configurations, thresholds, and fleet communication.

## Usage

```bash
gitmap settings                                  # View all current speed and AI settings
gitmap settings --json                           # Output settings as machine-readable JSON
gitmap settings set <key> <value>                # Update a configuration key
gitmap settings set ai_auto.default_name "<name>" # Update default AI bot user.name
gitmap settings set ai_auto.default_email "<email>" # Update default AI bot user.email
gitmap settings set ai_auto.default_agent "<agent>" # Update default fallback AI agent
gitmap settings set node_resource_threshold 30   # Set node benchmark resource threshold (30%)
gitmap settings set pae_recency_window_hours 48  # Set PAE commit recency lookback (48h)
```

## Configuration Keys

| Key | Description | Default |
|---|---|---|
| `ai_auto.default_name` | Global AI bot git `user.name` template | `gitmap automation bot - runs autonomously` |
| `ai_auto.default_email` | Global AI bot git `user.email` address | `info@rgitmap.com` |
| `ai_auto.default_agent` | Fallback AI agent identifier | `gemini` |
| `lap.default_hours` | Last Active Projects lookback window in hours | `24` |
| `account_switch.threshold` | Account credit remaining percentage threshold | `15` |
| `node_resource_threshold` | Hardware load limit percentage | `30` |
| `pae_recency_window_hours`| Pull Activity recency window in hours | `48` |
| `special_repos.secrets_name` | Local secrets repository alias | `repo-secrets` |
| `special_repos.cache_name` | Local cache repository alias | `repo-cache` |

## Web UI Settings

Settings are also editable through the interactive Web Settings dashboard at:
`http://127.0.0.1:4011/settings`

## Examples

```bash
# View active settings overview
gitmap settings

# Configure AI bot credentials
gitmap settings set ai_auto.default_name "gitmap automation bot - runs autonomously"
gitmap settings set ai_auto.default_email info@rgitmap.com
gitmap settings set ai_auto.default_agent muse

# Set PAE commit lookback to 72 hours
gitmap settings set pae_recency_window_hours 72
```
