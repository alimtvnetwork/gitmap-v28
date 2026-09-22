// Package cmdagy — agy_pins_edit.go handles editing pinned project display attributes.
package cmdagy

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var agyPinProjectsEditCmd = &cobra.Command{
	Use:     "edit [target] [new-name]",
	Aliases: []string{"rename", "set-name"},
	Short:   "Update display name or metadata for a pinned project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyPinsEdit(args)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeAgyRmArgs(toComplete)
	},
}

func runAgyPinsEdit(args []string) error {
	target, newName, err := parsePinsEditArgs(args)
	if err != nil {
		return err
	}
	return executePinsEdit(target, newName)
}

func parsePinsEditArgs(args []string) (string, string, error) {
	if len(args) < 2 {
		return "", "", apperror.NewSimple("requires target (seq/id/slug) and new name", "E9000")
	}
	return args[0], strings.TrimSpace(strings.Join(args[1:], " ")), nil
}

func executePinsEdit(target, newName string) error {
	store, loadErr := loadPinnedProjectsStore()
	if loadErr != nil {
		return loadErr
	}
	index, p := findPinnedIndex(store, target)
	if index == -1 || p == nil {
		return apperror.NewSimple(fmt.Sprintf("pinned project %q not found", target), "E9000")
	}
	return applyPinsEdit(store, p, newName)
}

func applyPinsEdit(store *PinnedProjectsStore, p *PinnedProject, newName string) error {
	oldName := p.Name
	p.Name = newName
	store.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if saveErr := savePinnedProjectsStore(store); saveErr != nil {
		return saveErr
	}
	fmt.Printf("  %s✓ Pinned project '%s' renamed to '%s'%s\n",
		constants.ColorGreen, oldName, newName, constants.ColorReset)
	return nil
}
