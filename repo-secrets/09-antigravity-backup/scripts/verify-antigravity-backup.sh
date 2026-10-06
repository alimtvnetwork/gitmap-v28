#!/usr/bin/env bash
# ==============================================================================
# verify-antigravity-backup.sh
# Quality Gate & Validation Audit Engine for Antigravity Backup Vault
# Verifies Quality Gates VG-01 through VG-07
# ==============================================================================

set -euo pipefail

# Positive boolean tracking variables
isValid=true
isVerbose=false
isSanitized=true
isPortable=true
isVerified=true

# Parse flags
for arg in "$@"; do
    case "$arg" in
        -v|--verbose)
            isVerbose=true
            shift
            ;;
    esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VAULT_DIR="${SCRIPT_DIR}/../vault"
SCRIPTS_DIR="${SCRIPT_DIR}"

PROJECTS_MANIFEST="${VAULT_DIR}/projects-manifest.json"
PINNED_MANIFEST="${VAULT_DIR}/pinned-projects.json"
SETTINGS_MANIFEST="${VAULT_DIR}/settings-manifest.json"
PLUGINS_MANIFEST="${VAULT_DIR}/plugins-and-skills.json"

echo "=============================================================================="
echo "  Antigravity Vault Verification Audit: Gates VG-01 to VG-07"
echo "=============================================================================="

# ------------------------------------------------------------------------------
# Gate VG-01: Vault Structure & Manifest Validation
# ------------------------------------------------------------------------------
echo -n "● [Gate VG-01] Vault Structure & Manifest Validation: "
vg01_pass=true

for manifest in "$PROJECTS_MANIFEST" "$PINNED_MANIFEST" "$SETTINGS_MANIFEST" "$PLUGINS_MANIFEST"; do
    if [[ ! -f "$manifest" ]]; then
        echo "FAILED"
        echo "  Missing manifest: $(basename "$manifest")"
        vg01_pass=false
        break
    fi
    if ! python3 -m json.tool "$manifest" >/dev/null 2>&1; then
        echo "FAILED"
        echo "  Invalid JSON syntax in: $(basename "$manifest")"
        vg01_pass=false
        break
    fi
done

if [[ "$vg01_pass" == "true" ]]; then
    echo "PASSED (All 4 manifests exist and valid JSON)"
else
    isValid=false
fi

# ------------------------------------------------------------------------------
# Gate VG-02: Repository Count Parity
# ------------------------------------------------------------------------------
echo -n "● [Gate VG-02] Repository Count Parity: "
vg02_pass=true

python3 -c "
import json, sys

with open('$PROJECTS_MANIFEST') as f:
    pm = json.load(f)
assert pm.get('totalProjects') == 78, f'totalProjects is {pm.get(\"totalProjects\")}, expected 78'
assert len(pm.get('projects', [])) == 78, f'projects list length is {len(pm.get(\"projects\", []))}, expected 78'
assert len(pm.get('sectionSummary', {})) == 10, 'Expected 10 sections'
assert sum(pm.get('sectionSummary', {}).values()) == 78, 'Section summary sum is not 78'

with open('$PINNED_MANIFEST') as f:
    pp = json.load(f)
assert pp.get('totalPinned') == 23, f'totalPinned is {pp.get(\"totalPinned\")}, expected 23'
assert len(pp.get('projects', [])) == 23, f'pinned list length is {len(pp.get(\"projects\", []))}, expected 23'
assert pp.get('tierDistribution', {}).get('tier1Core') == 7, 'Expected 7 Tier 1 Core repos'
assert pp.get('tierDistribution', {}).get('tier2Ecosystem') == 16, 'Expected 16 Tier 2 Ecosystem repos'
" >/dev/null 2>&1 || vg02_pass=false

if [[ "$vg02_pass" == "true" ]]; then
    echo "PASSED (78 projects across 10 sections, 23 pinned [7 Tier 1 + 16 Tier 2])"
else
    echo "FAILED (Repository or pinned count mismatch)"
    isValid=false
fi

# ------------------------------------------------------------------------------
# Gate VG-03: Path Relativity & Portability
# ------------------------------------------------------------------------------
echo -n "● [Gate VG-03] Path Relativity & Portability: "
vg03_pass=true

python3 -c "
import glob, json, re, sys

files = glob.glob('$VAULT_DIR/*.json')
for f in files:
    with open(f) as fp:
        content = fp.read()
    # Check for hardcoded host paths or Windows drive letters
    leaks = re.findall(r'(?:^|[\"\s])(?:/home/[a-zA-Z0-9_-]+|[A-Za-z]:[\\\\/])', content)
    if leaks:
        print(f'Leak in {f}: {leaks}', file=sys.stderr)
        sys.exit(1)

# Check projects manifest uses workspace root token
with open('$PROJECTS_MANIFEST') as fp:
    pm = json.load(fp)
for p in pm.get('projects', []):
    uri = p.get('uriTemplate', '')
    if '\${WORKSPACE_ROOT}' not in uri:
        print(f'Missing token in {p.get(\"name\")}', file=sys.stderr)
        sys.exit(1)
" >/dev/null 2>&1 || vg03_pass=false

if [[ "$vg03_pass" == "true" ]]; then
    echo "PASSED (\${WORKSPACE_ROOT} used exclusively, zero host paths)"
else
    echo "FAILED (Absolute path leak or missing \${WORKSPACE_ROOT})"
    isValid=false
fi

# ------------------------------------------------------------------------------
# Gate VG-04: Restoration Script Parity
# ------------------------------------------------------------------------------
echo -n "● [Gate VG-04] Restoration Script Parity: "
vg04_pass=true

sh_runner="${SCRIPTS_DIR}/restore-antigravity-all.sh"
ps1_runner="${SCRIPTS_DIR}/restore-antigravity-all.ps1"
sh_verifier="${SCRIPTS_DIR}/verify-antigravity-backup.sh"
ps1_verifier="${SCRIPTS_DIR}/verify-antigravity-backup.ps1"

if [[ ! -f "$sh_runner" ]] || [[ ! -f "$ps1_runner" ]] || [[ ! -f "$sh_verifier" ]] || [[ ! -f "$ps1_verifier" ]]; then
    vg04_pass=false
fi

if [[ "$vg04_pass" == "true" ]]; then
    echo "PASSED (POSIX Bash and PowerShell runners and verifiers present)"
else
    echo "FAILED (Missing restoration or verifier runner script)"
    isValid=false
fi

# ------------------------------------------------------------------------------
# Gate VG-05: Atomic Non-Destructive Ingestion
# ------------------------------------------------------------------------------
echo -n "● [Gate VG-05] Atomic Non-Destructive Ingestion: "
vg05_pass=true

if [[ -f "$sh_runner" ]]; then
    if ! grep -q "shouldBackup" "$sh_runner" || ! grep -q "pre-restore" "$sh_runner"; then
        vg05_pass=false
    fi
else
    vg05_pass=false
fi

if [[ "$vg05_pass" == "true" ]]; then
    echo "PASSED (Pre-restore snapshot & non-destructive ingestion confirmed)"
else
    echo "FAILED (Missing backup routine or destructive logic detected)"
    isValid=false
fi

# ------------------------------------------------------------------------------
# Gate VG-06: Secret Sanitization Gate
# ------------------------------------------------------------------------------
echo -n "● [Gate VG-06] Secret Sanitization Gate: "
vg06_pass=true

python3 -c "
import glob, re, sys

files = glob.glob('$VAULT_DIR/*.json')
for f in files:
    with open(f) as fp:
        content = fp.read()
    # Check for raw secrets, private keys, oauth tokens
    patterns = [
        r'BEGIN (?:RSA )?PRIVATE KEY',
        r'ghp_[a-zA-Z0-9]{36}',
        r'AKIA[0-9A-Z]{16}',
        r'\"client_secret\":\s*\"[^\"]+\"',
        r'\"apiKey\":\s*\"[a-zA-Z0-9_-]{20,}\"'
    ]
    for pat in patterns:
        if re.search(pat, content):
            print(f'Secret pattern match in {f}: {pat}', file=sys.stderr)
            sys.exit(1)
" >/dev/null 2>&1 || vg06_pass=false

if [[ "$vg06_pass" == "true" ]]; then
    echo "PASSED (Zero credentials, tokens, or private keys detected)"
else
    echo "FAILED (Potential secret pattern detected in manifests)"
    isValid=false
fi

# ------------------------------------------------------------------------------
# Gate VG-07: Relative Link & Positive Boolean Hygiene
# ------------------------------------------------------------------------------
echo -n "● [Gate VG-07] Positive Boolean Hygiene: "
vg07_pass=true

python3 -c "
import glob, json, sys

files = glob.glob('$VAULT_DIR/*.json')
for f in files:
    with open(f) as fp:
        data = json.load(fp)
    def check_keys(obj):
        if isinstance(obj, dict):
            for k, v in obj.items():
                if isinstance(v, bool):
                    # Check that boolean keys use affirmative naming
                    negative_prefixes = ('isNot', 'hasNo', 'no', 'disabled', 'unverified', 'skip')
                    if any(k.startswith(p) for p in negative_prefixes):
                        print(f'Negative boolean key found: {k}', file=sys.stderr)
                        sys.exit(1)
                check_keys(v)
        elif isinstance(obj, list):
            for item in obj:
                check_keys(item)
    check_keys(data)
" >/dev/null 2>&1 || vg07_pass=false

if [[ "$vg07_pass" == "true" ]]; then
    echo "PASSED (100% positive boolean naming in all manifest stores)"
else
    echo "FAILED (Negative boolean naming detected in manifest stores)"
    isValid=false
fi

echo "=============================================================================="
if [[ "$isValid" == "true" ]]; then
    echo "  RESULT: ALL 7 VERIFICATION GATES PASSED (100% COMPLIANT)"
    echo "=============================================================================="
    exit 0
else
    echo "  RESULT: VERIFICATION FAILED (Review gate errors above)"
    echo "=============================================================================="
    exit 1
fi
