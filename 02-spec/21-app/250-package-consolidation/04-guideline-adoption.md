# Spec 250 — Guideline adoption verdicts

> Candidate patterns/packages from `coding-guidelines-v24/04-code/golang/pkg/`,
> evaluated against gitmap-v28's actual needs. Rule: no churn for churn's sake
> — adopt only where the code has a real gap the guideline fills.

| Pattern / package | Location in coding-guidelines-v24 | Verdict | Reason |
|---|---|---|---|
| `streamwriter` | `04-code/golang/pkg/streamwriter/` | READ AS REFERENCE | Reference material for the stdout design in `03-stdout-and-stacktrace.md` (filter-chain composition ideas). NOT vendored — too heavyweight (async writer, payload converters, reflection) for what 250 needs. |
| `errcmd` | `04-code/golang/pkg/errcmd/` | DEFER | Genuine candidate — structured shell runner with `CommandResult`; gitmap shells out ad-hoc in hundreds of places. Too big for program 250; needs its own follow-up program with a migration plan. |
| `applogger` / `logger` | `04-code/golang/pkg/applogger/`, `04-code/golang/pkg/logger/` | DEFER | Genuine candidate — sink-based logging; gitmap has no central logger. Adopting either is a program-sized change (sink wiring, call-site migration), not 250 scope. |
| enum pattern | `04-code/golang/pkg/enum/` (+ `baseenumer`) | SKIP | Already covered — gitmap has its own enum pattern. |
| `appfault` | `04-code/golang/pkg/appfault/` | SKIP | Already covered — `apperror` in gitmap serves the same role. |
| `fileutil` | `04-code/golang/pkg/fileutil/` | SKIP | Already covered — atomic writes already exist in gitmap. |
| `lazyonce` | `04-code/golang/pkg/lazyonce/` | SKIP | Already covered — `≈ sync.Once` usage already in place. |
| result monads | `04-code/golang/pkg/result/` | SKIP | Already covered — partial result-monad usage already exists. |

## Notes

- DEFER ≠ reject: `errcmd` and `applogger`/`logger` are the two genuine
  follow-up candidates. Each deserves its own program with a scoped migration
  plan, not a drive-by adoption inside a consolidation program.
- SKIP means "no gap": adopting any of these would be churn for churn's sake.
- `streamwriter` stays a design reference only; the 250 `cli/output` package
  is deliberately lighter (synchronous `FilterWriter`, no async machinery).
