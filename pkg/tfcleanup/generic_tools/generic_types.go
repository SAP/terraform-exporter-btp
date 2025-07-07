package generictools

const EmptyString = "null"
const EmptyJson = "jsonencode({})"
const EmptyMap = "{}"
const ParentIdentifier = "parent_id"

type VariableInfo struct {
	Description  string
	DefaultValue string
	Type         string
}

type EntitlementKey struct {
	ServiceName string
	PlanName    string
}

type EntitlementInfo struct {
	Amount  int
	Address string
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

type DependencyAddresses struct {
	SubaccountAddress  string
	DirectoryAddress   string
	SpaceAddress       map[string]string
	EntitlementAddress map[EntitlementKey]EntitlementInfo
	RoleAddress        map[RoleKey]string
	DataSourceInfo     []DataSourceInfo
	BlocksToRemove     []BlockSpecifier
}

func NewDependencyAddresses() DependencyAddresses {
	return DependencyAddresses{
		EntitlementAddress: make(map[EntitlementKey]EntitlementInfo),
		RoleAddress:        make(map[RoleKey]string),
		SpaceAddress:       make(map[string]string),
	}
}
