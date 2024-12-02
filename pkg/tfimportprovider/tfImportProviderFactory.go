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
	case tfutils.CmdCfRouteParameter:
		return newCloudfoundryRouteImportProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported resource provided")
	}

}

func getEntitlementImportProviderByLevel(level string) (ITfImportProvider, error) {
	if level == tfutils.SubaccountLevel {
		return newSubaccountEntitlementImportProvider(), nil
	} else if level == tfutils.DirectoryLevel {
		return newDirectoryEntitlementImportProvider(), nil
	} else {
		return nil, fmt.Errorf("unsupported level provided")
	}
}

func getRoleImportProviderByLevel(level string) (ITfImportProvider, error) {
	if level == tfutils.SubaccountLevel {
		return newSubaccountRoleImportProvider(), nil
	} else if level == tfutils.DirectoryLevel {
		return newDirectoryRoleImportProvider(), nil
	} else {
		return nil, fmt.Errorf("unsupported level provided")
	}
}

func getRoleCollectionImportProviderByLevel(level string) (ITfImportProvider, error) {
	if level == tfutils.SubaccountLevel {
		return newSubaccountRoleCollectionImportProvider(), nil
	} else if level == tfutils.DirectoryLevel {
		return newDirectoryRoleCollectionImportProvider(), nil
	} else {
		return nil, fmt.Errorf("unsupported level provided")
	}
}
