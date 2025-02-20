package generictools

const EmptyString = "null"
const EmptyJson = "jsonencode({})"
const EmptyMap = "{}"
const ParentIdentifier = "parent_id"

type VariableInfo struct {
	Description string
	Value       string
}

type EntitlementKey struct {
	ServiceName string
	PlanName    string
}

type RoleKey struct {
	AppId            string
	Name             string
	RoleTemplateName string
}

type BlockSpecifier struct {
	BlockIdentifier string
	ResourceAddress string
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
	SpaceAddress       map[string]string
	EntitlementAddress map[EntitlementKey]string
	RoleAddress        map[RoleKey]string
	DataSourceInfo     []DataSourceInfo
	BlocksToRemove     []BlockSpecifier
}

func NewDepedendcyAddresses() DepedendcyAddresses {
	return DepedendcyAddresses{
		EntitlementAddress: make(map[EntitlementKey]string),
		RoleAddress:        make(map[RoleKey]string),
		SpaceAddress:       make(map[string]string),
	}
}
