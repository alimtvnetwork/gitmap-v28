package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type SSHTarget struct {
	Username string
	IP       string
	Port     int
}

func (t *SSHTarget) String() string {
	return fmt.Sprintf("%s@%s", t.Username, t.IP)
}

func ParseSSHTarget(raw string, defaultUser string, defaultPort int) (*SSHTarget, error) {
	if raw == "" {
		return nil, &apperror.AppError{
			Op:    "ParseSSHTarget",
			Code:  "E_INTERNAL_ERROR",
			Ctx:   map[string]any{"raw": raw},
			Cause: fmt.Errorf("empty target"),
		}
	}

	parts := strings.Split(raw, "@")
	var user, ip string

	if len(parts) == 1 {
		user = defaultUser
		ip = parts[0]
	} else if len(parts) == 2 {
		p1, p2 := parts[0], parts[1]
		if net.ParseIP(p1) != nil && net.ParseIP(p2) == nil {
			ip = p1
			user = p2
		} else {
			user = p1
			ip = p2
		}
	} else {
		return nil, &apperror.AppError{
			Op:    "ParseSSHTarget",
			Code:  "E_INTERNAL_ERROR",
			Ctx:   map[string]any{"raw": raw},
			Cause: fmt.Errorf("invalid target format"),
		}
	}

	return &SSHTarget{
		Username: user,
		IP:       ip,
		Port:     defaultPort,
	}, nil
}
