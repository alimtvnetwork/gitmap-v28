# `gitmap dedupe`

Detect identical repos cloned under different folders by comparing the
SHA of each repo's `HEAD^{tree}`. Groups of 2+ repos with the same tree
hash are reported.

## Flags

```
--root=DIR      scan root directory (default ".")
```

## Examples

```
gitmap dedupe
gitmap dedupe --root=~/code
```
