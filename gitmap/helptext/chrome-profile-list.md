# gitmap chrome-profile-list

List Chrome profiles discovered on disk plus every profile gitmap
has tracked in its local SQLite store (with export counts and the
last-seen timestamp).

## Alias

`cpl`

## Usage

    gitmap chrome-profile-list
    gitmap cpl

No arguments. Reads Chrome's User Data root for the current OS and
joins the result with the `ChromeProfile` table.

## Examples

    $ gitmap cpl
    Chrome profiles (C:\Users\me\AppData\Local\Google\Chrome\User Data):
      Default
      Profile 1
      Profile 2

    Tracked in gitmap DB:
      - Default                       exports=3   last=19-Jun-2026 09:14 PM (UTC)
      - Profile 1                     exports=1   last=14-Jun-2026 06:02 PM (UTC)

## See also

- [chrome-profile-copy](chrome-profile-copy.md)
- [chrome-profile-export](chrome-profile-export.md)
- [chrome-profile-delete](chrome-profile-delete.md)

## Examples

```bash
gitmap chrome-profile-list
```
