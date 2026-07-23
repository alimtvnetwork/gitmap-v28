# gitmap doctor

One-shot health check for every external dependency gitmap relies on.

## Usage

```
gitmap doctor [--json]
```

## What it checks

| Probe   | Pass means                                            | Fix hint on fail                                       |
|---------|-------------------------------------------------------|--------------------------------------------------------|
| git     | `git --version` runs                                  | Install git from https://git-scm.com/downloads         |
| ssh     | `ssh -V` runs                                         | Enable OpenSSH client (Windows optional feature)       |
| chrome  | Chrome binary located on disk                         | Install Chrome from https://www.google.com/chrome/     |
| PATH    | `gitmap` resolves on PATH                             | Run `gitmap self-install`                              |
| sqlite  | pure-Go driver embedded                               | Ensure `.gitmap/` is writable                          |
| disk    | cwd is writable                                       | Free disk space (gitmap needs ~50MB headroom)          |

Each failed probe prints a fix recipe inline. Exit code is non-zero
when any probe fails, so CI pipelines can gate on it.

## Examples

```
gitmap doctor
gitmap doctor --json | jq '.data.checks[] | select(.ok==false)'
```
