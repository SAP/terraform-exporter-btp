package tfimportprovider

type TfImportProvider struct {
	resourceType string
}

func (tf *TfImportProvider) GetResourceType() string {
	return tf.resourceType
}
