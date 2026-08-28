package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/crypto"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

var sjCommand = "sj"

type sjOptions struct {
	ImportFile string
	ExportFile string
	List       bool
	Args       []string
}

func parseSJFlags(args []string) sjOptions {
	fs := flag.NewFlagSet(sjCommand, flag.ExitOnError)
	var opts sjOptions
	fs.StringVar(&opts.ImportFile, "import", "", "Import SSH connections from JSON file")
	fs.StringVar(&opts.ExportFile, "export", "", "Export SSH connections to JSON file")
	fs.BoolVar(&opts.List, "list", false, "List SSH connections")
	fs.Parse(reorderFlagsBeforeArgs(args))

	argsAfterParse := fs.Args()
	if len(argsAfterParse) > 0 && (argsAfterParse[0] == "ls" || argsAfterParse[0] == "list") {
		opts.List = true
	}
	opts.Args = argsAfterParse
	return opts
}

func runSSHJoinLegacy(args []string) error {
	opts := parseSJFlags(args)

	if opts.List {
		runSSHJoinLs()
		return nil
	}
	if opts.ImportFile != "" {
		runSSHJoinImport(opts.ImportFile)
		return nil
	}
	if opts.ExportFile != "" {
		runSSHJoinExport(opts.ExportFile)
		return nil
	}

	runSSHJoinInteractive(opts.Args)
	return nil
}

func runSSHJoinLs() error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		return nil
	}
	defer dbConn.Close()

	conns, err := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if err != nil {
		fmt.Printf("Failed to get connections: %v\n", err)
		return nil
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9")).Render("SSH Joined Machines:"))
	for _, c := range conns {
		fmt.Printf("  Alias: %-15s IP: %-15s User: %-10s OS: %-8s\n", c.Alias, c.IPAddress, c.Username, c.OS)
	}
	return nil
}

func getEncryptionKey() []byte {
	// 32-byte key for AES-256
	return []byte("gitmap-ssh-secret-key-0123456789")
}

func runSSHJoinInteractive(args []string) error {
	reader := bufio.NewReader(os.Stdin)

	alias := promptInput(reader, "Machine Alias")
	ip := promptInput(reader, "IP Address")
	user := promptInput(reader, "Username")
	password := promptInput(reader, "Password (or leave blank for SSH key)")
	keyPath := ""
	if password == "" {
		keyPath = promptInput(reader, "SSH Key Path")
	}
	osType := promptInput(reader, "OS (windows/unix)")

	var encPass string
	if password != "" {
		var err error
		encPass, err = crypto.Encrypt([]byte(password), getEncryptionKey())
		if err != nil {
			fmt.Printf("Encryption failed: %v\n", err)
			return nil
		}
	}

	saveSSHConnection(alias, ip, user, encPass, keyPath, osType)
	return nil
}

func promptInput(reader *bufio.Reader, prompt string) string {
	fmt.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func saveSSHConnection(alias, ip, user, encPass, keyPath, osType string) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		return
	}
	defer dbConn.Close()

	conn := db.SSHConnection{
		Alias:             alias,
		IPAddress:         ip,
		Username:          user,
		EncryptedPassword: encPass,
		KeyPath:           keyPath,
		OS:                osType,
		CreatedAt:         time.Now(),
	}

	if err := db.InsertOrUpdateSSHConnection(dbConn.Context(), dbConn.SQL(), conn); err != nil {
		fmt.Printf("Failed to save connection: %v\n", err)
		return
	}
	fmt.Printf("Successfully joined %s (%s)\n", alias, ip)
}

func runSSHJoinImport(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Read file error: %v\n", err)
		return nil
	}

	var conns []db.SSHConnection
	if err := json.Unmarshal(data, &conns); err != nil {
		fmt.Printf("JSON unmarshal error: %v\n", err)
		return nil
	}

	dbConn, err := store.OpenDefault()
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		return nil
	}
	defer dbConn.Close()

	for _, c := range conns {
		if err := db.InsertOrUpdateSSHConnection(dbConn.Context(), dbConn.SQL(), c); err != nil {
			fmt.Printf("Failed to import %s: %v\n", c.Alias, err)
		} else {
			fmt.Printf("Imported %s\n", c.Alias)
		}
	}
	return nil
}

func runSSHJoinExport(file string) error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		return nil
	}
	defer dbConn.Close()

	conns, err := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if err != nil {
		fmt.Printf("Failed to get connections: %v\n", err)
		return nil
	}

	data, err := json.MarshalIndent(conns, "", "  ")
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
		return nil
	}

	if err := os.WriteFile(file, data, 0600); err != nil {
		fmt.Printf("Write file error: %v\n", err)
		return nil
	}
	fmt.Printf("Exported %d connections to %s\n", len(conns), file)
	return nil
}
