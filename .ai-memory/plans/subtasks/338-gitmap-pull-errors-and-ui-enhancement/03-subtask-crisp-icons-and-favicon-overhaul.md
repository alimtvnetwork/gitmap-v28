# Subtask 03: Crisp Icons and Favicon Overhaul
Parent Plan: 338-gitmap-pull-errors-and-ui-enhancement

## Scope
- Overhaul `GitMapLogoSVG` in `cli/cmdlogviewer/favicon.go`:
  - Enforce pure transparent background (`alpha = 0`) without chunky outer border cards.
  - Implement 8% quiet zone and crisp modern vector branching nodes (`#89b4fa`, `#a6e3a1`, `#cdd6f4`).
- Update embedded favicons and UI headers.
