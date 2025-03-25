package tfimportprovider

import (
	"fmt"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

func GetImportBlockProvider(cmdResourceName string, level string) (ITfImportProvider, error) {

	switch cmdResourceName {
	case tfutils.CmdSubaccountParameter:
		return newSubaccountImportProvider(), nil
	case tfutils.CmdEntitlementParameter:
		return getEntitlementImportProviderByLevel(level)
	case tfutils.CmdEnvironmentInstanceParameter:
		return newSubaccountEnvInstanceImportProvider(), nil
	case tfutils.CmdSubscriptionParameter:
		return newSubaccountSubscriptionImportProvider(), nil
	case tfutils.CmdTrustConfigurationParameter:
		return newSubaccountTrustConfigImportProvider(), nil
	case tfutils.CmdRoleParameter:
		return getRoleImportProviderByLevel(level)
	case tfutils.CmdRoleCollectionParameter:
		return getRoleCollectionImportProviderByLevel(level)
	case tfutils.CmdServiceInstanceParameter:
		return newSubaccountServiceInstanceImportProvider(), nil
	case tfutils.CmdServiceBindingParameter:
		return newSubaccountServiceBindingImportProvider(), nil
	case tfutils.CmdSecuritySettingParameter:
		return newSubaccountSecuritySettingImportProvider(), nil
	case tfutils.CmdDirectoryParameter:
		return newDirectoryImportProvider(), nil
	case tfutils.CmdCfSpaceParameter:
		return newcloudfoundrySpaceImportProvider(), nil
	case tfutils.CmdCfUserParameter:
		return newcloudfoundryUserImportProvider(), nil
	case tfutils.CmdCfDomainParamater:
		return newCloudfoundryDomainImportProvider(), nil
	case tfutils.CmdCfOrgRoleParameter:
		return newCloudfoundryOrgRolesImportProvider(), nil
	case tfutils.CmdCfRouteParameter:
		return newCloudfoundryRouteImportProvider(), nil
	case tfutils.CmdCfSpaceQuotaParameter:
		return newcloudfoundrySpaceQuotaImportProvider(), nil
	case tfutils.CmdCfServiceInstanceParameter:
		return newCloudfoundryServiceInstanceImportProvider(), nil
	case tfutils.CmdCfSpaceRoleParameter:
		return newCloudfoundrySpaceRolesImportProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported resource provided")
	}

}

func getEntitlementImportProviderByLevel(level string) (ITfImportProvider, error) {
	switch level {
	case tfutils.SubaccountLevel:
		return newSubaccountEntitlementImportProvider(), nil
	case tfutils.DirectoryLevel:
		return newDirectoryEntitlementImportProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported level provided")
	}
}

func getRoleImportProviderByLevel(level string) (ITfImportProvider, error) {
	switch level {
	case tfutils.SubaccountLevel:
		return newSubaccountRoleImportProvider(), nil
	case tfutils.DirectoryLevel:
		return newDirectoryRoleImportProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported level provided")
	}
}

func getRoleCollectionImportProviderByLevel(level string) (ITfImportProvider, error) {
	switch level {
	case tfutils.SubaccountLevel:
		return newSubaccountRoleCollectionImportProvider(), nil
	case tfutils.DirectoryLevel:
		return newDirectoryRoleCollectionImportProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported level provided")
	}
}
