package constants

// Tier-2 shared literals (plan 325 sweep B). These literals were verified
// as standalone code values (filenames, host defaults, error substrings)
// across many packages; embedded uses inside longer messages, URLs, paths,
// and help text were deliberately left alone. GoModFile and ReadmeFile
// already exist elsewhere in this package and cover "go.mod" and
// "README.md", so they are not redefined here.
const (
	// GitIgnoreFile is the repository ignore filename.
	GitIgnoreFile = ".gitignore"
	// Localhost is the loopback hostname used as a default host value
	// and in local-host comparisons.
	Localhost = "localhost"
	// AlreadyExistsSubstring matches the "already exists" error substring
	// used by string-based collision/benign-error checks.
	AlreadyExistsSubstring = "already exists"
)
