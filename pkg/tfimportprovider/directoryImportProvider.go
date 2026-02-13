package tfimportprovider

import (
	"fmt"
	"log"
	"strings"

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

func (tf *directoryImportProvider) GetImportBlock(data map[string]any, LevelId string, filterValues []string) (string, int, error) {

	directoryId := LevelId
	resourceDoc, err := tfutils.GetDocByResourceName(tfutils.ResourcesKind, tfutils.DirectoryType, tfutils.DirectoryLevel)
	if err != nil {
		return "", 0, err
	}

	importBlock, err := createDirectoryImportBlock(data, directoryId, filterValues, resourceDoc)
	if err != nil {
		return "", 0, err
	}

	//We only export one directory at a time
	return importBlock, 1, nil
}

func createDirectoryImportBlock(data map[string]any, directoryId string, filterValues []string, resourceDoc tfutils.EntityDocs) (importBlock string, err error) {
	if len(filterValues) != 0 {
		if filterValues[0] != fmt.Sprintf("%v", data["name"]) {
			err := fmt.Errorf("directory %s not found. Please adjust it in the provided file", filterValues[0])
			log.Println("Error:", err)
			return "", err
		}
	}

	template := strings.ReplaceAll(resourceDoc.Import, "<resource_name>", "directory_0")
	template = strings.ReplaceAll(template, "<directory_id>", directoryId)
	importBlock += template + "\n"

	return importBlock, nil
}
