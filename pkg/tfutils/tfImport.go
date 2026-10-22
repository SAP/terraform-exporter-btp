package tfutils

const BtpProviderVersion = "v1.27.0"
const CfProviderVersion = "v1.18.0"

const (
	SubaccountLevel   = "subaccountLevel"
	DirectoryLevel    = "directoryLevel"
	OrganizationLevel = "organizationLevel"
	SpaceLevel        = "spaceLevel"
)

const (
	CmdDirectoryParameter           = "directory"
	CmdSubaccountParameter          = "subaccount"
	CmdEntitlementParameter         = "entitlements"
	CmdEnvironmentInstanceParameter = "environment-instances"
	CmdSubscriptionParameter        = "subscriptions"
	CmdTrustConfigurationParameter  = "trust-configurations"
	CmdRoleParameter                = "roles"
	CmdRoleCollectionParameter      = "role-collections"
	CmdServiceInstanceParameter     = "service-instances"
	CmdServiceBindingParameter      = "service-bindings"
	CmdSecuritySettingParameter     = "security-settings"
	CmdCfSpaceParameter             = "spaces"
	CmdCfUserParameter              = "users"
	CmdCfDomainParamater            = "domains"
	CmdCfOrgRoleParameter           = "org-roles"
	CmdCfRouteParameter             = "routes"
	CmdCfSpaceQuotaParameter        = "space-quotas"
	CmdCfServiceInstanceParameter   = "cf-service-instances"
	CmdCfSpaceRoleParameter         = "space-roles"
)

const (
	SubaccountType                    = "btp_subaccount"
	SubaccountEntitlementType         = "btp_subaccount_entitlement"
	SubaccountEnvironmentInstanceType = "btp_subaccount_environment_instance"
	SubaccountSubscriptionType        = "btp_subaccount_subscription"
	SubaccountTrustConfigurationType  = "btp_subaccount_trust_configuration"
	SubaccountRoleType                = "btp_subaccount_role"
	SubaccountRoleCollectionType      = "btp_subaccount_role_collection"
	SubaccountServiceInstanceType     = "btp_subaccount_service_instance"
	SubaccountServiceBindingType      = "btp_subaccount_service_binding"
	SubaccountSecuritySettingType     = "btp_subaccount_security_settings"
	DirectoryType                     = "btp_directory"
	DirectoryEntitlementType          = "btp_directory_entitlement"
	DirectoryRoleType                 = "btp_directory_role"
	DirectoryRoleCollectionType       = "btp_directory_role_collection"
	CfSpaceType                       = "cloudfoundry_space"
	CfUserType                        = "cloudfoundry_user"
	CfUserExportType                  = "cloudfoundry_user_cf"
	CfOrgRoleType                     = "cloudfoundry_org_role"
	CfDomainType                      = "cloudfoundry_domain"
	CfRouteType                       = "cloudfoundry_route"
	CfSpaceQuotaType                  = "cloudfoundry_space_quota"
	CfServiceInstanceType             = "cloudfoundry_service_instance"
	CfSpaceRoleType                   = "cloudfoundry_space_role"
)

const (
	DirectoryFeatureDefault      = "DEFAULT"
	DirectoryFeatureEntitlements = "ENTITLEMENTS"
	DirectoryFeatureRoles        = "AUTHORIZATIONS"
)

type BtpResource struct {
	Name   string
	Values []string
}
type BtpResources struct{ BtpResources []BtpResource }

func GetExecutionLevelAndId(subaccountID, directoryID, organizationID, spaceID string) (string, string) {
	if subaccountID != "" {
		return SubaccountLevel, subaccountID
	}
	if directoryID != "" {
		return DirectoryLevel, directoryID
	}
	if organizationID != "" {
		if spaceID != "" {
			return SpaceLevel, spaceID
		}
		return OrganizationLevel, organizationID
	}
	return "", ""
}
