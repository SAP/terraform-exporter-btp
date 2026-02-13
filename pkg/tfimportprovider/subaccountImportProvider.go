package tfimportprovider

import (
	"fmt"
	"log"
	"strings"

	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type subaccountImportProvider struct {
	TfImportProvider
}

func newSubaccountImportProvider() ITfImportProvider {
	return &subaccountImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.SubaccountType,
		},
	}
}

func (tf *subaccountImportProvider) GetImportBlock(data map[string]any, levelId string, filterValues []string) (string, int, error) {

	subaccountId := levelId

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountType, tfutils.SubaccountLevel)
	if err != nil {
		return "", 0, err
	}

	importBlock, err := createSubaccountImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", 0, err
	}
	//We only export one subaccount at a time
	return importBlock, 1, nil
}

func createSubaccountImportBlock(data map[string]any, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {
	if len(filterValues) != 0 {
		if filterValues[0] != fmt.Sprintf("%v", data["name"]) {
			err := fmt.Errorf("subaccount %s not found. Please adjust it in the provided file", filterValues[0])
			log.Println("Error:", err)
			return "", err
		}
	}

	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "subaccount_0")
	template = strings.ReplaceAll(template, "<subaccount_id>", subaccountId)
	importBlock += template + "\n"

	return importBlock, nil
}
