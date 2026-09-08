# Subtask 02 — Chrome Refresh Token Vault & Cipher Variations

**Parent Plan:** [.lovable/plans/pending/64-remediation-fix-and-chrome-token-export.md](../../pending/64-remediation-fix-and-chrome-token-export.md)  
**Status:** COMPLETED  
**Files:**
- `gitmap/cmd/chromeprofile_tokens.go` [NEW]
- `gitmap/cmd/chromeprofile_tokens_test.go` [NEW]

---

## Objective
Implement a dedicated Chrome refresh token extraction and cipher transformation module capable of reading `token_service` from `Web Data`, generating double Base64 and Caesar cipher ("chespercypher") variations, and providing full mathematical reversibility.

## Requirements
1. Extract tokens from `<profile>/Web Data` table `token_service` using read-only SQLite with copy-to-temp fallback.
2. Implement `EncodeDoubleBase64` & `DecodeDoubleBase64` (double base64 encoding/decoding).
3. Implement `EncodeCaesarCipher` & `DecodeCaesarCipher` (text-based Caesar cipher with configurable shift).
4. Implement `EncodeCaesarByteShift` & `DecodeCaesarByteShift` (byte-level Caesar shift modulo 256).
5. Build `ChromeTokenVault` struct storing token variations (`doubleBase64`, `caesarCipher`, `caesarByteShift`) and explicit revert instructions.
6. Verify 100% bi-directional decoding equality in unit tests.
