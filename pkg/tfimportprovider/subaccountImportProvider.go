package tfimportprovider

import (
	"fmt"
	"log"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
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

func (tf *subaccountImportProvider) GetImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error) {

	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.SubaccountType)
	if err != nil {
		return "", err
	}

	importBlock, err := createSubaccountImportBlock(data, subaccountId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func createSubaccountImportBlock(data map[string]interface{}, subaccountId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {
	if len(filterValues) != 0 {
		if filterValues[0] != fmt.Sprintf("%v", data["name"]) {
			err := fmt.Errorf("subaccount %s not found. Please adjust it in the provided file", filterValues[0])
			log.Println("Error:", err)
			return "", err
		}
	}

	saName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", data["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", saName, -1)
	template = strings.Replace(template, "<subaccount_id>", subaccountId, -1)
	importBlock += template + "\n"

	return importBlock, nil
}
