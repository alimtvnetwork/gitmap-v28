# gitmap cp

Copy files across cluster fleet nodes or copy text to OS clipboard and memory buffer.

## Alias of

`copy`

## Usage

```bash
gitmap cp <local-src> <node>:<dest>             # Local to single remote node
gitmap cp <local-src> <dest> -t <node>          # Local to node via flag
gitmap cp <local-src> all:<dest>                # Broadcast to all online nodes
gitmap cp <node>:<remote-src> <local-dest>      # Remote pull to local
gitmap cp <node1>:<src> <node2>:<dest>          # Node to node transfer
gitmap cp <node1>:<src> all:<dest>              # Node to all online nodes
gitmap cp -r <local-dir> <node>:<dest-dir>      # Recursive directory copy
```

See [copy.md](copy.md) for full documentation and options.
