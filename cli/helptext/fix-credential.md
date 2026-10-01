# gitmap fix-credential

Repair git credential helper on Windows by configuring it to use 'manager' instead of unsupported 'cache'.

## Alias

fc

## Usage

    gitmap fix-credential [flags]
    gitmap fc [flags]

## Description

Resolves the `Can not use the 'cache' credential store on Windows due to lack of UNIX socket support` error by resetting `credential.helper` to Git Credential Manager (`manager`).

## Examples

### Example 1: Run credential store repair

    gitmap fix-credential

**Output:**

    Executing Git credential store fix for Windows...
    ✔ Git credential helper successfully updated to 'manager'.
      The 'cache' credential store issue on Windows is now resolved.

## See Also

- [fix](fix.md) — Universal repository fix dispatcher
