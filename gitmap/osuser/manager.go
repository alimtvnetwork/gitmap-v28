package osuser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// AddUser creates a new OS user across Windows and Linux.
func AddUser(username, password string) error {
	switch runtime.GOOS {
	case "windows":
		return addWindowsUser(username, password)
	case "linux":
		return addLinuxUser(username, password)
	default:
		return fmt.Errorf("unsupported operating system for user management: %s", runtime.GOOS)
	}
}

func addWindowsUser(username, password string) error {
	var cmd *exec.Cmd
	if password != "" {
		cmd = exec.Command("net", "user", username, password, "/ADD")
	} else {
		cmd = exec.Command("net", "user", username, "/ADD")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("windows user add failed: %w\nOutput: %s", err, out)
	}

	return nil
}

func addLinuxUser(username, password string) error {
	cmd := exec.Command("useradd", "-m", username)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("linux useradd failed: %w\nOutput: %s", err, out)
	}

	if password != "" {
		return setLinuxUserPassword(username, password)
	}

	return nil
}

func setLinuxUserPassword(username, password string) error {
	chCmd := exec.Command("chpasswd")
	stdin, err := chCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to pipe chpasswd: %w", err)
	}

	if err := chCmd.Start(); err != nil {
		return fmt.Errorf("failed to start chpasswd: %w", err)
	}

	if _, err := fmt.Fprintf(stdin, "%s:%s", username, password); err != nil {
		return fmt.Errorf("failed to write to chpasswd: %w", err)
	}

	stdin.Close()
	if err := chCmd.Wait(); err != nil {
		return fmt.Errorf("chpasswd failed: %w", err)
	}

	return nil
}

// RemoveUser deletes an OS user and their profile/home directory.
func RemoveUser(username string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("net", "user", username, "/DELETE")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("windows user remove failed: %w\nOutput: %s", err, out)
		}

	case "linux":
		// -r removes the home directory and mail spool
		cmd := exec.Command("userdel", "-r", username)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("linux userdel failed: %w\nOutput: %s", err, out)
		}

	default:
		return fmt.Errorf("unsupported operating system for user management: %s", runtime.GOOS)
	}

	return nil
}
