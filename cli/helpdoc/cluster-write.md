# gitmap cluster write

Write content or upload a local file to a remote cluster node.

## Usage

```bash
gitmap cluster write --id <node-id> --dest <remote-path> <local-file>
```

## Examples

```bash
gitmap cluster write --id 2 --dest "C:\app\config.json" ./local-config.json
```

See also: `gitmap cluster cat`, `gitmap cluster`
