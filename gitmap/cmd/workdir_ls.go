// Package cmd — workdir_ls.go displays registered work directories table.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/charmbracelet/lipgloss"
)

func runWorkDirLs() error {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return errDB
	}
	defer db.Close()

	dirs, errList := db.ListWorkDirs()
	if errList != nil {
		return errList
	}

	if len(dirs) == 0 {
		fmt.Println("No work directories registered. Run `gitmap workdir add <path>` to register one.")
		return nil
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	fmt.Println(titleStyle.Render("Registered Work Directories:"))
	fmt.Printf("  %-4s  %-35s  %-20s  %s\n", "ID", "PATH", "LABEL", "DEFAULT")
	fmt.Println("  --------------------------------------------------------------------------------")

	for _, d := range dirs {
		defMarker := ""
		if d.IsDefault {
			defMarker = "★ [DEFAULT]"
		}
		label := d.Label
		if label == "" {
			label = "-"
		}
		fmt.Printf("  %-4d  %-35s  %-20s  %s\n", d.ID, d.AbsolutePath, label, defMarker)
	}
	return nil
}
