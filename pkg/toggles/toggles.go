package toggles

import (
	"os"
	"testing"
)

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

func IsPlanUniqueIdentifierRemovalDeactivated() bool {
	return os.Getenv("BTPTF_SKIP_PLAN_UNIQUE_IDENTIFIER_REMOVAL") != ""
}
