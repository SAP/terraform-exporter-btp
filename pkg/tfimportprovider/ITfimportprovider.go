package tfimportprovider

type ITfImportProvider interface {
	GetResourceType() string
	GetImportBlock(data map[string]any, levelId string, filterValues []string) (string, int, error)
}
