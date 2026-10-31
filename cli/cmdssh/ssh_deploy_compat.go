// Package cmdssh — ssh_deploy_compat.go handles binary format detection and OS compatibility filtering.
package cmdssh

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// BinaryFormatType classifies executable binary machine architectures.
type BinaryFormatType string

const (
	BinaryFormatPE      BinaryFormatType = "PE"
	BinaryFormatELF     BinaryFormatType = "ELF"
	BinaryFormatMachO   BinaryFormatType = "Mach-O"
	BinaryFormatScript  BinaryFormatType = "Script"
	BinaryFormatUnknown BinaryFormatType = "Unknown"
)

const (
	SkipReasonOSIncompatible = "SkipReasonOSIncompatible"
	SkipReasonNodeOffline    = "SkipReasonNodeOffline"
	SkipReasonDryRun         = "SkipReasonDryRun"
)

// NodeLivenessResult stores reachability results from concurrent TCP probes.
type NodeLivenessResult struct {
	IsOnline bool
	Reason   string
}

// DetectBinaryFormat reads magic bytes to classify the candidate executable binary format.
func DetectBinaryFormat(filePath string) (BinaryFormatType, *apperror.AppError) {
	header, appErr := readBinaryHeader(filePath)
	if appErr != nil {
		return BinaryFormatUnknown, appErr
	}

	format := parseMagicBytes(header)

	return format, nil
}

func readBinaryHeader(filePath string) ([]byte, *apperror.AppError) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open binary for detection")
	}

	defer file.Close()
	buf := make([]byte, 16)
	numBytes, _ := file.Read(buf)

	return buf[:numBytes], nil
}

func parseMagicBytes(header []byte) BinaryFormatType {
	if isPEMagic(header) {
		return BinaryFormatPE
	}

	if isELFMagic(header) {
		return BinaryFormatELF
	}

	if isMachOMagic(header) {
		return BinaryFormatMachO
	}

	return parseScriptOrUnknown(header)
}

func isPEMagic(h []byte) bool {
	hasMZ := len(h) >= 2 && h[0] == 0x4D && h[1] == 0x5A

	return hasMZ
}

func isELFMagic(h []byte) bool {
	hasELF := len(h) >= 4 && h[0] == 0x7F && h[1] == 'E' && h[2] == 'L' && h[3] == 'F'

	return hasELF
}

func isMachOMagic(h []byte) bool {
	if len(h) < 4 {
		return false
	}

	is32 := (h[0] == 0xfe && h[1] == 0xed && h[2] == 0xfa && h[3] == 0xce) || (h[0] == 0xce && h[1] == 0xfa && h[2] == 0xed && h[3] == 0xfe)
	is64 := (h[0] == 0xfe && h[1] == 0xed && h[2] == 0xfa && h[3] == 0xcf) || (h[0] == 0xcf && h[1] == 0xfa && h[2] == 0xed && h[3] == 0xfe)
	isFat := (h[0] == 0xca && h[1] == 0xfe && h[2] == 0xba && h[3] == 0xbe) || (h[0] == 0xbe && h[1] == 0xba && h[2] == 0xfe && h[3] == 0xca)

	return is32 || is64 || isFat
}

func parseScriptOrUnknown(h []byte) BinaryFormatType {
	if len(h) >= 2 && h[0] == '#' && h[1] == '!' {
		return BinaryFormatScript
	}

	return BinaryFormatUnknown
}

// ResolveTargetNodeOS resolves OS metadata from connection record or active SSH probe.
func ResolveTargetNodeOS(node db.SSHConnection, client *ssh.Client) string {
	targetOS := strings.TrimSpace(node.OS)
	if targetOS != "" {
		return targetOS
	}

	if client != nil {
		return probeRemoteOSType(client)
	}

	return "linux"
}

func isWindowsTarget(targetOS string) bool {
	lower := strings.ToLower(strings.TrimSpace(targetOS))
	if strings.Contains(lower, "darwin") {
		return false
	}

	return strings.Contains(lower, "win")
}

func isDarwinTarget(targetOS string) bool {
	lower := strings.ToLower(strings.TrimSpace(targetOS))
	isMac := strings.Contains(lower, "darwin") || strings.Contains(lower, "macos") || strings.Contains(lower, "osx")

	return isMac
}

// CheckBinaryCompatibility verifies whether candidate binary format can execute on target node OS.
func CheckBinaryCompatibility(format BinaryFormatType, targetOS string) (bool, string) {
	if format == BinaryFormatPE {
		return checkPECompatibility(targetOS)
	}

	if format == BinaryFormatELF {
		return checkELFCompatibility(targetOS)
	}

	if format == BinaryFormatMachO {
		return checkMachOCompatibility(targetOS)
	}

	return true, ""
}

func checkPECompatibility(targetOS string) (bool, string) {
	isWin := isWindowsTarget(targetOS)
	if isWin {
		return true, ""
	}

	return false, fmt.Sprintf("binary is Windows PE (.exe) but target is %s", targetOS)
}

func checkELFCompatibility(targetOS string) (bool, string) {
	isWin := isWindowsTarget(targetOS)
	if isWin {
		return false, fmt.Sprintf("binary is Linux ELF but target is %s", targetOS)
	}

	return true, ""
}

func checkMachOCompatibility(targetOS string) (bool, string) {
	isMac := isDarwinTarget(targetOS)
	if isMac {
		return true, ""
	}

	return false, fmt.Sprintf("binary is macOS Mach-O but target is %s", targetOS)
}

// FormatIncompatibleRationale generates clear skip diagnostics for mismatched binary/OS combinations.
func FormatIncompatibleRationale(alias, targetOS string, format BinaryFormatType) string {
	nodeDesc := alias
	if targetOS != "" {
		nodeDesc = fmt.Sprintf("%s (%s)", alias, targetOS)
	}

	return fmt.Sprintf("Skipped %s node: binary is %s", nodeDesc, FormatBinaryTypeName(format))
}

// FormatBinaryTypeName returns human-friendly binary format description.
func FormatBinaryTypeName(format BinaryFormatType) string {
	if format == BinaryFormatPE {
		return "Windows PE (.exe)"
	}

	if format == BinaryFormatELF {
		return "Linux ELF"
	}

	if format == BinaryFormatMachO {
		return "macOS Mach-O"
	}

	return string(format)
}

// ProbeNodesLivenessConcurrent performs non-blocking concurrent TCP probes across target nodes.
func ProbeNodesLivenessConcurrent(nodes []db.SSHConnection, timeout time.Duration) map[string]NodeLivenessResult {
	results := make(map[string]NodeLivenessResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, node := range nodes {
		wg.Add(1)
		go probeNodeWorker(node, timeout, &mu, &wg, results)
	}

	wg.Wait()

	return results
}

func probeNodeWorker(node db.SSHConnection, timeout time.Duration, mu *sync.Mutex, wg *sync.WaitGroup, results map[string]NodeLivenessResult) {
	defer wg.Done()
	isOnline, reason := probeNodeTCP(node.IPAddress, 22, timeout)
	mu.Lock()
	results[node.Alias] = NodeLivenessResult{IsOnline: isOnline, Reason: reason}
	results[node.IPAddress] = NodeLivenessResult{IsOnline: isOnline, Reason: reason}
	mu.Unlock()
}

func probeNodeTCP(ip string, port int, timeout time.Duration) (bool, string) {
	if timeout <= 0 {
		timeout = 1200 * time.Millisecond
	}

	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false, err.Error()
	}

	_ = conn.Close()

	return true, "online"
}
