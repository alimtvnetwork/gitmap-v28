# copy

Copy files across cluster fleet nodes or copy text/items to the OS system clipboard and persistent repository memory buffer.

## Aliases

`cp`, `cp-mem`, `copy-to-memory`

## Usage

```bash
# Cluster Node File Copy (8 Syntax Variants)
gitmap cp <local-src> <node>:<dest>             # Local to single remote node
gitmap cp <local-src> <dest> -t <node>          # Local to node via flag
gitmap cp <local-src> all:<dest>                # Broadcast to all online nodes
gitmap cp <node>:<remote-src> <local-dest>      # Remote pull to local
gitmap cp <node1>:<src> <node2>:<dest>          # Node to node transfer
gitmap cp <node1>:<src> all:<dest>              # Node to all online nodes
gitmap cp -r <local-dir> <node>:<dest-dir>      # Recursive directory copy
gitmap copy <text...>                           # Clipboard & memory copy
gitmap copy --file <path>                       # File content to clipboard
```

## Description

The `copy` / `cp` command provides unified copying:
1. **Cluster Fleet File Transport:** Transfers files and directories between the local workstation and remote SSH fleet nodes, across remote nodes, or broadcasted to all online nodes using streaming tar sessions.
2. **Clipboard & Memory Buffer:** When called with plain text or without a destination node, writes content directly to the operating system clipboard and persists it in `.gitmap/memory/clipboard.txt`.

## Examples

```bash
# Copy file to remote worker node
gitmap cp ./build.tar.gz worker-1:/tmp/build.tar.gz

# Broadcast config to all online fleet nodes
gitmap cp ./config.json all:/etc/config.json

# Pull remote logs to local workspace
gitmap cp worker-1:/var/log/app.log ./logs/app.log

# Copy entire directory tree recursively
gitmap cp -r ./dist worker-1:/var/www/dist

# Copy text to clipboard and memory
gitmap copy "npm run build && npm run test"
```
