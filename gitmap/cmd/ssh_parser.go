package cmd

import "fmt"

type SSHTarget struct {
	Username string
	IP       string
	Port     int
}

func (t *SSHTarget) String() string {
	return fmt.Sprintf("%s@%s", t.Username, t.IP)
}
