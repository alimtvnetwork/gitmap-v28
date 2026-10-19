package constants

// Supported tool names — Databases.
const (
	ToolMySQL         = "mysql"
	ToolMariaDB       = "mariadb"
	ToolPostgreSQL    = "postgresql"
	ToolSQLite        = "sqlite"
	ToolMongoDB       = "mongodb"
	ToolCouchDB       = "couchdb"
	ToolRedis         = "redis"
	ToolCassandra     = "cassandra"
	ToolNeo4j         = "neo4j"
	ToolElasticsearch = "elasticsearch"
	ToolDuckDB        = "duckdb"
	ToolSupabase      = "supabase"
)

// Supported tool names — Languages & Runtimes.
const (
	ToolRust    = "rust"
	ToolCargo   = "cargo"
	ToolDotnet  = "dotnet"
	ToolJava    = "java"
	ToolFlutter = "flutter"
	ToolLaravel = "laravel"
)

// Supported tool names — Local AI.
const (
	ToolOllama     = "ollama"
	ToolLlamaCpp   = "llama-cpp"
	ToolPythonLibs = "python-libs"
	ToolMuse       = "muse"
	ToolMetaMuse   = "meta-muse"
)

// Meta Muse official distribution endpoints and binary names.
const (
	MetaMuseInstallUrlWindows = "https://dev.meta.ai/install.ps1"
	MetaMuseInstallUrlUnix    = "https://dev.meta.ai/install.sh"
	MetaMuseBinaryNameWindows = "muse.exe"
	MetaMuseBinaryNameUnix    = "muse"
)

// Supported tool names — DevOps & Containers.
const (
	ToolDocker      = "docker"
	ToolKubernetes  = "kubernetes"
	ToolJenkins     = "jenkins"
	ToolNginx       = "nginx"
	ToolVMware      = "vmware"
	ToolOpenVmTools = "open-vm-tools"
)

// Supported tool names — Terminal & Utilities.
const (
	ToolZsh         = "zsh"
	ToolFlameshot   = "flameshot"
	ToolConemu      = "conemu"
	ToolVLC         = "vlc"
	ToolQBittorrent = "qbittorrent"
	ToolUTorrent    = "utorrent"
	ToolJq          = "jq"
	ToolYq          = "yq"
	ToolZellij      = "zellij"
)

// Supported tool names — Additional & Profiles Parity Tools.
const (
	ToolUbuntuFont      = "ubuntu-font"
	ToolWhatsApp        = "whatsapp"
	ToolOneNote         = "onenote"
	ToolLightshot       = "lightshot"
	ToolWindowsTerminal = "windows-terminal"
	ToolAria2           = "aria2"
	Tool7Zip            = "7zip"
	ToolWinRAR          = "winrar"
	ToolXMind           = "xmind"
	ToolWordWeb         = "wordweb"
	ToolBeyondCompare   = "beyondcompare"
	ToolVcRedist        = "vcredist"
	ToolDirectX         = "directx"
	ToolDirectXSdk      = "directx-sdk"
	ToolStarship        = "starship"
	ToolOhMyPosh        = "oh-my-posh"
	ToolScoop           = "scoop"
	ToolTortoiseGit     = "tortoisegit"
)

// Tool categories.
const (
	ToolCategoryCore      = "Core Tools"
	ToolCategoryDatabase  = "Databases"
	ToolCategoryLanguages = "Languages & Runtimes"
	ToolCategoryAI        = "Local AI"
	ToolCategoryDevOps    = "DevOps & Containers"
	ToolCategoryUtilities = "Terminal & Utilities"
)

// Tool display names for --list output.
var InstallToolDescriptions = map[string]string{
	ToolVSCode:           "Visual Studio Code editor",
	ToolNodeJS:           "Node.js JavaScript runtime",
	ToolYarn:             "Yarn package manager",
	ToolBun:              "Bun JavaScript runtime",
	ToolPnpm:             "pnpm package manager",
	ToolPython:           "Python programming language",
	ToolGo:               "Go programming language",
	ToolGit:              "Git version control",
	ToolGitLFS:           "Git Large File Storage",
	ToolGHCLI:            "GitHub CLI",
	ToolGitHubDesktop:    "GitHub Desktop application",
	ToolCPP:              "C++ compiler (MinGW/g++)",
	ToolPHP:              "PHP programming language",
	ToolPowerShell:       "PowerShell shell",
	ToolChocolatey:       "Chocolatey package manager",
	ToolWinget:           "Winget package manager",
	ToolMySQL:            "MySQL relational database",
	ToolMariaDB:          "MariaDB (MySQL-compatible fork)",
	ToolPostgreSQL:       "PostgreSQL relational database",
	ToolSQLite:           "SQLite embedded database",
	ToolMongoDB:          "MongoDB document database",
	ToolCouchDB:          "CouchDB document database (REST API)",
	ToolRedis:            "Redis in-memory key-value store",
	ToolCassandra:        "Apache Cassandra wide-column NoSQL",
	ToolNeo4j:            "Neo4j graph database",
	ToolElasticsearch:    "Elasticsearch search and analytics",
	ToolDuckDB:           "DuckDB analytical columnar database",
	ToolSupabase:         "[?] Supabase open-source backend & database platform (Planned / To-Do)",
	ToolDbeaver:          "DBeaver database management tool",
	ToolStickyNotes:      "Microsoft Sticky Notes",
	ToolLiteDB:           "LiteDB embedded NoSQL database for .NET",
	ToolOBS:              "OBS Studio screen recorder and streamer",
	ToolVSCodeCtx:        "Add VS Code to Windows right-click context menu",
	ToolPwshCtx:          "Add PowerShell to Windows right-click context menu",
	ToolCtx:              ToolCtxDesc,
	ToolAllDevTools:      "Install all core developer tools at once",
	ToolNpp:              "NPP + Settings -- Notepad++ with settings",
	ToolNppSettings:      "NPP Settings -- Notepad++ settings sync only",
	ToolNppInstall:       "Install NPP -- Notepad++ install only (no settings)",
	ToolVSCodeSync:       "VS Code Settings -- sync VS Code settings and extensions",
	ToolOBSSync:          "OBS Settings -- sync OBS Studio profiles and scenes",
	ToolWTSync:           "WT Settings -- sync Windows Terminal settings.json",
	ToolScripts:          "Clone gitmap scripts to local folder",
	ToolGitmapOneliner:   "Print the Windows + macOS install-gitmap one-liners",
	ToolScriptsFixer:     "Install Scripts Fixer tool",
	ToolCodingGuidelines: "Install Coding Guidelines per repo",
	ToolMacroAhk:         "Install Macro AHK scripts",
	ToolAgManager:        "Antigravity Manager GUI and tools (latest GitHub release)",
	ToolAgCtx:            "Add Antigravity to Windows right-click context menu",
	ToolAntigravity:      "Google Antigravity Desktop IDE (AI-powered development environment)",
	ToolAgy:              "Antigravity CLI autonomous coding assistant (agy-cli)",
	ToolBuildEssential:   "Ubuntu build-essential compiler toolchain and common dev packages",
	ToolUbuntuCommon:     "Ubuntu common developer profile (compiler toolchain, editors, utilities)",
	ToolChrome:           "Google Chrome web browser",
	ToolGoogleChrome:     "Google Chrome web browser",
	ToolRust:             "Rust programming language and compiler toolchain",
	ToolCargo:            "Cargo package manager and Rust build system",
	ToolDotnet:           ".NET SDK and developer platform",
	ToolJava:             "Java OpenJDK development kit and runtime",
	ToolFlutter:          "Flutter cross-platform application SDK",
	ToolOllama:           "Ollama local large language model runner",
	ToolLlamaCpp:         "llama.cpp LLM inference engine in C/C++",
	ToolPythonLibs:       "Python AI, ML, data science, and web libraries",
	ToolDocker:           "Docker container platform and runtime",
	ToolKubernetes:       "Kubernetes container orchestration CLI (kubectl)",
	ToolJenkins:          "Jenkins automation and CI/CD server",
	ToolZsh:              "Zsh interactive shell and command interpreter",
	ToolFlameshot:        "Flameshot screenshot capture and annotation tool",
	ToolConemu:           "ConEmu Windows console emulator with tabs and splits",
	ToolVLC:              "VLC media player cross-platform multimedia player",
	ToolNginx:            "Nginx high-performance HTTP server and reverse proxy",
	ToolWordPress:        "WordPress web publishing platform and CMS / WP-CLI",
	ToolLaravel:          "Laravel PHP web application framework and installer CLI",
	ToolVMware:           "VMware guest tools and virtualization utilities",
	ToolGitmap:           "Gitmap CLI repository orchestration and dev tools manager",
	ToolComposer:         "Composer dependency manager for PHP",
	ToolWpCli:            "WP-CLI command-line interface for WordPress management",
	ToolOpenVmTools:      "Open Virtual Machine Tools for VMware guest optimization",
	ToolQBittorrent:      "qBittorrent free and open-source BitTorrent client",
	ToolUTorrent:         "uTorrent lightweight BitTorrent client",
	ToolUbuntuFont:       "Ubuntu font family for system and terminals",
	ToolWhatsApp:         "WhatsApp desktop messaging client",
	ToolOneNote:          "Microsoft OneNote desktop note-taking application",
	ToolLightshot:        "Lightshot screen capture tool",
	ToolWindowsTerminal:  "Windows Terminal modern tabbed console emulator",
	ToolAria2:            "aria2 lightweight multi-protocol download utility",
	Tool7Zip:             "7-Zip file archiver with high compression ratio",
	ToolWinRAR:           "WinRAR file compression and archive manager",
	ToolXMind:            "XMind mind mapping and brainstorming tool",
	ToolWordWeb:          "WordWeb international dictionary and thesaurus",
	ToolBeyondCompare:    "Beyond Compare data and directory comparison tool",
	ToolTortoiseGit:      "TortoiseGit Windows Explorer interface to Git",
	ToolVcRedist:         "Visual C++ Redistributable runtime libraries",
	ToolDirectX:          "DirectX end-user runtimes",
	ToolDirectXSdk:       "DirectX Software Development Kit",
	ToolStarship:         "Starship cross-shell customizable prompt",
	ToolOhMyPosh:         "Oh My Posh prompt theme engine for any shell",
	ToolScoop:            "Scoop command-line package installer for Windows",
	ToolGitCompact:       "git-compact repository compact tool and cleaner",
	ToolJq:               "jq command-line JSON processor",
	ToolYq:               "yq command-line YAML/JSON/XML processor",
	ToolZellij:           "Zellij terminal workspace and multiplexer",
	ToolMuse:             "Meta Muse autonomous AI developer agent",
	ToolMetaMuse:         "Meta Muse autonomous AI developer agent (alias)",
}

// InstallToolCategories groups tools by category for display.
var InstallToolCategories = map[string][]string{
	ToolCategoryCore: {
		ToolGitmap, ToolVSCode, ToolNodeJS, ToolYarn, ToolBun, ToolPnpm,
		ToolPython, ToolGo, ToolGit, ToolGitLFS, ToolGHCLI,
		ToolGitHubDesktop, ToolGitCompact, ToolCPP, ToolPHP, ToolPowerShell,
		ToolChocolatey, ToolWinget, ToolScoop, ToolDbeaver, ToolOBS,
		ToolStickyNotes, ToolVSCodeCtx, ToolPwshCtx, ToolCtx,
		ToolNpp, ToolNppSettings, ToolNppInstall,
		ToolVSCodeSync, ToolOBSSync, ToolWTSync,
		ToolScripts, ToolGitmapOneliner, ToolAllDevTools,
		ToolScriptsFixer, ToolCodingGuidelines, ToolMacroAhk, ToolAgManager, ToolAgCtx, ToolAntigravity, ToolAgy,
		ToolBuildEssential, ToolUbuntuCommon, ToolChrome, ToolGoogleChrome,
		ToolWordPress, ToolWpCli,
	},
	ToolCategoryDatabase: {
		ToolMySQL, ToolMariaDB, ToolPostgreSQL, ToolSQLite,
		ToolMongoDB, ToolCouchDB, ToolRedis, ToolCassandra,
		ToolNeo4j, ToolElasticsearch, ToolDuckDB, ToolLiteDB, ToolSupabase,
	},
	ToolCategoryLanguages: {
		ToolRust, ToolCargo, ToolDotnet, ToolJava, ToolFlutter, ToolLaravel, ToolComposer,
	},
	ToolCategoryAI: {
		ToolOllama, ToolLlamaCpp, ToolPythonLibs, ToolMuse, ToolMetaMuse,
	},
	ToolCategoryDevOps: {
		ToolDocker, ToolKubernetes, ToolJenkins, ToolNginx, ToolVMware, ToolOpenVmTools,
	},
	ToolCategoryUtilities: {
		ToolZsh, ToolFlameshot, ToolConemu, ToolVLC, ToolQBittorrent, ToolUTorrent,
		ToolUbuntuFont, ToolWhatsApp, ToolOneNote, ToolLightshot, ToolWindowsTerminal,
		ToolAria2, Tool7Zip, ToolWinRAR, ToolXMind, ToolWordWeb, ToolBeyondCompare,
		ToolTortoiseGit, ToolVcRedist, ToolDirectX, ToolDirectXSdk, ToolStarship, ToolOhMyPosh,
		ToolJq, ToolYq, ToolZellij,
	},
}
