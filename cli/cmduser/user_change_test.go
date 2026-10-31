package cmduser

import (
	"testing"
)

func TestParseAngleBracketIdentity_Standard(t *testing.T) {
	name, email, hasMatch := parseAngleBracketIdentity("AI automation bot <local@local>")
	if !hasMatch || name != "AI automation bot" || email != "local@local" {
		t.Errorf("unexpected angle bracket parse: %s / %s (match=%v)", name, email, hasMatch)
	}

	name2, email2, hasMatch2 := parseAngleBracketIdentity("\"Jane Doe <jane@example.com>\"")
	if !hasMatch2 || name2 != "Jane Doe" || email2 != "jane@example.com" {
		t.Errorf("unexpected quoted angle parse: %s / %s", name2, email2)
	}
}

func TestParseAngleBracketIdentity_Invalid(t *testing.T) {
	if _, _, hasMatch := parseAngleBracketIdentity("No angle brackets"); hasMatch {
		t.Errorf("expected no match for string without brackets")
	}

	if _, _, hasMatch := parseAngleBracketIdentity("Unclosed <bracket"); hasMatch {
		t.Errorf("expected no match for unclosed bracket")
	}

	if _, _, hasMatch := parseAngleBracketIdentity("Empty <>"); hasMatch {
		t.Errorf("expected no match for empty bracket")
	}
}

func TestResolveDualPositionalIdentity_Standard(t *testing.T) {
	name, email := resolveDualPositionalIdentity("Jane Doe", "jane@example.com")
	if name != "Jane Doe" || email != "jane@example.com" {
		t.Errorf("unexpected standard dual positional: %s / %s", name, email)
	}
}

func TestResolveDualPositionalIdentity_InvertedAndQuoted(t *testing.T) {
	nameInv, emailInv := resolveDualPositionalIdentity("jane@example.com", "Jane Doe")
	if nameInv != "Jane Doe" || emailInv != "jane@example.com" {
		t.Errorf("unexpected inverted dual positional: %s / %s", nameInv, emailInv)
	}

	nameQuoted, emailQuoted := resolveDualPositionalIdentity("\"Bob\"", "'bob@builder.com'")
	if nameQuoted != "Bob" || emailQuoted != "bob@builder.com" {
		t.Errorf("unexpected quoted dual positional: %s / %s", nameQuoted, emailQuoted)
	}
}

func TestResolveSinglePositionalIdentity_EmailAndName(t *testing.T) {
	nameEmail, emailOnly := resolveSinglePositionalIdentity("alice@example.com")
	if nameEmail != "" || emailOnly != "alice@example.com" {
		t.Errorf("unexpected email only: %s / %s", nameEmail, emailOnly)
	}

	nameOnly, emailEmpty := resolveSinglePositionalIdentity("Alice Smith")
	if nameOnly != "Alice Smith" || emailEmpty != "" {
		t.Errorf("unexpected name only: %s / %s", nameOnly, emailEmpty)
	}
}

func TestResolveSinglePositionalIdentity_Angle(t *testing.T) {
	nameAngle, emailAngle := resolveSinglePositionalIdentity("Alice <alice@work.com>")
	if nameAngle != "Alice" || emailAngle != "alice@work.com" {
		t.Errorf("unexpected angle in single positional: %s / %s", nameAngle, emailAngle)
	}
}

func TestParseChangeUserArgs_AngleBracket(t *testing.T) {
	opts, err := ParseChangeUserArgs([]string{"Jane Doe <jane@example.com>", "--global"}, true)
	if err != nil {
		t.Fatalf("unexpected error parsing angle bracket: %v", err)
	}

	if opts.Name != "Jane Doe" || opts.Email != "jane@example.com" {
		t.Errorf("unexpected identity: %s / %s", opts.Name, opts.Email)
	}

	if !opts.IsGlobal || opts.Scope != ScopeGlobal {
		t.Errorf("expected global scope, got: %s", opts.Scope)
	}
}

func TestParseChangeUserArgs_DualPositional(t *testing.T) {
	opts, err := ParseChangeUserArgs([]string{"Jane Doe", "jane@example.com", "--local"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Name != "Jane Doe" || opts.Email != "jane@example.com" {
		t.Errorf("unexpected identity: %s / %s", opts.Name, opts.Email)
	}

	if !opts.IsLocal || opts.Scope != ScopeLocal {
		t.Errorf("expected local scope, got: %s", opts.Scope)
	}
}

func TestParseChangeUserArgs_ExplicitFlags(t *testing.T) {
	opts, err := ParseChangeUserArgs([]string{"--name", "Alice", "--email", "alice@corp.com"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Name != "Alice" || opts.Email != "alice@corp.com" {
		t.Errorf("unexpected identity from flags: %s / %s", opts.Name, opts.Email)
	}
}

func TestParseChangeUserArgs_ShortFlags(t *testing.T) {
	optsShort, errShort := ParseChangeUserArgs([]string{"-n", "Bob", "-e", "bob@corp.com", "-g"}, false)
	if errShort != nil || optsShort.Name != "Bob" || optsShort.Email != "bob@corp.com" || !optsShort.IsGlobal {
		t.Errorf("unexpected short flag result: %v", optsShort)
	}
}

func TestParseChangeUserArgs_EmptyEmail(t *testing.T) {
	opts, err := ParseChangeUserArgs([]string{"Alice", "--empty-email"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Name != "Alice" || opts.Email != "" || !opts.HasEmptyEmail {
		t.Errorf("unexpected empty email options: %+v", opts)
	}
}

func TestParseChangeUserArgs_MissingIdentity(t *testing.T) {
	_, err := ParseChangeUserArgs([]string{})
	if err == nil {
		t.Errorf("expected error for empty args, got nil")
	}
}

func TestParseAutobotArgs_Defaults(t *testing.T) {
	opts, err := ParseAutobotArgs([]string{}, true)
	if err != nil {
		t.Fatalf("unexpected error parsing autobot defaults: %v", err)
	}

	if opts.Name != "GitMap automation" || opts.Email != "info@rgitmap.com" {
		t.Errorf("unexpected autobot defaults: %s / %s", opts.Name, opts.Email)
	}

	if !opts.IsAutobot || opts.Role != "automation agent" {
		t.Errorf("unexpected autobot flags: %+v", opts)
	}
}

func TestParseAutobotArgs_CustomEmail(t *testing.T) {
	opts, err := ParseAutobotArgs([]string{"custom@corp.internal"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Name != "gitmap automation" || opts.Email != "custom@corp.internal" {
		t.Errorf("unexpected custom email autobot: %s / %s", opts.Name, opts.Email)
	}
}

func TestParseAutobotArgs_EmptyEmail(t *testing.T) {
	opts, err := ParseAutobotArgs([]string{"--empty-email"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.Name != "gitmap automation" || opts.Email != "" || !opts.HasEmptyEmail {
		t.Errorf("unexpected empty email autobot: %+v", opts)
	}
}

func TestParseAutobotArgs_ScopeFlags(t *testing.T) {
	optsGlobal, errG := ParseAutobotArgs([]string{"--global"}, true)
	if errG != nil || !optsGlobal.IsGlobal || optsGlobal.Scope != ScopeGlobal {
		t.Errorf("expected global autobot scope: %+v", optsGlobal)
	}

	optsProj, errP := ParseAutobotArgs([]string{"--project"}, true)
	if errP != nil || !optsProj.IsProject || optsProj.Scope != ScopeProject {
		t.Errorf("expected project autobot scope: %+v", optsProj)
	}
}

func TestIsChangeUserCommand(t *testing.T) {
	validChange := []string{"change-user", "change-login", "chuser", "chlogin", "cu", "cl"}
	for _, cmd := range validChange {
		if !isChangeUserCommand(cmd) {
			t.Errorf("expected isChangeUserCommand(%s) = true", cmd)
		}
	}

	if isChangeUserCommand("other") {
		t.Errorf("expected isChangeUserCommand(other) = false")
	}
}

func TestIsAutobotCommand(t *testing.T) {
	validBot := []string{"autobot", "bot"}
	for _, cmd := range validBot {
		if !isAutobotCommand(cmd) {
			t.Errorf("expected isAutobotCommand(%s) = true", cmd)
		}
	}

	if isAutobotCommand("other") {
		t.Errorf("expected isAutobotCommand(other) = false")
	}
}

func TestDispatchUserSubcommand_ChangeRoutes(t *testing.T) {
	errChange := dispatchUserSubcommand("change", []string{"Test User", "test@example.com"})
	if errChange != nil {
		t.Errorf("unexpected error on dispatch change: %v", errChange)
	}

	errBot := dispatchUserSubcommand("autobot", []string{})
	if errBot != nil {
		t.Errorf("unexpected error on dispatch autobot: %v", errBot)
	}
}
