package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

func runUserList(args []string) error {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading profiles: %v\n", err)
		return nil
	}

	boundProf, _ := usercontext.ResolveProjectUser(".")
	rows := buildProfileTableRows(cfg.Profiles, cfg.Active, cfg.Default, boundProf)
	renderUserProfilesTable(rows)

	return nil
}

func renderUserProfilesTable(rows []termtable.Row) {
	termtable.PrintTable(termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "ACT", Align: termtable.AlignLeft, MinWidth: 4},
			{Title: "ALIAS", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "NAME", Align: termtable.AlignLeft, MinWidth: 20},
			{Title: "EMAIL", Align: termtable.AlignLeft, MinWidth: 26},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 14},
		},
		Rows: rows,
	})
}

func buildProfileTableRows(profiles []model.GitProfile, active, def string, bound *model.GitProfile) []termtable.Row {
	rows := make([]termtable.Row, 0, len(profiles))
	for _, p := range profiles {
		row := buildSingleProfileRow(p, active, def, bound)
		rows = append(rows, row)
	}

	return rows
}

func buildSingleProfileRow(p model.GitProfile, active, def string, bound *model.GitProfile) termtable.Row {
	mark, status := resolveProfileMarkers(p, active, def, bound)

	return termtable.Row{
		Cells: []string{mark, p.ID, p.Name, p.Email, status},
	}
}

func resolveProfileMarkers(p model.GitProfile, active, def string, bound *model.GitProfile) (string, string) {
	isActive := p.ID == active || p.Name == active
	isDefault := p.ID == def || p.Name == def || p.IsDefault
	isBound := bound != nil && (p.ID == bound.ID || p.Name == bound.Name)

	if isBound {
		return "*", "bound"
	}
	if isActive {
		return "*", "active"
	}
	if isDefault {
		return "*", "default"
	}

	return " ", "configured"
}

func runUserSwitch(args []string) error {
	alias := extractFirstPositional(args)
	if len(alias) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap user switch <alias> [--global] [--project]\n")
		return nil
	}

	cfg, err := store.LoadGitProfiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading profiles: %v\n", err)
		return nil
	}

	prof, isFound := findProfileInConfig(cfg.Profiles, alias)
	if !isFound {
		fmt.Fprintf(os.Stderr, "Error: profile %q not found\n", alias)
		return nil
	}

	return executeProfileSwitch(&cfg, prof, args)
}

func findProfileInConfig(profiles []model.GitProfile, alias string) (model.GitProfile, bool) {
	for _, p := range profiles {
		if p.ID == alias || strings.EqualFold(p.Name, alias) {
			return p, true
		}
	}

	return model.GitProfile{}, false
}

func executeProfileSwitch(cfg *model.GitProfileConfig, p model.GitProfile, args []string) error {
	isGlobal := hasUserFlag(args, "--global")
	isProject := hasUserFlag(args, "--project")

	if isGlobal {
		return applyGlobalSwitch(cfg, p)
	}
	if isProject {
		return applyProjectSwitch(p)
	}

	return applyDefaultSwitch(cfg, p)
}

func applyGlobalSwitch(cfg *model.GitProfileConfig, p model.GitProfile) error {
	cfg.Active = p.ID
	store.SaveGitProfiles(*cfg)
	usercontext.SetGitIdentity(p.Name, p.Email, true)
	fmt.Printf("✔ Switched global Git identity to %s <%s>\n", p.Name, p.Email)

	return nil
}

func applyProjectSwitch(p model.GitProfile) error {
	usercontext.BindProject(".", p.ID)
	usercontext.SetGitIdentity(p.Name, p.Email, false)
	fmt.Printf("✔ Bound current project to profile %q (%s <%s>)\n", p.ID, p.Name, p.Email)

	return nil
}

func applyDefaultSwitch(cfg *model.GitProfileConfig, p model.GitProfile) error {
	cfg.Active = p.ID
	store.SaveGitProfiles(*cfg)
	usercontext.BindProject(".", p.ID)
	usercontext.SetGitIdentity(p.Name, p.Email, false)
	fmt.Printf("✔ Switched active profile to %q (%s <%s>)\n", p.ID, p.Name, p.Email)

	return nil
}

func runUserAddProfile(args []string) error {
	alias := extractFirstPositional(args)
	if len(alias) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap user add <alias> --name <name> --email <email>\n")
		return nil
	}

	name := extractUserArgFlagDefault(args, "--name", alias)
	email := extractUserArgFlag(args, "--email")
	cfg, _ := store.LoadGitProfiles()
	prof := createNewGitProfile(alias, name, email, hasUserFlag(args, "--default"))

	cfg.Profiles = appendOrUpdateProfile(cfg.Profiles, prof)
	if prof.IsDefault || len(cfg.Active) == 0 {
		cfg.Active = prof.ID
	}

	store.SaveGitProfiles(cfg)
	fmt.Printf("✔ Added Git profile %q (%s <%s>)\n", alias, name, email)

	return nil
}

func createNewGitProfile(alias, name, email string, isDefault bool) model.GitProfile {
	return model.GitProfile{
		ID:         alias,
		Name:       name,
		Email:      email,
		Provider:   "git",
		Type:       "user",
		AuthMethod: "config",
		IsDefault:  isDefault,
		LastUsedAt: time.Now(),
	}
}

func appendOrUpdateProfile(profiles []model.GitProfile, p model.GitProfile) []model.GitProfile {
	for idx := range profiles {
		if profiles[idx].ID == p.ID {
			profiles[idx] = p
			return profiles
		}
	}

	return append(profiles, p)
}

func extractFirstPositional(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			return a
		}
	}

	return ""
}
