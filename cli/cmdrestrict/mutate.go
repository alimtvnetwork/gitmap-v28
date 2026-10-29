package cmdrestrict

import (
	"fmt"
	"strings"
)

const (
	flagName         = "--name"
	flagEmail        = "--email"
	flagTimezone     = "--tz"
	flagBlockedAfter = "--blocked-after"
	flagWeekendOpen   = "--weekend-open"
	flagWeekendClosed = "--weekend-closed"
	flagPermitStay    = "--permit-stay"
)

func addRestriction(args []string) error {
	item := restrictionFromArgs(args)
	if item.Name == "" {
		return fmt.Errorf("restrict add requires --name")
	}
	file, err := loadRestrictions()
	if err != nil {
		return err
	}
	file.Items = append(file.Items, item)
	if errSave := saveRestrictions(file); errSave != nil {
		return errSave
	}
	fmt.Printf("Added restriction %s\n", item.Name)

	return nil
}

func editRestriction(args []string) error {
	patch := restrictionFromArgs(args)
	if patch.Name == "" {
		return fmt.Errorf("restrict edit requires --name")
	}
	file, err := loadRestrictions()
	if err != nil {
		return err
	}
	next, isFound := replaceNamed(file.Items, patch, args)
	if !isFound {
		return fmt.Errorf("restriction %s was not found", patch.Name)
	}
	file.Items = next
	if errSave := saveRestrictions(file); errSave != nil {
		return errSave
	}
	fmt.Printf("Edited restriction %s\n", patch.Name)

	return nil
}

func replaceNamed(items []Restriction, patch Restriction, args []string) ([]Restriction, bool) {
	var next []Restriction
	isFound := false
	for _, item := range items {
		next = append(next, mergedItem(item, patch, args, &isFound))
	}

	return next, isFound
}

func mergedItem(item, patch Restriction, args []string, isFound *bool) Restriction {
	if item.Name != patch.Name {
		return item
	}
	*isFound = true

	return applyPatch(item, patch, args)
}

func applyPatch(item, patch Restriction, args []string) Restriction {
	item.Email = firstNonEmpty(patch.Email, item.Email)
	item.Timezone = firstNonEmpty(patch.Timezone, item.Timezone)
	item.BlockedAfter = firstNonEmpty(patch.BlockedAfter, item.BlockedAfter)
	item.IsWeekendOpen = weekendValue(item.IsWeekendOpen, args)
	item.IsPermitStay = permitValue(item.IsPermitStay, args)

	return item
}

func weekendValue(current bool, args []string) bool {
	if hasFlag(args, flagWeekendClosed) {
		return false
	}
	if hasFlag(args, flagWeekendOpen) {
		return true
	}

	return current
}

func permitValue(current bool, args []string) bool {
	if hasFlag(args, flagPermitStay) {
		return true
	}

	return current
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}

	return false
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}

func restrictionFromArgs(args []string) Restriction {
	item := defaultRestriction()
	item.Name = ""
	item.IsWeekendOpen = false
	for i := 0; i < len(args); i++ {
		item = applyFlag(item, args, &i)
	}

	return item
}

func applyFlag(item Restriction, args []string, index *int) Restriction {
	switch args[*index] {
	case flagName:
		item.Name = nextValue(args, index)
	case flagEmail:
		item.Email = nextValue(args, index)
	case flagTimezone:
		item.Timezone = nextValue(args, index)
	case flagBlockedAfter:
		item.BlockedAfter = nextValue(args, index)
	case flagWeekendOpen:
		item.IsWeekendOpen = true
	case flagWeekendClosed:
		item.IsWeekendOpen = false
	case flagPermitStay:
		item.IsPermitStay = true
	default:
	}

	return item
}

func nextValue(args []string, index *int) string {
	if *index+1 >= len(args) {
		return ""
	}
	*index = *index + 1

	return args[*index]
}
