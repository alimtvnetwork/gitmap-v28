---
name: constants-ownership
description: Constants must be organized by owning package/domain, not forced artificial prefixes
type: constraint
---

# Constants ownership rule

Do not force artificial naming prefixes just to organize constants.

- If a constant belongs to CLI routing and is shared broadly, it may live in `gitmap/constants`.
- If a constant belongs only to one feature package, keep it in that package.
- If a constant belongs only to one cmd flow, keep it in `cmd` as file-local/package-local state.

Avoid reintroducing rules that require every new constant to use canonical prefixes like `Cmd*`, `Msg*`, `Err*`, `Flag*`, or `Default*`.

Why: ownership and package boundaries are the real organization mechanism; forced prefixes create noisy names and push feature-local literals into the wrong package.