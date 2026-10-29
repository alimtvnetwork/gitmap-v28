package cmdrestrict

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/timeutil"
)

const defaultBlockedAfter = "08:00"
const defaultTimezone = "Local"
const statusTimeLayout = "02-Jan-06 03:04 PM"

func printStatus() error {
	file, err := loadRestrictions()
	if err != nil {
		return err
	}
	item := activeRestriction(file)
	loc, errLoc := timeutil.ResolveLocation(item.Timezone)
	if errLoc != nil {
		return errLoc
	}

	return writeStatus(item, loc)
}

func writeStatus(item Restriction, loc *time.Location) error {
	now := time.Now().In(loc)
	fmt.Println("GitMap restricted status")
	fmt.Printf("Name: %s\n", item.Name)
	fmt.Printf("Time: %s\n", now.Format(statusTimeLayout))
	fmt.Printf("Timezone: %s\n", loc.String())
	fmt.Printf("Blocked after: %s\n", item.BlockedAfter)
	fmt.Printf("Email: %s\n", item.Email)
	fmt.Printf("Weekend open: %t\n", item.IsWeekendOpen)
	fmt.Printf("Permit stay: %t\n", item.IsPermitStay)
	fmt.Printf("Window: %s\n", windowLabel(now, loc, item))

	return nil
}

func activeRestriction(file RestrictionFile) Restriction {
	if len(file.Items) == 0 {
		return defaultRestriction()
	}

	return file.Items[0]
}

func defaultRestriction() Restriction {
	return Restriction{
		Name:          "weekday-morning",
		BlockedAfter:  defaultBlockedAfter,
		Timezone:      defaultTimezone,
		IsWeekendOpen: true,
	}
}

func windowLabel(now time.Time, loc *time.Location, item Restriction) string {
	if item.IsWeekendOpen && timeutil.IsWeekend(now) {
		return "open (weekend)"
	}
	if timeutil.IsCommitWindowAllowed(now, loc) {
		return "open"
	}
	next := timeutil.CalculateNextWindow(now, loc)

	return "restricted until " + next.Format(statusTimeLayout)
}
