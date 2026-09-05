# regexnew — Lazy-Compiled Regular Expression Engine

Package `regexnew` provides a **lazy-loaded, thread-safe** regular expression caching engine following the New Creator Pattern. Patterns are compiled at most once, cached globally in synchronized registries (`regexMaps`, `lazyRegexOnceMap`), and reused across goroutines. All methods handle `nil` `*LazyRegex` receivers gracefully without panicking.

## Architecture

```
04-code/golang/pkg/regexnew/
├── consts.go                 # Default capacity, formatting constants, and callback types
├── vars.go                   # Package-level singletons: New, regexMaps, lazyRegexOnceMap, mutexes
├── new_creator.go            # New.* — Lazy, LazyLock, Default, DefaultLock entry points
├── new_lazy_regex_creator.go # New.LazyRegex.* — New, NewLock, TwoLock, ManyUsingLock
├── lazy_regex.go             # LazyRegex struct: lazy-compiled, cached, nil-safe receiver methods
├── lazy_regex_map.go         # Global pattern -> *LazyRegex cache with mutex synchronization
├── compile_helpers.go        # Create, CreateLock, CreateMust, NewMustLock helpers
├── match_helpers.go          # IsMatchLock, IsMatchFailed, MatchError, validation errors
├── pretty_json.go            # Indented JSON serializer for LazyRegex state
├── regexnew_test.go          # Full unit and concurrency test suite
└── readme.md                 # Technical architecture and usage guide
```

## New Creator Structure

```
regexnew.New (newCreator)
├── Lazy(pattern)              -> *LazyRegex (package var-level, lockless)
├── LazyLock(pattern)          -> *LazyRegex (method-level, locked)
├── Default(pattern)           -> (*regexp.Regexp, error)
├── DefaultLock(pattern)       -> (*regexp.Regexp, error)
├── DefaultLockIf(bool, str)   -> (*regexp.Regexp, error)
├── DefaultApplicableLock(str) -> (*regexp.Regexp, error, bool)
└── LazyRegex (newLazyRegexCreator)
    ├── New(pattern)           -> *LazyRegex
    ├── NewLock(pattern)       -> *LazyRegex
    ├── NewLockIf(bool, str)   -> *LazyRegex
    ├── TwoLock(p1, p2)        -> (first, second *LazyRegex)
    ├── ManyUsingLock(ps...)   -> map[string]*LazyRegex
    └── AllPatternsMap()       -> map[string]*LazyRegex
```

## Thread Safety

| Creator | Lock Mechanism | Usage Guideline |
| :--- | :--- | :--- |
| `New.Lazy(pattern)` | None | Package-level `var` declarations (Go runtime guarantees single-goroutine init). |
| `New.LazyLock(pattern)` | `sync.Mutex` | Inside methods and concurrent handlers — safe for simultaneous calls. |
| `New.LazyRegex.TwoLock(p1, p2)` | Single `sync.Mutex` | Batch registration of two patterns under one mutex lock. |
| `New.LazyRegex.ManyUsingLock(ps...)` | Single `sync.Mutex` | Batch registration of multiple patterns under one mutex lock. |

All `LazyRegex` instances are deduplicated via `lazyRegexOnceMap` and compiled `*regexp.Regexp` pointers are cached in `regexMaps`, eliminating redundant allocations and CPU-intensive regex recompilations.
