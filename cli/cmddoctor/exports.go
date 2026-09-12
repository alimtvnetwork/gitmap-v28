package cmddoctor

// RunDoctorCmd executes the doctor command.
func RunDoctorCmd(args []string) error {
	return runDoctor(args)
}

// RunCleanCorrupted executes the clean-corrupted command.
func RunCleanCorrupted(args []string) error {
	return runCleanCorrupted(args)
}
