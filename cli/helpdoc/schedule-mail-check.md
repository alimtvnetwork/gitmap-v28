# schedule mail-check

Inspect configured email mailboxes, evaluate unread message metrics, and verify inbox communication health on a schedule.

## Aliases

`schedule-mail`, `sched-check`

---

## Usage

```bash
gitmap schedule mail-check [--all] [--account <alias>] [--unread-only] [--show-full-email] [--json]
```

---

## Flags

| Flag | Description |
|---|---|
| `--all` | Check all configured mailbox profiles (default) |
| `--account <alias>`, `-a <alias>` | Check a specific mailbox account profile |
| `--unread-only`, `-u` | Filter output table to accounts that have unread messages |
| `--show-full-email` | Reveal unmasked email addresses in table output |
| `--json`, `-j` | Format mailbox metrics as Universal JSON Response Envelope |

---

## Examples

```bash
# Check all mailboxes with masked address rendering
gitmap schedule mail-check

# Check mailboxes with unread messages only
gitmap schedule mail-check --unread-only

# Unmask addresses in output
gitmap schedule mail-check --show-full-email
```
