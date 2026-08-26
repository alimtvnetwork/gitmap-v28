# [PROJECT_NAME] - Error Codes

> **Project Prefix:** `XX`  
> **Range:** X000-X999  
> **Last Updated:** YYYY-MM-DD
**Version:** 3.2.0  
**Last Updated:** 2026-04-16  

---

## Overview

This document defines all error codes for [PROJECT_NAME].

---

## XX-000: Initialization (X000-X099)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-000-01 | | | |
| XX-000-02 | | | |

---

## XX-100: Authentication (X100-X199)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-100-01 | | | 401 |
| XX-100-02 | | | 401 |

---

## XX-200: Authorization (X200-X299)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-200-01 | | | 403 |
| XX-200-02 | | | 403 |

---

## XX-300: Validation (X300-X399)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-300-01 | | | 400 |
| XX-300-02 | | | 400 |

---

## XX-400: Business Logic (X400-X499)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-400-01 | | | 422 |
| XX-400-02 | | | 422 |

---

## XX-500: Database (X500-X599)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-500-01 | | | 500 |
| XX-500-02 | | | 500 |

---

## XX-600: External Services (X600-X699)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-600-01 | | | 502 |
| XX-600-02 | | | 503 |

---

## XX-700: File System (X700-X799)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-700-01 | | | 500 |
| XX-700-02 | | | 404 |

---

## XX-800: Network (X800-X899)

| Code | Name | Message | HTTP |
|------|------|---------|------|
| XX-800-01 | | | 503 |
| XX-800-02 | | | 504 |

---

## Usage

```go
// Go
return errors.New(ErrXX00001, "Description")
```

```typescript
// TypeScript
throw new AppError(ErrorCodes.XX_000_01, "Description");
```

---

## Registration

This error code range is registered in:
`spec/03-error-manage/03-error-code-registry/01-registry.md`
