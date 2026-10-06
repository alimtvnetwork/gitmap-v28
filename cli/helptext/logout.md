# gitmap logout

Remove the stored GitHub credential.

## Usage

```bash
gitmap logout
```

## How it works

Clears the token stored by `gitmap login` (git global config `github.token`) and rejects any
credential cached by the git credential helper for github.com.

## Examples

```bash
gitmap logout
gitmap login --status   # verify: should report "Not logged in"
```

## See also

- `gitmap login` — log in to GitHub
- `gitmap login --status` — show current login state
