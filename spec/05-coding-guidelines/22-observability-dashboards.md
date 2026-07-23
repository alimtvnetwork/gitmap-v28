# Observability Dashboards

Universal guidelines for dashboard layout, metric visualization, alert integration, and SLO tracking across all projects.

---

## 1. Dashboard Layout

### Structure

Every observability dashboard follows a consistent top-down layout:

```
┌─────────────────────────────────────────┐
│  Title Bar  (service name + time range) │
├─────────────────────────────────────────┤
│  SLO Summary  (error budget + status)   │
├──────────────────┬──────────────────────┤
│  Key Metrics     │  Active Alerts       │
│  (golden signals)│  (firing + recent)   │
├──────────────────┴──────────────────────┤
│  Detailed Panels  (request rate, errors,│
│  latency percentiles, saturation)       │
├─────────────────────────────────────────┤
│  Drill-Down Links  (logs, traces, deps) │
└─────────────────────────────────────────┘
```

### Rules

| Rule | Detail |
|---|---|
| Top-down priority | Most critical information at the top |
| Golden signals first | Rate, errors, latency, saturation — always visible |
| Consistent time range | All panels use the same time window |
| No decoration | Remove chart borders, gridlines, and legends that don't add clarity |
| Mobile-aware | Key metrics must be readable on narrow viewports |

### Naming

| Convention | Example |
|---|---|
| Dashboard title | `<Service> — Overview` |
| Panel title | `<Metric> — <Unit>` (e.g., `Request Rate — req/s`) |
| Variable names | `$service`, `$environment`, `$region` |

---

## 2. Metric Visualization

### Chart Type Selection

| Data Shape | Chart Type | Example |
|---|---|---|
| Time series (rate) | Line chart | Request rate over time |
| Time series (count) | Stacked area | Error count by type |
| Current value | Stat panel / gauge | Current CPU usage |
| Distribution | Histogram / heatmap | Latency distribution |
| Comparison | Bar chart | Error rate by endpoint |
| Status | Traffic light / badge | Service health (up/down) |

### Golden Signals

Every service dashboard must include these four panels:

```
1. Rate      — requests per second (line chart)
2. Errors    — error rate as percentage (line + threshold)
3. Latency   — p50, p95, p99 percentiles (multi-line)
4. Saturation — CPU, memory, connections (gauges)
```

### Visualization Rules

| Rule | Detail |
|---|---|
| Fixed Y-axis minimum | Always start Y-axis at 0 for rate/count charts |
| Percentile lines | Show p50, p95, p99 — never just "average" |
| Color consistency | Green = healthy, yellow = warning, red = critical |
| Threshold lines | Draw horizontal lines at SLO/alert thresholds |
| Unit labels | Always include units (ms, req/s, %, bytes) |
| No pie charts | Use bar charts or stat panels instead |

### Go (Prometheus-Style Metrics)

```go
// metrics.go
var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests by method and status",
        },
        []string{"method", "status", "endpoint"},
    )

    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request latency in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)
```

### TypeScript (Metric Reporting)

```typescript
interface MetricPoint {
  name: string;
  value: number;
  labels: Record<string, string>;
  timestamp: number;
}

function recordMetric(point: MetricPoint): void {
  // Send to metrics backend
  metricsClient.record(point);
}

recordMetric({
  name: "http_request_duration_ms",
  value: elapsed,
  labels: { method: "GET", endpoint: "/api/users", status: "200" },
  timestamp: Date.now(),
});
```

---

## 3. Alert Integration

### Alert Severity Levels

| Level | Response Time | Example |
|---|---|---|
| Critical (P1) | Immediate (< 5 min) | Service down, data loss |
| High (P2) | Within 30 min | Error rate > SLO, high latency |
| Warning (P3) | Within 4 hours | Disk usage > 80%, elevated errors |
| Info (P4) | Next business day | Deployment completed, cert renewal |

### Alert Rules

| Rule | Detail |
|---|---|
| Alert on symptoms, not causes | "Error rate > 1%" not "CPU > 90%" |
| Include runbook link | Every alert links to a resolution guide |
| Set meaningful thresholds | Base on SLOs and historical baselines |
| Avoid alert fatigue | < 5 actionable alerts per service per week |
| Use alert grouping | Group related alerts into a single notification |
| Require ownership | Every alert has an assigned team/owner |

### Alert Template

```yaml
alert: HighErrorRate
expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.01
for: 5m
labels:
  severity: high
  team: backend
annotations:
  summary: "Error rate above 1% for {{ $labels.service }}"
  runbook: "https://runbooks.internal/high-error-rate"
  dashboard: "https://dashboards.internal/{{ $labels.service }}"
```

### Dashboard–Alert Linking

- Every alert panel on the dashboard shows current firing alerts.
- Clicking an alert navigates to the corresponding detailed panel.
- Alert threshold lines are drawn on the relevant metric chart.
- Alert history panel shows recent firings with resolution times.

---

## 4. SLO Tracking

### Definitions

| Term | Definition |
|---|---|
| SLI (Service Level Indicator) | A measurable metric (e.g., request success rate) |
| SLO (Service Level Objective) | Target value for an SLI (e.g., 99.9% success rate) |
| Error Budget | Allowed failure = `1 - SLO` (e.g., 0.1% of requests can fail) |
| Burn Rate | How fast the error budget is being consumed |

### SLO Dashboard Panel

Every service dashboard includes an SLO summary at the top:

```
┌──────────────────────────────────────────┐
│  SLO: 99.9%  │  Current: 99.95%  │  ✅  │
│  Budget: 43.2 min remaining (30d window) │
│  Burn Rate: 0.5x (healthy)              │
└──────────────────────────────────────────┘
```

### Burn Rate Alerts

| Burn Rate | Window | Severity | Meaning |
|---|---|---|---|
| 14.4x | 1 hour | Critical | Budget exhausted in ~1 hour |
| 6x | 6 hours | High | Budget exhausted in ~5 hours |
| 1x | 3 days | Warning | On track to exhaust budget |

### Implementation

```go
// slo.go
type SLO struct {
    Name       string
    Target     float64 // e.g., 0.999
    Window     time.Duration
    Indicator  string // metric name
}

func (s SLO) ErrorBudget() float64 {
    return 1.0 - s.Target
}

func (s SLO) BurnRate(currentErrorRate float64) float64 {
    budget := s.ErrorBudget()
    if budget == 0 {
        return 0
    }
    return currentErrorRate / budget
}
```

```typescript
interface SLO {
  name: string;
  target: number; // e.g., 0.999
  windowDays: number;
  indicator: string;
}

function errorBudget(slo: SLO): number {
  return 1.0 - slo.target;
}

function burnRate(slo: SLO, currentErrorRate: number): number {
  const budget = errorBudget(slo);
  return budget === 0 ? 0 : currentErrorRate / budget;
}

function remainingBudgetMinutes(slo: SLO, currentErrorRate: number): number {
  const totalMinutes = slo.windowDays * 24 * 60;
  const budgetMinutes = totalMinutes * errorBudget(slo);
  const consumedMinutes = totalMinutes * currentErrorRate;
  return Math.max(0, budgetMinutes - consumedMinutes);
}
```

---

## 5. Dashboard Standards

### Required Dashboards

| Dashboard | Content |
|---|---|
| Service overview | Golden signals + SLO summary + active alerts |
| Infrastructure | CPU, memory, disk, network per host/container |
| Dependency map | Upstream/downstream health and latency |
| Deployment | Deploy frequency, rollback rate, lead time |

### Maintenance Rules

| Rule | Detail |
|---|---|
| Review quarterly | Remove unused panels, update thresholds |
| Version control | Dashboard definitions stored as code (JSON/YAML) |
| No manual dashboards | All dashboards reproducible from config |
| Template variables | Use `$service`, `$env` — never hardcode values |
| Annotations | Mark deployments and incidents on timeline charts |

### Offline / Self-Contained Dashboards

For CLI-generated dashboards (see [Dashboard Command](../01-app/60-help-dashboard.md)):

| Rule | Detail |
|---|---|
| No CDN dependencies | Bundle all CSS/JS inline |
| Canvas-based charts | Use HTML5 `<canvas>` for visualizations |
| Static data snapshot | Embed JSON data directly in the HTML |
| Single-file output | One `.html` file, fully self-contained |

---

## Constraints

| Constraint | Detail |
|---|---|
| No pie charts | Use bar charts or stat panels |
| No average-only latency | Always show percentiles (p50, p95, p99) |
| No alert without runbook | Every alert must link to resolution steps |
| No dashboard without SLO | Every service dashboard includes SLO tracking |
| No hardcoded thresholds | Use variables or config for alert thresholds |
| No unowned alerts | Every alert has an assigned team |

---

## References

- [Monitoring & Alerting](./15-monitoring-alerting.md)
- [Logging & Observability](./07-logging-observability.md)
- [Resilience Patterns](./18-resilience-patterns.md)
- [CI/CD Patterns](./17-cicd-patterns.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
