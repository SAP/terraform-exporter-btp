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
	DirectoyAddress    string
	EntitlementAddress map[EntilementKey]string
}
