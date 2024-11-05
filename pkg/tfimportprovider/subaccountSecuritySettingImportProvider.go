package tfimportprovider

import (
	"fmt"
	"log"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountSecuritySettingImportProvider struct {
	TfImportProvider
}

func newSubaccountSecuritySettingImportProvider() ITfImportProvider {
	return &subaccountSecuritySettingImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountSecuritySettingType,
		},
	}
}

func (tf *subaccountSecuritySettingImportProvider) GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, int, error) {

	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountSecuritySettingType, tfutils.SubaccountLevel)
	if err != nil {
		fmt.Print("\r\n")
		log.Fatalf("read doc failed!")
		return "", 0, err
	}

	importBlock, err := createSecuritySettingImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", 0, err
	}

	// There is only one security setting per level
	return importBlock, 1, nil
}

func createSecuritySettingImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {

	if len(filterValues) != 0 {
		if filterValues[0] != fmt.Sprintf("%v", data["subaccount_id"]) {
			err := fmt.Errorf("security setting %s not found. Please adjust it in the provided file", filterValues[0])
			log.Println("Error:", err)
			return "", err
		}

	}

	template := strings.Replace(resourceDoc.Import, "<resource_name>", "sec_setting", -1)
	template = strings.Replace(template, "'<subaccount_id>'", subaccountId, -1)
	importBlock += template + "\n"

	return importBlock, nil

}
