# Error Resolution Cross-Reference Diagram

> **Generated:** 2026-03-09  
**Version:** 3.2.0  
> **Status:** Complete  
> **Total Connections:** 18 outbound, 19 inbound

---

## Visual Architecture

```
                                    ┌─────────────────────────────────────────────┐
                                    │         ERROR RESOLUTION                     │
                                    │  spec/03-error-manage/01-error-resolution/  │
                                    │                                              │
                                    │  ┌──────────────────────────────────────┐   │
                                    │  │ 00-overview.md                       │   │
                                    │  │ 01-retrospectives/                   │   │
                                    │  │ 02-verification-patterns/            │   │
                                    │  │ 03-debugging-guides/                 │   │
                                    │  │   ├── 01-debugging-php.md            │   │
                                    │  │   ├── 02-debugging-go.md             │   │
                                    │  │   └── 03-debugging-typescript.md     │   │
                                    │  └──────────────────────────────────────┘   │
                                    └─────────────────────────────────────────────┘
                                                           │
                                                           │
            ┌──────────────────────────────────────────────┼──────────────────────────────────────────────┐
            │                                              │                                               │
            ▼                                              ▼                                               ▼
┌───────────────────────────┐              ┌───────────────────────────┐               ┌───────────────────────────┐
│    ARCHITECTURE SPECS     │              │       CLI TOOLS           │               │    WORDPRESS SPECS        │
│                           │              │                           │               │                           │
│ ┌───────────────────────┐ │              │ ┌───────────────────────┐ │               │ ┌───────────────────────┐ │
│ │ 04-split-db           │◄├──────────────┤►│ 09-gsearch-cli        │ │               │ │ 13-wp-plugin/         │ │
│ │   architecture/       │ │              │ │   [Go + TS debug]     │ │               │ │   exam-manager/       │ │
│ └───────────────────────┘ │              │ └───────────────────────┘ │               │ │   [PHP + Go debug]    │ │
│                           │              │                           │               │ └───────────────────────┘ │
│ ┌───────────────────────┐ │              │ ┌───────────────────────┐ │               │                           │
│ │ 05-seedable-config    │◄├──────────────┤►│ 10-brun-cli           │ │               │ ┌───────────────────────┐ │
│ │   architecture/       │ │              │ │   [Go + TS debug]     │ │               │ │ 13-wp-plugin/         │ │
│ └───────────────────────┘ │              │ └───────────────────────┘ │               │ │   link-manager/       │ │
│                           │              │                           │               │ │   [PHP + Go debug]    │ │
│ ┌───────────────────────┐ │              │ ┌───────────────────────┐ │               │ └───────────────────────┘ │
│ │ 06-powershell         │◄├──────────────┤►│ 11-ai-bridge-cli      │ │               │                           │
│ │   integration-v2/     │ │              │ │   [Go + TS debug]     │ │               │ ┌───────────────────────┐ │
│ └───────────────────────┘ │              │ └───────────────────────┘ │               │ │ 13-wp-plugin/         │ │
│                           │              │                           │               │ │   wp-plugin-publish/  │ │
│ ┌───────────────────────┐ │              │ ┌───────────────────────┐ │               │ │   [PHP debug]         │ │
│ │ 07-error-code         │◄├──────────────┤►│ 12-nexus-flow-cli     │ │               │ └───────────────────────┘ │
│ │   registry/           │ │              │ │   [Go + TS debug]     │ │               │                           │
│ └───────────────────────┘ │              │ └───────────────────────┘ │               │ ┌───────────────────────┐ │
│                           │              │                           │               │ │ 14-wp-plugin-builder/ │ │
└───────────────────────────┘              │ ┌───────────────────────┐ │               │ │   [PHP + Go debug]    │ │
                                           │ │ 15-spec-reverse-cli   │ │               │ └───────────────────────┘ │
                                           │ │   [Go + TS debug]     │ │               │                           │
                                           │ └───────────────────────┘ │               └───────────────────────────┘
                                           │                           │
                                            │ ┌───────────────────────┐ │
                                            │ │ 21-wp-seo-publish-cli │ │
                                            │ │   [Go + TS debug]     │ │
                                            │ └───────────────────────┘ │
                                           │                           │
                                           │ ┌───────────────────────┐ │
                                           │ │ 16-ai-transcribe-cli  │ │
                                           │ │   [Go + TS debug]     │ │
                                           │ └───────────────────────┘ │
                                           │                           │
                                           └───────────────────────────┘

            ┌───────────────────────────┐
            │    SUPPORTING SPECS       │
            │                           │
            │ ┌───────────────────────┐ │
            │ │ 01-general-spec/      │ │        Legend:
            │ │   [Reference only]    │ │        ────────
            │ └───────────────────────┘ │        ◄► = Bidirectional reference
            │                           │        → = Outbound reference from Error Resolution
            │ ┌───────────────────────┐ │        [Go debug] = References 02-debugging-go.md
            │ │ 03-shared-cli         │◄├────────[PHP debug] = References 01-debugging-php.md
            │ │   frontend/           │ │        [TS debug] = References 03-debugging-typescript.md
            │ └───────────────────────┘ │
            │                           │
            └───────────────────────────┘
```

---

## Connection Matrix

| Source Spec | Debugging Guide | Verification Pattern | Status |
|-------------|-----------------|---------------------|--------|
| **Architecture Specs** | | | |
| 05-split-db-architecture | Go | ✓ | ✅ Connected |
| 06-seedable-config-architecture | Go | ✓ | ✅ Connected |
| 06-powershell-integration | — | ✓ | ✅ Connected |
| 07-error-code-registry | PHP + Go | ✓ | ✅ Connected |
| **CLI Tools** | | | |
| 09-gsearch-cli | Go + TS | ✓ | ✅ Connected |
| 10-brun-cli | Go + TS | ✓ | ✅ Connected |
| 11-ai-bridge-cli | Go + TS | ✓ | ✅ Connected |
| 12-nexus-flow-cli | Go + TS | ✓ | ✅ Connected |
| 15-spec-reverse-cli | Go + TS | ✓ | ✅ Connected |
| 21-wp-seo-publish-cli | Go + TS | ✓ | ✅ Connected |
| 16-ai-transcribe-cli | Go + TS | ✓ | ✅ Connected |
| **WordPress Specs** | | | |
| 13-wp-plugin/exam-manager | PHP + Go | ✓ | ✅ Connected |
| 13-wp-plugin/link-manager | PHP + Go | ✓ | ✅ Connected |
| 13-wp-plugin/wp-plugin-publish | PHP | ✓ | ✅ Connected |
| 14-wp-plugin-builder | PHP + Go | ✓ | ✅ Connected |
| **Supporting Specs** | | | |
| 01-general-spec | — | — | ⚪ Reference only |
| 02-spec-management-software | — | — | ⚪ Reference only |
| 20-shared-cli-frontend | TS | ✓ | ✅ Connected |

---

## Coverage Summary

```
┌────────────────────────────────────────┐
│         COVERAGE STATISTICS            │
├────────────────────────────────────────┤
│                                        │
│  Total Specs in Project:    18         │
│  Connected to Error Resolution: 16     │
│  Reference Only (no debug):  2         │
│                                        │
│  Coverage: 100%                        │
│                                        │
├────────────────────────────────────────┤
│  Debugging Guides Used:                │
│    - PHP:        5 specs               │
│    - Go:        14 specs               │
│    - TypeScript: 8 specs               │
│                                        │
└────────────────────────────────────────┘
```

---

## Specs Not Requiring References

These specs are **reference-only** and don't require direct error resolution links:

| Spec | Reason |
|------|--------|
| 01-general-spec | High-level standards library (contains own error patterns) |
| 02-spec-management-software | Application spec (references general-spec for errors) |
| 32-shared-preset-data | Data only, no code |

---

*This diagram provides a visual overview of all error resolution cross-references.*
