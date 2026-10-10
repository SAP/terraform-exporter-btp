package tfutils

import "testing"

func TestEscapeHclString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "feature-flags-dashboard", "feature-flags-dashboard"},
		{"uuid", "8D363D57-86F1-4FF4-B3F8-C1FB7FFD34AC", "8D363D57-86F1-4FF4-B3F8-C1FB7FFD34AC"},
		{"composite", "1234,feature-flags-dashboard,dashboard", "1234,feature-flags-dashboard,dashboard"},
		{"double_quote", `a"b`, `a\"b`},
		{"backslash", `a\b`, `a\\b`},
		{"newline", "a\nb", `a\n` + "b"},
		{"interpolation", "${var.x}", "$${var.x}"},
		{"template_directive", "%{if true}", "%%{if true}"},
		{"breakout_attempt", `x","y`, `x\",\"y`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeHclString(tt.in); got != tt.want {
				t.Errorf("EscapeHclString(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSanitizeHclLabel(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "cloudfoundry", "cloudfoundry"},
		{"underscores", "role_organization_user_0", "role_organization_user_0"},
		{"guid_prefixed", "role_23456_space_manager_0", "role_23456_space_manager_0"},
		{"space", "Space Developer", "Space_Developer"},
		{"leading_digit", "123foo", "_123foo"},
		{"special", `a.b"c`, "a_b_c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeHclLabel(tt.in); got != tt.want {
				t.Errorf("SanitizeHclLabel(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
