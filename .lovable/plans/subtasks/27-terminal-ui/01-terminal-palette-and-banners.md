# Subtask 27-01: ANSI Palette & Super-Category Banners

## Goal

Standardize terminal color definitions to bright bold 9X ANSI sequences and establish structured box-drawing banners:
- `ColorGreen = "\033[1;92m"`
- `ColorRed = "\033[1;91m"`
- `ColorYellow = "\033[1;93m"`
- `ColorCyan = "\033[1;96m"`
- `ColorMagenta = "\033[1;95m"`
- `printSuperCategory` with `━━ SECTION ━━━━━━━━━━`

## Status: DONE

- Implemented in `gitmap/constants/constants_colors.go` and `gitmap/cmd/rootusage.go`.
