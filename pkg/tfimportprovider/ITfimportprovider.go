package tfimportprovider

type ITfImportProvider interface {
	GetResourceType() string
	GetImportBlock(data map[string]interface{}, levelId string, filterValues []string) (string, error)
}
