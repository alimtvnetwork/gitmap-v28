package cmdsupabase

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSupabaseCLI routes and executes Supabase subcommands.
func RunSupabaseCLI(args []string) error {
	return Run(args)
}

// Run is the main dispatcher for the Supabase management CLI.
func Run(args []string) error {
	if len(args) == 0 {
		return printSupabaseHelp()
	}

	subcmd := strings.ToLower(args[0])
	subArgs := args[1:]

	return dispatchSubcommand(subcmd, subArgs)
}

func dispatchSubcommand(subcmd string, subArgs []string) error {
	switch subcmd {
	case "add":
		return handleAdd(subArgs)
	case "list", "ls":
		return handleList(subArgs)
	case "get":
		return handleGet(subArgs)
	case "use":
		return handleUse(subArgs)
	case "ping", "test":
		return handleTest(subArgs)
	case "remove", "rm":
		return handleRemove(subArgs)
	case "env", "export-env":
		return handleExportEnv(subArgs)
	case "help", "-h", "--help":
		return printSupabaseHelp()
	default:
		return apperror.NewValidationError("unknown supabase subcommand: " + subcmd)
	}
}

func handleAdd(args []string) error {
	positional, flags := parseAddFlags(args)
	if len(positional) < 4 {
		return apperror.NewValidationError("usage: gitmap supabase add <alias> <url> <anon_key> <service_key> [db_url]")
	}

	alias, apiUrl, anonKey, serviceKey := positional[0], positional[1], positional[2], positional[3]
	if !isValidSupabaseUrl(apiUrl) {
		return apperror.NewValidationError("invalid supabase api url: must begin with http:// or https://")
	}

	dbUrl := extractOptionalDbUrl(positional)
	return persistEncryptedProject(alias, apiUrl, anonKey, serviceKey, dbUrl, flags)
}

func parseAddFlags(args []string) ([]string, map[string]string) {
	positional := make([]string, 0, len(args))
	flags := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--desc" && i+1 < len(args) {
			flags["desc"] = args[i+1]
			i++
			continue
		}
		if arg == "--project-ref" && i+1 < len(args) {
			flags["project-ref"] = args[i+1]
			i++
			continue
		}
		if !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
		}
	}

	return positional, flags
}

func extractOptionalDbUrl(positional []string) string {
	if len(positional) >= 5 {
		return positional[4]
	}

	return ""
}

func isValidSupabaseUrl(u string) bool {
	lower := strings.ToLower(u)

	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func persistEncryptedProject(alias, apiUrl, anonKey, serviceKey, dbUrl string, flags map[string]string) error {
	anonEnc, err := EncryptSecret(anonKey)
	if err != nil {
		return err
	}
	servEnc, err := EncryptSecret(serviceKey)
	if err != nil {
		return err
	}
	dbEnc, err := EncryptSecret(dbUrl)
	if err != nil {
		return err
	}

	rec := SupabaseProjectRecord{
		Alias: alias, ProjectRef: flags["project-ref"], ApiUrl: apiUrl,
		AnonKeyEnc: anonEnc, ServiceKeyEnc: servEnc, DbUrlEnc: dbEnc,
		Status: StatusActive, Description: flags["desc"],
	}

	if err := InsertProject(rec); err != nil {
		return err
	}
	fmt.Printf("✓ Encrypted and registered Supabase project [%s] into vault.\n", alias)

	return nil
}

func handleList(_ []string) error {
	projects, err := ListAllProjects()
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		fmt.Println("No Supabase projects configured in vault. Use 'gitmap supabase add' to register one.")
		return nil
	}

	renderProjectsTable(projects)

	return nil
}

func renderProjectsTable(projects []SupabaseProjectRecord) {
	fmt.Printf("\n%-16s %-32s %-10s %-20s %-12s %-20s\n",
		"ALIAS", "API URL", "STATUS", "ANON KEY (MASKED)", "HAS DB URL", "UPDATED AT")
	fmt.Println(strings.Repeat("-", 114))

	for _, p := range projects {
		hasDb := p.DbUrlEnc != ""
		decAnon, _ := DecryptSecret(p.AnonKeyEnc)
		masked := MaskSecret(decAnon)
		updatedStr := p.UpdatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("%-16s %-32s %-10s %-20s %-12t %-20s\n",
			p.Alias, p.ApiUrl, p.Status, masked, hasDb, updatedStr)
	}
	fmt.Println()
}

func handleGet(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap supabase get <alias> [--decrypt]")
	}

	alias := args[0]
	hasDecrypt := hasDecryptFlag(args)
	rec, err := GetProjectByAlias(alias)
	if err != nil {
		return err
	}

	return renderProjectDetails(rec, hasDecrypt)
}

func hasDecryptFlag(args []string) bool {
	for _, a := range args {
		if a == "--decrypt" || a == "-d" {
			return true
		}
	}

	return false
}

func renderProjectDetails(p *SupabaseProjectRecord, isDecrypted bool) error {
	anon, serv, dbUrl := resolveSecretDisplay(p, isDecrypted)

	fmt.Printf("\nSupabase Project: %s\n", p.Alias)
	fmt.Printf("  Project Ref:   %s\n", p.ProjectRef)
	fmt.Printf("  API URL:       %s\n", p.ApiUrl)
	fmt.Printf("  Status:        %s\n", p.Status)
	fmt.Printf("  Description:   %s\n", p.Description)
	fmt.Printf("  Anon Key:      %s\n", anon)
	fmt.Printf("  Service Key:   %s\n", serv)
	fmt.Printf("  Database URL:  %s\n", dbUrl)
	fmt.Printf("  Updated At:    %s\n\n", p.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func resolveSecretDisplay(p *SupabaseProjectRecord, isDecrypted bool) (string, string, string) {
	if !isDecrypted {
		return MaskSecret(p.AnonKeyEnc), MaskSecret(p.ServiceKeyEnc), MaskSecret(p.DbUrlEnc)
	}

	anon, _ := DecryptSecret(p.AnonKeyEnc)
	serv, _ := DecryptSecret(p.ServiceKeyEnc)
	dbUrl, _ := DecryptSecret(p.DbUrlEnc)

	return anon, serv, dbUrl
}

func handleUse(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap supabase use <alias>")
	}

	alias := args[0]
	if err := UpdateProjectStatus(alias, StatusActive); err != nil {
		return err
	}
	fmt.Printf("✓ Switched active Supabase project to: %s\n", alias)

	return nil
}

func handleTest(args []string) error {
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
	apiKey, err := DecryptSecret(p.AnonKeyEnc)
	if err != nil || apiKey == "" {
		apiKey, _ = DecryptSecret(p.ServiceKeyEnc)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	reqUrl := strings.TrimRight(p.ApiUrl, "/") + "/rest/v1/"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return apperror.WrapSimple(err, "executePingProbe_NewRequest")
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	start := time.Now()
	resp, reqErr := client.Do(req)
	latency := time.Since(start)

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

func handleRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap supabase remove <alias>")
	}

	alias := args[0]
	if err := DeleteProjectByAlias(alias); err != nil {
		return err
	}
	fmt.Printf("✓ Removed Supabase project [%s] from vault.\n", alias)

	return nil
}

func handleExportEnv(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("usage: gitmap supabase env <alias>")
	}

	alias := args[0]
	p, err := GetProjectByAlias(alias)
	if err != nil {
		return err
	}

	anon, _ := DecryptSecret(p.AnonKeyEnc)
	serv, _ := DecryptSecret(p.ServiceKeyEnc)
	dbUrl, _ := DecryptSecret(p.DbUrlEnc)

	fmt.Printf("export SUPABASE_URL=%q\n", p.ApiUrl)
	fmt.Printf("export SUPABASE_ANON_KEY=%q\n", anon)
	fmt.Printf("export SUPABASE_SERVICE_ROLE_KEY=%q\n", serv)
	if dbUrl != "" {
		fmt.Printf("export DATABASE_URL=%q\n", dbUrl)
	}

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
