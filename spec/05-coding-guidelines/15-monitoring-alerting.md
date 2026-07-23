# Monitoring & Alerting

## Overview

Standards for health checks, metrics collection, log aggregation, and
incident response patterns. Monitoring exists to detect problems before
users do — make it reliable, actionable, and low-noise.

---

## 1. Health Checks

### Endpoint Convention

Every service exposes a `/healthz` endpoint returning structured status.

```json
{
  "status": "ok",
  "checks": {
    "database": { "status": "ok", "latency_ms": 12 },
    "storage":  { "status": "ok", "latency_ms": 45 },
    "network":  { "status": "degraded", "detail": "high latency" }
  },
  "version": "2.36.7",
  "uptime_s": 84200
}
```

### Rules

| Rule | Rationale |
|------|-----------|
| Health checks are read-only | Never mutate state during a check |
| Include dependency status | Database, network, storage — each separate |
| Return HTTP 200 for healthy, 503 for degraded | Load balancers rely on status codes |
| Include version and uptime | Aids debugging without extra queries |
| Timeout individual checks at 3 seconds | Prevent cascading slowdowns |

### CLI Health Checks (Doctor Pattern)

CLI tools use the `doctor` command pattern:

```
✓ Config file       OK — parsed 6 fields
✓ Database          OK — 42 repos, 12 releases
✗ Git binary        FAIL — git not found in PATH
```

- One line per check — icon, name, status, detail.
- Exit code reflects overall health (0 = pass, 1 = any fail).
- See `spec/05-coding-guidelines/07-logging-observability.md` §5.

---

## 2. Metrics Collection

### What to Measure

| Category | Metrics | Example |
|----------|---------|---------|
| Latency | p50, p95, p99 response times | `http_request_duration_ms{p99=240}` |
| Traffic | Request rate, throughput | `http_requests_total{method="GET"}` |
| Errors | Error rate, error ratio | `http_errors_total{status="500"}` |
| Saturation | CPU, memory, disk, connections | `system_memory_used_percent{host="api-1"}` |

### Naming Convention

```
<namespace>_<subsystem>_<metric>_<unit>
```

Examples:
- `gitmap_scan_duration_seconds`
- `gitmap_clone_errors_total`
- `gitmap_db_connections_active`

### Rules

- Use counters for cumulative values (requests, errors).
- Use gauges for current-state values (connections, queue depth).
- Use histograms for latency distributions.
- Label with dimensions, not embedded in metric names.
- Never use high-cardinality labels (user IDs, request IDs).

---

## 3. Log Aggregation

### Structured Log Format

All logs must be structured JSON in production:

```json
{
  "timestamp": "2025-01-15T14:30:22.001Z",
  "level": "error",
  "component": "clone",
  "message": "git clone failed",
  "repo": "api-gateway",
  "error": "authentication required",
  "duration_ms": 3400,
  "trace_id": "abc123"
}
```

### Required Fields

| Field | When |
|-------|------|
| `timestamp` | Always — ISO 8601 with milliseconds |
| `level` | Always — error, warn, info, debug |
| `component` | Always — subsystem identifier |
| `message` | Always — human-readable summary |
| `trace_id` | On request-scoped operations |
| `error` | On failures — original error string |
| `duration_ms` | On timed operations |

### Aggregation Rules

- Ship logs to a central store (not local files in production).
- Retain error/warn logs for 90 days minimum.
- Retain info logs for 30 days.
- Retain debug logs for 7 days (if enabled).
- Never log secrets, tokens, or PII.
- See `spec/05-coding-guidelines/07-logging-observability.md` for format details.

---

## 4. Alerting Strategy

### Severity Levels

| Severity | Response Time | Example |
|----------|--------------|---------|
| Critical | ≤15 minutes | Service down, data loss risk |
| High | ≤1 hour | Error rate >5%, degraded performance |
| Medium | Next business day | Elevated latency, disk >80% |
| Low | Next sprint | Warning trends, non-critical deprecations |

### Alert Rules

- Alert on **symptoms** (error rate, latency), not causes (CPU, memory).
- Every alert must have a **runbook** link.
- Suppress duplicate alerts — deduplicate within a 5-minute window.
- Require **two consecutive failures** before firing (avoid flaps).
- Include context in the alert: service, metric value, threshold, dashboard link.

### Alert Format

```
[CRITICAL] gitmap-api: error rate 12.3% (threshold: 5%)
Dashboard: https://monitoring.example.com/d/gitmap-v28
Runbook: https://wiki.example.com/runbooks/high-error-rate
```

---

## 5. Incident Response

### Severity Classification

| Level | Criteria | Response |
|-------|----------|----------|
| SEV-1 | Service outage, data loss | All-hands, war room |
| SEV-2 | Major degradation | On-call + backup |
| SEV-3 | Minor degradation | On-call investigates |
| SEV-4 | Cosmetic or low-impact | Normal prioritization |

### Incident Lifecycle

1. **Detect** — Alert fires or user reports.
2. **Acknowledge** — On-call acknowledges within SLA.
3. **Triage** — Classify severity, assemble responders.
4. **Mitigate** — Restore service (rollback, scale, failover).
5. **Resolve** — Root cause fixed and deployed.
6. **Post-mortem** — Written within 48 hours of resolution.

### Post-Mortem Template

Every SEV-1 and SEV-2 incident requires a post-mortem:

- **Summary**: One-paragraph description.
- **Timeline**: Timestamped sequence of events.
- **Root Cause**: Technical explanation.
- **Impact**: Users affected, duration, data loss.
- **Action Items**: Concrete tasks with owners and deadlines.
- **Lessons Learned**: What worked, what didn't.

---

## 6. Dashboard Standards

- Every service has a dedicated dashboard.
- Dashboards include: error rate, latency percentiles, throughput, saturation.
- Use consistent time ranges (1h, 6h, 24h, 7d).
- Mark deployment events on timeline graphs.
- Dashboard links are included in alert notifications.

---

## Constraints

- Health check endpoints respond in <500ms.
- Metric names follow the `namespace_subsystem_metric_unit` convention.
- All production logs are structured JSON.
- Every alert has a runbook.
- Post-mortems are blameless and completed within 48 hours.
- No high-cardinality labels in metrics.

---

## References

- Logging Standards: `spec/05-coding-guidelines/07-logging-observability.md`
- Error Handling: `spec/05-coding-guidelines/04-error-handling.md`
- Doctor Pattern: `spec/01-app/18-compliance-audit.md`
