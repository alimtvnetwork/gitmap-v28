Created At: 2026-10-24T20:23:39+08:00
Completed At: 2026-10-24T20:30:00+08:00
# Subtask 320-04: cmdlogin Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdlogin` (6 files)
- Leaf Packages Created:
  - `cli/cmdlogin/enums/`
    - `login_enums.go`
    - `enums_test.go`
  - `cli/cmdlogin/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Implemented
1. `LoginAuthMethodType`: Auth methods (`ssh`, `token`, `browser`, `status`).
2. `LoginInteractiveChoiceType`: Menu options (`1`, `2`, `3`, `4`).
3. `GitAuthProviderType`: Providers (`github`, `gitlab`).
4. `ClipboardToolType`: OS tools (`clip`, `pbcopy`, `wl-copy`, `xclip`, `xsel`).

Each enum follows:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- <= 15 line function limits with strict blank line formatting.
- Leaf package: imports standard library only (`fmt`, `strings`).

## Specific Constants Implemented
- Subcommand names: `CmdLogin`, `CmdLogout`, `SubCmdLoginSSH`, `SubCmdLoginBrowser`, `SubCmdLoginBrowserFlow`, `SubCmdLoginToken`, `SubCmdLoginPAT`, `SubCmdLoginStatus`.
- Flags: `--web`, `--browser`, `--token`, `--no-verify`, `--yes`, `--force`, `--dir`, `--help`, `--status`.
- URLs & endpoints: `GithubUserAPI`, `GithubNewTokenURL`, `GithubNewSSHKeyURL`, `GitlabNewSSHKeyURL`.
- Git config keys: `GitConfigGithubToken`, `GitConfigCoreSSHCommand`, `GitConfigGlobalFlag`, `GitConfigUnsetFlag`, `GitCredentialRejectPayload`.
- Tools & platform constants: `BinGh`, `BinSshKeygen`, `ClipWindows`, `ClipDarwin`, `ClipWlCopy`, `ClipXclip`, `ClipXsel`, `BannerDivider`.

## Consuming Callers Refactored
- `cli/cmdlogin/login_cmd.go`
- `cli/cmdlogin/login_browser.go`
- `cli/cmdlogin/login_ssh.go`
- `cli/cmdlogin/login_status.go`
- `cli/cmdlogin/login_token.go`
- Added comprehensive unit tests in `cli/cmdlogin/login_test.go`.

## Verification Status
- Unit tests: `go test -count=1 -v ./cmdlogin/...` PASS (100% pass across all subpackages).
- Static analysis: `go vet ./cmdlogin/...` clean (zero issues).
- Binary build: `go build .` clean.
- Status: COMPLETED.
