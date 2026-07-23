package gitutil

import (
	"fmt"
	"net"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// CheckOnline verifies network connectivity by dialing a remote host.
func CheckOnline() error {
	conn, err := net.DialTimeout(
		constants.NetworkProto,
		constants.NetworkCheckHost,
		time.Duration(constants.NetworkTimeoutSec)*time.Second,
	)
	if err != nil {
		return fmt.Errorf(constants.ErrOffline)
	}

	conn.Close()

	return nil
}

// IsOnline returns true if network connectivity is available.
func IsOnline() bool {
	err := CheckOnline()

	return err == nil
}

// PrintOfflineWarning prints a user-friendly offline message.
func PrintOfflineWarning() {
	fmt.Print(constants.MsgOfflineWarning)
	fmt.Print(constants.MsgOfflineHint)
}
