// Package cmd — doctor_run.go: flag-aware entry point for `gitmap doctor`.
//
//	--json   emit machine-readable JSON instead of the colorized text report
//	--fix    attempt safe auto-fixes (create .gitmap/, suggest exact recipes)
//
// Exit status: 0 when every probe passes, 1 otherwise.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDoctor is invoked by the root dispatcher.
func runDoctor(args []string) {
	// Sub-command dispatch (v6.80.1+): `gitmap doctor fix-repo [...]`
	// routes to the fix-repo → gofmt probe suite instead of the
	// generic dependency checks. Alias `fr` matches CmdFixRepoAlias.
	if len(args) > 0 && (args[0] == constants.CmdFixRepo || args[0] == constants.CmdFixRepoAlias) {
		runDoctorFixRepo(args[1:])

		return
	}
	wantJSON, wantFix := false, false
	for _, a := range args {
		switch a {
		case "--json":
			wantJSON = true
		case "--fix":
			wantFix = true
		}
	}
	checks := defaultDoctorChecks()
	results := make([]DoctorResult, 0, len(checks))
	failed := 0
	for _, c := range checks {
		ok, detail := c.Run()
		if !ok && wantFix {
			ok, detail = applyDoctorFix(c, detail)
		}
		results = append(results, DoctorResult{Name: c.Name, OK: ok, Detail: detail, FixHint: c.FixHint})
		if !ok {
			failed++
		}
	}
	if wantJSON {
		emitDoctorJSON(results, failed)
	} else {
		emitDoctorText(results, failed)
	}
	if failed > 0 {
		os.Exit(1)
	}
}

func emitDoctorJSON(results []DoctorResult, failed int) {
	payload := struct {
		Failed  int            `json:"failed"`
		Results []DoctorResult `json:"results"`
	}{Failed: failed, Results: results}
	buf, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(buf))
}

func emitDoctorText(results []DoctorResult, failed int) {
	for _, r := range results {
		mark := "[ok]  "
		if !r.OK {
			mark = "[fail]"
		}
		fmt.Printf("%s %-10s %s\n", mark, r.Name, r.Detail)
		if !r.OK && r.FixHint != "" {
			fmt.Printf("           fix: %s\n", r.FixHint)
		}
	}
	if failed > 0 {
		fmt.Printf("\n%d check(s) failed. Re-run with --fix to attempt auto-repair.\n", failed)
		return
	}
	fmt.Println("\nAll systems nominal.")
}

// applyDoctorFix performs safe, idempotent repairs for known check names
// and re-runs the probe. Anything it cannot fix is annotated with a
// copy-pasteable next-step recipe.
func applyDoctorFix(c DoctorCheck, detail string) (bool, string) {
	switch c.Name {
	case "config":
		_ = os.MkdirAll(".gitmap", 0o755)
		return c.Run()
	case "gh-token":
		return false, detail + "\n           run: gitmap secrets set GITHUB_TOKEN  (or export GITHUB_TOKEN=<pat>)"
	case "PATH":
		return false, detail + "\n           run: gitmap self-install"
	default:
		return false, detail
	}
}
