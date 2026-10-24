package cmdsupabase

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func handleTest(args []string) error {
	if IsSupabaseErrorSuppressed() {
		return apperror.NewSimple("supabase connection suppressed (5m backoff active)", "E_SB_SUPPRESSED")
	}

	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap supabase ping <alias>")
	}

	alias := args[0]
	p, err := GetProjectByAlias(alias)
	if err != nil {
		return err
	}

	return executePingProbe(p)
}

func executePingProbe(p *SupabaseProjectRecord) error {
	if IsSupabaseErrorSuppressed() {
		// Skip connection attempt for 5 minutes
		return apperror.NewSimple("supabase connection suppressed (5m backoff active)", "E_SB_SUPPRESSED")
	}

	apiKey, err := DecryptSecret(p.AnonKeyEnc)
	if err != nil || apiKey == "" {
		apiKey, _ = DecryptSecret(p.ServiceKeyEnc)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	reqUrl := strings.TrimRight(p.ApiUrl, "/") + "/rest/v1/"
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		RecordSupabaseError(err)
		return apperror.WrapSimple(err, "executePingProbe_NewRequest")
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", constants.AuthBearerPrefix+apiKey)

	start := time.Now()
	resp, reqErr := client.Do(req)
	latency := time.Since(start)

	if reqErr != nil {
		RecordSupabaseError(reqErr)
	} else if resp != nil && (resp.StatusCode < 200 || resp.StatusCode >= 400) {
		RecordSupabaseError(apperror.NewExecutionError(fmt.Sprintf("supabase ping received HTTP %d", resp.StatusCode)))
	}

	return reportPingResult(p.Alias, p.ApiUrl, resp, latency, reqErr)
}

func reportPingResult(alias, url string, resp *http.Response, latency time.Duration, reqErr error) error {
	if reqErr != nil {
		fmt.Printf("✖ Ping failed for [%s] (%s): %v\n", alias, url, reqErr)
		return reqErr
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		fmt.Printf("✓ Ping succeeded for [%s] HTTP %d (latency: %v)\n", alias, resp.StatusCode, latency)
		return nil
	}

	fmt.Printf("✖ Ping received non-2xx for [%s] HTTP %d (latency: %v)\n", alias, resp.StatusCode, latency)
	return nil
}

func printSupabaseHelp() error {
	fmt.Printf(`%sGitMap Multi-Supabase Database Vault%s

Usage:
  gitmap supabase <command> [arguments]

Commands:
  add <alias> <url> <anon> <service> [db]  Register encrypted Supabase project
  list, ls                                List all configured Supabase projects
  get <alias> [--decrypt]                 Show details for a project
  use <alias>                             Set active status for project
  ping, test <alias>                      Probe connectivity to project REST API
  env, export-env <alias>                 Print shell environment exports
  remove, rm <alias>                      Delete project from vault
  help                                    Show this help message

Options:
  --desc <description>                    Project description for add
  --project-ref <ref>                     Supabase project reference for add
  --decrypt                               Show decrypted credentials in get

Zero Cleartext Invariant:
  All service keys and database URLs are stored using authenticated AES-256-GCM encryption.
`, constants.ColorCyan, constants.ColorReset)

	return nil
}
