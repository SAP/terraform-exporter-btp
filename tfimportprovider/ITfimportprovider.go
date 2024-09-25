package tfimportprovider

type ITfImportProvider interface {
	GetResourceType() string
	GetImportBlock(data map[string]interface{}, subaccountId string, filterValues []string) (string, error)
}
