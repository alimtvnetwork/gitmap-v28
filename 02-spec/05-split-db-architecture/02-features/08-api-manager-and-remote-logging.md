# API Manager & Remote Logging Specification

## 1. Overview

This specification outlines the remote logging sink architecture (`ApiSink` / `ApiManager`) that serializes structured log events and transmits them over HTTP to centralized log aggregators or microservice endpoints.

## 2. Architecture & Contracts

### 2.1 Pluggable Sender Function

```go
type ApiLogSenderFunc func(ctx context.Context, endpoint string, headers map[string]string, payload []byte) *appfault.AppError
```
Allows swapping the underlying network transport (e.g., standard HTTP, gRPC bridge, or in-memory mock during tests) without modifying the sink core.

### 2.2 Configuration & Lifecycle

- `Endpoint`: Target URL for log submission.
- `Headers`: Custom HTTP headers (e.g. Authorization, Tenant-ID).
- `BatchSize`: Maximum entries accumulated before automatic flush.
- `FlushInterval`: Time threshold triggering periodic buffer flushing.
- `IsCompress`: Optional gzip compression for network bandwidth reduction.

### 2.3 Error Handling

All transmission failures are returned as `*appfault.AppError` with category `errtype.Network` or `errtype.External`.
