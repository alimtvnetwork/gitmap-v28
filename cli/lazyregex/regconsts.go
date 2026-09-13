package lazyregex

// Canonical regular expression pattern constants harvested from core and gitmap.
const (
	Alpha                                  = "^[a-zA-Z]+$"
	AlphaDash                              = "^[a-zA-Z0-9_\\-]+$"
	AlphaSpace                             = "^[-a-zA-Z0-9_ ]+$"
	AlphaNumeric                           = "^[a-zA-Z0-9]+$"
	CreditCard                             = "^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$"
	Coordinate                             = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?),\\s*[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	CSSColor                               = "^(#([\\da-f]{3}){1,2}|(rgb|hsl)a\\((\\d{1,3}%?,\\s?){3}(1|0?\\.\\d+)\\)|(rgb|hsl)\\(\\d{1,3}%?(,\\s?\\d{1,3}%?){2}\\))$"
	Date                                   = "^(((19|20)([2468][048]|[13579][26]|0[48])|2000)[/-]02[/-]29|((19|20)[0-9]{2}[/-](0[469]|11)[/-](0[1-9]|[12][0-9]|30)|(19|20)[0-9]{2}[/-](0[13578]|1[02])[/-](0[1-9]|[12][0-9]|3[01])|(19|20)[0-9]{2}[/-]02[/-](0[1-9]|1[0-9]|2[0-8])))$"
	DateDDMMYY                             = "^(0?[1-9]|[12][0-9]|3[01])[\\/\\-](0?[1-9]|1[012])[\\/\\-]\\d{4}$"
	MacAddress                             = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
	Numeric                                = "^[-+]?[0-9]+$"
	IntegerOnly                            = "^[0-9]+"
	SignedOrInteger                        = "^[-+]?[0-9]+"
	Path                                   = `[/\\]{0,2}[_-]*[Aa0-zZ9. ]+[_-]*[Aa0-zZ9. ]*(?:\!\@\#\$\%\&)?[/\\]?|^/$`
	IpSimple                               = `(\d{1,3}\.){3}\d{1,3}[^\r\n]*`
	FloatingPointNumberOnly                = `^[-+]?[0-9]+[\.][0-9]+`
	SignOrIntegerOrFloatingPointNumber     = `^[-+]?[0-9]+[\.]?[0-9]+`
	AnyIdentifier                          = `[a-zA-Z][a-zA-Z_\d\-]*`
	StringDoubleQuote                      = `"(?:\\.|[^"])*"`
	StringSingleQuote                      = `'((?:\\.|[^'])*'|"(?:\\.|[^"])*")`
	StringSingleOrDoubleQuote              = `'((?:\\.|[^'])*'|"(?:\\.|[^"])*")`
	Url                                    = "^(?:http(s)?:\\/\\/)?[\\w.-]+(?:\\.[\\w\\.-]+)+[\\w\\-\\._~:/?#[\\]@!\\$&'\\(\\)\\*\\+,;=.]+$"
	UUID                                   = "^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}$"
	UUIDAny                                = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	UUID3                                  = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	UUID4                                  = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID5                                  = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	ISBN10                                 = "^(?:[0-9]{9}X|[0-9]{10})$"
	ISBN13                                 = "^(?:[0-9]{13})$"
	Alphanumeric                           = "^[a-zA-Z0-9]+$"
	Int                                    = `^[-+]?(?:[1-9][0-9]*)$`
	Hexadecimal                            = "^[0-9a-fA-F]+$"
	Hexcolor                               = "^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
	RGBcolor                               = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	ASCII                                  = "^[\x00-\x7F]+$"
	Multibyte                              = "[^\x00-\x7F]"
	FullWidth                              = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	HalfWidth                              = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	Base64                                 = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	PrintableASCII                         = "^[\x20-\x7E]+$"
	DataUri                                = "^data:.+\\/(.+);base64$"
	MagnetUri                              = "^magnet:\\?xt=urn:[a-zA-Z0-9]+:[a-zA-Z0-9]{32,40}&dn=.+&tr=.+$"
	DNSName                                = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	UrlSchema                              = `((ftp|tcp|udp|wss?|https?):\/\/)`
	UrlUsername                            = `(\S+(:\S*)?@)`
	UrlPath                                = `((\/|\?|#)[^\s]*)`
	UrlPort                                = `(:(\d{1,5}))`
	UrlIP                                  = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	UrlSubdomain                           = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	SSN                                    = `^\d{3}[- ]?\d{2}[- ]?\d{4}$`
	WinPathWithOrWithoutUnc                = `(\\\\\?\\)?[a-zA-Z]:(\\|/)[_-]*[Aa0-zZ9. ]+[_-]*[Aa0-zZ9. ]*(?:\!\@\#\$\%\&)?(\\|/)?[_-]*[Aa0-zZ9. ]+[_-]*[Aa0-zZ9. ]*`
	UnixPathWithoutSpace                   = `/[Aa0-zZ9.]+[_-]*[Aa0-zZ9.]*|^/$`
	Semver                                 = "^v?(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$"
	HasLowerCase                           = ".*[[:lower:]]"
	HasUpperCase                           = ".*[[:upper:]]"
	HasWhitespace                          = ".*[[:space:]]"
	HasWhitespaceOnly                      = "^[[:space:]]+$"
	IMEI                                   = "^[0-9a-f]{14}$|^\\d{15}$|^\\d{18}$"
	Identifier                             = AnyIdentifier
	UnderscoreIdentifier                   = `_*` + Identifier
	ThirdBracketHeader                     = `\[[a-zA-Z][a-zA-Z_\d\-]*\]`
	IMSI                                   = "^\\d{14,15}$"
	HashComment                            = "#[^\\n]*"
	SlashComment                           = `\\[^\n]*`
	HashSemicolonComment                   = "[#;][^\\n]*"
	HashCommentWithSpaceOptional           = "\\s*#[^\\n]*"
	HashOrSemicolonComment                 = "\\s*[#;][^\\n]*"
	AllWhitespaces                         = "\\s+"
	AllWhitespacesOrPipe                   = "\\s+|\\|+"
	EachWordsWithDollarSymbolDefinition    = `(\$\{` + AnyIdentifier + `\}|\$` + AnyIdentifier + `)+`
	EachWordsWithinPercentSymbolDefinition = `(\%\{` + AnyIdentifier + `\}|\%` + AnyIdentifier + `)+`
	PrettyName                             = `^PRETTY_NAME=(.*)$`
	ExactIdFieldMatching                   = `^ID=(.*)$`
	ExactVersionIdFieldMatching            = `^VERSION_ID=(.*)$`
	UbuntuNameChecker                      = `[\( ]([\d\.]+)`
	CentOsNameChecker                      = `^CentOS( Linux)? release ([\d\.]+) `
	RedHatNameChecker                      = `[\( ]([\d\.]+)`
	FirstNumberAnyWhere                    = `(\d+){1}`
	WindowsVersionNumberChecker            = FirstNumberAnyWhere
	IpEthernet                             = `(\d{1,3}\.){3}\d{1,3}`
	IpWithSubnet                           = `(\d{1,3}\.){3}\d{1,3}\/\d{1,2}`
	Ipv6WithSubnet                         = `(([0-9a-fA-F]{0,4}:){1,7}[0-9a-fA-F]{0,4})`
	IpWithPort                             = `(\d{1,3}\.){3}\d{1,3}:\d{1,6}`
	ValidMac                               = `^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`
	SemiColonSeparatedEqualKeyValSplitter  = `^(.+=(.*);)+$`
	CommaSeparatedValue                    = `^(.+)(,\s*.+)+$`

	// GitMap and CLI-specific common patterns.
	NumberPrefix      = `^(\d+)[-_](.*)$`
	SlugSanitize      = `[^a-z0-9_-]+`
	RepoVersionSuffix = `^(.+)-v(\d+)$`
	RepoVersionCi     = `(?i)^(.*)-v(\d+)$`
	SemverVersion     = `^(v?\d+\.\d+\.\d+([.\-+].+)?|v\+\+|v\+\d+)$`
	SemverLoose       = `v?(\d+\.\d+(?:\.\d+)*)`
	IdentifierExact   = `^[A-Za-z_][A-Za-z0-9_]*$`
	RemoteUrlScp      = `^[^@]+@([^:]+):([^/]+)/(.+)$`
	RemoteUrlHttp     = `^https?://([^/]+)/([^/]+)/(.+)$`
	RemoteUrlSsh      = `^ssh://[^@]+@([^/:]+)(?::\d+)?/([^/]+)/(.+)$`
	AnsiEscape        = `(?:\x1b|\^\[)\[[0-9;]*[a-zA-Z]`
	FileUri           = `file:///[^\x00-\x1f\x7f-\xff"'\s]+`
	SecretPattern     = `\b[A-Z][A-Z0-9_]{2,}(?:KEY|TOKEN|SECRET|PASSWORD|PAT|API)[A-Z0-9_]*\b|\b(?:API|GITMAP|GH|GITHUB|OPENAI)_[A-Z0-9_]+\b`
	MdHeader          = `^(#{1,6})\s+(.+?)\s*$`
	MdTableDivider    = `^\s*\|(\s*:?-{2,}:?\s*\|)+\s*$`
	MdLink            = `\[([^\]]+)\]\(([^)]+)\)`
	MdBold            = `\*\*([^*]+?)\*\*`
	MdInlineCode      = "`([^`]+?)`"

	// Language function declarations (funcintel).
	FuncIntelGo         = `^func\s+(?:\([^)]*\)\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*\(`
	FuncIntelJava       = `^\s*(?:public|protected|private)?\s*(?:static\s+)?[A-Za-z_][A-Za-z0-9_<>\[\]]*\s+([A-Za-z_][A-Za-z0-9_]*)\s*\([^)]*\)\s*\{`
	FuncIntelJavascript = `^(?:export\s+)?(?:async\s+)?function\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\(`
	FuncIntelPhp        = `^(?:(?:public|protected|private)\s+)?function\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`
	FuncIntelPython     = `^def\s+([a-z_][A-Za-z0-9_]*)\s*\(`
	FuncIntelRust       = `^(?:pub(?:\([^)]*\))?\s+)?fn\s+([a-z_][A-Za-z0-9_]*)\s*[<(]`
	FuncIntelTs         = `^(?:export\s+)?const\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*=\s*(?:async\s*)?\(`
)
