# Blue-Green Deployment Standards

## 1. Zero Downtime Requirement

All critical production services MUST be deployed using a Blue-Green or Canary deployment strategy. In-place updates that cause service interruption are strictly prohibited.

## 2. Infrastructure Setup

- **Blue Environment:** The current active production environment.
- **Green Environment:** The idle environment where the new version is deployed.
- **Router/Load Balancer:** Must support instantaneous traffic switching from Blue to Green.

## 3. Deployment Flow

1. Deploy the new artifact to the Green environment.
2. Run automated smoke tests and health checks against the Green environment.
3. If tests pass, switch the router to send 100% of traffic to Green (or gradually shift traffic if doing Canary).
4. Monitor Error Budgets and SLIs for anomalies.
5. If anomalies are detected within the cool-down period, instantly rollback by switching traffic back to Blue.
6. Once stable, tear down the Blue environment to save costs (it becomes the new Green for the next deployment).
