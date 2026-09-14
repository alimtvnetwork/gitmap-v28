package cmdos

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runOSCron(args []string) error {
	if len(args) == 0 || isOSHelpArg(args[0]) {
		printOSCronUsage()

		return nil
	}

	sub := strings.ToLower(args[0])
	rest := args[1:]

	return dispatchOSCronSubcommand(sub, rest)
}

func dispatchOSCronSubcommand(sub string, rest []string) error {
	switch sub {
	case "ls", "list", "show":
		return runOSCronList()
	case "add":
		return runOSCronAdd(rest)
	case "rm", "delete", "del":
		return runOSCronRemove(rest)
	case "clear":
		return runOSCronClear()
	case "edit":
		return runOSCronEdit()
	default:
		return apperror.NewSimple("unknown os cron subcommand: "+sub, "E_OS_CRON_INVALID")
	}
}

func runOSCronList() error {
	if runtime.GOOS == constants.OSWindows {
		fmt.Println("▶ Crontab is a POSIX utility. Showing Windows Scheduled Tasks:")
		cmd := exec.Command("schtasks", "/query", "/fo", "TABLE")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}

	cmd := exec.Command("crontab", "-l")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runOSCronAdd(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("cron job specification required", "E_OS_CRON_JOB_REQUIRED")
	}

	job := strings.Join(args, " ")
	if runtime.GOOS == constants.OSWindows {
		fmt.Printf("ℹ Windows detected. Recommending 'gitmap schedule add' for job: %q\n", job)

		return nil
	}

	return appendCronJobPOSIX(job)
}

func appendCronJobPOSIX(job string) error {
	out, _ := exec.Command("crontab", "-l").Output()
	current := string(out)
	if strings.Contains(current, job) {
		fmt.Printf("ℹ Job already exists in crontab: %q\n", job)

		return nil
	}

	newCron := current + "\n" + job + "\n"
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(newCron)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "write crontab")
	}

	fmt.Printf("✔ Cron job added successfully: %q\n", job)

	return nil
}

func runOSCronRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("pattern or job required", "E_OS_CRON_PATTERN_REQUIRED")
	}

	pattern := args[0]
	if runtime.GOOS == constants.OSWindows {
		fmt.Printf("ℹ Windows detected. To manage scheduled tasks, use: schtasks /delete /tn %q\n", pattern)

		return nil
	}

	return removeCronJobPOSIX(pattern)
}

func removeCronJobPOSIX(pattern string) error {
	out, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		return apperror.WrapSimple(err, "read crontab")
	}

	lines := strings.Split(string(out), "\n")
	var kept []string
	removedCount := 0
	for _, l := range lines {
		if strings.Contains(l, pattern) {
			removedCount++
			continue
		}
		kept = append(kept, l)
	}

	return writeFilteredCrontab(kept, removedCount, pattern)
}

func writeFilteredCrontab(kept []string, removedCount int, pattern string) error {
	newCron := strings.Join(kept, "\n")
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(newCron)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "write filtered crontab")
	}

	fmt.Printf("✔ Removed %d matching cron job(s) for pattern %q\n", removedCount, pattern)

	return nil
}

func runOSCronClear() error {
	if runtime.GOOS == constants.OSWindows {
		fmt.Println("ℹ Crontab clear is only applicable on POSIX systems.")

		return nil
	}

	cmd := exec.Command("crontab", "-r")
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "crontab -r")
	}

	fmt.Println("✔ Crontab cleared successfully")

	return nil
}

func runOSCronEdit() error {
	if runtime.GOOS == constants.OSWindows {
		fmt.Println("ℹ Crontab edit is only applicable on POSIX systems.")

		return nil
	}

	cmd := exec.Command("crontab", "-e")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func printOSCronUsage() {
	fmt.Println(`Usage: gitmap os cron [subcommand] [args]

Commands:
  ls, list           List current user crontab jobs
  add "<job>"        Append a new job line to crontab
  rm "<pattern>"     Remove lines matching pattern from crontab
  clear              Remove all user crontab entries
  edit               Open crontab in default interactive editor`)
}
