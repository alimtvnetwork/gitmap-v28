import re

# 1. Fix race-detector.yml
with open(".github/workflows/race-detector.yml", "r", encoding="utf-8") as f:
    rc = f.read()
rc = rc.replace("- uses: ./.github/actions/setup-go-cached\n\n      - name: File-size lint",
                "- uses: ./.github/actions/setup-go-cached\n        with:\n          cache-suffix: 'race'\n\n      - name: File-size lint")
with open(".github/workflows/race-detector.yml", "w", encoding="utf-8") as f:
    f.write(rc)

# 2. Fix ci.yml
with open(".github/workflows/ci.yml", "r", encoding="utf-8") as f:
    ci = f.read()

# Fix smoke test missing run
ci = ci.replace("""      - name: Run installer smoke (source mode)\n        if: needs.sha-check.outputs.already-built != 'true'\n        shell: pwsh\n\n  e2e-cross-platform:""",
                """      - name: Run installer smoke (source mode)\n        if: needs.sha-check.outputs.already-built != 'true'\n        shell: pwsh\n        run: ./.github/scripts/smoke-installer.ps1 source\n\n  e2e-cross-platform:""")

# Fix [email protected] -> BYOB@v1.3.0
ci = ci.replace("uses: RubbaBoy/[email protected]", "uses: RubbaBoy/BYOB@v1.3.0")

with open(".github/workflows/ci.yml", "w", encoding="utf-8") as f:
    f.write(ci)

print("done")
