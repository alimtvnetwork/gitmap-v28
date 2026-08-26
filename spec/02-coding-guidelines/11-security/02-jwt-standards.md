# JWT Standards

## 1. Token Lifecycles

- **Access Tokens:** Must have a short expiration time (e.g., 15 minutes).
- **Refresh Tokens:** Used to obtain new access tokens. Must be rotated upon use to detect and prevent replay attacks.

## 2. Storage

- **Access Tokens:** Should be stored in memory or in a short-lived closure. Never store access tokens in `localStorage`.
- **Refresh Tokens:** MUST be stored in `HttpOnly`, `Secure`, `SameSite=Strict` cookies. They should never be accessible via JavaScript.

## 3. Payload Constraints

- Never include PII (Personally Identifiable Information) or sensitive secrets in the JWT payload. The payload is Base64 encoded and can be read by anyone. Include only necessary identifiers like `UserId` and role claims.
