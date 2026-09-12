package cmdschedule

var (
	CheckHelpFn func(subcmd string, args []string)
)

// RunSchedule handles gitmap schedule commands.
func RunSchedule(args []string) error {
	return runSchedule(args)
}
