package tfutils

import "strings"

// EscapeHclString escapes a raw string so it can be embedded safely inside an
// HCL double-quoted string literal. For values that contain no HCL
// metacharacters the result is identical to the input, so generated output is
// unchanged for well-formed BTP/CF identifiers and names; only values that
// contain quotes, backslashes, control characters or interpolation sequences
// are neutralised so they cannot terminate the literal or inject
// interpolation/template expressions.
func EscapeHclString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	// Neutralise HCL interpolation (${...}) and template (%{...}) sequences.
	s = strings.ReplaceAll(s, "${", "$${")
	s = strings.ReplaceAll(s, "%{", "%%{")
	return s
}

// SanitizeHclLabel maps a raw string to a safe HCL identifier usable as a
// resource label, keeping only [A-Za-z0-9_-] and replacing every other
// character with '_'. A leading digit is prefixed with '_' so the result is a
// valid identifier. For labels that already consist of safe characters the
// result is identical to the input.
func SanitizeHclLabel(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_', r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	res := b.String()
	if res != "" && res[0] >= '0' && res[0] <= '9' {
		res = "_" + res
	}
	return res
}
