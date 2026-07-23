# Monitoring & Alerting

Universal guidelines for metric collection, alert thresholds, dashboard design, and incident response across all projects.

---

## 1. Metric Collection

### Golden Signals

Every service must expose these four categories:

| Signal | What to Measure | Example Metric |
|---|---|---|
| Traffic | Request rate | `http_requests_total` |
| Errors | Failure rate / ratio | `http_errors_total`, `error_rate_percent` |
| Latency | Response time distribution | `http_request_duration_seconds` |
| Saturation | Resource utilization | `cpu_usage_percent`, `memory_usage_bytes` |

### Metric Types

| Type | Use Case | Example |
|---|---|---|
| Counter | Monotonically increasing totals | Requests served, errors occurred |
| Gauge | Current point-in-time value | Active connections, queue depth |
| Histogram | Value distributions (percentiles) | Request latency, payload size |
| Summary | Pre-calculated percentiles | Client-side latency (when server histograms are unavailable) |

### Naming Conventions

| Rule | Example |
|---|---|
| Snake_case | `http_requests_total` |
| Include unit suffix | `_seconds`, `_bytes`, `_total` |
| Prefix with subsystem | `db_query_duration_seconds` |
| Counters end with `_total` | `http_requests_total` |
| No metric name collisions | Unique per service namespace |

### Go (Prometheus)

```go
var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests by method and status.",
        },
        []string{"method", "status"},
    )

    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Request latency in seconds.",
            Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(requestsTotal, requestDuration)
}
```

### TypeScript

```typescript
interface Metric {
  name: string;
  type: "counter" | "gauge" | "histogram";
  help: string;
  labels: string[];
}

const metrics = {
  requestsTotal: createCounter({
    name: "http_requests_total",
    help: "Total HTTP requests",
    labels: ["method", "status"],
  }),

  requestDuration: createHistogram({
    name: "http_request_duration_seconds",
    help: "Request latency in seconds",
    labels: ["method", "endpoint"],
    buckets: [0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10],
  }),
};
```

### Collection Rules

| Rule | Detail |
|---|---|
| Instrument at boundaries | HTTP handlers, DB calls, external API calls |
| Use labels sparingly | High-cardinality labels (user IDs, URLs) cause metric explosion |
| Max 5–7 labels per metric | Beyond this, use logs for detail |
| Collect at source | Services emit their own metrics — no external scraping of logs |
| Standard endpoints | Expose `/metrics` (Prometheus) or push to collector |

---

## 2. Alert Thresholds

### Severity Levels

| Level | Response Time | Notification | Example |
|---|---|---|---|
| Critical (P1) | Immediate (< 5 min) | Page on-call | Service down, data corruption |
| High (P2) | Within 30 min | Page + ticket | Error rate > SLO, degraded latency |
| Warning (P3) | Within 4 hours | Ticket only | Disk > 80%, elevated error rate |
| Info (P4) | Next business day | Dashboard only | Deploy completed, cert renewal due |

### Threshold Rules

| Rule | Detail |
|---|---|
| Alert on symptoms | "Error rate > 1%" not "CPU > 90%" |
| Use rate over raw count | `rate(errors[5m]) > 0.01` not `errors > 100` |
| Require sustained duration | `for: 5m` — avoid firing on transient spikes |
| Base on SLOs | Thresholds derive from service-level objectives |
| Review quarterly | Adjust based on traffic growth and historical data |

### Alert Definition

```yaml
# alert_rules.yml
groups:
  - name: service_health
    rules:
      - alert: HighErrorRate
        expr: >
          rate(http_requests_total{status=~"5.."}[5m])
          / rate(http_requests_total[5m]) > 0.01
        for: 5m
        labels:
          severity: high
          team: backend
        annotations:
          summary: "Error rate above 1% for {{ $labels.service }}"
          runbook: "https://runbooks.internal/high-error-rate"
          dashboard: "https://dashboards.internal/{{ $labels.service }}"

      - alert: HighLatency
        expr: >
          histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1.0
        for: 10m
        labels:
          severity: high
          team: backend
        annotations:
          summary: "p95 latency above 1s for {{ $labels.service }}"
          runbook: "https://runbooks.internal/high-latency"
```

### Anti-Patterns

| Anti-Pattern | Problem | Fix |
|---|---|---|
| Alerting on CPU alone | Symptom-blind, causes false positives | Alert on error rate or latency instead |
| No `for` duration | Fires on momentary spikes | Add `for: 5m` minimum |
| Too many alerts | Alert fatigue, missed real issues | < 5 actionable alerts per service per week |
| Missing runbook | On-call has no guidance | Every alert must link to a runbook |
| Hardcoded thresholds | Break when traffic scales | Use SLO-derived or percentile-based thresholds |

---

## 3. Dashboard Design

### Layout Hierarchy

```
Title + Time Range Selector
├── SLO Status (error budget, burn rate)
├── Golden Signals (rate, errors, latency, saturation)
├── Active Alerts (currently firing)
├── Detailed Panels (per-endpoint, per-dependency)
└── Drill-Down Links (logs, traces)
```

### Panel Rules

| Rule | Detail |
|---|---|
| Critical info at top | SLO and golden signals are always first |
| Consistent time range | All panels share one time window |
| Y-axis starts at 0 | For rate and count charts — never auto-scale from minimum |
| Show percentiles | p50, p95, p99 — never "average" alone |
| Draw threshold lines | SLO and alert thresholds visible on charts |
| Use color meaningfully | Green = healthy, yellow = warning, red = critical |

### Chart Selection

| Data | Chart Type |
|---|---|
| Rate over time | Line chart |
| Error breakdown | Stacked area |
| Current value | Stat panel or gauge |
| Latency distribution | Heatmap or histogram |
| Comparison across services | Bar chart |
| Up/down status | Traffic light badge |

### Dashboard Standards

| Standard | Detail |
|---|---|
| Dashboard-as-code | JSON/YAML definitions in version control |
| Template variables | `$service`, `$environment`, `$region` — never hardcoded |
| Annotations | Mark deployments and incidents on timelines |
| No manual dashboards | All reproducible from config |
| Review quarterly | Remove stale panels, update thresholds |

---

## 4. Incident Response

### Severity Classification

| Severity | Impact | Example | Response |
|---|---|---|---|
| SEV-1 | Full outage | Service unreachable | All-hands, status page update |
| SEV-2 | Major degradation | > 50% error rate | On-call + team lead |
| SEV-3 | Minor degradation | Elevated latency, partial errors | On-call investigates |
| SEV-4 | No user impact | Internal metric anomaly | Track in ticket |

### Response Workflow

```
Alert fires → Acknowledge (< 5 min) → Triage severity → Investigate
→ Mitigate → Communicate → Resolve → Post-mortem (within 48 hours)
```

### On-Call Rules

| Rule | Detail |
|---|---|
| Acknowledge within 5 min | Auto-escalate if unacknowledged |
| Mitigate before root-cause | Restore service first, investigate later |
| Communicate early | Status page update within 15 min for SEV-1/2 |
| Escalate when stuck | If no progress in 30 min, bring in next responder |
| No blame | Focus on systems and processes, not individuals |

### Post-Mortem Template

```markdown
# Incident: [Title]

## Summary
- **Duration**: Start time – End time (total minutes)
- **Severity**: SEV-X
- **Impact**: What users experienced
- **Detection**: How the incident was detected (alert name, user report)

## Timeline
| Time | Event |
|------|-------|
| HH:MM | Alert fired |
| HH:MM | Acknowledged by [responder] |
| HH:MM | Root cause identified |
| HH:MM | Mitigation applied |
| HH:MM | Fully resolved |

## Root Cause
[Description of what caused the incident]

## Action Items
| Action | Owner | Due Date |
|--------|-------|----------|
| [Fix] | [Name] | [Date] |
| [Prevention] | [Name] | [Date] |
| [Monitoring improvement] | [Name] | [Date] |
```

### Runbook Standards

Every alert must link to a runbook containing:

| Section | Content |
|---|---|
| Alert meaning | What triggered and why it matters |
| Impact assessment | How to determine user impact |
| Diagnostic steps | Commands/queries to run (copy-pasteable) |
| Mitigation steps | How to restore service quickly |
| Escalation path | Who to contact if stuck |
| Related dashboards | Links to relevant monitoring views |

---

## Constraints

| Constraint | Detail |
|---|---|
| No alert without runbook | Every alert links to resolution steps |
| No average-only metrics | Always include percentiles for latency |
| No high-cardinality labels | Max 5–7 labels per metric |
| No manual dashboards | All dashboard definitions in version control |
| No unacknowledged alerts | Auto-escalate after 5 min |
| No blameful post-mortems | Focus on systems, not people |
| Post-mortem within 48 hours | For SEV-1 and SEV-2 incidents |

---

## References

- [Observability Dashboards](./22-observability-dashboards.md)
- [Logging & Observability](./07-logging-observability.md)
- [Resilience Patterns](./18-resilience-patterns.md)
- [Security & Secrets](./08-security-secrets.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
