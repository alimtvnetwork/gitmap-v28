package cmduser

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// UserChangeScopeType defines the configuration boundary.
type UserChangeScopeType string

const (
	ScopeLocal   UserChangeScopeType = "local"
	ScopeGlobal  UserChangeScopeType = "global"
	ScopeProject UserChangeScopeType = "project"
)

// ParsedIdentity holds extracted author credentials.
type ParsedIdentity struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HasName        bool   `json:"hasName"`
	HasEmail       bool   `json:"hasEmail"`
	IsAngleBracket bool   `json:"isAngleBracket"`
}

// UserChangeOptions specifies parameters for identity modification.
type UserChangeOptions struct {
	Name             string              `json:"name"`
	Email            string              `json:"email"`
	Scope            UserChangeScopeType `json:"scope"`
	Role             string              `json:"role,omitempty"`
	ProfileAlias     string              `json:"profileAlias,omitempty"`
	IsGlobal         bool                `json:"isGlobal"`
	IsLocal          bool                `json:"isLocal"`
	IsProject        bool                `json:"isProject"`
	IsAutobot        bool                `json:"isAutobot"`
	HasEmptyEmail    bool                `json:"hasEmptyEmail"`
	HasExplicitName  bool                `json:"hasExplicitName"`
	HasExplicitEmail bool                `json:"hasExplicitEmail"`
}

// ChangeUserOptions is an alias for UserChangeOptions.
type ChangeUserOptions = UserChangeOptions

type buildUserChangeInput struct {
	PosName       string
	PosEmail      string
	FlagName      string
	FlagEmail     string
	Scope         UserChangeScopeType
	IsAutobot     bool
	HasEmptyEmail bool
}

func parseAngleBracketIdentity(raw string) (string, string, bool) {
	trimmed := strings.Trim(strings.TrimSpace(raw), `"'`)
	openIdx := strings.IndexByte(trimmed, '<')
	closeIdx := strings.LastIndexByte(trimmed, '>')
	if openIdx < 0 || closeIdx <= openIdx+1 {
		return "", "", false
	}

	name := strings.Trim(strings.TrimSpace(trimmed[:openIdx]), `"'`)
	email := strings.Trim(strings.TrimSpace(trimmed[openIdx+1:closeIdx]), `"'`)

	return name, email, true
}

func stripWrappingQuotes(s string) string {
	trimmed := strings.TrimSpace(s)
	hasDouble := len(trimmed) >= 2 && trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"'
	if hasDouble {
		return strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}

	hasSingle := len(trimmed) >= 2 && trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\''
	if hasSingle {
		return strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}

	return trimmed
}

func resolveDualPositionalIdentity(first, second string) (string, string) {
	cleanFirst := stripWrappingQuotes(first)
	cleanSecond := stripWrappingQuotes(second)
	hasFirstAt := strings.Contains(cleanFirst, "@")
	hasSecondAt := strings.Contains(cleanSecond, "@")
	if hasFirstAt && !hasSecondAt {
		return cleanSecond, cleanFirst
	}

	return cleanFirst, cleanSecond
}

func resolveSinglePositionalIdentity(arg string) (string, string) {
	clean := stripWrappingQuotes(arg)
	if name, email, hasAngle := parseAngleBracketIdentity(clean); hasAngle {
		return name, email
	}

	if strings.Contains(clean, "@") {
		return "", clean
	}

	return clean, ""
}

func resolveMultiPositionalIdentity(pos []string) (string, string) {
	lastIdx := len(pos) - 1
	if strings.Contains(pos[lastIdx], "@") {
		return strings.Join(pos[:lastIdx], " "), pos[lastIdx]
	}

	if strings.Contains(pos[0], "@") {
		return strings.Join(pos[1:], " "), pos[0]
	}

	return strings.Join(pos, " "), ""
}

func resolvePositionalIdentityFromArgs(pos []string) (string, string) {
	if len(pos) == 0 {
		return "", ""
	}

	if name, email, hasAngle := checkJoinedAngleBracket(pos); hasAngle {
		return name, email
	}

	return resolveDiscretePositionals(pos)
}

func checkJoinedAngleBracket(pos []string) (string, string, bool) {
	joined := strings.Join(pos, " ")

	return parseAngleBracketIdentity(joined)
}

func resolveDiscretePositionals(pos []string) (string, string) {
	if len(pos) == 1 {
		return resolveSinglePositionalIdentity(pos[0])
	}

	if len(pos) == 2 {
		return resolveDualPositionalIdentity(pos[0], pos[1])
	}

	return resolveMultiPositionalIdentity(pos)
}

func isValuedFlag(arg string) bool {
	return arg == "--name" || arg == "-n" || arg == "--email" || arg == "-e"
}

func isIgnoredSubcommandToken(arg string) bool {
	lower := strings.ToLower(arg)

	return lower == "change" || lower == "change-user" || lower == "change-login" || lower == "autobot" || lower == "bot"
}

func isChangePositionalToken(arg string) bool {
	return !strings.HasPrefix(arg, "-") && !isIgnoredSubcommandToken(arg)
}

func extractChangePositionalArgs(args []string) []string {
	var pos []string
	hasSkip := false
	for _, a := range args {
		pos, hasSkip = stepChangePositionalArg(pos, a, hasSkip)
	}

	return pos
}

func stepChangePositionalArg(pos []string, a string, hasSkip bool) ([]string, bool) {
	if hasSkip {
		return pos, false
	}
	if isValuedFlag(a) {
		return pos, true
	}
	if isChangePositionalToken(a) {
		return append(pos, a), false
	}

	return pos, false
}

func extractFlagValue(args []string, longFlag, shortFlag string) string {
	for i := 0; i < len(args); i++ {
		isMatch := args[i] == longFlag || (shortFlag != "" && args[i] == shortFlag)
		if isMatch && i+1 < len(args) {
			return stripWrappingQuotes(args[i+1])
		}
	}

	return ""
}

func extractFlagValueWithEqual(args []string, longFlag, shortFlag string) string {
	val := extractFlagValue(args, longFlag, shortFlag)
	if len(val) > 0 {
		return val
	}

	prefix := longFlag + "="
	for _, a := range args {
		if strings.HasPrefix(a, prefix) {
			return stripWrappingQuotes(strings.TrimPrefix(a, prefix))
		}
	}

	return ""
}

func hasAnyFlag(args []string, flags ...string) bool {
	for _, a := range args {
		for _, f := range flags {
			if a == f {
				return true
			}
		}
	}

	return false
}

func resolveScopeFromFlags(args []string, isInsideRepo bool) UserChangeScopeType {
	if scope, hasExplicit := resolveExplicitScope(args); hasExplicit {
		return scope
	}

	return resolveInferredScope(isInsideRepo)
}

func resolveExplicitScope(args []string) (UserChangeScopeType, bool) {
	if hasAnyFlag(args, "--global", "-g") {
		return ScopeGlobal, true
	}
	if hasAnyFlag(args, "--project", "-p") {
		return ScopeProject, true
	}
	if hasAnyFlag(args, "--local", "-l") {
		return ScopeLocal, true
	}

	return ScopeLocal, false
}

func resolveInferredScope(isInsideRepo bool) UserChangeScopeType {
	if isInsideRepo {
		return ScopeLocal
	}

	return ScopeGlobal
}

func detectInsideRepo(override ...bool) bool {
	if len(override) > 0 {
		return override[0]
	}
	if hasGitDirectory() {
		return true
	}

	return probeGitRevParseWorkTree()
}

func hasGitDirectory() bool {
	_, err := os.Stat(".git")

	return err == nil
}

func probeGitRevParseWorkTree() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

// ParseChangeUserArgs parses one-liner arguments for user change.
func ParseChangeUserArgs(args []string, isInside ...bool) (UserChangeOptions, error) {
	in := collectChangeUserInput(args, isInside...)

	return buildUserChangeOptions(in)
}

func collectChangeUserInput(args []string, isInside ...bool) buildUserChangeInput {
	scope := resolveScopeFromFlags(args, detectInsideRepo(isInside...))
	posName, posEmail := resolvePositionalIdentityFromArgs(extractChangePositionalArgs(args))

	return buildUserChangeInput{
		PosName:       posName,
		PosEmail:      posEmail,
		FlagName:      extractFlagValueWithEqual(args, "--name", "-n"),
		FlagEmail:     extractFlagValueWithEqual(args, "--email", "-e"),
		Scope:         scope,
		IsAutobot:     hasAnyFlag(args, "--autobot", "-b", "autobot", "bot"),
		HasEmptyEmail: hasAnyFlag(args, "--empty-email", "--no-email"),
	}
}

func buildUserChangeOptions(in buildUserChangeInput) (UserChangeOptions, error) {
	if in.IsAutobot {
		return buildAutobotChangeOptions(in), nil
	}

	return buildStandardUserChangeOptions(in)
}

func buildStandardUserChangeOptions(in buildUserChangeInput) (UserChangeOptions, error) {
	name, email := resolveEffectiveChangePair(in)
	if len(name) == 0 && len(email) == 0 && !in.HasEmptyEmail {
		return UserChangeOptions{}, apperror.WrapSimple(os.ErrInvalid, "missing identity: specify name, email, or --autobot")
	}

	return constructUserChangeOptions(name, email, in.Scope, false, in.HasEmptyEmail, ""), nil
}

func resolveEffectiveChangePair(in buildUserChangeInput) (string, string) {
	name := resolveEffectiveIdentityField(in.FlagName, in.PosName)
	if in.HasEmptyEmail {
		return name, ""
	}

	return name, resolveEffectiveIdentityField(in.FlagEmail, in.PosEmail)
}

func resolveEffectiveIdentityField(flagVal, posVal string) string {
	if len(flagVal) > 0 {
		return flagVal
	}

	return posVal
}

func constructUserChangeOptions(name, email string, scope UserChangeScopeType, isAutobot, hasEmptyEmail bool, role string) UserChangeOptions {
	opts := UserChangeOptions{
		Name: name, Email: email, Scope: scope, Role: role,
		ProfileAlias: resolveProfileAlias(isAutobot),
		IsGlobal:     scope == ScopeGlobal, IsLocal: scope == ScopeLocal, IsProject: scope == ScopeProject,
		IsAutobot: isAutobot, HasEmptyEmail: hasEmptyEmail,
		HasExplicitName: len(name) > 0, HasExplicitEmail: len(email) > 0,
	}

	return opts
}

func resolveProfileAlias(isAutobot bool) string {
	if isAutobot {
		return "autobot"
	}

	return ""
}

func buildAutobotChangeOptions(in buildUserChangeInput) UserChangeOptions {
	name := resolveAutobotName(in)
	email := resolveAutobotEmail(in)
	hasEmpty := in.HasEmptyEmail || (len(email) == 0 && in.HasEmptyEmail)

	return constructUserChangeOptions(name, email, in.Scope, true, hasEmpty, "automation agent")
}

func resolveAutobotName(in buildUserChangeInput) string {
	if len(in.FlagName) > 0 {
		return in.FlagName
	}

	if len(in.PosName) > 0 && !strings.Contains(in.PosName, "@") && in.PosName != "autobot" && in.PosName != "bot" {
		return in.PosName
	}

	return "AI automation bot"
}

func resolveAutobotEmail(in buildUserChangeInput) string {
	if in.HasEmptyEmail {
		return ""
	}

	if explicitEmail := resolveExplicitAutobotEmail(in); len(explicitEmail) > 0 {
		return explicitEmail
	}

	return "local@local"
}

func resolveExplicitAutobotEmail(in buildUserChangeInput) string {
	if len(in.FlagEmail) > 0 {
		return in.FlagEmail
	}
	if len(in.PosEmail) > 0 {
		return in.PosEmail
	}
	if strings.Contains(in.PosName, "@") {
		return in.PosName
	}

	return ""
}

// ParseAutobotArgs parses options tailored for the autobot command.
func ParseAutobotArgs(args []string, isInside ...bool) (UserChangeOptions, error) {
	newArgs := append([]string{"--autobot"}, args...)

	return ParseChangeUserArgs(newArgs, isInside...)
}

func applyUserChangeIdentity(opts UserChangeOptions) error {
	if err := applyUserChangeName(opts); err != nil {
		return err
	}

	return applyUserChangeEmail(opts)
}

func applyUserChangeName(opts UserChangeOptions) error {
	if len(opts.Name) == 0 {
		return nil
	}

	return writeGitConfigKey("user.name", opts.Name, opts.IsGlobal)
}

func applyUserChangeEmail(opts UserChangeOptions) error {
	if opts.HasEmptyEmail || len(opts.Email) == 0 {
		return unsetGitConfigKey("user.email", opts.IsGlobal)
	}

	return writeGitConfigKey("user.email", opts.Email, opts.IsGlobal)
}

func writeGitConfigKey(key, value string, isGlobal bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	args := buildConfigSetArgs(key, value, isGlobal)
	out, err := exec.CommandContext(ctx, "git", args...).CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, "git config set failed: "+strings.TrimSpace(string(out)))
	}

	return nil
}

func buildConfigSetArgs(key, value string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", key, value}
	}

	return []string{"config", key, value}
}

func unsetGitConfigKey(key string, isGlobal bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	args := buildConfigUnsetArgs(key, isGlobal)
	_ = exec.CommandContext(ctx, "git", args...).Run()

	return nil
}

func buildConfigUnsetArgs(key string, isGlobal bool) []string {
	if isGlobal {
		return []string{"config", "--global", "--unset", key}
	}

	return []string{"config", "--unset", key}
}

func renderChangeSuccess(opts UserChangeOptions) {
	scopeName := resolveScopeDisplayName(opts)
	if opts.HasEmptyEmail || len(opts.Email) == 0 {
		fmt.Printf("✔ Changed %s Git author identity to: %s (email unset)\n", scopeName, opts.Name)
	} else {
		fmt.Printf("✔ Changed %s Git author identity to: %s <%s>\n", scopeName, opts.Name, opts.Email)
	}

	if len(opts.Role) > 0 {
		fmt.Printf("  Role: %s\n", opts.Role)
	}
}

func resolveScopeDisplayName(opts UserChangeOptions) string {
	if opts.IsGlobal {
		return "Global"
	}

	if opts.IsProject {
		return "Project"
	}

	return "Local"
}

// RunChangeUser changes the Git user identity based on arguments.
func RunChangeUser(args []string) error {
	opts, err := ParseChangeUserArgs(args)
	if err != nil {
		return logUserChangeError("Error parsing user change options: %v\n", err)
	}

	if err := applyUserChangeIdentity(opts); err != nil {
		return logUserChangeError("Error updating Git author identity: %v\n", err)
	}

	renderChangeSuccess(opts)

	return nil
}

// RunUserChange is an alias for RunChangeUser.
func RunUserChange(args []string) error {
	return RunChangeUser(args)
}

// RunChangeLogin is an alias for RunChangeUser.
func RunChangeLogin(args []string) error {
	return RunChangeUser(args)
}

// RunUserAutobot switches to the AI automation bot identity.
func RunUserAutobot(args []string) error {
	opts, err := ParseAutobotArgs(args)
	if err != nil {
		return logUserChangeError("Error parsing autobot options: %v\n", err)
	}

	if err := applyUserChangeIdentity(opts); err != nil {
		return logUserChangeError("Error applying autobot identity: %v\n", err)
	}

	renderChangeSuccess(opts)

	return nil
}

// RunAutobot is an alias for RunUserAutobot.
func RunAutobot(args []string) error {
	return RunUserAutobot(args)
}

func logUserChangeError(format string, err error) error {
	fmt.Fprintf(os.Stderr, format, err)

	return err
}
