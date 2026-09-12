package cmd

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/power"
)

func runPowerStatus() error {
	mgr, err := power.NewManager()
	if err != nil {
		return err
	}

	current, err := mgr.GetStatus()
	if err != nil {
		return err
	}

	dbSetting := power.Settings{}
	if db, err := openDB(); err == nil && db != nil {
		defer db.Close()
		dbSetting, _ = db.GetActivePowerSetting()
	}

	printPowerStatus(current, dbSetting)

	return nil
}

func runPowerNeverSleep() error {
	mgr, err := power.NewManager()
	if err != nil {
		return err
	}

	current, _ := mgr.GetStatus()
	recordPowerTransition("never-sleep", current, 0, 0, true)

	if err := mgr.SetNeverSleep(); err != nil {
		return err
	}

	fmt.Println("✓ Power setting updated: Never-Sleep mode enabled (display & standby timeouts disabled).")

	return nil
}

func runPowerSet(args []string) error {
	disp, sleep, err := parseSetArgs(args)
	if err != nil {
		return err
	}

	mgr, err := power.NewManager()
	if err != nil {
		return err
	}

	current, _ := mgr.GetStatus()
	isNever := disp == 0 && sleep == 0
	recordPowerTransition("set", current, disp, sleep, isNever)

	if err := mgr.SetTimeouts(disp, sleep); err != nil {
		return err
	}

	fmt.Printf("✓ Power setting updated: Display timeout=%dm, Sleep timeout=%dm.\n", disp, sleep)

	return nil
}

func parseSetArgs(args []string) (int, int, error) {
	fs := flag.NewFlagSet("set", flag.ContinueOnError)
	disp := fs.Int("display", -1, "Display timeout in minutes")
	sleep := fs.Int("sleep", -1, "Sleep timeout in minutes")

	if err := fs.Parse(args); err != nil {
		return 0, 0, apperror.WrapSimple(err, "parseSetArgs")
	}

	if *disp >= 0 || *sleep >= 0 {
		d, s := normalizeFlags(*disp, *sleep)

		return d, s, nil
	}

	val, err := parseFirstPositional(fs)
	if err == nil && val >= 0 {
		return val, val, nil
	}

	return 0, 0, apperror.NewSimple("usage: gitmap power set <minutes> OR --display <m> --sleep <m>", "E_INVALID_ARGS")
}

func parseFirstPositional(fs *flag.FlagSet) (int, error) {
	if fs.NArg() == 0 {
		return -1, fmt.Errorf("no positional args")
	}

	return strconv.Atoi(fs.Arg(0))
}

func normalizeFlags(disp, sleep int) (int, int) {
	if disp < 0 {
		disp = 10
	}

	if sleep < 0 {
		sleep = 30
	}

	return disp, sleep
}

func recordPowerTransition(action string, current power.Settings, disp, sleep int, isNever bool) {
	db, err := openDB()
	if err != nil || db == nil {
		return
	}

	defer db.Close()

	_ = db.SavePowerProfile("previous", current, false)
	newS := power.Settings{
		Platform:              current.Platform,
		DisplayTimeoutMinutes: disp,
		SleepTimeoutMinutes:   sleep,
		IsNeverSleep:          isNever,
		Source:                action,
	}

	_ = db.SavePowerProfile("current", newS, true)
	_ = db.RecordPowerHistory(action, newS, fmt.Sprintf("Action=%s, display=%d, sleep=%d", action, disp, sleep))
}

func runPowerReset() error {
	mgr, err := power.NewManager()
	if err != nil {
		return err
	}

	target := power.Settings{
		Platform:              mgr.Platform(),
		DisplayTimeoutMinutes: 10,
		SleepTimeoutMinutes:   30,
	}

	target = loadPreviousOrDefault(target)

	if err := mgr.ApplySettings(target); err != nil {
		return err
	}

	fmt.Printf("✓ Power settings reset: Display=%s, Sleep=%s.\n",
		formatTimeoutMinutes(target.DisplayTimeoutMinutes),
		formatTimeoutMinutes(target.SleepTimeoutMinutes))

	return nil
}

func loadPreviousOrDefault(def power.Settings) power.Settings {
	db, err := openDB()
	if err != nil || db == nil {
		return def
	}

	defer db.Close()

	prev, err := db.GetPowerProfile("previous")
	if err == nil {
		def = prev
	}

	_ = db.RecordPowerHistory("reset", def, "Reset to previous configuration")

	return def
}

func runPowerHistory(args []string) error {
	db, err := openDB()
	if err != nil || db == nil {
		return apperror.NewSimple("database unavailable for power history", "E_DB_UNAVAILABLE")
	}

	defer db.Close()

	records, err := db.ListPowerHistory(30)
	if err != nil {
		return err
	}

	printPowerHistoryTable(records)

	return nil
}
