# Error Management — Acceptance Criteria

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## AC-01: Structured Error Response

**GIVEN** any CLI backend encounters an error during request processing  
**WHEN** the error response is generated  
**THEN** it contains: `Code` (numeric), `Message` (human-readable), `Details` (technical), and `Stack` (up to 40 frames)  
**AND** the error code falls within the tool's assigned range per the Error Code Registry

**Edge Cases:**
- **GIVEN** the error originates from a third-party library **WHEN** the stack trace is captured **THEN** both the library frames and the application frames are included with clear delineation
- **GIVEN** the error code is not registered in the Error Code Registry **WHEN** it is returned **THEN** a fallback generic code within the tool's range is used and a warning is logged
- **GIVEN** the error `Details` field contains sensitive data (file paths, credentials) **WHEN** the response is generated **THEN** sensitive values are redacted before sending to the client

---

## AC-02: Frontend-Backend Verification Protocol

**GIVEN** the frontend receives an error response from the backend  
**WHEN** the error is displayed in the error modal  
**THEN** the user can see the backend error code, the frontend component that triggered the request, and the timestamp  
**AND** "Copy All" copies both frontend and backend context

**Edge Cases:**
- **GIVEN** the clipboard API is unavailable **WHEN** "Copy All" is clicked **THEN** a fallback textarea is shown with the content pre-selected
- **GIVEN** the backend returns an error with no `Code` field **WHEN** the frontend processes it **THEN** a synthetic code `GEN-1000` is assigned and a parsing warning is logged

---

## AC-03: Retrospective Document Structure

**GIVEN** a production bug has been resolved  
**WHEN** a retrospective document is created  
**THEN** it contains: Root Cause, Timeline, Resolution Steps, Prevention Measures, and Related Error Codes

---

## AC-04: Verification Pattern Application

**GIVEN** a developer implements a fix for a known error pattern  
**WHEN** they consult the verification patterns documentation  
**THEN** they find step-by-step verification instructions specific to the error category

---

## AC-05: Debugging Guide Coverage

**GIVEN** a developer encounters a backend error  
**WHEN** they follow the language-specific debugging guide  
**THEN** they can identify common issues and each issue links to the relevant specification

---

## AC-06: Quick Resolution

**GIVEN** a common error scenario  
**WHEN** the developer consults the cheat sheet  
**THEN** they find a 3-step resolution procedure: Identify → Diagnose → Fix  
**AND** each step includes the exact command or code to run

---

## Cross-References

- [Overview](./00-overview.md)
