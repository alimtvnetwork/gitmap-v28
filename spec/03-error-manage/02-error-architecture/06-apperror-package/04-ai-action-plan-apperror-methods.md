# Action Plan for AI: Implementing Contact, Errors, and Msg in AppError

> **Target:** Apperror package (e.g., gitlab.com/auk-go/core/apperror)
> **Goal:** Extend AppError with new fields and fluent setters as requested by the spec.

When an AI agent is instructed to implement the new AppError methods (contact, errors, msg), the agent MUST follow these exact steps to ensure spec compliance.

---

## 1. Modify the AppError Struct

**File to edit:** Apperror.go (or wherever AppError is defined)

Add the new fields to the AppError struct, matching the JSON tags required by the spec:

```go
type AppError struct {
    Code         string            json:"Code"
    Message      string            json:"Message"
    DisplayError string            json:"DisplayError,omitempty"
    Details      string            json:"Details,omitempty"
    Contact      string            json:"Contact,omitempty"      // NEW
    Values       map[string]string json:"Values,omitempty"
    Diagnostic   ErrorDiagnostic   json:"Diagnostic,omitempty"
    Errors       []error           json:"Errors,omitempty"       // NEW
    Stack        StackTrace        json:"Stack"
    Cause        error             json:"-"
}
```
*(Note: Ensure you do not accidentally overwrite existing JSON tags if they are explicitly declared in the codebase.)*

## 2. Implement the Fluent Setters

**File to edit:** Apperror_setters.go (or append to Apperror.go)

Create the fluent setter methods for AppError:

```go
// WithContact sets support contact information (email, URL, Slack channel)
// and returns the modified AppError for chaining.
func (e *AppError) WithContact(contact string) *AppError {
    if e == nil {
        return nil
    }
    e.Contact = contact
    return e
}

// WithErrors appends a list of sub-errors (e.g., batch validation failures)
// and returns the modified AppError for chaining.
func (e *AppError) WithErrors(errs ...error) *AppError {
    if e == nil {
        return nil
    }
    e.Errors = append(e.Errors, errs...)
    return e
}

// WithMsg overrides the developer-facing Message on the fly
// and returns the modified AppError for chaining.
func (e *AppError) WithMsg(msg string) *AppError {
    if e == nil {
        return nil
    }
    e.Message = msg
    return e
}
```

## 3. Update Display Methods (String, FullString)

**File to edit:** Apperror_display.go (or equivalent stringification file)

Ensure that Contact and Errors are rendered when the error is printed to logs via FullString().

**Checklist for FullString() / ToClipboard():**

- [ ] If e.Contact != "", append it to the output string (e.g., `fmt`.Sprintf(" Contact: %s", e.Contact)).
- [ ] If len(e.Errors) > 0, iterate over e.Errors and print each one, indenting them under the main error.

*Example stringification update:*
```go
if e.Contact != "" {
    sb.WriteString(fmt.Sprintf(" | Contact: %s", e.Contact))
}
if len(e.Errors) > 0 {
    sb.WriteString("\n  Sub-errors:")
    for i, subErr := range e.Errors {
        sb.WriteString(fmt.Sprintf("\n    %d. %s", i+1, subErr.Error()))
    }
}
```

## 4. Run Unit Tests & Linting

- Write unit tests in Apperror_test.go to ensure chaining WithContact("foo").WithMsg("bar") works.
- Verify json.Marshal(appErr) correctly includes "Contact" and "Errors" if set, and omits them if empty.
- Run go run linter-scripts/validate-guidelines.go --path . to ensure no CODE-RED violations were introduced.

# Action Plan for AI: Implementing Apperror.New Namespace

> **Goal:** Implement the new creator struct namespace to elegantly wrap errors and allow early
il returns.

## 1. Create the creator struct and New variable

**File to edit:** Apperror_creator.go (new file) or Apperror.go

```go
type creator struct{}

// New exposes the fluent creator methods for AppError.
var New = creator{}
```

## 2. Implement the Creator Methods

**File to edit:** Apperror_creator.go

You MUST implement the early if err == nil { return nil } check for all methods that take an error.

```go
// Error wraps an existing error. If no error, nothing is created.
func (c creator) Error(errType apperrtype.ErrorType, err error) *AppError {
    if err == nil {
        return nil
    }
    return WrapType(err, errType)
}

// UsingErrorMsg wraps an error and overrides the default message.
func (c creator) UsingErrorMsg(errType apperrtype.ErrorType, err error, msg string) *AppError {
    if err == nil {
        return nil
    }
    return WrapTypeMsg(err, errType, msg)
}

// UsingMsg creates a new error from scratch (always returns an error).
func (c creator) UsingMsg(errType apperrtype.ErrorType, msg string) *AppError {
    appErr := NewType(errType)
    appErr.Message = msg
    return appErr
}

// ErrorVar wraps an error and injects a single variable.
func (c creator) ErrorVar(errType apperrtype.ErrorType, err error, varName string, varValue any) *AppError {
    if err == nil {
        return nil
    }
    // Note: Assuming WithValue is implemented or can accept stringified values
    return WrapType(err, errType).WithValue(varName, fmt.Sprintf("%v", varValue))
}

// ErrorVars wraps an error and injects multiple variables.
func (c creator) ErrorVars(errType apperrtype.ErrorType, err error, vars map[string]any) *AppError {
    if err == nil {
        return nil
    }
    appErr := WrapType(err, errType)
    for k, v := range vars {
        appErr = appErr.WithValue(k, fmt.Sprintf("%v", v))
    }
    return appErr
}
```

## 3. Verify nil-safety in tests

Write tests ensuring `apperror.New.Error(t, nil)` strictly returns a typeless
`nil` or typed `(*AppError)(nil)` so that `err != nil` checks do not falsely trigger in Go.

---

# Action Plan for AI: Implementing Human & Logger Output Methods

> **Goal:** Extend `AppError` to support segregated output formats for humans (UI), loggers (JSON), and developers (Console).

## 1. Implement `HumanString()`

**File to edit:** `apperror_display.go`

This method must act as a firewall, ensuring no technical details leak to the end user.

```go
// HumanString returns the safe, user-facing error message.
func (e *AppError) HumanString() string {
    if e == nil {
        return ""
    }
    if e.DisplayError != "" {
        return e.DisplayError
    }
    // Strict fallback if the developer forgot to set a DisplayError
    return "An unexpected error occurred. Please contact support if the issue persists."
}
```

## 2. Implement `LogMap()`

**File to edit:** `apperror_display.go`

This method translates the error into a flat/nested map for structured loggers (e.g., Zap, Logrus).

```go
// LogMap returns structured fields for JSON loggers.
func (e *AppError) LogMap() map[string]any {
    if e == nil {
        return nil
    }

    m := map[string]any{
        "error_code": e.Code,
        "error_msg":  e.Message,
    }

    if e.Details != "" {
        m["error_details"] = e.Details
    }
    if len(e.Values) > 0 {
        m["error_values"] = e.Values
    }
    // You must map the Diagnostic fields if they exist
    // You must include the stack trace array
    m["error_stack"] = e.Stack.Frames // Assuming StackTrace has a slice of frames

    if len(e.Errors) > 0 {
        var subErrs []string
        for _, sub := range e.Errors {
            subErrs = append(subErrs, sub.Error())
        }
        m["error_sub_errors"] = subErrs
    }

    return m
}
```

## 3. Implement `ConsoleString()`

**File to edit:** `apperror_display.go`

Alias this to the existing `FullString()` method, or replace `FullString()` with `ConsoleString()`.

```go
// ConsoleString returns the full, multi-line representation of the error
// designed for developer terminal output.
func (e *AppError) ConsoleString() string {
    return e.FullString()
}
```

---

# Action Plan for AI: Constructing Enums


## AI Implementation Examples & Non-Negotiable Checklist (`auk-go/errorwrapper`)

To ensure AI agents properly construct the `errtype` structure without hallucination, you MUST use the following files from `auk-go/errorwrapper` as your canonical reference. **Do not just copy URLs; you must read them and strictly implement the following non-negotiable checklist for each file.**

### 1. The Core Enum & Bounds
*   **`consts.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/consts.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/consts.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Must define the numerical bounds for the domain (e.g., Min, Max ranges). No hardcoded magic numbers outside of the `iota` block.
*   **`vars.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/vars.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/vars.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Declare global variables or defaults cleanly.

### 2. The Data Structures
*   **`VariantStructure.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/VariantStructure.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/VariantStructure.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Must define the rich metadata struct (`Name`, `Code`, `Message`, `Contact`). This struct must exactly match the JSON properties required by the API.
*   **`JsonModel.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/JsonModel.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/JsonModel.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Must implement strict JSON marshaling/unmarshaling without leaking internal stack traces to the public UI.

### 3. The Registries & Lookups
*   **`errTypesMap.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/errTypesMap.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/errTypesMap.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Every single enum variant MUST be registered here with its corresponding `VariantStructure`. Missing registrations are a CODE-RED violation.
*   **`StringToVariantMap.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/StringToVariantMap.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/StringToVariantMap.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Must provide a flawless reverse-lookup (`map[string]Variation`) to allow parsing from external systems.

### 4. The Methods & APIs
*   **`variation-functions.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/variation-functions.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/variation-functions.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Must implement `String()`, `Code()`, and `Message()` on the enum directly, pulling data from the `errTypesMap` registry.
*   **`istype.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/istype.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/istype.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Must implement type guard methods (e.g., `IsX`, `IsAny`) to allow clean boolean checks in business logic.
*   **`New.go`** ([View](https://gitlab.com/auk-go/errorwrapper/-/blob/develop/errtype/New.go?ref_type=heads) | [Raw](https://gitlab.com/auk-go/errorwrapper/-/raw/develop/errtype/New.go?ref_type=heads))
    *   [ ] **Non-negotiable:** Must implement the fluent constructor namespace (`New.Error`, `New.UsingMsg`) with early `if err == nil { return nil }` safety checks.

**AI Instruction:** Use your network reading tools (like `read_url_content` or `curl`) to download the `Raw` URLs above before you attempt to scaffold the Error Enum structure.
