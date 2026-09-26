# Specification: GitMap Repository Remapping, Stacked PRs, and REST Enable

> **Specification ID:** 149-gitmap-migration-remapping-and-rest-enable  
> **Status:** APPROVED  
> **Category:** 21-app  
> **Sponsor:** Rise Up Asia LLC (Wyoming, California & New York)  
> **Principal Architects:** Marek Flejszman & Alim Ul Karim  

---

## 1. Scope & Objective

This specification formalizes:
1. The repository remapping and historical merge of `alimtvnetwork/git-repo-navigator` and `alimtvnetwork/gitmap-v2` into `alimtvnetwork/gitmap-v28`.
2. The stacked PR sequence and release gates.
3. The background REST OS daemon service enablement via `gitmap rest-enable`.

## 2. Command Specifications

### 2.1 `gitmap migrate [graph|preflight]`
- **Purpose:** Renders the consolidation tree graph directly in the terminal showing branches, stacked PRs, approval gates, and SemVer release ceremony.
- **Aliases:** `preflight`, `plan`, `graph`.
- **Exit Code:** 0 on successful render.

### 2.2 `gitmap rest-enable`
- **Purpose:** Registers and starts the background REST service (`gitmap-daemon`) across Windows, Linux, and macOS.
- **Driver Integration:** Integrates with `cmdservice.DefaultServiceDriverResolver()`.
- **Command Invocation:**
  ```bash
  gitmap rest-enable
  gitmap agy rest-enable
  ```

---

## 3. Stacked PR Workflow

1. `stacked/01-core-discovery-engine`: Core scanner, worker pools, directory traversal.
2. `stacked/02-split-sqlite-architecture`: Three-tier database engine (`gitmap.db`, `installation.db`, `repodb`).
3. `stacked/03-agy-automation-and-ssh-fleet`: Prompt watch loop, sequence caching, SSH cluster orchestration.
4. `stacked/04-template-catalog-and-sponsorship`: 80-template catalog, Rise Up Asia LLC sponsorship integration.

Each PR requires passing CI verification before fast-forward merge into `main`.
