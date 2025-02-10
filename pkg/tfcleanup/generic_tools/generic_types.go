package generictools

const EmptyString = "null"
const EmptyJson = "jsonencode({})"

type VariableInfo struct {
	Description string
	Value       string
}

type EntilementKey struct {
	ServiceName string
	PlanName    string
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
