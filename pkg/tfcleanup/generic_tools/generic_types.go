package generictools

const EmptyString = "null"
const EmptyJson = "jsonencode({})"
const EmptyMap = "{}"

type VariableInfo struct {
	Description string
	Value       string
}

type EntitlementKey struct {
	ServiceName string
	PlanName    string
}

type DataSourceInfo struct {
	DatasourceAddress  string
	SubaccountAddress  string
	OfferingName       string
	Name               string
	EntitlementAddress string
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
	EntitlementAddress map[EntitlementKey]string
	DataSourceInfo     []DataSourceInfo
}

func NewDepedendcyAddresses() DepedendcyAddresses {
	return DepedendcyAddresses{
		EntitlementAddress: make(map[EntitlementKey]string),
	}
}
