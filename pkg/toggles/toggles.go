package toggles

import "os"

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
