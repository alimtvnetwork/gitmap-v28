package cmdinstall

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

func buildInstallDatabaseSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Databases & Cache Servers",
		Entries: []termhelp.CommandEntry{
			{Command: "mysql", Description: "MySQL relational database server"},
			{Command: "mariadb", Description: "MariaDB SQL database (MySQL fork)"},
			{Command: "postgresql (pg)", Description: "PostgreSQL object-relational database"},
			{Command: "sqlite", Description: "SQLite embedded SQL engine"},
			{Command: "mongodb (mongo)", Description: "MongoDB NoSQL document database"},
			{Command: "redis", Description: "In-memory key-value store and cache"},
			{Command: "neo4j", Description: "Neo4j graph database platform"},
		},
	}
}

func buildInstallAISection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Local AI & Coding Suite",
		Entries: []termhelp.CommandEntry{
			{Command: "antigravity (agy)", Description: "Autonomous AI coding assistant CLI"},
			{Command: "ag-manager (agm)", Description: "Antigravity Manager GUI desktop app"},
			{Command: "ollama", Description: "Local LLM runner and model daemon"},
			{Command: "llama-cpp", Description: "High-performance local inference in C++"},
			{Command: "python-libs", Description: "Core AI stack (numpy, torch, transformers)"},
		},
	}
}

func buildInstallProfileSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Workstation Profiles",
		Entries: []termhelp.CommandEntry{
			{Command: "profile [name]", Description: "Install curated bundle (minimal, dev, ai...)", HasSubcommands: true},
			{Command: "profile --list", Description: "List all available profiles and tool counts"},
			{Command: "dev", Description: "Standard dev workstation with AI suite"},
			{Command: "minimal (min)", Description: "Minimal dev workstation (code, git, node, py)"},
			{Command: "ubuntu", Description: "Ubuntu developer workstation bundle"},
			{Command: "ai", Description: "AI/ML stack (Python, Ollama, torch, agy)"},
			{Command: "backend", Description: "Backend workstation (DBs, Docker, Go, .NET)"},
			{Command: "fullstack", Description: "Full-stack web developer workstation"},
		},
	}
}

func buildInstallCustomSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Custom Scripts & Archives",
		Entries: []termhelp.CommandEntry{
			{Command: "tar <file>", Description: "Install .tar.gz/.zip archive with auto-detect"},
			{Command: "scripts-fixer", Description: "GitMap scripts and path fixer suite"},
			{Command: "coding-guidelines (cg)", Description: "AlimTV coding guidelines compliance (v24)"},
			{Command: "macro-ahk", Description: "AutoHotkey v2 macro automation suite"},
		},
	}
}

func buildInstallMgmtSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Management, Logs & Config",
		Entries: []termhelp.CommandEntry{
			{Command: "list (ls)", Description: "List all supported tools with installed status"},
			{Command: "logs [tool]", Description: "View installer execution and error logs"},
			{Command: "add <name> <ver>", Description: "Interactively register custom installer"},
			{Command: "export <tool>", Description: "Export installer configs to YAML/JSON"},
			{Command: "import <file>", Description: "Import custom installer configurations"},
		},
	}
}
