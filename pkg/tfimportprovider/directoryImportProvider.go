package tfimportprovider

import (
	"fmt"
	"log"
	"strings"

	"github.com/SAP/terraform-exporter-btp/pkg/output"
	tfutils "github.com/SAP/terraform-exporter-btp/pkg/tfutils"
)

type directoryImportProvider struct {
	TfImportProvider
}

func newDirectoryImportProvider() ITfImportProvider {
	return &directoryImportProvider{
		TfImportProvider: TfImportProvider{
			resourceType: tfutils.DirectoryType,
		},
	}
}

func (tf *directoryImportProvider) GetImportBlock(data map[string]interface{}, LevelId string, filterValues []string) (string, error) {

	directoryId := LevelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.DirectoryType)
	if err != nil {
		return "", err
	}

	importBlock, err := createDirectoryImportBlock(data, directoryId, filterValues, resourceDoc)
	if err != nil {
		return "", err
	}

	return importBlock, nil
}

func createDirectoryImportBlock(data map[string]interface{}, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {
	if len(filterValues) != 0 {
		if filterValues[0] != fmt.Sprintf("%v", data["name"]) {
			err := fmt.Errorf("directory %s not found. Please adjust it in the provided file", filterValues[0])
			log.Println("Error:", err)
			return "", err
		}
	}

	saName := output.FormatResourceNameGeneric(fmt.Sprintf("%v", data["name"]))
	template := strings.Replace(resourceDoc.Import, "<resource_name>", saName, -1)
	template = strings.Replace(template, "<directory_id>", directoryId, -1)
	importBlock += template + "\n"

	return importBlock, nil
}
