package cmdpull

// pullOptions holds parsed pull flags.
type pullOptions struct {
	slug          string
	group         string
	all           bool
	verbose       bool
	stopOnFail    bool
	parallel      int
	workers       int
	hands         int
	isWWOH        bool
	isAutoScale   bool
	isHighPerf    bool
	isLowCPU      bool
	onlyAvailable bool
	autoFix       bool
	yes           bool
	noFix         bool
	isRaw         bool
	useSSH        bool
	useHTTPS      bool
	showStatus    bool
	isJSON        bool
	isProbe       bool
}
