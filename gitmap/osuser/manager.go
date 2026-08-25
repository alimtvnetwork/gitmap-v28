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
		var cmd *exec.Cmd
		if password != "" {
			cmd = exec.Command("net", "user", username, password, "/ADD")
		} else {
			cmd = exec.Command("net", "user", username, "/ADD")
		}
		
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("windows user add failed: %v\nOutput: %s", err, out)
		}
		
	case "linux":
		// useradd is the standard low-level utility across Debian, Ubuntu, and Fedora.
		cmd := exec.Command("useradd", "-m", username)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("linux useradd failed: %v\nOutput: %s", err, out)
		}
		
		if password != "" {
			// use chpasswd to set the password non-interactively
			chCmd := exec.Command("chpasswd")
			stdin, err := chCmd.StdinPipe()
			if err != nil {
				return fmt.Errorf("failed to pipe chpasswd: %v", err)
			}
			
			if err := chCmd.Start(); err != nil {
				return fmt.Errorf("failed to start chpasswd: %v", err)
			}
			
			// write username:password format
			if _, err := stdin.Write([]byte(fmt.Sprintf("%s:%s", username, password))); err != nil {
				return fmt.Errorf("failed to write to chpasswd: %v", err)
			}
			stdin.Close()
			
			if err := chCmd.Wait(); err != nil {
				return fmt.Errorf("chpasswd failed: %v", err)
			}
		}

	default:
		return fmt.Errorf("unsupported operating system for user management: %s", runtime.GOOS)
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
			return fmt.Errorf("windows user remove failed: %v\nOutput: %s", err, out)
		}
		
	case "linux":
		// -r removes the home directory and mail spool
		cmd := exec.Command("userdel", "-r", username)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("linux userdel failed: %v\nOutput: %s", err, out)
		}

	default:
		return fmt.Errorf("unsupported operating system for user management: %s", runtime.GOOS)
	}
	
	return nil
}
