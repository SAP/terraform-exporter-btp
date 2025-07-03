package toggles

import (
	"os"
	"testing"
)

func GetIacTool() string {
	return os.Getenv("BTPTF_IAC_TOOL")
}

func IsCodeCleanupDeactivated() bool {
	return os.Getenv("BTPTF_SKIP_CODECLEANUP") != ""
}

func IsRoleCollectionFilterDeactivated() bool {
	return os.Getenv("BTPTF_SKIP_RCFILTER") != ""
}

func IsRoleFilterDeactivated() bool {
	return os.Getenv("BTPTF_SKIP_ROLEFILTER") != ""
}

func IsEntitlementFilterDeactivated() bool {
	return os.Getenv("BTPTF_SKIP_ENTITLEMENTFILTER") != ""
}

func IsEntitlementModuleGenerationDeactivated() bool {
	// For Unit tests we rely on the existing logic
	if testing.Testing() {
		return true
	}

	return os.Getenv("BTPTF_ADD_ENTITLEMENTMODULE") == ""
}
