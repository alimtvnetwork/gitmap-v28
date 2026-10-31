package constants

import "time"

// Command and subcommand names for cmdlogin.
const (
	CmdLogin               = "login"
	CmdLogout              = "logout"
	SubCmdLoginSSH         = "ssh"
	SubCmdLoginBrowser     = "browser"
	SubCmdLoginBrowserFlow = "browser-flow"
	SubCmdLoginToken       = "token"
	SubCmdLoginPAT         = "pat"
	SubCmdLoginStatus      = "status"
	SubCmdLoginSet         = "set"
)

// Command flag identifiers and arguments for cmdlogin.
const (
	FlagGlobalArg      = "--global"
	FlagGlobalShortArg = "-g"
	FlagHelp           = "help"
	FlagHelpArg        = "--help"
	FlagHelpShortArg   = "-h"
	FlagStatus         = "status"
	FlagStatusArg      = "--status"
	FlagWeb            = "web"
	FlagWebArg         = "--web"
	FlagBrowser        = "browser"
	FlagBrowserArg     = "--browser"
	FlagNoVerify       = "no-verify"
	FlagNoVerifyArg    = "--no-verify"
	FlagToken          = "token"
	FlagTokenArg       = "--token"
	FlagTokenPrefix    = "--token="
	FlagYes            = "yes"
	FlagYesArg         = "--yes"
	FlagYesShortArg    = "-y"
	FlagForce          = "force"
	FlagForceArg       = "--force"
	FlagForceShortArg  = "-f"
	FlagDir            = "dir"
	FlagDirArg         = "--dir"
	FlagDirShortArg    = "-d"
	FlagDirPrefix      = "--dir="
)

// Git configuration keys and credential commands.
const (
	GitBin                     = "git"
	GitRevParse                = "rev-parse"
	GitInsideWorkTree          = "--is-inside-work-tree"
	GitConfig                  = "config"
	GitConfigGlobalFlag        = "--global"
	GitConfigUnsetFlag         = "--unset"
	GitConfigGithubToken       = "github.token"
	GitConfigCoreSSHCommand    = "core.sshCommand"
	GitCredential              = "credential"
	GitCredentialReject        = "reject"
	GitCredentialRejectPayload = "protocol=https\nhost=github.com\n\n"
)

// Remote service URLs and endpoints.
const (
	GithubUserAPI      = "https://api.github.com/user"
	GithubNewTokenURL  = "https://github.com/settings/tokens/new"
	GithubNewSSHKeyURL = "https://github.com/settings/ssh/new"
	GitlabNewSSHKeyURL = "https://gitlab.com/-/user_settings/ssh_keys"
)

// HTTP headers and media types for API requests.
const (
	HeaderAuthorization             = "Authorization"
	HeaderAuthorizationBearerPrefix = "Bearer "
	HeaderAccept                    = "Accept"
	HeaderGithubJSON                = "application/vnd.github+json"
	TokenValidationTimeout          = 15 * time.Second
)

// External tools, command arguments, and platform strings.
const (
	BinGh                 = "gh"
	GhCmdAuth             = "auth"
	GhCmdLogin            = "login"
	GhHostFlag            = "-h"
	GhHostDomain          = "github.com"
	GhWebFlag             = "-w"
	BrowserOpenDarwin     = "open"
	BrowserOpenWindows    = "rundll32"
	BrowserWindowsDllArg  = "url.dll,FileProtocolHandler"
	BrowserOpenLinux      = "xdg-open"
	PlatformDarwin        = "darwin"
	PlatformWindows       = "windows"
	BinSshKeygen          = "ssh-keygen"
	SshKeygenTypeFlag     = "-t"
	SshKeygenCommentFlag  = "-C"
	SshKeygenFileFlag     = "-f"
	SshKeygenPassFlag     = "-N"
	DefaultSSHKeyType     = "ed25519"
	DefaultSSHKeyBaseName = "id_ed25519"
	DefaultSSHKeyPrefix   = "id_ed25519_"
	DefaultSSHComment     = "gitmap"
	DefaultSSHDir         = ".ssh"
	PubExt                = ".pub"
	SSHCommandFormat      = "ssh -i %s -F /dev/null -o IdentitiesOnly=yes"
)

// System clipboard tools and arguments.
const (
	ClipWindows            = "clip"
	ClipDarwin             = "pbcopy"
	ClipWlCopy             = "wl-copy"
	ClipXclip              = "xclip"
	ClipXsel               = "xsel"
	ClipboardSelectionFlag = "-selection"
	ClipboardClipboardArg  = "clipboard"
	ClipboardFlag          = "--clipboard"
	ClipboardInputFlag     = "--input"
)

// GitProfile defaults.
const (
	DefaultProfileName    = "default"
	DefaultProfileType    = "user"
	ProfileProviderGitHub = "github"
	ProfileAuthMethodSSH  = "ssh"
	TokenMaskFull         = "********"
	TokenScopeRepo        = "repo"
	DefaultAIAutoBotName  = "gitmap automation bot - runs autonomously"
	DefaultAIAutoBotEmail = "info@rgitmap.com"
	DefaultAIAutoBotAlias = "autobot"
)

// UI messages, labels, and dividers.
const (
	BannerDivider = "================================================================================"
)
