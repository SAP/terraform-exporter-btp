package tfimportprovider

import (
	"fmt"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

func GetImportBlockProvider(cmdResourceName string) (ITfImportProvider, error) {

	switch cmdResourceName {
	case tfutils.CmdSubaccountParameter:
		return newSubaccountImportProvider(), nil
	case tfutils.CmdEntitlementParameter:
		return newSubaccountEntitlementImportProvider(), nil
	case tfutils.CmdEnvironmentInstanceParameter:
		return newSubaccountEnvInstanceImportProvider(), nil
	case tfutils.CmdSubscriptionParameter:
		return newSubaccountSubscriptionImportProvider(), nil
	case tfutils.CmdTrustConfigurationParameter:
		return newSubaccountTrustConfigImportProvider(), nil
	case tfutils.CmdRoleParameter:
		return newSubaccountRoleImportProvider(), nil
	case tfutils.CmdRoleCollectionParameter:
		return newSubaccountRoleCollectionImportProvider(), nil
	case tfutils.CmdServiceInstanceParameter:
		return newSubaccountServiceInstanceImportProvider(), nil
	case tfutils.CmdServiceBindingParameter:
		return newSubaccountServiceBindingImportProvider(), nil
	case tfutils.CmdSecuritySettingParameter:
		return newSubaccountSecuritySettingImportProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported resource provided")
	}

}
