# AI Code Review Guide — Naming, Signatures, Whitespace

> **Audience:** an AI assistant performing a code review.
> **Scope:** Go, TypeScript/JavaScript, Python, Rust, Java, C#.
> **Mode of operation:** do not merely comment. **Rewrite the code** and output the corrected
> version (or a unified diff), a rule-numbered findings list, and a self-verification checklist.

---

## 0. How the AI must apply this guide

For every function, method, constructor, struct, and call site you touch:

1. Run the **Quick checklist** (§1) top to bottom.
2. For each violation, produce the rule number, a one-line reason, and the corrected code.
3. Never change behaviour while applying a formatting or naming rule. Naming, splitting,
   struct-grouping, and whitespace are mechanical refactors — keep logic identical unless a rule
   explicitly says otherwise (R6 unused parameters, R7 swallowed errors).
4. If a rule conflicts with a formatter the repo already runs (`gofmt`, `prettier`, `black`,
   `rustfmt`, `spotless`, `dotnet format`), the formatter wins **only where it actually conflicts**,
   and only on whitespace it controls. Naming, parameter counts, struct grouping, error handling,
   and the blank-line rules R13–R20 (which no formatter enforces) always apply.
5. Apply rules to **definitions and call sites together**. A split signature with a 200-char call
   site is still a violation.
6. Before printing the rewritten code, run §18 self-verification and emit a pass/fail line per rule.

---

## 1. Quick checklist (run on every function)

| #   | Check                                                                                             | Must fix?  |
| --- | ------------------------------------------------------------------------------------------------- | ---------- |
| R1  | Acronyms (`Id`, `Json`, `Url`, `Ip`, `Http`, `Api`, `Tls`, …) are **Pascal case**, never all-caps | Yes        |
| R2  | JSON / serialization keys are **Pascal case** (`"NewIp"`, `"ApiUrl"`)                             | Yes        |
| R3  | Every boolean starts with `is` / `has` (`Is` / `Has` when Pascal-cased)                           | Yes        |
| R4  | > 3 parameters **or** signature line > 100 chars → one parameter per line                         | Yes        |
| R5  | > 4 parameters **or** 2+ adjacent same-typed parameters → params struct / options object          | Yes        |
| R6  | Every parameter is used, or explicitly discarded with a comment                                   | Yes        |
| R7  | Every error/failure is propagated with context; nothing swallowed                                 | Yes        |
| R8  | No magic literals passed as arguments — extract named constants                                   | Yes        |
| R9  | Call > 100 chars or > 4 args → one argument per line (mirrors R4 at the call site)                | Yes        |
| R10 | No boolean positional parameters → named field / enum                                             | Suggestion |
| R11 | `context.Context` (Go) / cancellation token (C#) placed per language convention                   | Yes        |
| R12 | Return values documented; multi-return meanings unambiguous                                       | Suggestion |
| R13 | One blank line before every `return` / `throw`, unless it is the only statement in the block      | Yes        |
| R14 | One blank line after a closing `}`, unless the next line is `}`, `else`, `case`, or `catch`       | Yes        |
| R15 | Never two blank lines in a row, anywhere                                                          | Yes        |
| R16 | No blank line immediately after `{` or immediately before `}`                                     | Yes        |
| R17 | One blank line between top-level declarations                                                     | Yes        |
| R18 | Imports grouped stdlib / third-party / first-party absolute / first-party relative                | Yes        |
| R19 | Trailing newline at end of file; no trailing whitespace on any line                               | Yes        |
| R20 | No section-separator blank lines inside a function — refactor instead                             | Suggestion |
| R21 | No per-field comments inside structs unless the field is genuinely not self-explanatory           | Yes        |

---

## 2. R1 — Acronym casing is Pascal case

**Acronyms are written as ordinary words: first letter capitalized, the rest lowercase.**
`IP` → `Ip`. `JSON` → `Json`. `URL` → `Url`. `ID` → `Id`. `HTTP` → `Http`. `TLS` → `Tls`.
`API` → `Api`. `UUID` → `Uuid`. `SQL` → `Sql`. `HTML` → `Html`.

**All-uppercase acronyms are a violation**, in Go and TypeScript especially. Do **not** write
`UserID`, `HTTPClient`, `JSONPayload`, `TLSConfig`, `parseJSONURL`, `swapIPWindows`.
Write `UserId`, `HttpClient`, `JsonPayload`, `TlsConfig`, `parseJsonUrl`, `swapIpWindows`.

Half-capitalized forms are equally wrong: never `Jsonurl`, `iD`, `uRL`.

**Acronym list to normalize:**
`Acl Api Ascii Cpu Css Db Dns Eof Guid Html Http Https Id Ip Json Lhs Qps Ram Rhs Rpc Sla Smtp Sql Ssh Tcp Tls Ttl Udp Ui Uid Uri Url Utf8 Uuid Vm Xml Xsrf Xss`

### Per-language mapping

| Language        | Exported / public                               | Unexported / local                | Notes                                                                       |
| --------------- | ----------------------------------------------- | --------------------------------- | ----------------------------------------------------------------------------- |
| Go              | `UserId`, `ParseJsonUrl`, `SwapIpWindows`       | `userId`, `parseJsonUrl`, `newIp` | Unexported = same word shape, lowercase first letter: `ipAddress`, `jsonBody` |
| TypeScript / JS | `type UserId`, `class HttpClient`               | `userId`, `apiUrl`                | Types/classes PascalCase: `HttpClient` — never `HTTPClient`                   |
| Java            | `getUserId()`, `class HttpClient`               | `userId`, `apiUrl`                | `JsonParser`, not `JSONParser`                                                |
| C#              | `UserId`, `ParseJsonUrl()`                      | `userId`, `_httpClient`           | PascalCase members with Pascal-cased acronyms                                 |
| Python          | `user_id`, `parse_json_url`, `class JsonParser` | same                              | snake_case functions/vars; Pascal-cased acronyms in class names               |
| Rust            | `user_id`, `parse_json_url`, `struct JsonParser`| same                              | snake_case items; Pascal-cased acronyms in type names                         |

### Before / after

```go
// BEFORE
func swapIPWindows(ctx context.Context, interfaceName string, newIP string) error
type userJSONURL struct{ ID string; URL string }
type TLSConfig struct{ CACert string }
```

```go
// AFTER
func swapIpWindows(ctx context.Context, params SwapIpParams) error

type UserJsonUrl struct {
	Id  string
	Url string
}

type TlsConfig struct {
	CaCert string
}
```

```ts
// BEFORE
function fetchUserJSONURL(userID: string, apiURL: string) {}
class HTTPClient {}
interface JSONPayload {}
```

```ts
// AFTER
function fetchUserJsonUrl(userId: string, apiUrl: string) {}

class HttpClient {}

interface JsonPayload {}
```

```python

# BEFORE

def fetchUserJsonUrl(userId, apiUrl): ...
class JSONParser: ...
```

```python

# AFTER

def fetch_user_json_url(user_id: str, api_url: str) -> None: ...

class JsonParser: ...
```

```rust
// BEFORE
fn fetch_user_json_url(user_id: &str, api_url: &str) -> Result<UserJSONURL, Error>

// AFTER
fn fetch_user_json_url(user_id: &str, api_url: &str) -> Result<UserJsonUrl, Error>
```

```java
// BEFORE
public String getUserID(String apiURL) {}

// AFTER
public String getUserId(String apiUrl) {}
```

```csharp
// BEFORE
public string GetUserID(string apiURL) { }

// AFTER
public string GetUserId(string apiUrl) { }
```

---

## 3. R2 — JSON / serialization keys are Pascal case

Serialized keys use **Pascal case**, matching the Go field name exactly: `Id`, `NewIp`, `ApiUrl`,
`IsEnabled`. Never `id`, `new_ip`, `api_url`, `newIP`.

```go
// BEFORE
type User struct {
	ID     string `json:"id"`
	APIURL string `json:"api_url"`
	Active bool   `json:"active"`
}
```

```go
// AFTER
type User struct {
	Id       string `json:"Id"`
	ApiUrl   string `json:"ApiUrl"`
	IsActive bool   `json:"IsActive"`
}
```

Rules:

- Apply Pascal case across the **whole payload**. Never mix `Id` and `id` in one schema.
- Nested structs, slices of structs, and map keys used as field names follow the same rule.
- Custom marshalers, query-parameter binding, and OpenAPI/schema files are updated in the same
  change so the wire contract stays internally consistent.
- Other languages default to Pascal case for serialized keys as well, until a language-specific
  convention is explicitly decided for that service.

```ts
// AFTER
interface UserDto {
  Id: string;
  ApiUrl: string;
  IsActive: boolean;
}
```

```python

# AFTER — serialization alias is Pascal case even though Python fields are snake_case

class UserDto(BaseModel):
    id: str = Field(alias="Id")
    api_url: str = Field(alias="ApiUrl")
    is_active: bool = Field(alias="IsActive")
```

```rust
// AFTER
#[derive(Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct UserDto {
    pub id: String,
    pub api_url: String,
    pub is_active: bool,
}
```

```csharp
// AFTER — PascalCase is the default; do not apply a camelCase naming policy
public sealed record UserDto(string Id, string ApiUrl, bool IsActive);
```

---

## 4. R3 — Boolean naming: `is` or `has`

Every boolean — variable, parameter, struct field, JSON key, and boolean-returning function —
**must** start with `is` or `has` (`Is` / `Has` when the identifier is Pascal-cased).

- State / classification → `is`: `isEnabled`, `isStatic`, `isDryRun`, `isTlsVerified`.
- Possession / presence → `has`: `hasSubnetMask`, `hasRetries`, `hasAdminRole`.

Banned: `enabled`, `active`, `legacy`, `flag`, `verifyPeer`, `useTls`, `ok`, `done`, `valid`,
`retries` (when boolean), and any bare adjective or noun.

```go
// BEFORE
type Config struct {
	Enabled    bool `json:"enabled"`
	Legacy     bool `json:"legacy"`
	VerifyPeer bool `json:"verify_peer"`
}

func (c Config) Valid() bool
```

```go
// AFTER
type Config struct {
	IsEnabled      bool `json:"IsEnabled"`
	IsLegacy       bool `json:"IsLegacy"`
	IsPeerVerified bool `json:"IsPeerVerified"`
	HasSubnetMask  bool `json:"HasSubnetMask"`
}

func (c Config) IsValid() bool
```

```ts
// BEFORE
const enabled = true;
function valid(user: User): boolean {}
```

```ts
// AFTER
const isEnabled = true;

function isValid(user: User): boolean {}

function hasAdminRole(user: User): boolean {}
```

```python

# AFTER

is_enabled: bool = True

def has_admin_role(user: User) -> bool: ...
```

```rust
// AFTER
let is_enabled = true;

pub fn has_admin_role(user: &User) -> bool {}
```

```java
// AFTER
private boolean isEnabled;

public boolean hasAdminRole() {}
```

```csharp
// AFTER
public bool IsEnabled { get; init; }

public bool HasAdminRole() { }
```

---

## 5. R4 — Split the parameter list

**Trigger:** more than **3** parameters, **or** the signature line exceeds **100 characters**.
**Action:** one parameter per line, closing paren and return type on their own line, trailing comma
where the language allows it.

```go
// BEFORE (>100 chars, 5 params)
func swapIpWindows(ctx context.Context, interfaceName string, oldIp string, newIp string, subnetMask string) error {
```

```go
// AFTER
func swapIpWindows(
	ctx context.Context,
	interfaceName string,
	oldIp string,
	newIp string,
	subnetMask string,
) error {
```

Note: in Go, do **not** collapse to `oldIp, newIp, subnetMask string` when splitting — one name per
line, each with its explicit type.

```ts
// AFTER
export async function swapIpWindows(
  signal: AbortSignal,
  interfaceName: string,
  oldIp: string,
  newIp: string,
  subnetMask: string,
): Promise<void> {}
```

```python

# AFTER

def swap_ip_windows(
    ctx: Context,
    interface_name: str,
    old_ip: str,
    new_ip: str,
    subnet_mask: str,
) -> None: ...
```

```rust
// AFTER
pub fn swap_ip_windows(
    ctx: &Context,
    interface_name: &str,
    old_ip: &str,
    new_ip: &str,
    subnet_mask: &str,
) -> Result<(), AppError> {}
```

```java
// AFTER
public void swapIpWindows(
        Context ctx,
        String interfaceName,
        String oldIp,
        String newIp,
        String subnetMask) throws AppException {
}
```

```csharp
// AFTER
public Task SwapIpWindowsAsync(
    string interfaceName,
    string oldIp,
    string newIp,
    string subnetMask,
    CancellationToken cancellationToken)
{
}
```

---

## 6. R5 — Group parameters into a struct / options object

**Trigger:** more than **4** parameters, **or** two or more **adjacent parameters of the same type**
(they can be swapped at a call site with no compiler error — a silent bug).

`swapIp(ctx, interfaceName, oldIp, newIp string)` has three adjacent `string` parameters.
`swapIp(ctx, oldIp, interfaceName, newIp)` compiles and is wrong. Group them.

```go
// AFTER
type SwapIpParams struct {
	InterfaceName string `json:"InterfaceName"`
	OldIp         string `json:"OldIp"`
	NewIp         string `json:"NewIp"`
	SubnetMask    string `json:"SubnetMask"`
	IsDryRun      bool   `json:"IsDryRun"`
}

func SwapIp(ctx context.Context, params SwapIpParams) error {}

err := SwapIp(ctx, SwapIpParams{
	InterfaceName: defaultInterfaceName,
	OldIp:         current,
	NewIp:         desired,
	SubnetMask:    defaultSubnetMask,
	IsDryRun:      false,
})
```

```ts
interface SwapIpParams {
  interfaceName: string;
  oldIp: string;
  newIp: string;
  subnetMask: string;
  isDryRun: boolean;
}

export async function swapIp(signal: AbortSignal, params: SwapIpParams): Promise<void> {}
```

```python
@dataclass(frozen=True)
class SwapIpParams:
    interface_name: str
    old_ip: str
    new_ip: str
    subnet_mask: str
    is_dry_run: bool = False

def swap_ip(ctx: Context, params: SwapIpParams) -> None: ...

# Python alternative: keyword-only parameters

def swap_ip(ctx: Context, *, interface_name: str, old_ip: str, new_ip: str) -> None: ...
```

```rust
pub struct SwapIpParams<'a> {
    pub interface_name: &'a str,
    pub old_ip: &'a str,
    pub new_ip: &'a str,
    pub subnet_mask: &'a str,
    pub is_dry_run: bool,
}

pub fn swap_ip(ctx: &Context, params: SwapIpParams<'_>) -> Result<(), AppError> {}
```

```java
public record SwapIpParams(
        String interfaceName,
        String oldIp,
        String newIp,
        String subnetMask,
        boolean isDryRun) {}

public void swapIp(Context ctx, SwapIpParams params) {}
```

```csharp
public sealed record SwapIpParams(
    string InterfaceName,
    string OldIp,
    string NewIp,
    string SubnetMask,
    bool IsDryRun);

public Task SwapIpAsync(SwapIpParams parameters, CancellationToken cancellationToken) { }
```

---

## 7. R6 — No unused parameters

An unused parameter is either dead weight or a forgotten feature. Choose one:

1. **Remove it** and update all call sites (preferred).
2. **Use it** — if it was meant to be used, that is the actual bug.
3. **Discard explicitly** — only when an interface fixes the signature:

```go
func (h handler) Handle(_ context.Context, req Request) error { // ctx unused: interface-required
```

```python
def handle(self, _ctx: Context, req: Request) -> None:  # ctx unused: Protocol-required
```

```rust
fn handle(&self, _ctx: &Context, req: Request) -> Result<(), AppError> {}
```

```csharp
public Task HandleAsync(Request req, CancellationToken _) { } // token unused: interface-required
```

Always add the trailing comment explaining *why* the signature cannot shrink.

---

## 8. R7 — Propagate and wrap every error

Never swallow. Never return a bare error that loses the operation name. Always return early.

```go
// BEFORE
cmd.Run()
return nil
```

```go
// AFTER
if err := cmd.Run(); err != nil {
	return apperror.Wrap(err, "swapIpWindows", map[string]any{
		"InterfaceName": params.InterfaceName,
		"NewIp":         params.NewIp,
	})
}

return nil
```

Cross-language equivalents:

```ts
try {
  await runCommand(cmd);
} catch (cause) {
  throw new AppError("swapIpWindows failed", {
    cause,
    interfaceName: params.interfaceName,
  });
}
```

```python
try:
    run_command(cmd)
except OSError as exc:
    raise AppError("swap_ip_windows failed", interface_name=params.interface_name) from exc
```

```rust
run_command(&cmd)
    .with_context(|| format!("swap_ip_windows: iface={}", params.interface_name))?;
```

```java
try {
    runCommand(cmd);
} catch (IOException cause) {
    throw new AppException("swapIpWindows failed: iface=" + params.interfaceName(), cause);
}
```

```csharp
try {
    await RunCommandAsync(cmd, cancellationToken);
} catch (Exception cause) {
    throw new AppException($"SwapIpWindows failed: iface={parameters.InterfaceName}", cause);
}
```

Also under R7: a wrapper that drops a caller's argument on the floor (see §16, where `swapIp`
receives `oldIp` and never forwards it) is a defect, not a style issue.

---

## 9. R8 — No magic literals as arguments

```go
// BEFORE
netshExecutor(ctx, "netsh", "interface", "ip", "set", "address", "name="+interfaceName, "static", newIp, "255.255.255.0")
```

```go
// AFTER
const (
	defaultInterfaceName = "Ethernet"
	defaultSubnetMask    = "255.255.255.0"
	netshBinary          = "netsh"
	addressModeStatic    = "static"
)
```

The same applies in every language: `const`, `final`, `enum`, module-level constants, or a config
struct — never a bare literal at the call site.

---

## 10. R9 — Split long argument lists at call sites

R4 and R9 are the same rule applied at two places. **A split definition with a packed call site is a
violation, and so is the reverse.**

**Trigger:** the call spans more than 100 characters, or passes more than 4 arguments.
**Action:** one argument per line, closing paren on its own line, trailing comma where the language
allows it. Never two arguments on one line of a split call.

```go
cmd := netshExecutor(
	ctx,
	netshBinary,
	"interface",
	"ip",
	"set",
	"address",
	"name="+params.InterfaceName,
	addressModeStatic,
	params.NewIp,
	params.SubnetMask,
)
```

When the arguments are a homogeneous list, prefer a slice/array variable — the literal's elements
may then be grouped by meaning:

```go
args := []string{
	"interface",
	"ip",
	"set",
	"address",
	"name=" + params.InterfaceName,
	addressModeStatic,
	params.NewIp,
	params.SubnetMask,
}

cmd := netshExecutor(ctx, netshBinary, args...)
```

```ts
await swapIpWindows(
  signal,
  interfaceName,
  oldIp,
  newIp,
  subnetMask,
);
```

```python
swap_ip_windows(
    ctx,
    interface_name,
    old_ip,
    new_ip,
    subnet_mask,
)
```

```rust
swap_ip_windows(
    ctx,
    interface_name,
    old_ip,
    new_ip,
    subnet_mask,
)?;
```

```java
swapIpWindows(
        ctx,
        interfaceName,
        oldIp,
        newIp,
        subnetMask);
```

```csharp
await SwapIpWindowsAsync(
    interfaceName,
    oldIp,
    newIp,
    subnetMask,
    cancellationToken);
```

---

## 11. R10 — No boolean positional parameters

`configure(host, true, false)` is unreadable. Use a named field or an enum — and the field still
obeys R3 (`is` / `has`).

```go
type AddressMode int

const (
	AddressModeStaticValue AddressMode = iota
	AddressModeDhcp
)

func Configure(ctx context.Context, params ConfigureParams) error
```

```ts
configure({ host, isTlsEnabled: true, isPeerVerified: false });
```

```python
configure(host, is_tls_enabled=True, is_peer_verified=False)
```

---

## 12. R11 — Context / cancellation position

- Go: `ctx context.Context` is always the first parameter, never stored in a struct field.
- C#: `CancellationToken` is the **last** parameter by .NET convention — apply the language
  convention, not Go's, and be consistent within the codebase.
- Rust / Python / TS: pass the cancellation handle (`&Context`, `AbortSignal`) explicitly rather
  than reading a global.

---

## 13. R12 — Document returns

Multi-value returns must be unambiguous. If two returns share a type, name them (Go named results in
the doc comment, TS a named object, Python a `NamedTuple`, Rust a struct).

```go
// SwapIp applies params.NewIp to the named interface.
// It returns a wrapped *apperror.Error when the netsh command fails, nil on success.
func SwapIp(ctx context.Context, params SwapIpParams) error
```

---

## 14. R13–R20 — Line gaps and whitespace

These rules govern blank lines. No mainstream formatter enforces them, so they never conflict with
`gofmt`, `prettier`, `black`, `rustfmt`, `spotless`, or `dotnet format`.

### R13 — Blank line before `return` / `throw`

One blank line before every `return`, `throw`, `raise`, or early exit — **unless it is the only
statement in the block**. A single-statement `if`, `for`, `match`, or function body keeps its
`return` tight against the brace.

```go
// AFTER
func classify(params SwapIpParams) string {
	if params.NewIp == "" {
		return classificationEmpty
	}

	mode := detectMode(params)
	if mode == AddressModeDhcp {
		return classificationDhcp
	}

	return classificationStatic
}
```

`classificationEmpty` and `classificationDhcp` are the only statements in their blocks, so no blank
line precedes them. The final `return` follows other statements, so it gets one.

```ts
export function classify(params: SwapIpParams): string {
  if (!params.newIp) {
    return CLASSIFICATION_EMPTY;
  }

  const mode = detectMode(params);
  if (mode === AddressMode.Dhcp) {
    return CLASSIFICATION_DHCP;
  }

  return CLASSIFICATION_STATIC;
}
```

```python
def classify(params: SwapIpParams) -> str:
    if not params.new_ip:
        return CLASSIFICATION_EMPTY

    mode = detect_mode(params)
    if mode is AddressMode.DHCP:
        return CLASSIFICATION_DHCP

    return CLASSIFICATION_STATIC
```

Python equivalents of `throw`: `raise`. Rust: the `?` operator counts as an early exit only when it
stands on its own statement line preceded by other statements; an expression-position tail value
needs no blank line if it is the block's only statement.

```rust
pub fn classify(params: &SwapIpParams) -> Result<Classification, AppError> {
    if params.new_ip.is_empty() {
        return Ok(Classification::Empty);
    }

    let mode = detect_mode(params)?;

    Ok(Classification::from(mode))
}
```

```java
public String classify(SwapIpParams params) {
    if (params.newIp().isEmpty()) {
        return CLASSIFICATION_EMPTY;
    }

    AddressMode mode = detectMode(params);

    return mode == AddressMode.DHCP ? CLASSIFICATION_DHCP : CLASSIFICATION_STATIC;
}
```

```csharp
public string Classify(SwapIpParams parameters)
{
    if (string.IsNullOrEmpty(parameters.NewIp))
    {
        return ClassificationEmpty;
    }

    AddressMode mode = DetectMode(parameters);

    return mode == AddressMode.Dhcp ? ClassificationDhcp : ClassificationStatic;
}
```

### R14 — Blank line after a closing `}`

One blank line after a closing brace, **unless** the next line is another `}`, `else`, `case`, or
`catch`. Chained constructs stay chained; independent statements get separated.

```go
// WRONG
if err := validate(params); err != nil {
	return err
}
cmd := build(params)
```

```go
// RIGHT
if err := validate(params); err != nil {
	return err
}

cmd := build(params)
```

```go
// RIGHT — no blank line before `else`
if params.IsDryRun {
	logPlan(cmd)
} else {
	runPlan(cmd)
}
```

Python has no braces: the equivalent is one blank line after a dedent that ends a block, unless the
next line is `else`, `elif`, `except`, `finally`, or `case`.

### R15 — Never two blank lines in a row

Anywhere, in any language. This overrides PEP 8's two-blank-line convention between top-level
Python definitions: use exactly one.

### R16 — No padded braces

No blank line immediately after `{`, and none immediately before `}`. The same applies to a Python
block's first and last lines.

```go
// WRONG
func run() error {

	return doWork()

}
```

```go
// RIGHT
func run() error {
	return doWork()
}
```

### R17 — One blank line between top-level declarations

Functions, classes, types, interfaces, and exported constants are separated by exactly one blank
line. Grouped `const`/`var` blocks count as one declaration.

### R18 — Import grouping

Exactly one blank line between groups, in this order, with no mixing:

1. Standard library
2. Third-party
3. First-party absolute
4. First-party relative

```go
import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"example.com/service/internal/apperror"

	"./netsh"
)
```

```ts
import { readFile } from "node:fs/promises";

import { z } from "zod";

import { AppError } from "@/lib/app-error";

import { netshExecutor } from "./netsh";
```

```python
import os
from contextlib import suppress

import httpx

from service.apperror import AppError

from .netsh import netsh_executor
```

Rust groups as `std` / external crates / `crate::` / `self::`+`super::`. Java groups as `java.*` and
`javax.*` / third-party / first-party. C# groups as `System.*` / third-party / first-party.

### R19 — File hygiene

Every file ends with exactly one trailing newline. No trailing whitespace on any line, including
otherwise-blank lines inside blocks.

### R20 — No section separators inside a function

Blank lines inside a function exist only because R13 or R14 requires them. If you want a blank line
to mark "this next part does something else", the function is doing more than one thing — extract
the part into its own function instead of adding whitespace. Suggestion-level, but a repeated
pattern of it should be raised as a design finding.

---

## 15. R21 — No per-field comments inside structs

Fields carry no comment by default. A well-named field with an explicit type documents itself, and a
comment on every field is noise that hides the one field that actually needed explaining.

A field comment is allowed **only** when the meaning cannot be inferred from name and type:

- units or scale (`TimeoutMs`, `SizeBytes` where the name cannot carry it)
- invariants or valid ranges the type does not express
- a non-obvious default applied elsewhere
- a wire-contract quirk kept for backward compatibility

Doc comments on the **type itself** stay and are encouraged.

```go
// WRONG — every field commented, all of them obvious
type SwapIpParams struct {
	// InterfaceName is the OS interface name.
	InterfaceName string `json:"InterfaceName"`
	// NewIp is the static address to assign.
	NewIp string `json:"NewIp"`
	// IsDryRun reports whether the command is logged instead of executed.
	IsDryRun bool `json:"IsDryRun"`
}
```

```go
// RIGHT — type documented, fields speak for themselves, one earned comment
// SwapIpParams describes a single static-Ip reassignment on one network interface.
// Empty InterfaceName and SubnetMask fall back to the package defaults.
type SwapIpParams struct {
	InterfaceName string `json:"InterfaceName"`
	// OldIp is advisory only: it is recorded in error context, never asserted.
	OldIp      string `json:"OldIp"`
	NewIp      string `json:"NewIp"`
	SubnetMask string `json:"SubnetMask"`
	IsDryRun   bool   `json:"IsDryRun"`
}
```

The same applies to TypeScript interfaces, Python dataclasses and Pydantic models, Rust structs,
Java records, and C# records.

```ts
/** A single static-Ip reassignment on one network interface. */
interface SwapIpParams {
  interfaceName: string;
  /** Advisory only: recorded in error context, never asserted. */
  oldIp: string;
  newIp: string;
  subnetMask: string;
  isDryRun: boolean;
}
```

---

## 16. Worked example — the reviewed snippet

### Input

```go
func swapIPWindows(ctx context.Context, interfaceName string, newIP string) error {
	if interfaceName == "" {
		interfaceName = "Ethernet"
	}
	cmd := netshExecutor(ctx, "netsh", "interface", "ip", "set", "address", "name="+interfaceName, "static", newIP, "255.255.255.0")
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "swapIPWindows", map[string]any{"interface": interfaceName})
	}
	return nil
}

func swapIP(ctx context.Context, interfaceName, oldIP, newIP string) error {
	return swapIPWindows(ctx, interfaceName, newIP)
}
```

### Findings

| Rule | Location        | Finding                                                                                                              |
| ---- | --------------- | ---------------------------------------------------------------------------------------------------------------------- |
| R1   | whole file      | `swapIPWindows`, `newIP`, `oldIP`, `swapIP` use all-caps `IP`. Rename to `swapIpWindows`, `newIp`, `oldIp`, `SwapIp`.    |
| R2   | error context   | Context key `"interface"` is lowercase; serialization keys must be Pascal case (`"InterfaceName"`).                      |
| R5   | `swapIp`        | `interfaceName, oldIp, newIp string` — three adjacent same-typed params; trivially swappable. Group into `SwapIpParams`. |
| R6   | `swapIp`        | `oldIp` is accepted and never used or forwarded. Defect.                                                                 |
| R8   | `swapIpWindows` | `"Ethernet"`, `"255.255.255.0"`, `"netsh"`, `"static"` are magic literals.                                               |
| R9   | `swapIpWindows` | 130-char, 10-argument call on one line. Split, or build an args slice.                                                   |
| R12  | both            | No doc comment on the error contract.                                                                                    |
| R13  | `swapIpWindows` | Final `return nil` follows other statements with no blank line before it.                                                |
| R14  | `swapIpWindows` | No blank line after the closing `}` of the `if interfaceName == ""` block or the `cmd.Run()` block.                      |

### Output

```go
const (
	defaultInterfaceName = "Ethernet"
	defaultSubnetMask    = "255.255.255.0"
	netshBinary          = "netsh"
	addressModeStatic    = "static"
)

// SwapIpParams describes a single static-Ip reassignment on one network interface.
// Empty InterfaceName and SubnetMask fall back to the package defaults.
type SwapIpParams struct {
	InterfaceName string `json:"InterfaceName"`
	// OldIp is advisory only: it is recorded in error context, never asserted.
	OldIp      string `json:"OldIp"`
	NewIp      string `json:"NewIp"`
	SubnetMask string `json:"SubnetMask"`
	IsDryRun   bool   `json:"IsDryRun"`
}

func (p SwapIpParams) withDefaults() SwapIpParams {
	if p.InterfaceName == "" {
		p.InterfaceName = defaultInterfaceName
	}

	if p.SubnetMask == "" {
		p.SubnetMask = defaultSubnetMask
	}

	return p
}

// swapIpWindows assigns params.NewIp to params.InterfaceName via netsh.
// It returns a wrapped error when the netsh command fails, nil on success.
func swapIpWindows(ctx context.Context, params SwapIpParams) error {
	params = params.withDefaults()

	args := []string{
		"interface",
		"ip",
		"set",
		"address",
		"name=" + params.InterfaceName,
		addressModeStatic,
		params.NewIp,
		params.SubnetMask,
	}

	cmd := netshExecutor(ctx, netshBinary, args...)
	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, "swapIpWindows", map[string]any{
			"InterfaceName": params.InterfaceName,
			"OldIp":         params.OldIp,
			"NewIp":         params.NewIp,
		})
	}

	return nil
}

// SwapIp dispatches a static-Ip swap to the platform-specific implementation.
// It returns a wrapped error on failure, nil on success.
func SwapIp(ctx context.Context, params SwapIpParams) error {
	return swapIpWindows(ctx, params)
}
```

`SwapIp`'s body is a single statement, so its `return` keeps no blank line before it (R13).
`OldIp` is now carried through to the error context instead of being silently discarded — if the
intent was a precondition check, add it explicitly rather than leaving a dead parameter.

---

## 17. Per-language appendix — same violations, corrected

Each block shows the same function violating R1, R3, R4, R5, and R6, then corrected under the full
ruleset including R13–R21.

### Go

```go
// BEFORE
func createUser(ctx context.Context, ID string, name string, email string, role string, tenant string, legacy bool) error
```

```go
// AFTER
type CreateUserParams struct {
	UserId   string `json:"UserId"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Role     Role   `json:"Role"`
	Tenant   string `json:"Tenant"`
	IsLegacy bool   `json:"IsLegacy"`
}

func CreateUser(ctx context.Context, params CreateUserParams) error
```

### TypeScript / JavaScript

```ts
// BEFORE
function createUser(id, name, email, role, tenant, legacyFlag) {}
```

```ts
// AFTER
interface CreateUserParams {
  userId: string;
  name: string;
  email: string;
  role: Role;
  tenant: string;
  isLegacy: boolean;
}

export function createUser(params: CreateUserParams): Promise<User> {}
```

### Python

```python

# BEFORE

def create_user(id, name, email, role, tenant, legacy_flag): ...
```

```python

# AFTER

@dataclass(frozen=True)
class CreateUserParams:
    user_id: str
    name: str
    email: str
    role: Role
    tenant: str
    is_legacy: bool = False

def create_user(params: CreateUserParams) -> User: ...
```

### Rust

```rust
// BEFORE
fn create_user(id: &str, name: &str, email: &str, role: &str, tenant: &str, legacy: bool) -> Result<User, AppError>
```

```rust
// AFTER
pub struct CreateUserParams<'a> {
    pub user_id: &'a str,
    pub name: &'a str,
    pub email: &'a str,
    pub role: Role,
    pub tenant: &'a str,
    pub is_legacy: bool,
}

pub fn create_user(params: CreateUserParams<'_>) -> Result<User, AppError> {}
```

### Java

```java
// BEFORE
public User createUser(String ID, String name, String email, String role, String tenant, boolean legacy)
```

```java
// AFTER
public record CreateUserParams(
        String userId,
        String name,
        String email,
        Role role,
        String tenant,
        boolean isLegacy) {}

public User createUser(CreateUserParams params) {}
```

### C#

```csharp
// BEFORE
public Task<User> CreateUserAsync(string ID, string name, string email, string role, string tenant, bool legacy)
```

```csharp
// AFTER
public sealed record CreateUserParams(
    string UserId,
    string Name,
    string Email,
    Role Role,
    string Tenant,
    bool IsLegacy);

public Task<User> CreateUserAsync(
    CreateUserParams parameters,
    CancellationToken cancellationToken = default);
```

---

## 18. Linter specification — how the AI verifies its own output

This section is the machine-checkable form of the guide. After rewriting code, walk it top to
bottom and emit one pass/fail line per rule. Rules marked auto-fixable must be applied silently in
the rewrite; the rest are reported.

### Rule index

| Rule | Machine-checkable statement                                                                    | Severity   | Auto-fixable |
| ---- | ------------------------------------------------------------------------------------------------ | ---------- | ------------ |
| R1   | No identifier contains two or more consecutive uppercase letters                                 | Must fix   | Yes          |
| R2   | Every serialization key matches `^[A-Z][A-Za-z0-9]*$`                                            | Must fix   | Yes          |
| R3   | Every boolean identifier matches `^_?(is\|has\|Is\|Has)[A-Z0-9]`                                 | Must fix   | Yes          |
| R4   | Signature line ≤ 100 chars and ≤ 3 params, else one param per line                               | Must fix   | Yes          |
| R5   | ≤ 4 params and no two adjacent params share a type                                               | Must fix   | No           |
| R6   | Every param appears in the body, or is prefixed `_` with a trailing reason comment               | Must fix   | No           |
| R7   | No call whose error/exception result is discarded; every propagation adds operation context      | Must fix   | No           |
| R8   | No string/number literal passed as an argument outside a `const`/`enum` declaration              | Must fix   | No           |
| R9   | Call line ≤ 100 chars and ≤ 4 args, else exactly one argument per line                           | Must fix   | Yes          |
| R10  | No boolean literal in a positional argument slot                                                 | Suggestion | No           |
| R11  | Go: param 0 is `context.Context`. C#: last param is `CancellationToken`                          | Must fix   | Yes          |
| R12  | Exported function returning 2+ values or an error has a doc comment naming the contract          | Suggestion | No           |
| R13  | Line before a `return`/`throw`/`raise` is blank, unless that statement is alone in its block      | Must fix   | Yes          |
| R14  | Line after a closing `}` is blank, unless it is `}`, `else`, `case`, `catch`, or EOF              | Must fix   | Yes          |
| R15  | No two consecutive blank lines                                                                    | Must fix   | Yes          |
| R16  | Line after `{` is non-blank; line before `}` is non-blank                                        | Must fix   | Yes          |
| R17  | Exactly one blank line between consecutive top-level declarations                                | Must fix   | Yes          |
| R18  | Import block groups are ordered stdlib → third-party → first-party abs → relative, one blank line apart | Must fix | Yes    |
| R19  | File ends with exactly one `\n`; no line matches `[ \t]+$`                                       | Must fix   | Yes          |
| R20  | No blank line inside a function that is not required by R13 or R14                               | Suggestion | No           |
| R21  | No comment attached to a struct field unless it states units, invariants, defaults, or wire quirks | Must fix | No           |

### Machine-readable rules

```yaml
version: 2
rules:
  - id: R1
    kind: naming
    severity: must-fix
    autofix: true
    deny: '[A-Z]{2,}'
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R2
    kind: serialization
    severity: must-fix
    autofix: true
    require: '^[A-Z][A-Za-z0-9]*$'
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R3
    kind: naming
    severity: must-fix
    autofix: true
    require: '^_?(is|has|Is|Has)[A-Z0-9]'
    target: boolean-identifiers
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R4
    kind: structure
    severity: must-fix
    autofix: true
    max_params_inline: 3
    max_line_length: 100
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R5
    kind: structure
    severity: must-fix
    autofix: false
    max_params: 4
    deny_adjacent_same_type: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R6
    kind: structure
    severity: must-fix
    autofix: false
    require_used_or_underscore_with_comment: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R7
    kind: errors
    severity: must-fix
    autofix: false
    deny_discarded_errors: true
    require_context_on_propagation: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R8
    kind: literals
    severity: must-fix
    autofix: false
    deny_literal_arguments: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R9
    kind: structure
    severity: must-fix
    autofix: true
    max_args_inline: 4
    max_line_length: 100
    one_arg_per_line_when_split: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R10
    kind: structure
    severity: suggestion
    autofix: false
    deny_positional_boolean: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R11
    kind: structure
    severity: must-fix
    autofix: true
    go: context-first
    csharp: cancellation-token-last
    applies_to: [go, csharp]
  - id: R12
    kind: docs
    severity: suggestion
    autofix: false
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R13
    kind: whitespace
    severity: must-fix
    autofix: true
    blank_line_before: [return, throw, raise]
    except_when_only_statement_in_block: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R14
    kind: whitespace
    severity: must-fix
    autofix: true
    blank_line_after_block_end: true
    except_next_line_matches: ['}', 'else', 'case', 'catch', 'elif', 'except', 'finally', 'EOF']
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R15
    kind: whitespace
    severity: must-fix
    autofix: true
    max_consecutive_blank_lines: 1
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R16
    kind: whitespace
    severity: must-fix
    autofix: true
    deny_blank_after_open_brace: true
    deny_blank_before_close_brace: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R17
    kind: whitespace
    severity: must-fix
    autofix: true
    blank_lines_between_top_level_decls: 1
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R18
    kind: imports
    severity: must-fix
    autofix: true
    groups: [stdlib, third-party, first-party-absolute, first-party-relative]
    separator_blank_lines: 1
    deny_mixed_groups: true
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R19
    kind: whitespace
    severity: must-fix
    autofix: true
    require_trailing_newline: true
    deny: '[ \t]+$'
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R20
    kind: whitespace
    severity: suggestion
    autofix: false
    deny_blank_lines_not_required_by: [R13, R14]
    scope: function-body
    applies_to: [go, ts, js, python, rust, java, csharp]
  - id: R21
    kind: docs
    severity: must-fix
    autofix: false
    deny_field_comments_unless:
      [units, invariants, non-obvious-default, wire-contract-quirk]
    applies_to: [go, ts, js, python, rust, java, csharp]
formatter_precedence:
  note: >-
    gofmt, prettier, black, rustfmt, spotless, and dotnet format win only where they
    actually conflict, and only on whitespace they control. R13-R20 govern blank lines
    none of these formatters enforce, so they always apply.
```

### Self-verification procedure

1. Apply every auto-fixable rule to the rewritten code.
2. Re-read the rewritten code as if it were new input and walk R1 → R21 in order.
3. Emit `R{n}: pass` or `R{n}: fail — <location>` for each rule. A `fail` means go back to step 1.
4. Only after every must-fix rule passes, print the rewritten code.

---

## 19. Severity table

| Rule                                            | Severity   | Blocks review? |
| ----------------------------------------------- | ---------- | -------------- |
| R1 acronym casing (Pascal case, never all-caps) | Must fix   | Yes            |
| R2 Pascal-case serialization keys               | Must fix   | Yes            |
| R3 boolean `is` / `has` prefix                  | Must fix   | Yes            |
| R4 split parameters (>3 or >100 chars)          | Must fix   | Yes            |
| R5 params struct (>4 or adjacent same types)    | Must fix   | Yes            |
| R6 unused parameter                             | Must fix   | Yes            |
| R7 error propagation / wrapping                 | Must fix   | Yes            |
| R8 magic literals                               | Must fix   | Yes            |
| R9 split long call arguments                    | Must fix   | Yes            |
| R10 boolean positional parameter                | Suggestion | No             |
| R11 context / cancellation position             | Must fix   | Yes            |
| R12 documented returns                          | Suggestion | No             |
| R13 blank line before `return` / `throw`        | Must fix   | Yes            |
| R14 blank line after closing brace              | Must fix   | Yes            |
| R15 no double blank lines                       | Must fix   | Yes            |
| R16 no padded braces                            | Must fix   | Yes            |
| R17 blank line between top-level declarations   | Must fix   | Yes            |
| R18 import grouping                             | Must fix   | Yes            |
| R19 trailing newline / no trailing whitespace   | Must fix   | Yes            |
| R20 no in-function section separators           | Suggestion | No             |
| R21 no unnecessary struct field comments        | Must fix   | Yes            |

---

## 20. Review output format the AI should produce

```text

### Findings

- R1 (must fix) — `swapIPWindows`: all-caps acronym; rename to `swapIpWindows`, `newIP` → `newIp`.
- R5 (must fix) — `swapIp`: three adjacent `string` parameters; group into `SwapIpParams`.
- R6 (must fix) — `swapIp`: `oldIp` accepted but never used.
- R9 (must fix) — `swapIpWindows`: 130-char, 10-argument call to `netshExecutor`.
- R13 (must fix) — `swapIpWindows`: no blank line before the final `return nil`.
- R21 (must fix) — `SwapIpParams`: every field commented; keep only the `OldIp` note.

### Self-verification

R1: pass
R2: pass
...
R21: pass

### Rewritten code

<full corrected file or unified diff>
```
