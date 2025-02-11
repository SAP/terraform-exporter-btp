package generictools

const EmptyString = "null"
const EmptyJson = "jsonencode({})"
const EmptyMap = "{}"

type VariableInfo struct {
	Description string
	Value       string
}

type EntilementKey struct {
	ServiceName string
	PlanName    string
}

type LevelIds struct {
	SubaccountId string
	DirectoryId  string
	CfOrgId      string
}

type VariableContent map[string]VariableInfo

type DepedendcyAddresses struct {
	SubaccountAddress  string
	DirectoryAddress   string
	EntitlementAddress map[EntilementKey]string
}

func NewDepedendcyAddresses() DepedendcyAddresses {
	return DepedendcyAddresses{
		EntitlementAddress: make(map[EntilementKey]string),
	}
}
