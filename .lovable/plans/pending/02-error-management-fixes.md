# Error Management Migration Plan

This plan documents all required fixes to migrate the codebase to strict error management using `AppError`.

## Rules to Enforce:
1. **Never swallow errors.** Every catch/ignore logs the operation name and key inputs, then rethrows or returns a typed error.
2. **Wrap, do not lose.** Wrap the original error with an operation label and context (`apperror.Wrap(err, 'op', ctx)`).
3. **Typed errors only.** No bare `panic('msg')` or raw `fmt.Errorf` without wrapping.
4. Every variable needs to be captured in an error log/path.
5. All functions returning errors must return `*apperror.AppError`, not raw `error`.

## Subtasks Overview:
- [ ] `01-swallowed-error-part1.md`: Fix 62 `swallowed error` issues.
- [ ] `02-panic-part1.md`: Fix 5 `panic` issues.
- [ ] `03-raw-error-return-part1.md`: Fix 150 `raw error return` issues.
- [ ] `04-raw-error-return-part2.md`: Fix 150 `raw error return` issues.
- [ ] `05-raw-error-return-part3.md`: Fix 150 `raw error return` issues.
- [ ] `06-raw-error-return-part4.md`: Fix 150 `raw error return` issues.
- [ ] `07-raw-error-return-part5.md`: Fix 150 `raw error return` issues.
- [ ] `08-raw-error-return-part6.md`: Fix 150 `raw error return` issues.
- [ ] `09-raw-error-return-part7.md`: Fix 150 `raw error return` issues.
- [ ] `10-raw-error-return-part8.md`: Fix 150 `raw error return` issues.
- [ ] `11-raw-error-return-part9.md`: Fix 150 `raw error return` issues.
- [ ] `12-raw-error-return-part10.md`: Fix 85 `raw error return` issues.
- [ ] `13-fmtErrorf-part1.md`: Fix 150 `fmt.Errorf` issues.
- [ ] `14-fmtErrorf-part2.md`: Fix 150 `fmt.Errorf` issues.
- [ ] `15-fmtErrorf-part3.md`: Fix 150 `fmt.Errorf` issues.
- [ ] `16-fmtErrorf-part4.md`: Fix 150 `fmt.Errorf` issues.
- [ ] `17-fmtErrorf-part5.md`: Fix 150 `fmt.Errorf` issues.
- [ ] `18-fmtErrorf-part6.md`: Fix 150 `fmt.Errorf` issues.
- [ ] `19-fmtErrorf-part7.md`: Fix 45 `fmt.Errorf` issues.
