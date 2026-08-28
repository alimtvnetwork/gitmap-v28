# fix-seq-files

Re-sequence numbered files (e.g. 01-notes.md, 02-draft.md) systematically.

## ALIAS

fsf

## Examples

Order files by their filesystem modification time:
```bash
gitmap fix-seq-files /folder1 /folder2 -orderbytime
```

Order files alphabetically by their text suffix, keeping existing numbers where possible:
```bash
gitmap fix-seq-files /folder1 -orderbyaz -keep-old-order
```

Pin a specific filename base to a specific sequence number, letting others increment around it:
```bash
gitmap fix-seq-files /folder1 -pin "readme=00,intro=01"
```

## OVERVIEW

This command safely adjusts numeric prefixes in filenames. It defaults to sorting unpinned files by your chosen strategy and applying zero-padded sequence numbers (e.g., `01`, `02`).
